package misc_test

import (
	"testing"

	"github.com/setlog/trivrost/pkg/misc"
)

func TestWordWrapIgnoreNewLine(t *testing.T) {
	tests := []struct {
		text      string
		lineWidth int
		expected  string
	}{
		{"hello world", 100, "hello world"},
		{"hello world", 11, "hello world"},
		{"hello world", 10, "hello\nworld"},
		{"hello world", 7, "hello\nworld"},
		{"hello world", 6, "hello\nworld"},
		{"hello world", 5, "hello\nworld"},
		{"hello world", 4, "hello\nworld"},
		{"日 本 語 日 本 語", 6, "日 本 語\n日 本 語"},
		{"日 本 語 日 本 語", 5, "日 本 語\n日 本 語"},
		{"日 本 語 日 本 語", 4, "日 本\n語 日\n本 語"},
		{"日 本 語 日 本 語", 3, "日 本\n語 日\n本 語"},
		{"日\n本 語\n日 本 語", 5, "日 本 語\n日 本 語"},
	}
	for i, test := range tests {
		result := misc.WordWrapIgnoreNewLine(test.text, test.lineWidth)
		if result != test.expected {
			t.Errorf("Test #%d failed: WordWrapIgnoreNewLine(\"%s\", %d) yielded %s. Expected %s.", i+1, test.text, test.lineWidth, result, test.expected)
		}
	}
}

func TestRemoveLines(t *testing.T) {
	tests := []struct {
		text        string
		from, to    int
		expected    string
		expectedErr bool
	}{
		{"what\na\ndrag", 1, 2, "what\ndrag", false},
		{"what\na\ndrag", 1, 3, "what", false},
		{"what\na\ndrag", 1, 4, "", true},
	}
	for i, test := range tests {
		result, err := misc.RemoveLines(test.text, test.from, test.to)
		if result != test.expected {
			t.Errorf("Test #%d failed: RemoveLines(\"%s\", %d, %d) yielded %s. Expected %s.", i+1, test.text, test.from, test.to, result, test.expected)
		}
		if (err == nil) == test.expectedErr {
			t.Errorf("Test #%d failed: RemoveLines(\"%s\", %d, %d) yielded error %v.", i+1, test.text, test.from, test.to, err)
		}
	}
}

func TestShortString(t *testing.T) {
	tests := []struct {
		text                        string
		leadingCount, trailingCount int
		expected                    string
	}{
		{"what a stupid text", 7, 5, "what a ... text"},
		{"what a great text", 7, 7, "what a great text"},
		{"what a great text", 6, 7, "what a...at text"},
		{"日本語偽善者", 1, 2, "日本語偽善者"},
		{"日本語偽善者", 1, 1, "日...者"},
		{"ab日本語皮を被る", 1, 2, "a...被る"},
		{"", 0, 0, ""},
		{"0123456789", 4, 0, "0123..."},
		{"0123456789", 0, 4, "...6789"},
	}
	for i, test := range tests {
		result := misc.ShortString(test.text, test.leadingCount, test.trailingCount)
		if result != test.expected {
			t.Errorf("Test #%d failed: ShortString(\"%s\", %d, %d) yielded %s. Expected %s.", i+1, test.text, test.leadingCount, test.trailingCount, result, test.expected)
		}
	}
}

func TestSplitTrailing(t *testing.T) {
	tests := []struct {
		text, trailSet              string
		expectedLead, expectedTrail string
	}{
		{"Heeey\n\n\n", "\n", "Heeey", "\n\n\n"},
		{"Heeey\n\n\n", "", "Heeey\n\n\n", ""},
		{"Heeey", "\n", "Heeey", ""},
		{"Heeey", "", "Heeey", ""},
		{"\n\n\n", "\n", "", "\n\n\n"},
		{"\n\n\n", "", "\n\n\n", ""},
		{"", "\n", "", ""},
		{"", "", "", ""},
		{"THEQUICKBROWNFOX", "OWXFBRN", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "XBFRWNO", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "WXRONBF", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "OWNRFXB", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "OWXOFBRN", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "XBFORWNO", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "WXROONBF", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "OWNORFXB", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "OWXOFOBRN", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "XBFOROWNO", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "WXROOONBF", "THEQUICK", "BROWNFOX"},
		{"THEQUICKBROWNFOX", "OWNOROFXB", "THEQUICK", "BROWNFOX"},
	}
	for i, test := range tests {
		lead, trail := misc.SplitTrailing(test.text, test.trailSet)
		if lead != test.expectedLead || trail != test.expectedTrail {
			t.Errorf("Test #%d: misc.SplitTrailing(%v, %v) return %v, %v. Expected: %v, %v", i+1, test.text, test.trailSet, lead, trail, test.expectedLead, test.expectedTrail)
		}
	}
}
