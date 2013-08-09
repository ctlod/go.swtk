package displays

import "image"
import "image/draw"
import "image/color"
import "code.google.com/p/freetype-go/freetype"
import "github.com/ctlod/go.swtk"


type LabelDisplayPane struct {
	thePane       swtk.Pane
	im            draw.Image
	fg            image.Image
	renderChannel chan swtk.PaneImage
	drawChannel   chan int
	sizeChannel   chan image.Point
	closeChannel chan int
	stateChannel chan int
	ftContext *freetype.Context
	label string
}

func NewLabelDisplayPane(l string, c color.Color) *LabelDisplayPane {
	pn := new(LabelDisplayPane)
	pn.fg = image.NewUniform(c)
	pn.sizeChannel = make(chan image.Point, 1)
	pn.drawChannel = make(chan int, 1)
	pn.closeChannel = make(chan int, 1)
	pn.stateChannel = make(chan int, 1)
	
	pn.ftContext = freetype.NewContext()
	pn.ftContext.SetDPI(swtk.Dpi)
	pn.ftContext.SetFont(swtk.Font)
	pn.ftContext.SetFontSize(swtk.FontSize)
	pn.ftContext.SetSrc(pn.fg)

	pn.label = l
	return pn
}

func (pn *LabelDisplayPane) SetPane(p swtk.Pane) {
	pn.thePane = p
}

func (dp *LabelDisplayPane) SetSize() chan image.Point {
	return dp.sizeChannel
}

func (dp *LabelDisplayPane) CloseChannel() chan int {
	return dp.closeChannel
}

func (dp *LabelDisplayPane) setSize(s image.Point) {
	//simply set the correct size in the buffer
	if s.X == 0 || s.Y == 0 {
		dp.im = nil
	}
	dp.im = image.NewRGBA(image.Rect(0, 0, s.X, s.Y))
}

func (dp *LabelDisplayPane) DrawingHandler() {
	for {
		select {
		case _ = <-dp.drawChannel:
			dp.draw()
		case re := <-dp.sizeChannel:
			dp.setSize(re)
			dp.draw()
		case _ = <-dp.closeChannel:
			break
		}
	}
}

func (pn *LabelDisplayPane) Draw() chan int {
	return pn.drawChannel
}

func (dp *LabelDisplayPane) SetRenderChannel(rc chan swtk.PaneImage) {
	dp.renderChannel = rc
}

func (pn *LabelDisplayPane) draw() {
	if pn.im != nil {
		r := pn.im.Bounds()
		pn.ftContext.SetClip(r)
		pn.ftContext.SetDst(pn.im)

		textHeight := int(pn.ftContext.PointToFix32(swtk.FontSize)>>8)
		pt := freetype.Pt(0, textHeight)
		pt, _ = pn.ftContext.DrawString(pn.label, pt)
		lableLength := int(pt.X >> 8)
		pt = freetype.Pt((r.Dx() - lableLength)  / 2, (r.Dy() + textHeight) / 2 - 1)

		draw.Draw(pn.im, r, image.Transparent, image.ZP, draw.Src)
		pn.ftContext.DrawString(pn.label, pt)
	}
	pn.renderChannel <- swtk.PaneImage{pn.thePane, pn.im}
}
