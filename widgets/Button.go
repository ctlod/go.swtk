package widgets

import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/displays"
import "github.com/ctlod/go.swtk/inputs"
import "image/color"

type Button struct {
	*swtk.StandardPane
}

func NewButton(bg, fg color.Color, label string) *Button {
	bn := new(Button)
	bn.StandardPane = swtk.NewStandardPane()
	bn.SetDisplayPane(displays.NewButtonDisplayPane(bg, fg, label))
	bn.DisplayPane().SetPane(bn)
	bn.SetInputHandler(inputs.NewButtonInputHandler())
	bn.InputHandler().SetPane(bn)
	go bn.DisplayPane().DrawingHandler()
	go bn.InputHandler().InputHandler()
	return bn
}
