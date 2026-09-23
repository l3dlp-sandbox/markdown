package parser

import "testing"

func TestHelperFindEmphCharUsesBracketTable(t *testing.T) {
	buf := []byte("prefix [text*] tail*")
	p := &Parser{brackets: bracketTable{data: buf}}
	data := buf[len("prefix "):]

	if got := helperFindEmphChar(p, data, '*'); got != len("[text") {
		t.Fatalf("matched bracket: got %d, want %d", got, len("[text"))
	}

	buf = []byte("prefix [text*")
	p.brackets = bracketTable{data: buf}
	data = buf[len("prefix "):]
	if got := helperFindEmphChar(p, data, '*'); got != len("[text") {
		t.Fatalf("unmatched bracket: got %d, want %d", got, len("[text"))
	}
}
