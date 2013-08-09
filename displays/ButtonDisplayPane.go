package displays

import "image"
import "image/draw"
import "image/color"
import "code.google.com/p/freetype-go/freetype"
import "github.com/ctlod/go.swtk"

type ButtonDisplayPane struct {
	thePane       swtk.Pane
	im            draw.Image
	col           color.Color
	renderChannel chan swtk.PaneImage
	bTop          image.Image
	bBottom       image.Image
	mask          image.Image
	highlightWidth   int
	shadowWidth   int
	drawChannel   chan int
	sizeChannel   chan image.Point
	closeChannel chan int
	stateChannel chan int
	state int
	label string
}

func NewButtonDisplayPane(a, b int, c color.Color) *ButtonDisplayPane {
	pn := new(ButtonDisplayPane)
	pn.col = c
	pn.bTop = image.White
	pn.bBottom = image.Black
	pn.mask = image.NewUniform(color.Alpha{128})
	pn.highlightWidth = 2
	pn.sizeChannel = make(chan image.Point, 1)
	pn.drawChannel = make(chan int, 1)
	pn.closeChannel = make(chan int, 1)
	pn.stateChannel = make(chan int, 1)
	pn.state = 0
	pn.label = "Button"
	return pn
}

func (pn *ButtonDisplayPane) SetPane(p swtk.Pane) {
	pn.thePane = p
}

func (dp *ButtonDisplayPane) SetState() chan int {
	return dp.stateChannel
}

func (dp *ButtonDisplayPane) SetSize() chan image.Point {
	return dp.sizeChannel
}

func (dp *ButtonDisplayPane) CloseChannel() chan int {
	return dp.closeChannel
}

func (dp *ButtonDisplayPane) setSize(s image.Point) {
	//simply set the correct size in the buffer
	if s.X == 0 || s.Y == 0 {
		dp.im = nil
	}
	dp.im = image.NewRGBA(image.Rect(0, 0, s.X, s.Y))
}

func (dp *ButtonDisplayPane) DrawingHandler() {
	for {
		select {
		case _ = <-dp.drawChannel:
			dp.draw()
		case re := <-dp.sizeChannel:
			dp.setSize(re)
			dp.draw()
		case st := <-dp.stateChannel:
			dp.state = st
			dp.draw()
		case _ = <-dp.closeChannel:
			break
		}
	}
}

func (pn *ButtonDisplayPane) Draw() chan int {
	return pn.drawChannel
}

func (dp *ButtonDisplayPane) SetRenderChannel(rc chan swtk.PaneImage) {
	dp.renderChannel = rc
}

func (pn *ButtonDisplayPane) draw() {
	if pn.im != nil {
		r := pn.im.Bounds()
		t := freetype.NewContext()
		t.SetDPI(swtk.Dpi)
		t.SetFont(swtk.Font)
		t.SetFontSize(swtk.FontSize)
		t.SetClip(r)
		t.SetDst(pn.im)
		t.SetSrc(image.Black)

		fHeight := int(t.PointToFix32(swtk.FontSize)>>8)
		pt := freetype.Pt(0, fHeight)
		pt1, _ := t.DrawString(pn.label, pt)
		lableLength := int(pt1.X >> 8)

		pt1 = freetype.Pt((r.Dx() - lableLength)  / 2, (r.Dy() + fHeight) / 2 - 1)

		draw.Draw(pn.im, r, &image.Uniform{pn.col}, image.ZP, draw.Src)
		t.DrawString(pn.label, pt1)

		if pn.state > 0 {
			r0 := image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y + pn.highlightWidth)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//left
			r0 = image.Rect(r.Min.X, r.Min.Y + pn.highlightWidth, r.Min.X + pn.highlightWidth, r.Max.Y - pn.highlightWidth)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//bottom
			r0 = image.Rect(r.Min.X, r.Max.Y - pn.highlightWidth, r.Max.X, r.Max.Y)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//right
			r0 = image.Rect(r.Max.X - pn.highlightWidth, r.Min.Y + pn.highlightWidth, r.Max.X, r.Max.Y - pn.highlightWidth)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
		}
		if pn.state == 2 {
			//shade if pressed
			draw.DrawMask(pn.im, r, image.Black, image.ZP, pn.mask, image.ZP, draw.Over)
		}
	}
	pn.renderChannel <- swtk.PaneImage{pn.thePane, pn.im}
}
