package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
)

// ErrAborted is returned by Run when the user presses ctrl+c.
var ErrAborted = errors.New("aborted by user")

// Model is the Bubbletea model for the crux new prompt flow.
type Model struct {
	session        *prompt.Session
	serviceName    string
	current        *prompt.Question
	input          string
	cursor         int
	selected       map[int]bool
	err            string
	done           bool
	aborted        bool
	confirmingQuit bool
	reviewMode     bool // true when all questions answered — shows summary screen
	editPickMode   bool // true when user is picking which answer to edit
	reviewCursor   int  // cursor in review action list or edit-pick list
	showHelp       bool
	theme          Theme
	styles         Styles
	width          int
	height         int
	answers        map[string]prompt.Answer
}

// NewModel creates a TUI model for the given session and service name.
func NewModel(session *prompt.Session, serviceName string) Model {
	theme := LoadThemePreference()
	m := Model{
		session:     session,
		serviceName: serviceName,
		theme:       theme,
		styles:      theme.BuildStyles(),
		selected:    make(map[int]bool),
		answers:     make(map[string]prompt.Answer),
		width:       80,
		height:      24,
	}
	m.current = session.NextQuestion()
	return m
}

// Answers returns the collected answers after the model is done.
func (m *Model) Answers() map[string]prompt.Answer { return m.answers }

// Done reports whether all questions have been answered.
func (m *Model) Done() bool { return m.done || m.reviewMode }

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys processed before mode checks.
	switch msg.String() {
	case "ctrl+c":
		m.confirmingQuit = true
		return m, nil
	case "?":
		if !m.isTextInput() {
			m.showHelp = !m.showHelp
			return m, nil
		}
	case "t":
		if !m.isTextInput() {
			next := NextTheme(&m.theme)
			m.theme = next
			m.styles = next.BuildStyles()
			SaveThemePreference(&next)
			return m, nil
		}
	}

	if m.confirmingQuit {
		return m.handleQuitConfirm(msg)
	}
	if m.reviewMode {
		return m.updateReview(msg)
	}
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}
	if m.current == nil {
		return m, tea.Quit
	}

	switch m.current.Type {
	case prompt.QuestionTypeText, prompt.QuestionTypeNumber:
		return m.updateTextInput(msg)
	case prompt.QuestionTypeConfirm:
		return m.updateConfirm(msg)
	case prompt.QuestionTypeSelect:
		return m.updateSelect(msg)
	case prompt.QuestionTypeMultiSelect:
		return m.updateMultiSelect(msg)
	}
	return m, nil
}

func (m Model) handleQuitConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.aborted = true
		return m, tea.Quit
	default:
		m.confirmingQuit = false
	}
	return m, nil
}

func (m Model) isTextInput() bool {
	if m.current == nil {
		return false
	}
	return m.current.Type == prompt.QuestionTypeText || m.current.Type == prompt.QuestionTypeNumber
}

func (m Model) updateTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		raw := m.input
		if raw == "" {
			raw = m.current.Default
		}
		answer, err := prompt.Validate(m.current, raw)
		if err != nil {
			m.err = err.Error()
			return m, nil
		}
		m.err = ""
		m = m.recordAnswer(answer)
	case "backspace":
		if m.input != "" {
			m.input = m.input[:len(m.input)-1]
		}
		m.err = ""
	case "b":
		if m.input == "" {
			var err error
			m, err = m.goBack()
			if err != nil {
				return m, nil
			}
			return m, nil
		}
		m.input += "b"
		m.err = ""
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
			m.err = ""
		}
	}
	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		answer, _ := prompt.Validate(m.current, "y")
		m = m.recordAnswer(answer)
	case "n", "N":
		answer, _ := prompt.Validate(m.current, "n")
		m = m.recordAnswer(answer)
	case "enter":
		answer, _ := prompt.Validate(m.current, m.current.Default)
		m = m.recordAnswer(answer)
	case "b", "backspace":
		m, _ = m.goBack()
	}
	return m, nil
}

