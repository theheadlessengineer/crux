package prompt

import "errors"

// ErrAtFirstQuestion is returned when the user tries to navigate back from the first question.
var ErrAtFirstQuestion = errors.New("already at the first question, cannot go back further")

type historyEntry struct {
	question *Question
	answer   Answer
}

// Session drives the full prompt flow: ordered questions, history, back navigation.
type Session struct {
	graph   *DecisionGraph
	history []historyEntry
	answers map[string]Answer
}

// NewSession creates a Session for the given graph.
func NewSession(graph *DecisionGraph) *Session {
	return &Session{graph: graph, answers: make(map[string]Answer)}
}

// Answers returns a copy of the current answer map.
func (s *Session) Answers() map[string]Answer {
	out := make(map[string]Answer, len(s.answers))
	for k, v := range s.answers {
		out[k] = v
	}
	return out
}

// NextQuestion returns the next unanswered visible question, or nil when done.
func (s *Session) NextQuestion() *Question {
	for i := range s.graph.questions {
		q := &s.graph.questions[i]
		if _, answered := s.answers[q.ID]; answered {
			continue
		}
		if s.graph.IsVisible(q, s.answers) {
			return q
		}
	}
	return nil
}

// Record stores an answer and pushes it onto the history stack.
func (s *Session) Record(q *Question, a Answer) {
	s.answers[q.ID] = a
	s.history = append(s.history, historyEntry{question: q, answer: a})
}

// Back removes the last answer and clears any answers that are no longer visible.
// Returns ErrAtFirstQuestion if history is empty.
func (s *Session) Back() error {
	if len(s.history) == 0 {
		return ErrAtFirstQuestion
	}
	last := s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]
	delete(s.answers, last.question.ID)
	s.clearHidden()
	return nil
}

func (s *Session) clearHidden() {
	for id := range s.answers {
		idx, ok := s.graph.index[id]
		if !ok {
			continue
		}
		if !s.graph.IsVisible(&s.graph.questions[idx], s.answers) {
			delete(s.answers, id)
		}
	}
	filtered := s.history[:0]
	for i := range s.history {
		if _, ok := s.answers[s.history[i].question.ID]; ok {
			filtered = append(filtered, s.history[i])
		}
	}
	s.history = filtered
}
