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

type ResizeMsg struct {
	PaneMessage
	Size image.Point
	View image.Rectangle
	VisualMessage
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

