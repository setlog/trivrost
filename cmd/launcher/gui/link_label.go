package gui

import (
	"github.com/andlabs/ui"
)

// See https://docs.microsoft.com/en-us/windows/win32/uxguide/vis-fonts#fonts-and-colors
var linkColor = ui.TextColor{R: 0, G: float64(0x66) / 255.0, B: float64(0xCC) / 255.0, A: 1}
var linkColorHover = ui.TextColor{R: float64(0x33) / 255.0, G: float64(0x99) / 255.0, B: float64(0xFF) / 255.0, A: 1}

func newLinkLabel(labelText string, onClickFunc func()) *ui.Area {
	attrString := ui.NewAttributedString(labelText)
	attrString.SetAttribute(linkColor, 0, len(labelText))
	return ui.NewArea(&linkAreaHandler{attributedString: attrString, defaultFont: getDefaultFont(), onClickFunc: onClickFunc})
}

type linkAreaHandler struct {
	attributedString *ui.AttributedString
	defaultFont      *ui.FontDescriptor
	onClickFunc      func()
}

func getDefaultFont() *ui.FontDescriptor {
	// Ideally, the current platform's default text font. However, right now this caters to Windows.
	return &ui.FontDescriptor{
		Family:  ui.TextFamily("Segoe UI"),
		Size:    ui.TextSize(9),
		Weight:  ui.TextWeightNormal,
		Italic:  ui.TextItalicNormal,
		Stretch: ui.TextStretchNormal,
	}
}

func (ah *linkAreaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	tl := ui.DrawNewTextLayout(&ui.DrawTextLayoutParams{
		String:      ah.attributedString,
		DefaultFont: ah.defaultFont,
		Width:       p.AreaWidth,
		Align:       ui.DrawTextAlignRight,
	})
	defer tl.Free()
	p.Context.Text(tl, 0, -1)
}

func (ah *linkAreaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	if me.Down == 1 { // left mouse button
		ah.onClickFunc()
	}
}

func (ah *linkAreaHandler) MouseCrossed(a *ui.Area, left bool) {
	if left {
		ah.attributedString.SetAttribute(linkColor, 0, len(ah.attributedString.String()))
	} else {
		ah.attributedString.SetAttribute(linkColorHover, 0, len(ah.attributedString.String()))
	}
	a.QueueRedrawAll()
}

func (ah *linkAreaHandler) DragBroken(a *ui.Area) {}

func (ah *linkAreaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) { return false }
