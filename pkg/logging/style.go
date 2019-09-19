package logging

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type withoutStyleWriter struct {
	target    io.Writer
	isCutting bool
}

func (dw *withoutStyleWriter) Write(p []byte) (n int, err error) {
	decolored := make([]byte, 0)
	for _, b := range p {
		if b == '\x1B' {
			dw.isCutting = true
			continue
		} else if b == 'm' {
			if dw.isCutting {
				dw.isCutting = false
				continue
			}
		}
		if !dw.isCutting {
			decolored = append(decolored, b)
		}
	}

	n, err = dw.target.Write(decolored)
	if n == len(decolored) && err == nil {
		n = len(p) // This is a lie, but it is required to satisfy the io.Writer interface.
	}
	return
}

func colorize(text string, colorCode color.Attribute, clearAfter bool) string {
	if text == "" {
		if clearAfter {
			return clearStyle()
		}
		return text
	}
	if clearAfter {
		return fmt.Sprintf("%s%s%s", formatStyle(colorCode), text, clearStyle())
	}
	return fmt.Sprintf("%s%s", formatStyle(colorCode), text)
}

func formatStyle(colorCode color.Attribute) string {
	return fmt.Sprintf("\x1B[%dm", colorCode)
}

func clearStyle() string {
	return fmt.Sprintf("\x1B[%dm", color.Reset)
}
