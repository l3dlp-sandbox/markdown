package parser

import "sort"

// bracketTable maps each '[' in one Inline() buffer to its matching ']'.
// Built once per buffer so a run of unmatched '[' is O(n) (GHSA-85vw-wvf9-r522).
type bracketTable struct {
	data   []byte
	close  []int
	nested []bool
	emph   [3][]int
}

func (t *bracketTable) lookup(open int) (closeAt int, nested, ok bool) {
	t.ensure()
	if open < 0 || open >= len(t.close) {
		return 0, false, false
	}
	closeAt = t.close[open]
	if closeAt < 0 {
		return 0, false, false
	}
	return closeAt, t.nested[open], true
}

// firstEmphasis returns the first emphasis delimiter in [start, end). The
// positions are collected while building the bracket table, so an unmatched
// '[' does not require another scan to the end of the inline buffer.
func (t *bracketTable) firstEmphasis(start, end int, c byte) (int, bool) {
	t.ensure()
	var slot int
	switch c {
	case '*':
		slot = 0
	case '_':
		slot = 1
	case '~':
		slot = 2
	default:
		return 0, false
	}
	positions := t.emph[slot]
	i := sort.SearchInts(positions, start)
	if i == len(positions) || positions[i] >= end {
		return 0, false
	}
	return positions[i], true
}

func (t *bracketTable) ensure() {
	if t.close != nil {
		return
	}
	n := len(t.data)
	closeAt := make([]int, n)
	nested := make([]bool, n)
	for i := 0; i < n; i++ {
		closeAt[i] = -1
	}
	stack := make([]int, 0, 16)
	for i := 0; i < n; i++ {
		switch t.data[i] {
		case '*':
			t.emph[0] = append(t.emph[0], i)
		case '_':
			t.emph[1] = append(t.emph[1], i)
		case '~':
			t.emph[2] = append(t.emph[2], i)
		}
		// an odd run of backslashes escapes the bracket; an even run is
		// escaped backslashes and leaves the bracket live
		if isEscape(t.data, i) {
			continue
		}
		switch t.data[i] {
		case '[':
			stack = append(stack, i)
		case ']':
			if len(stack) == 0 {
				continue
			}
			open := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			closeAt[open] = i
			if len(stack) > 0 {
				nested[stack[len(stack)-1]] = true
			}
		}
	}
	t.close = closeAt
	t.nested = nested
}