func (m Model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	opts := m.current.Options
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(opts)-1 {
			m.cursor++
		}
	case "enter":
		if len(opts) > 0 {
			answer, _ := prompt.Validate(m.current, opts[m.cursor].Value)
			m = m.recordAnswer(answer)
		}
	case "b", "backspace":
		m, _ = m.goBack()
	}
	return m, nil
}

func (m Model) updateMultiSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	opts := m.current.Options
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(opts)-1 {
			m.cursor++
		}
	case " ":
		m.selected[m.cursor] = !m.selected[m.cursor]
	case "enter":
		var chosen []string
		for i, o := range opts {
			if m.selected[i] {
				chosen = append(chosen, o.Value)
			}
		}
		answer, _ := prompt.Validate(m.current, strings.Join(chosen, ","))
		m = m.recordAnswer(answer)
	case "b", "backspace":
		m, _ = m.goBack()
	}
	return m, nil
}

func (m Model) recordAnswer(answer prompt.Answer) Model {
	m.session.Record(m.current, answer)
	m.answers[m.current.ID] = answer
	m.input = ""
	m.cursor = 0
	m.selected = make(map[int]bool)
	m.err = ""
	m.current = m.session.NextQuestion()
	if m.current == nil {
		m.reviewMode = true
		m.reviewCursor = 0
	}
	return m
}

func (m Model) goBack() (Model, error) {
	if err := m.session.Back(); err != nil {
		m.err = err.Error()
		return m, err
	}
	for k := range m.answers {
		if _, ok := m.session.Answers()[k]; !ok {
			delete(m.answers, k)
		}
	}
	m.input = ""
	m.cursor = 0
	m.selected = make(map[int]bool)
	m.err = ""
	m.current = m.session.NextQuestion()
	return m, nil
}

// reviewActions are the fixed options shown at the bottom of the review screen.
var reviewActions = []string{"Confirm & create service", "Change plugin selection", "Edit an answer"}

func (m Model) updateReview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	answered := m.answeredQuestions()
	if m.editPickMode {
		return m.updateEditPick(msg, answered)
	}
	switch msg.String() {
	case "up", "k":
		if m.reviewCursor > 0 {
			m.reviewCursor--
		}
	case "down", "j":
		if m.reviewCursor < len(reviewActions)-1 {
			m.reviewCursor++
		}
	case "enter":
		switch m.reviewCursor {
		case 0:
			m.done = true
			return m, tea.Quit
		case 1:
			m.reviewMode = false
			m = m.navigateTo("_plugins")
		case 2:
			m.editPickMode = true
			m.reviewCursor = 0
		}
	}
	return m, nil
}

func (m Model) updateEditPick(msg tea.KeyMsg, answered []prompt.Question) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.reviewCursor > 0 {
			m.reviewCursor--
		}
	case "down", "j":
		if m.reviewCursor < len(answered)-1 {
			m.reviewCursor++
		}
	case "enter":
		if m.reviewCursor < len(answered) {
			m.editPickMode = false
			m.reviewMode = false
			m = m.navigateTo(answered[m.reviewCursor].ID)
		}
	case "esc", "b":
		m.editPickMode = false
		m.reviewCursor = 0
	}
	return m, nil
}

func (m Model) answeredQuestions() []prompt.Question {
	questions := m.session.Graph().Questions()
	answered := make([]prompt.Question, 0, len(questions))
	for _, q := range questions {
		if _, ok := m.answers[q.ID]; ok {
			answered = append(answered, q)
		}
	}
	return answered
}

// navigateTo walks the session back until the question with the given ID is current.
func (m Model) navigateTo(id string) Model {
	for {
		if next := m.session.NextQuestion(); next != nil && next.ID == id {
			m.current = next
			break
		}
		if err := m.session.Back(); err != nil {
			m.current = m.session.NextQuestion()
			break
		}
		for k := range m.answers {
			if _, ok := m.session.Answers()[k]; !ok {
				delete(m.answers, k)
			}
		}
	}
	m.input = ""
	m.cursor = 0
	m.selected = make(map[int]bool)
	m.err = ""
	return m
}

