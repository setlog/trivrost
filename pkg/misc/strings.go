package misc

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"log"
)

const ellipsis = "…"

// WordWrap performs word-wrapping on space-separated words within text, keeping existing newline
// characters intact. Any words which exceed lineWidth in length will move to a new line, but words
// themself will never be split.
func WordWrap(text string, lineWidth int) string {
	lines := strings.Split(text, "\n")
	wrappedString := ""
	lastLineIndex := len(lines) - 1
	for i, line := range lines {
		wrappedString += WordWrapIgnoreNewLine(line, lineWidth)
		if i != lastLineIndex {
			wrappedString += "\n"
		}
	}
	return wrappedString
}

// WordWrapIgnoreNewLine performs word-wrapping on space-separated words within text, whereas any
// newline characters are treated as spaces themselves. Any words which exceed lineWidth in
// length will move to a new line, but words themself will never be split.
func WordWrapIgnoreNewLine(text string, lineWidth int) string {
	words := strings.FieldsFunc(text, isBreakingSpace)
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - utf8.RuneCountInString(wrapped)
	for _, word := range words[1:] {
		wordRuneCount := utf8.RuneCountInString(word)
		if wordRuneCount >= spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - wordRuneCount
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + wordRuneCount
		}
	}

	return wrapped
}

func isBreakingSpace(r rune) bool {
	return unicode.IsSpace(r) && r != 0xA0
}

func ShortString(s string, leadingCount int, trailingCount int) string {
	maxLen := leadingCount + trailingCount + utf8.RuneCountInString(ellipsis)
	runeCount := utf8.RuneCountInString(s)
	if runeCount > maxLen {
		firstByteIndex := RuneIndexToByteIndex(s, leadingCount)
		omitByteIndex := RuneIndexToByteIndex(s, runeCount-trailingCount)
		s = s[:firstByteIndex] + ellipsis + s[omitByteIndex:]
	}
	return s
}

func RuneIndexToByteIndex(s string, runeIndex int) int {
	currentRuneIndex := 0
	for i := range s {
		if currentRuneIndex == runeIndex {
			return i
		}
		currentRuneIndex++
	}
	if currentRuneIndex == runeIndex {
		return len(s)
	}
	return -1
}

func ExtensionlessFileName(filePath string) string {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	return fileName[:len(fileName)-len(ext)]
}

func StripTrailingLineBreak(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s[:len(s)-1]
	}
	return s
}

func SplitTrailing(s string, trailSet string) (lead, trail string) {
	splitIndex := -1
	for i, r := range s {
		if strings.ContainsRune(trailSet, r) {
			if splitIndex == -1 {
				splitIndex = i
			}
		} else {
			splitIndex = -1
		}
	}
	if splitIndex >= 0 {
		return s[:splitIndex], s[splitIndex:]
	}
	return s, ""
}

func JoinNonEmpty(stringList []string, sep string) string {
	candidates := make([]string, 0, len(stringList))
	for _, s := range stringList {
		if s != "" {
			candidates = append(candidates, s)
		}
	}
	return strings.Join(candidates, sep)
}

func RemoveLines(s string, from, to int) (string, error) {
	if from < 0 {
		return "", fmt.Errorf("from (%d) < 0", from)
	}
	if to < from {
		return "", fmt.Errorf("to (%d) < from (%d)", to, from)
	}
	lines := strings.SplitN(s, "\n", -1)
	if to > len(lines) {
		return "", fmt.Errorf("to (%d) > linecount (%d)", to, len(lines))
	}
	return strings.Join(append(lines[:from], lines[to:]...), "\n"), nil
}

func TryRemoveLines(s string, from, to int) string {
	result, err := RemoveLines(s, from, to)
	if err != nil {
		log.Printf("Could not remove lines %d to %d from string \"%s\": %v\n", from, to, s, err)
		return s
	}
	return result
}
