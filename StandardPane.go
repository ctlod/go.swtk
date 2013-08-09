package swtk

import "image"

type StandardPane struct {
	displayPane  DisplayPane
	layoutPane   LayoutPane
	inputHandler InputHandler

	minSize image.Point
	maxSize image.Point

	closeChan      chan int
	resizeChan     chan image.Point
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
	lp.SetPane(pn)
}

func (pn *StandardPane) InputHandler() InputHandler {
	return pn.inputHandler
}

func (pn *StandardPane) SetInputHandler(ih InputHandler) {
	pn.inputHandler = ih
	ih.SetPane(pn)
	go ih.InputHandler()
}

func (pn *StandardPane) DisplayPane() DisplayPane {
	return pn.displayPane
}

func (pn *StandardPane) SetDisplayPane(dp DisplayPane) {
	pn.displayPane = dp
	dp.SetPane(pn)
	go dp.DrawingHandler()
}

func (pn *StandardPane) ResizeChannel() chan image.Point {
	return pn.resizeChan
}

func (pn *StandardPane) CloseChannel() chan int {
	return pn.closeChan
}

func (pn *StandardPane) MouseStateChannel() chan MouseState {
	return pn.mouseStateChan
}

func (pn *StandardPane) EventHandler() {
	for {
		select {
		case me := <-pn.mouseStateChan:
			if pn.inputHandler != nil {
				pn.inputHandler.HandleMouseState() <- me
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
				pn.displayPane.SetSize() <- re
			}
		case ce := <-pn.closeChan:
			if pn.layoutPane != nil {
				pn.layoutPane.HandleCloseEvent(ce)
			}
			break
		}
	}
}

func NewStandardPane() *StandardPane {
	pn := new(StandardPane)
	pn.minSize = image.Point{0, 0}
	pn.maxSize = image.Point{0, 0}
	pn.resizeChan = make(chan image.Point, 100)
	pn.closeChan = make(chan int)
	pn.mouseStateChan = make(chan MouseState)
	return pn
}
