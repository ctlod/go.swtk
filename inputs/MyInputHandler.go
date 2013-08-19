package inputs

import "image/color"
import "github.com/ctlod/go.swtk"
import "github.com/ctlod/go.swtk/displays"

type MyInputHandler struct {
	col1, col2 color.Color
	thePane    swtk.Pane
	ms         swtk.MouseState
	mouseStateChannel chan swtk.MouseState
	closeChannel chan int
}

func NewMyInputHandler(col1, col2 color.Color) *MyInputHandler {
	ih := new(MyInputHandler)
	ih.col1 = col1
	ih.col2 = col2
	ih.ms = swtk.MouseState{int8(-1), int16(-1), int16(-1)}
	ih.mouseStateChannel = make(chan swtk.MouseState, 1)
	return ih
}

func (ih *MyInputHandler) SetPane(pane swtk.Pane) {
	ih.thePane = pane
}

func (ih *MyInputHandler) InputHandler() {
	for {
		select {
			case _ = <- ih.closeChannel:
				break
			case mse := <- ih.mouseStateChannel:
				ih.handleMouseState(mse)
		}
	}
}

func (ih *MyInputHandler) HandleMouseState(ms swtk.MouseState) {
	ih.mouseStateChannel <- ms
}


func (ih *MyInputHandler) handleMouseState(mse swtk.MouseState) {
	if ih.thePane.DisplayPane() != nil {
		theDisplayPane, ok := ih.thePane.DisplayPane().(*displays.MyDisplayPane)
		if ok {
			//check if left click is released
			if ih.ms.B > int8(0) && (ih.ms.B&int8(1)) == int8(1) && (mse.B < int8(0) || (mse.B&int8(1)) == int8(0)) {
				theDisplayPane.SetColor() <- ih.col1
				theDisplayPane.Draw()
			}

			//check if left click is pressed (ignore if mouse enters pressed)
			if mse.B > int8(0) && (mse.B&int8(1) == int8(1)) && (ih.ms.B >= int8(0) && (ih.ms.B&int8(1)) == int8(0)) {
				theDisplayPane.SetColor() <- ih.col2
				theDisplayPane.Draw()
			}
		}
	}
	ih.ms = mse
}
