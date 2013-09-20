package inputs

import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/displays"

type ButtonInputHandler struct {
	thePane             swtk.Pane
	ms                  swtk.MouseState
	closeChannel        chan int
	mouseStateChannel   chan swtk.MouseState
	pointerStateChannel chan swtk.PointerState
}

func NewButtonInputHandler() *ButtonInputHandler {
	ih := new(ButtonInputHandler)
	ih.ms = swtk.MouseState{int8(-1), int16(-1), int16(-1)}
	ih.closeChannel = make(chan int, 1)
	ih.mouseStateChannel = make(chan swtk.MouseState, 1)
	ih.pointerStateChannel = make(chan swtk.PointerState, 1)
	return ih
}

func (ih *ButtonInputHandler) SetPane(pane swtk.Pane) {
	ih.thePane = pane
}

func (ih *ButtonInputHandler) InputHandler() {
	for {
		select {
		case _ = <-ih.closeChannel:
			break
		case mse := <-ih.mouseStateChannel:
			ih.handleMouseState(mse)
		case mse := <-ih.pointerStateChannel:
			ih.handlePointerState(mse)
		}
	}
}

func (ih *ButtonInputHandler) HandleMouseState(ms swtk.MouseState) {
	ih.mouseStateChannel <- ms
}

func (ih *ButtonInputHandler) handlePointerState(mse swtk.PointerState) {
	if ih.thePane.DisplayPane() != nil {
		theDisplayPane, ok := ih.thePane.DisplayPane().(*displays.ButtonDisplayPane)
		if ok {
			if mse.X < 0 {
				//an exit event need to check if another pointer still here
				theDisplayPane.SetState(0)
			} else {
				//an enter event need to check if another pointer already here
				theDisplayPane.SetState(1)
			}
		}
	}
}

func (ih *ButtonInputHandler) handleContanctState(mse swtk.ContactState) {
	if ih.thePane.DisplayPane() != nil {
		theDisplayPane, ok := ih.thePane.DisplayPane().(*displays.ButtonDisplayPane)
		if ok {
			if mse.X < 0 {
				//an exit event
				theDisplayPane.SetState(0)
			} else {
				//an enter event
				theDisplayPane.SetState(1)
			}
		}
	}
}

func (ih *ButtonInputHandler) handleMouseState(mse swtk.MouseState) {
	if ih.thePane.DisplayPane() != nil {
		theDisplayPane, ok := ih.thePane.DisplayPane().(*displays.ButtonDisplayPane)
		if ok {
			if mse.B < int8(0) {
				if mse.X < int16(0) {
					//an exit event
					theDisplayPane.SetState(0)
				} else {
					//an enter event
					theDisplayPane.SetState(1)
				}
			} else {
				if ih.ms.B > int8(0) && (ih.ms.B&int8(1)) == int8(1) && (mse.B < int8(0) || (mse.B&int8(1)) == int8(0)) {
					//check if left click is released
					theDisplayPane.SetState(1)
				} else if mse.B > int8(0) && (mse.B&int8(1) == int8(1)) && (ih.ms.B >= int8(0) && (ih.ms.B&int8(1)) == int8(0)) {
					//check if left click is pressed (ignore if mouse enters pressed)
					theDisplayPane.SetState(2)
				}
			}
		}
	}
	ih.ms = mse
}