const (
	headerHeight    = 3 // title + separator + blank
	statusBarHeight = 1
	rightColWidth   = 36 // answers panel fixed width
	panelPadding    = 2  // lipgloss border takes 2 cols each side
)

// View implements tea.Model.
func (m Model) View() string {
	if m.done {
		return ""
	}
	if m.confirmingQuit {
		return m.renderQuitConfirm()
	}
	if m.reviewMode {
		return m.renderReview()
	}
	if m.showHelp {
		return m.renderHelp()
	}

	// ── Row 1: header bar ────────────────────────────────────────────────────
	header := m.renderHeaderBar()

	// ── Row 2: 2-column body ────────────────────────────────────────────────
	bodyHeight := m.height - headerHeight - statusBarHeight - 2 // 2 for panel borders
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	leftW := m.width - rightColWidth - 1 // -1 for join gap
	if leftW < 20 {
		leftW = 20
	}

	leftPanel := m.renderQuestionPanel(leftW, bodyHeight)
	rightPanel := m.renderAnswersPanel(rightColWidth, bodyHeight)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// ── Row 3: status bar ────────────────────────────────────────────────────
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, statusBar)
}

// renderHeaderBar renders the top bar: title + progress + theme name.
func (m Model) renderHeaderBar() string {
	answered, total := m.session.Progress()
	progress := ""
	if total > 0 {
		progress = fmt.Sprintf(" (%d/%d)", answered+1, total)
	}

	s := m.styles
	left := s.Title.Render("  crux") +
		s.Muted.Render(" — "+m.serviceName) +
		s.Accent.Render(progress)

	right := s.Muted.Render("theme: ") + s.Accent.Render(m.theme.Name) +
		s.Muted.Render("  [t] cycle  [?] help")

	// Pad between left and right
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	line1 := left + strings.Repeat(" ", gap) + right
	line2 := s.Border.Render(strings.Repeat("─", m.width))

	style := lipgloss.NewStyle()
	if m.theme.BG != "" {
		style = style.Background(m.theme.BG)
	}
	return style.Render(line1) + "\n" + line2
}

// renderQuestionPanel renders the left bordered panel with the current question.
func (m Model) renderQuestionPanel(w, h int) string {
	innerW := w - panelPadding
	if innerW < 10 {
		innerW = 10
	}

	var b strings.Builder
	if m.current != nil {
		q := m.current
		b.WriteString(m.styles.Prompt.Render("? "+q.Prompt) + "\n")
		if q.Default != "" {
			b.WriteString(m.styles.Muted.Render("  default: "+q.Default) + "\n")
		}
		b.WriteString("\n")
		b.WriteString(m.renderQuestionInput(q))
		if m.err != "" {
			b.WriteString("\n" + m.styles.Error.Render("  ✘ "+m.err) + "\n")
		}
		if q.Help != "" {
			b.WriteString("\n" + m.styles.Border.Render(strings.Repeat("─", innerW)) + "\n")
			b.WriteString(wrapText(q.Help, innerW, m.styles.Muted))
		}
	}

	title := m.styles.Title.Render("Question")
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Width(innerW).
		Height(h).
		Padding(0, 1).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(m.theme.Title).Render(title) + "\n\n" + b.String(),
		)
}

