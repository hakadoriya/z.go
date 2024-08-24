package simplediffz

import (
	"strings"
)

// DefaultSeparator is the separator of lines.
//
//nolint:gochecknoglobals
var DefaultSeparator = "\n"

type (
	// DiffOperation represents a single operation in a diff.
	DiffOperation struct {
		Op   string // "+" for add, "-" for delete, " " for no change
		Text string
	}
	// DiffResult represents a collection of DiffOperation.
	DiffResult struct {
		Separator string
		Ops       []DiffOperation
	}
)

type (
	diffConfig struct {
		separator string
	}
	DiffOption interface {
		apply(c *diffConfig)
	}
)

type withDiffOptionSeparator struct{ separator string }

func (o withDiffOptionSeparator) apply(c *diffConfig) {
	c.separator = o.separator
}

// WithDiffOptionSeparator sets the separator of lines.
func WithDiffOptionSeparator(separator string) DiffOption {
	return withDiffOptionSeparator{separator: separator}
}

// Diff returns the diff between two strings.
func Diff(before, after string, opts ...DiffOption) *DiffResult {
	config := &diffConfig{separator: DefaultSeparator}
	for _, opt := range opts {
		opt.apply(config)
	}

	diffOps := diff(strings.Split(before, config.separator), strings.Split(after, config.separator))

	return &DiffResult{
		Separator: config.separator,
		Ops:       diffOps,
	}
}

// String returns the string representation of the diff result.
func (r *DiffResult) String() string {
	var result strings.Builder
	for i := range r.Ops {
		_, _ = result.WriteString(r.Ops[i].Op + r.Ops[i].Text + r.Separator)
	}
	return result.String()
}

//nolint:cyclop
func diff(a, b []string) []DiffOperation {
	m := len(a)
	n := len(b)
	diffs := []DiffOperation{}

	// Create a 2D slice to store the edit distance between slices
	edits := make([][]int, m+1)
	for i := range edits {
		edits[i] = make([]int, n+1)
	}

	// Fill the table
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			switch {
			case i == 0:
				edits[i][j] = j
			case j == 0:
				edits[i][j] = i
			case a[i-1] == b[j-1]:
				edits[i][j] = edits[i-1][j-1]
			default:
				edits[i][j] = min(edits[i-1][j]+1, edits[i][j-1]+1)
			}
		}
	}

	// Backtrack to find the diff
	for i, j := m, n; i > 0 || j > 0; {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			diffs = append([]DiffOperation{{" ", a[i-1]}}, diffs...)
			i--
			j--
		case j > 0 && (i == 0 || edits[i][j-1] <= edits[i-1][j]):
			diffs = append([]DiffOperation{{"+", b[j-1]}}, diffs...)
			j--
		case i > 0 && (j == 0 || edits[i][j-1] > edits[i-1][j]):
			diffs = append([]DiffOperation{{"-", a[i-1]}}, diffs...)
			i--
		}
	}

	return diffs
}
