package gui

import (
	"github.com/andlabs/ui"
	"github.com/setlog/trivrost/pkg/misc"
	"strings"
)

func textBox(box *ui.Box, message string, lineWidth int) (*ui.Box, int) {
	message = misc.WordWrap(message, lineWidth) // The ui library itself wraps nothing, not even spaces, resulting in very wide windows (>10000 pixels) without this.
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		box.Append(ui.NewLabel(line), false)
	}
	return box, len(lines)
}

func clearBox(box *ui.Box, n int) {
	for i := n - 1; i >= 0; i-- {
		box.Delete(i)
	}
}