// renderQuestionInput renders the interactive input area for the current question type.
func (m Model) renderQuestionInput(q *prompt.Question) string {
	var b strings.Builder
	switch q.Type {
	case prompt.QuestionTypeText, prompt.QuestionTypeNumber:
		display := m.input
		if display == "" {
			display = m.styles.Muted.Render(q.Default)
		}
		b.WriteString(m.styles.Cursor.Render("❯ ") + display + m.styles.Cursor.Render("█") + "\n")

	case prompt.QuestionTypeConfirm:
		yStyle, nStyle := m.styles.Muted, m.styles.Muted
		if q.Default == "y" {
			yStyle = m.styles.Selected
		} else {
			nStyle = m.styles.Selected
		}
		b.WriteString(yStyle.Render("  [ Y ] Yes") + "   " + nStyle.Render("[ N ] No") + "\n")

	case prompt.QuestionTypeSelect:
		for i, o := range q.Options {
			if i == m.cursor {
				b.WriteString(m.styles.Cursor.Render("❯ ") + m.styles.Selected.Render(o.Label) + "\n")
			} else {
				b.WriteString("  " + m.styles.Muted.Render(o.Label) + "\n")
			}
		}

	case prompt.QuestionTypeMultiSelect:
		count := 0
		for _, v := range m.selected {
			if v {
				count++
			}
		}
		b.WriteString(m.styles.Accent.Render(
			fmt.Sprintf("  %d selected — space to toggle, enter to confirm", count),
		) + "\n\n")
		for i, o := range q.Options {
			check := m.styles.Muted.Render("○")
			if m.selected[i] {
				check = m.styles.Success.Render("●")
			}
			pre, lbl := "  ", m.styles.Muted.Render(o.Label)
			if i == m.cursor {
				pre = m.styles.Cursor.Render("❯ ")
				lbl = m.styles.Prompt.Render(o.Label)
			}
			b.WriteString(pre + check + " " + lbl + "\n")
		}
	}
	return b.String()
}

// renderAnswersPanel renders the right bordered panel with all answers so far.
func (m Model) renderAnswersPanel(w, h int) string {
	innerW := w - panelPadding
	if innerW < 8 {
		innerW = 8
	}

	var b strings.Builder
	if len(m.answers) == 0 {
		b.WriteString(m.styles.Muted.Render("  no answers yet") + "\n")
	} else {
		// Walk questions in order so answers appear in question order.
		for _, q := range m.session.Graph().Questions() {
			a, ok := m.answers[q.ID]
			if !ok {
				continue
			}
			key := m.styles.Muted.Render(truncate(q.ID, innerW-2))
			val := lipgloss.NewStyle().Foreground(m.theme.Accent).
				Render(truncate(fmt.Sprintf("%v", a.Value), innerW-2))
			b.WriteString(key + "\n")
			b.WriteString("  " + val + "\n\n")
		}
	}

	title := m.styles.Title.Render("answers")
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Width(innerW).
		Height(h).
		Padding(0, 1).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(m.theme.Title).Render(title) + "\n\n" + b.String(),
		)
}

// renderReview renders the full-screen review/confirmation screen.
func (m Model) renderReview() string {
	s := m.styles
	questions := m.session.Graph().Questions()

	answered := make([]prompt.Question, 0, len(questions))
	for _, q := range questions {
		if _, ok := m.answers[q.ID]; ok {
			answered = append(answered, q)
		}
	}

	var b strings.Builder
	b.WriteString(s.Title.Render("  Review your configuration") + "\n")
	b.WriteString(s.Border.Render("  "+strings.Repeat("─", m.width-4)) + "\n\n")

	for i, q := range answered {
		a := m.answers[q.ID]
		var val string
		switch v := a.Value.(type) {
		case []string:
			if len(v) == 0 {
				val = s.Muted.Render("(none)")
			} else {
				val = s.Accent.Render(strings.Join(v, ", "))
			}
		default:
			val = s.Accent.Render(fmt.Sprintf("%v", v))
		}
		if m.editPickMode && i == m.reviewCursor {
			b.WriteString(s.Cursor.Render("  ❯ ") + s.Selected.Render(q.Prompt+": ") + val + "\n")
		} else {
			b.WriteString("    " + s.Muted.Render(q.Prompt+": ") + val + "\n")
		}
	}

	b.WriteString("\n" + s.Border.Render("  "+strings.Repeat("─", m.width-4)) + "\n\n")

	if m.editPickMode {
		b.WriteString("  " + s.Prompt.Render("Select a question to edit — enter to confirm, esc to cancel") + "\n")
	} else {
		for i, action := range reviewActions {
			if i == m.reviewCursor {
				b.WriteString(s.Cursor.Render("  ❯ ") + s.Selected.Render(action) + "\n")
			} else {
				b.WriteString("    " + s.Muted.Render(action) + "\n")
			}
		}
	}

	var statusKeys string
	if m.editPickMode {
		statusKeys = " ↑/↓ select answer  enter edit  esc cancel"
	} else {
		statusKeys = " ↑/↓ navigate  enter select  ctrl+c abort"
	}
	statusBar := s.StatusBar.Width(m.width).Render(statusKeys)
	body := lipgloss.NewStyle().
		Background(m.theme.BG).
		Width(m.width).
		Height(m.height - statusBarHeight).
		Render(b.String())
	return lipgloss.JoinVertical(lipgloss.Left, body, statusBar)
}

