package displays

/*
import "image"
import "image/draw"
import "image/color"
import "code.google.com/p/freetype-go/freetype"
import "github.com/ctlod/go.swtk"
import "log"

type LabelDisplay struct {
	thePane       swtk.Pane
	im            draw.Image
	fg            image.Image
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

func NewLabelDisplay(pane swtk.Pane, l string, c color.Color) *LabelDisplay {
	pn := new(LabelDisplay)
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

	pn.thePane = pane

	pn.label = l
	return pn
}

func (dp *LabelDisplay) SetSize(size swtk.ResizeEvent) {
	dp.sizeChannel <- size
}

func (dp *LabelDisplay) Stop() {
	dp.closeChannel <- 1
}

func (dp *LabelDisplay) setSize(s swtk.ResizeEvent) {
	//set the correct size in the buffer
	dp.size = s.Size
	dp.view = s.View
	if s.Size.X == 0 || s.Size.Y == 0 {
		dp.im = nil
	}
	dp.im = image.NewRGBA(image.Rect(0, 0, s.View.Dx(), s.View.Dy()))
}

func (dp *LabelDisplay) Display() {
	log.Println("LabelDisplay Started")
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

func (pn *LabelDisplay) Draw() {
	pn.drawChannel <- 1
}

func (pn *LabelDisplay) draw() {
	log.Println("LabelDisplay Drawing")
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
	pn.thePane.Renderer().SetAspect(swtk.PaneImage{pn.thePane, pn.im})
}
*/
