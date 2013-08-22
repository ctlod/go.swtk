package widgets

import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/displays"
import "image/color"

type Label struct {
	*swtk.StandardPane
}

func NewLabel(label string, fg color.Color) *Label {
	bn := new(Label)
	bn.StandardPane = swtk.NewStandardPane()
	bn.SetDisplayPane(displays.NewLabelDisplayPane(label, fg))
	bn.DisplayPane().SetPane(bn)
	go bn.DisplayPane().DrawingHandler()
	return bn
}