package swtk

import "image"

type Pane interface {
	//return minimum and maximum sizes desired
	//0 means not applicable.
	//minimum may not be respected.
	MinMax() (image.Point, image.Point)
	SetMinMax(min, max image.Point)

	SetSize(size ResizeEvent)
	Close()

	//Event handling
	SetMouseState(ms MouseState)
	//SetPointerState(ps PointerState)
	//SetContactState(cs ContactState)

	//Panes only need to know about focus, the inputHandler will take care of keyboard events
	//SetFocusState(fs FocusState)

	SetLayoutPane(lp LayoutPane)
	LayoutPane() LayoutPane

	SetDisplayPane(dp DisplayPane)
	DisplayPane() DisplayPane

	SetInputHandler(ih InputHandler)
	InputHandler() InputHandler

	PaneHandler()
}

type DisplayPane interface {
	SetPane(pane Pane)
	SetRenderer(r Renderer)

	Draw()
	SetSize(size ResizeEvent)
	Close()

	DrawingHandler()
}

type Alignmenter interface {
	Alignment() alignment
}

// This handles children size and location
type LayoutPane interface {
	SetPane(pane Pane)
	HandleResizeEvent(re ResizeEvent)
	HandleCloseEvent()
	RegisterRenderer(wr Renderer)

	//This initializes pane
	//no further setup should be possible afterwards
	AddPane(pane Pane, x, y int)
}

type Renderer interface {
	RegisterPane(pane Pane, parentPane Pane)
	SetAspect(pi PaneImage)
	SetLocation(pc PaneCoords)
	SetBasePane(pane Pane)
	Run()
	BackEndRun()

	//RequestFocus(pn Pane)
}

type InputHandler interface {
	SetPane(pn Pane)
	HandleMouseState(ms MouseState)
	//	HandlePointerState(ps PointerState)
	//	HandleContactState(cs ContactState)
	InputHandler()
}

