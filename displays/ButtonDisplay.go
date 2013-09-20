package displays

/*
import "image"
import "image/draw"
import "image/color"
import "code.google.com/p/freetype-go/freetype"
import "github.com/ctlod/go.swtk"

type ButtonDisplay struct {
	thePane       	swtk.Pane
	im            	draw.Image
	ftim						draw.Image
	fg     	     		image.Image
	bg       				image.Image
	mask          	image.Image
	highlightWidth  int
	shadowWidth   	int
	drawChannel   	chan int
	sizeChannel   	chan swtk.ResizeEvent
	closeChannel 		chan int
	stateChannel 		chan int
	state 					int
	ftContext 			*freetype.Context
	label 					string
	size						image.Point
	view						image.Rectangle
}

func NewButtonDisplay(pane swtk.Pane, label string, fgc, bgc color.Color) *ButtonDisplay {
	pn := new(ButtonDisplay)
	pn.bg = image.NewUniform(bgc)
	pn.fg = image.NewUniform(fgc)
	pn.mask = image.NewUniform(color.Alpha{128})
	pn.highlightWidth = 2
	pn.sizeChannel = make(chan swtk.ResizeEvent, 1)
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

	pn.ftim = image.NewRGBA(image.Rect(0,0,0,0))

	pn.thePane = pane

	return pn
}

func (dp *ButtonDisplay) SetState(s int) {
	dp.stateChannel <- s
}

func (dp *ButtonDisplay) SetSize(size swtk.ResizeEvent) {
	dp.sizeChannel <- size
}

func (dp *ButtonDisplay) Stop() {
	dp.closeChannel <- 1
}

func (dp *ButtonDisplay) setSize(s swtk.ResizeEvent) {
	//set the correct size in the buffer
	dp.size = s.Size
	dp.view = s.View
	if s.View.Dx() == 0 || s.View.Dy() == 0 {
		dp.im = nil
	}
	dp.im = image.NewRGBA(image.Rect(0, 0, s.View.Dx(), s.View.Dy()))
}

func (dp *ButtonDisplay) Display() {
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

func (pn *ButtonDisplay) Draw() {
	pn.drawChannel <- 1
}

func (pn *ButtonDisplay) draw() {
	if pn.im != nil {
		r := pn.im.Bounds()
		draw.Draw(pn.im, r, pn.bg, image.ZP, draw.Src)

		pn.ftContext.SetClip(image.Rect(0,0,pn.size.X,pn.size.Y))
		pn.ftContext.SetDst(pn.ftim)
		textHeight := int(pn.ftContext.PointToFix32(swtk.FontSize)>>8)
		pt := freetype.Pt(0, textHeight)
		pt, _ = pn.ftContext.DrawString(pn.label, pt)
		labelLength := int(pt.X >> 8)
		pt = freetype.Pt((pn.size.X - labelLength)  / 2 - pn.view.Min.X, (pn.size.Y + textHeight) / 2 - 1 - pn.view.Min.Y)
		pn.ftContext.SetDst(pn.im)
		pn.ftContext.DrawString(pn.label, pt)

		if pn.state > 0 {
			r0 := image.Rect(0, 0, pn.size.X, pn.highlightWidth).Intersect(pn.view).Sub(pn.view.Min)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//left
			r0 = image.Rect(0, pn.highlightWidth, pn.highlightWidth, pn.size.Y - pn.highlightWidth).Intersect(pn.view).Sub(pn.view.Min)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//bottom
			r0 = image.Rect(0, pn.size.Y - pn.highlightWidth, pn.size.X, pn.size.Y).Intersect(pn.view).Sub(pn.view.Min)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
			//right
			r0 = image.Rect(pn.size.X - pn.highlightWidth, pn.highlightWidth, pn.size.X, pn.size.Y - pn.highlightWidth).Intersect(pn.view).Sub(pn.view.Min)
			draw.Draw(pn.im, r0, image.Black, image.ZP, draw.Src)
		}
		if pn.state == 2 {
			//shade if pressed
			draw.DrawMask(pn.im, r, image.Black, image.ZP, pn.mask, image.ZP, draw.Over)
		}
	}
	pn.thePane.Renderer().SetAspect(swtk.PaneImage{pn.thePane, pn.im})
}
*/
