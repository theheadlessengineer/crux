package template

import (
	"io/fs"

	domain "github.com/theheadlessengineer/crux/internal/domain/template"
)

// NewFromFS is exported for testing only.
func NewFromFS(fsys fs.FS) (domain.Engine, error) {
	return newFromFS(fsys)
}
