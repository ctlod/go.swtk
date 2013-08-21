package displays

import "image"
import "image/draw"
import "image/color"
import "code.google.com/p/freetype-go/freetype"
import "github.com/ctlod/go.swtk"

type ButtonDisplayPane struct {
	thePane       swtk.Pane
	im            draw.Image
	renderer      swtk.Renderer
	fg          image.Image
	bg       image.Image
	mask          image.Image
	highlightWidth   int
	shadowWidth   int
	drawChannel   chan int
	sizeChannel   chan image.Point
	closeChannel chan int
	stateChannel chan int
	state int
	ftContext *freetype.Context
	label string
}

func NewButtonDisplayPane(bgc, fgc color.Color, label string) *ButtonDisplayPane {
	pn := new(ButtonDisplayPane)
	pn.bg = image.NewUniform(bgc)
	pn.fg = image.NewUniform(fgc)
	pn.mask = image.NewUniform(color.Alpha{128})
	pn.highlightWidth = 2
	pn.sizeChannel = make(chan image.Point, 1)
	pn.drawChannel = make(chan int, 1)
	pn.closeChannel = make(chan int, 1)
	pn.stateChannel = make(chan int, 1)
	pn.state = 0
	pn.label = label
	
	pn.ftContext = freetype.NewContext()
	pn.ftContext.SetDPI(swtk.Dpi)
	pn.ftContext.SetFont(swtk.Font)
	pn.ftContext.SetFontSize(swtk.FontSize)
	pn.ftContext.SetSrc(pn.fg)

	return pn
}

func (pn *ButtonDisplayPane) SetPane(p swtk.Pane) {
	pn.thePane = p
}

func (dp *ButtonDisplayPane) SetState(s int) {
	dp.stateChannel <- s
}

func (dp *ButtonDisplayPane) SetSize(size image.Point) {
	dp.sizeChannel <- size
}

func (dp *ButtonDisplayPane) Close() {
	dp.closeChannel <- 1
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

func (pn *ButtonDisplayPane) Draw() {
	pn.drawChannel <- 1
}

func (dp *ButtonDisplayPane) SetRenderer(r swtk.Renderer) {
	dp.renderer = r
}

func (pn *ButtonDisplayPane) draw() {
	if pn.im != nil {
		r := pn.im.Bounds()
		pn.ftContext.SetClip(r)
		pn.ftContext.SetDst(pn.im)

		textHeight := int(pn.ftContext.PointToFix32(swtk.FontSize)>>8)
		pt := freetype.Pt(0, textHeight)
		pt, _ = pn.ftContext.DrawString(pn.label, pt)
		lableLength := int(pt.X >> 8)
		pt = freetype.Pt((r.Dx() - lableLength)  / 2, (r.Dy() + textHeight) / 2 - 1)

		draw.Draw(pn.im, r, pn.bg, image.ZP, draw.Src)
		pn.ftContext.DrawString(pn.label, pt)

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
	pn.renderer.SetAspect(swtk.PaneImage{pn.thePane, pn.im})
}
