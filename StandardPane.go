package swtk

import "image"

type StandardPane struct {
	displayPane  DisplayPane
	layoutPane   LayoutPane
	inputHandler InputHandler

	minSize image.Point
	maxSize image.Point

	closeChan      chan int
	resizeChan     chan ResizeEvent
	mouseStateChan chan MouseState
}

func (pn *StandardPane) MinMax() (image.Point, image.Point) {
	return pn.minSize, pn.maxSize
}

func (pn *StandardPane) SetMinMax(min, max image.Point) {
	pn.minSize = min
	pn.maxSize = max
}

func (pn *StandardPane) LayoutPane() LayoutPane {
	return pn.layoutPane
}

func (pn *StandardPane) SetLayoutPane(lp LayoutPane) {
	pn.layoutPane = lp
}

func (pn *StandardPane) InputHandler() InputHandler {
	return pn.inputHandler
}

func (pn *StandardPane) SetInputHandler(ih InputHandler) {
	pn.inputHandler = ih
}

func (pn *StandardPane) DisplayPane() DisplayPane {
	return pn.displayPane
}

func (pn *StandardPane) SetDisplayPane(dp DisplayPane) {
	pn.displayPane = dp
}

func (pn *StandardPane) SetSize(size ResizeEvent) {
	pn.resizeChan <- size
}

func (pn *StandardPane) Close() {
	pn.closeChan <- 1
}

func (pn *StandardPane) SetMouseState(ms MouseState) {
	pn.mouseStateChan <- ms
}

func (pn *StandardPane) PaneHandler() {
	for {
		select {
		case me := <-pn.mouseStateChan:
			if pn.inputHandler != nil {
				pn.inputHandler.HandleMouseState(me)
			}
		case re := <-pn.resizeChan:
			//only treat last resize in a batch
			cond := true
			for cond {
				select {
				case re = <-pn.resizeChan:
				default:
					cond = false
				}
			}
			if pn.layoutPane != nil {
				pn.layoutPane.HandleResizeEvent(re)
			}
			if pn.displayPane != nil {
				pn.displayPane.SetSize(re)
			}
		case _ = <-pn.closeChan:
			if pn.layoutPane != nil {
				pn.layoutPane.HandleCloseEvent()
			}
			break
		}
	}
}

func NewStandardPane() *StandardPane {
	pn := new(StandardPane)
	pn.minSize = image.Point{0, 0}
	pn.maxSize = image.Point{0, 0}
	pn.resizeChan = make(chan ResizeEvent, 8)
	pn.closeChan = make(chan int)
	pn.mouseStateChan = make(chan MouseState)
	return pn
}
