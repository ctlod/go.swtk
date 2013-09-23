package main

import "image"
import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/renderers"
import "github.com/ctlod/go.swtk/layouts"
import "github.com/ctlod/go.swtk/displays"

func main() {

	window := renderers.NewWdeRenderer("Test Window", image.White)

	p := swtk.NewStandardPane()
	p.PaneMsgChan() <- swtk.SetLayouterMsg{Layouter: layouts.NewGridLayout()}

	window.SetBasePane(p)

	go window.Run()

	p2 := swtk.NewStandardPane()
	p2.PaneMsgChan() <- swtk.SetVisualerMsg{Visualer: displays.NewSimpleColor(image.Black)}
	p.Layouter().LayoutMsgChan() <- swtk.AddPaneMsg{p2, 1, 2}

	p3 := swtk.NewStandardPane()
	p3.PaneMsgChan() <- swtk.SetVisualerMsg{Visualer: displays.NewSimpleColor(image.Black)}
	p.Layouter().LayoutMsgChan() <- swtk.AddPaneMsg{p3, 2, 1}

	p.PaneMsgChan() <- swtk.ResizeMsg{Size: image.Point{400, 300}, View: image.Rect(0, 0, 400, 300)}

	window.BackEndRun()
}
