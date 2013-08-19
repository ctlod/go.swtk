package swtk

import "image"

type PaneCoords struct {
	Pane   Pane
	Coords image.Point
	Size   image.Point
}

type PaneImage struct {
	Pane  Pane
	Image image.Image
}

type Pane interface {
	//return minimum and maximum sizes desired
	//0 means not applicable.
	//minimum may not be respected.
	MinMax() (image.Point, image.Point)
	SetMinMax(min, max image.Point)

	SetSize(size image.Point)
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
	SetSize(size image.Point)
	Close()

	DrawingHandler()
}

// This handles children size and location
type LayoutPane interface {
	SetPane(pane Pane)
	HandleResizeEvent(re image.Point)
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

type MouseState struct {
	B    int8
	X, Y int16
}

//This is where on screen the pointer is
//There may be multiple pointers from multiple Devices
type PointerState struct {
	Device int
	Id     int
	X, Y   int
}

//This is where on screen has been 'touched'
// ie: with finger, or mouse button down
//There will certainly be multiple contacts from multiple Devices
type ContactState struct {
	Device int
	Id     int
	X, Y   int
}
