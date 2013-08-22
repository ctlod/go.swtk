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
	renderer  swtk.Renderer
	drawChannel   chan int
	sizeChannel   chan swtk.ResizeEvent
	closeChannel chan int
	stateChannel chan int
	ftContext *freetype.Context
	label string
	size						image.Point
	view						image.Rectangle
	ftim						draw.Image
}

func NewLabelDisplayPane(l string, c color.Color) *LabelDisplayPane {
	pn := new(LabelDisplayPane)
	pn.fg = image.NewUniform(c)
	pn.sizeChannel = make(chan swtk.ResizeEvent, 1)
	pn.drawChannel = make(chan int, 1)
	pn.closeChannel = make(chan int, 1)
	pn.stateChannel = make(chan int, 1)
	
	pn.ftContext = freetype.NewContext()
	pn.ftContext.SetDPI(swtk.Dpi)
	pn.ftContext.SetFont(swtk.Font)
	pn.ftContext.SetFontSize(swtk.FontSize)
	pn.ftContext.SetSrc(pn.fg)

	pn.ftim = image.NewRGBA(image.Rect(0,0,0,0))

	pn.label = l
	return pn
}

func (pn *LabelDisplayPane) SetPane(p swtk.Pane) {
	pn.thePane = p
}

func (dp *LabelDisplayPane) SetSize(size swtk.ResizeEvent) {
	dp.sizeChannel <- size
}

func (dp *LabelDisplayPane) Close() {
	dp.closeChannel <- 1
}

func (dp *LabelDisplayPane) setSize(s swtk.ResizeEvent) {
	//set the correct size in the buffer
	dp.size = s.Size
	dp.view = s.View
	if s.Size.X == 0 || s.Size.Y == 0 {
		dp.im = nil
	}
	dp.im = image.NewRGBA(image.Rect(0, 0, s.View.Dx(), s.View.Dy()))
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

func (pn *LabelDisplayPane) Draw() {
	pn.drawChannel <- 1
}

func (dp *LabelDisplayPane) SetRenderer(r swtk.Renderer) {
	dp.renderer = r
}

func (pn *LabelDisplayPane) draw() {
	if pn.im != nil {
		r := pn.im.Bounds()
		draw.Draw(pn.im, r, image.Transparent, image.ZP, draw.Src)

		pn.ftContext.SetClip(image.Rect(0,0,pn.size.X,pn.size.Y))
		pn.ftContext.SetDst(pn.ftim)
		textHeight := int(pn.ftContext.PointToFix32(swtk.FontSize)>>8)
		pt := freetype.Pt(0, textHeight)
		pt, _ = pn.ftContext.DrawString(pn.label, pt)
		labelLength := int(pt.X >> 8)
		pt = freetype.Pt((pn.size.X - labelLength)  / 2 - pn.view.Min.X, (pn.size.Y + textHeight) / 2 - 1 - pn.view.Min.Y)
		pn.ftContext.SetDst(pn.im)
		pn.ftContext.DrawString(pn.label, pt)
	}
	pn.renderer.SetAspect(swtk.PaneImage{pn.thePane, pn.im})
}
