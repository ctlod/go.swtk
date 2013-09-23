package swtk

import "image"

type PaneMessage struct {
}

func (PaneMessage) PaneMsg() {
}

type VisualMessage struct {
}

func (vs VisualMessage) VisualMsg() {
}

type LayoutMessage struct {
}

func (ly LayoutMessage) LayoutMsg() {
}

type ResizeMsg struct {
	Size image.Point
	View image.Rectangle
	PaneMessage
	VisualMessage
	LayoutMessage
}

type AddPaneMsg struct {
	Pane Pane
	X, Y int
}

func (msg AddPaneMsg) LayoutMsg() {
}

type SetLayouterMsg struct {
	Layouter Layouter
	PaneMessage
}

type SetVisualerMsg struct {
	Visualer Visualer
	PaneMessage
}

type SetActionerMsg struct {
	Actioner Actioner
	PaneMessage
}

type SetRendererMsg struct {
	Renderer Renderer
	PaneMessage
}

type SetPaneMsg struct {
	Pane Pane
	VisualMessage
	LayoutMessage
}
