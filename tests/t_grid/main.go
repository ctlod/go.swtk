package main

import "image"
import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/renderers"
import "github.com/ctlod/go.swtk/layouts"
import "github.com/ctlod/go.swtk/displays"

func main() {

	r := renderers.NewWdeRenderer()

	w1 := r.NewWindowActor("Test Window", image.White, 400, 300)
	w1.PaneMsgChan() <- swtk.SetLayouterMsg{Layouter: layouts.NewGridLayoutActor()}

	p2 := swtk.NewStandardPaneActor()
	p2.PaneMsgChan() <- swtk.SetVisualerMsg{Visualer: displays.NewSimpleColorActor(image.Black)}
	w1.Layouter().LayoutMsgChan() <- swtk.AddPaneMsg{p2, 1, 2}
	p3 := swtk.NewStandardPaneActor()
	p3.PaneMsgChan() <- swtk.SetVisualerMsg{Visualer: displays.NewSimpleColorActor(image.Black)}
	w1.Layouter().LayoutMsgChan() <- swtk.AddPaneMsg{p3, 2, 1}
	w1.PaneMsgChan() <- swtk.ResizeMsg{Size: image.Point{400, 300}, View: image.Rect(0, 0, 400, 300)}

	r.NewWindowActor("Test Window 2", image.Black, 200, 200)

	r.BackEndRun()
}
