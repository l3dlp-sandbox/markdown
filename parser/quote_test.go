package parser

import "testing"

func parseWithFlags(input string, flags Flags) string {
	p := NewWithExtensions(CommonExtensions)
	p.Opts.Flags = flags
	return astToString(p.Parse([]byte(input)))
}

func TestBlockQuoteBlankLine(t *testing.T) {
	input := "> test\n\n> test 2\n"
	classic := "BlockQuote\n  Paragraph\n    Text 'test'\n  Paragraph\n    Text 'test 2'\n"
	split := "BlockQuote\n  Paragraph\n    Text 'test'\nBlockQuote\n  Paragraph\n    Text 'test 2'\n"

	if got := parseWithFlags(input, FlagsNone); got != classic {
		t.Errorf("default\nExpected[%#v]\nGot     [%#v]\n%s", classic, got, got)
	}
	if got := parseWithFlags(input, CommonMark); got != split {
		t.Errorf("commonmark\nExpected[%#v]\nGot     [%#v]\n%s", split, got, got)
	}
}

func TestBlockQuoteCommonMarkParagraphs(t *testing.T) {
	// A '>' on the blank line keeps both paragraphs in one quote.
	input := "> test\n>\n> test 2\n"
	want := "BlockQuote\n  Paragraph\n    Text 'test'\n  Paragraph\n    Text 'test 2'\n"
	if got := parseWithFlags(input, CommonMark); got != want {
		t.Errorf("\nExpected[%#v]\nGot     [%#v]\n%s", want, got, got)
	}
}

func TestBlockQuoteCommonMarkLazyContinuation(t *testing.T) {
	input := "> foo\nbar\n"
	want := "BlockQuote\n  Paragraph\n    Text 'foo\\nbar'\n"
	if got := parseWithFlags(input, CommonMark); got != want {
		t.Errorf("\nExpected[%#v]\nGot     [%#v]\n%s", want, got, got)
	}
}

func TestBlockQuoteCommonMarkNested(t *testing.T) {
	input := "> > foo\n>\n> > bar\n"
	want := "BlockQuote\n  BlockQuote\n    Paragraph\n      Text 'foo'\n  BlockQuote\n    Paragraph\n      Text 'bar'\n"
	if got := parseWithFlags(input, CommonMark); got != want {
		t.Errorf("\nExpected[%#v]\nGot     [%#v]\n%s", want, got, got)
	}
}

func TestBlockQuoteCommonMarkFence(t *testing.T) {
	input := "> ```\n> a\n>\n> b\n> ```\n"
	want := "BlockQuote\n  CodeBlock: 'a\\n\\nb\\n'\n"
	if got := parseWithFlags(input, CommonMark); got != want {
		t.Errorf("\nExpected[%#v]\nGot     [%#v]\n%s", want, got, got)
	}
}
