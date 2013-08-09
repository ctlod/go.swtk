package swtk

import "image"

type PaneCoords struct {
	Pane   Pane
	Coords image.Point
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

	//resize events are treated as absolute, no negotiating on size
	ResizeChannel() chan image.Point
	CloseChannel() chan int

	//Event handling
	MouseStateChannel() chan MouseState
	//GetKeyboardInputChan() chan int

	SetLayoutPane(lp LayoutPane)
	LayoutPane() LayoutPane

	SetDisplayPane(dp DisplayPane)
	DisplayPane() DisplayPane

	SetInputHandler(ih InputHandler)
	InputHandler() InputHandler

	EventHandler()
}

type DisplayPane interface {
	SetPane(pane Pane)
	SetRenderChannel(chan PaneImage)

	Draw() chan int
	SetSize() chan image.Point
	CloseChannel() chan int

	DrawingHandler()
}

// This handles children size and location
type LayoutPane interface {
	SetPane(pane Pane)
	HandleResizeEvent(re image.Point)
	HandleCloseEvent(ce int)
	RegisterRenderer(wr Renderer)

	//This initializes pane
	//no further setup should be possible afterwards
	AddPane(pane Pane, x, y int)
}

type Renderer interface {
	RegisterPane(pane Pane, parentPane Pane)
	UpdateNotifyChannel() chan PaneImage
	RefreshLocationChannel() chan PaneCoords
	SetBasePane(pane Pane)
	Run()
	BackEndRun()
}

type InputHandler interface {
	SetPane(pn Pane)
	HandleMouseState() chan MouseState
	InputHandler()
}

type MouseState struct {
	B    int8
	X, Y int16
}
