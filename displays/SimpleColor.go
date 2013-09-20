package displays

import "image"
import "image/draw"
import "image/color"
import "github.com/ctlod/go.swtk"
import "log"

type SimpleColor struct {
	thePane       	swtk.Pane
	canvas         	draw.Image
	bg       				image.Image
	mask						image.Image

  vchan chan swtk.VisualMsger
	
	view image.Rectangle
}

func NewSimpleColor(bgc color.Color) *SimpleColor {
	sc := new(SimpleColor)
	sc.bg = image.NewUniform(bgc)
	sc.mask = image.NewUniform(color.Alpha{128})
	sc.view = image.Rectangle{image.ZP, image.ZP}
	sc.vchan = make(chan swtk.VisualMsger)
	go swtk.VisualActor(sc)
	return sc
}

func (sc *SimpleColor) VisualMsgChan() chan swtk.VisualMsger {
	return sc.vchan
}

func (sc *SimpleColor) Pane() swtk.Pane {
	return sc.thePane
}

func (sc *SimpleColor) SetPane(p swtk.Pane) {
	sc.thePane = p
}

func (sc *SimpleColor) CreateVisualMsgChan() int {
  if sc.vchan != nil {
  	return 1
  }
  sc.vchan = make(chan swtk.VisualMsger, 0)
  return 0
}

func (sc *SimpleColor) OtherVisualMsg(msg swtk.VisualMsger) {
}

func (sc *SimpleColor) ResizeCanvas(rs swtk.ResizeMsg) {
  log.Println("Received msg: ", rs)
	sc.view = rs.View
	if rs.View.Dx() == 0 || rs.View.Dy() == 0 {
		sc.canvas = nil
	} else {
		sc.canvas = image.NewRGBA(image.Rect(0, 0, rs.View.Dx(), rs.View.Dy()))
	}
}

func (sc *SimpleColor) Draw() {
	if sc.canvas != nil {
		r := sc.canvas.Bounds()
		draw.Draw(sc.canvas, r, sc.bg, image.ZP, draw.Src)
  }
 	sc.thePane.Renderer().SetAspect(swtk.PaneImage{sc.thePane, sc.canvas})
}
