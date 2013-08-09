package displays

import "image"
import "image/draw"
import "image/color"
import "github.com/ctlod/go.swtk"


type MyDisplayPane struct {
	thePane       swtk.Pane
	im            draw.Image
	col           color.Color
	renderChannel chan swtk.PaneImage
	drawChannel   chan int
	sizeChannel   chan image.Point
	closeChannel chan int
	colorChannel chan color.Color
}

func (pn *MyDisplayPane) SetColor() chan color.Color {
	return pn.colorChannel
}

func (pn *MyDisplayPane) setColor(newCol color.Color) {
	pn.col = newCol
}

func NewMyDisplayPane(a, b int, c color.Color) *MyDisplayPane {
	pn := new(MyDisplayPane)
	pn.col = c
	pn.drawChannel = make(chan int, 1)
	pn.sizeChannel = make(chan image.Point, 1)
	pn.closeChannel = make(chan int, 1)
	return pn
}

func (pn *MyDisplayPane) CloseChannel () chan int {
	return pn.closeChannel
}

func (pn *MyDisplayPane) SetPane(p swtk.Pane) {
	pn.thePane = p
}

func (dp *MyDisplayPane) SetSize() chan image.Point {
	return dp.sizeChannel
}

func (dp *MyDisplayPane) setSize(s image.Point) {
	if s.X == 0 || s.Y == 0 {
		dp.im = nil
	}
	if (dp.im == nil || dp.im.Bounds().Dx() != s.X || dp.im.Bounds().Dy() != s.Y) {
		dp.im = image.NewRGBA(image.Rect(0, 0, s.X, s.Y))
	}
}

func (dp *MyDisplayPane) DrawingHandler() {
	for {
		select {
		case _ = <- dp.drawChannel:
			dp.draw()
		case re := <- dp.sizeChannel:
			dp.setSize(re)
			dp.draw()
		case cl := <- dp.colorChannel:
			dp.setColor(cl)
			dp.draw()
		case _ = <- dp.closeChannel:
			break
		}
	}
}

func (dp *MyDisplayPane) SetRenderChannel(rc chan swtk.PaneImage) {
	dp.renderChannel = rc
}

func (pn *MyDisplayPane) Draw() chan int {
	return pn.drawChannel
}

func (pn *MyDisplayPane) draw() {
	if pn.im != nil {
		draw.Draw(pn.im, pn.im.Bounds(), &image.Uniform{pn.col}, image.ZP, draw.Src)
	}
	pn.renderChannel <- swtk.PaneImage{pn.thePane, pn.im}
}