// renderQuitConfirm renders a centered confirmation overlay.
func (m Model) renderQuitConfirm() string {
	msg := m.styles.Prompt.Render("Abort service creation?") + "\n\n" +
		m.styles.Error.Render("  [ Y ] Yes, abort") + "   " +
		m.styles.Success.Render("[ any key ] Continue") + "\n"
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Warning).
		Padding(1, 3).
		Render(msg)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// renderStatusBar renders the full-width inverted status bar.
func (m Model) renderStatusBar() string {
	var keys string
	switch {
	case m.current != nil && m.current.Type == prompt.QuestionTypeMultiSelect:
		keys = " ↑/↓ navigate  space toggle  enter confirm  b back  t theme  ctrl+c quit"
	case m.current != nil && m.isTextInput():
		keys = " type answer  enter confirm  backspace delete  b back (when empty)  ctrl+c quit"
	default:
		keys = " ↑/↓ navigate  enter confirm  b back  t theme  ? help  ctrl+c quit"
	}
	return m.styles.StatusBar.Width(m.width).Render(keys)
}

// renderHelp renders the full-screen help overlay.
func (m Model) renderHelp() string {
	s := m.styles
	lines := []string{
		s.Title.Render("  Keyboard Shortcuts"),
		s.Border.Render("  " + strings.Repeat("─", 36)),
		"",
		s.Prompt.Render("  Navigation"),
		s.Help.Render("  ↑ / k          move up"),
		s.Help.Render("  ↓ / j          move down"),
		s.Help.Render("  b / backspace  go back to previous question"),
		s.Help.Render("  enter          confirm answer"),
		"",
		s.Prompt.Render("  Multi-select"),
		s.Help.Render("  space          toggle selection"),
		s.Help.Render("  enter          confirm selections"),
		"",
		s.Prompt.Render("  Global"),
		s.Help.Render("  t              cycle theme (" + strings.Join(ThemeNames(), " → ") + ")"),
		s.Help.Render("  ?              toggle this help"),
		s.Help.Render("  ctrl+c         quit"),
		"",
		s.Muted.Render("  Press any key to close"),
	}
	content := strings.Join(lines, "\n")
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Width(m.width-2).
		Padding(1, 2).
		Render(content)
	return lipgloss.JoinVertical(lipgloss.Left, panel, m.renderStatusBar())
}

// truncate shortens s to max runes, appending "…" if needed.
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}

// wrapText wraps text to maxWidth columns, rendering each line with style.
func wrapText(text string, maxWidth int, style lipgloss.Style) string {
	if maxWidth < 10 {
		return style.Render(text) + "\n"
	}
	words := strings.Fields(text)
	var b strings.Builder
	line := ""
	for _, w := range words {
		if line == "" {
			line = w
		} else if len(line)+1+len(w) <= maxWidth {
			line += " " + w
		} else {
			b.WriteString(style.Render("  "+line) + "\n")
			line = w
		}
	}
	if line != "" {
		b.WriteString(style.Render("  "+line) + "\n")
	}
	return b.String()
}

// Run executes the TUI for the given session and returns the collected answers.
// Returns ErrAborted if the user pressed ctrl+c.
func Run(session *prompt.Session, serviceName string) (map[string]prompt.Answer, error) {
	m := NewModel(session, serviceName)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	fm, ok := final.(Model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}
	if fm.aborted {
		return nil, ErrAborted
	}
	return fm.session.Answers(), nil
}

// IsTTY reports whether stdin is an interactive terminal.
func IsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
