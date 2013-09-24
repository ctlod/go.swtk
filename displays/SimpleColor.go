package displays

import "image"
import "image/draw"
import "image/color"
import "github.com/ctlod/go.swtk"
import "log"

type simpleColor struct {
	thePane swtk.Pane
	canvas  draw.Image
	bg      image.Image
	mask    image.Image

	vchan chan swtk.VisualMsger

	view image.Rectangle
}

func NewSimpleColor(bgc color.Color) *simpleColor {
	sc := new(simpleColor)
	sc.bg = image.NewUniform(bgc)
	sc.mask = image.NewUniform(color.Alpha{128})
	sc.view = image.Rectangle{image.ZP, image.ZP}
	sc.vchan = make(chan swtk.VisualMsger)
	return sc
}

func NewSimpleColorActor(bgc color.Color) *simpleColor {
	sc := NewSimpleColor(bgc)
	go swtk.VisualActor(sc)
	return sc
}

func (sc *simpleColor) VisualMsgChan() chan swtk.VisualMsger {
	return sc.vchan
}

func (sc *simpleColor) Pane() swtk.Pane {
	return sc.thePane
}

func (sc *simpleColor) SetPane(p swtk.Pane) {
	sc.thePane = p
}

func (sc *simpleColor) CreateVisualMsgChan() int {
	if sc.vchan != nil {
		return 1
	}
	sc.vchan = make(chan swtk.VisualMsger, 0)
	return 0
}

func (sc *simpleColor) OtherVisualMsg(msg swtk.VisualMsger) {
}

func (sc *simpleColor) ResizeCanvas(rs swtk.ResizeMsg) {
	log.Println("Received msg: ", rs)
	sc.view = rs.View
	if rs.View.Dx() == 0 || rs.View.Dy() == 0 {
		sc.canvas = nil
	} else {
		sc.canvas = image.NewRGBA(image.Rect(0, 0, rs.View.Dx(), rs.View.Dy()))
	}
}

func (sc *simpleColor) Draw() {
	if sc.canvas != nil {
		r := sc.canvas.Bounds()
		draw.Draw(sc.canvas, r, sc.bg, image.ZP, draw.Src)
	}
	sc.thePane.Windower().SetAspect(swtk.PaneImage{sc.thePane, sc.canvas})
}
