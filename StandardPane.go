package swtk

import "image"

type standardPane struct {
	minSize image.Point
	maxSize image.Point
	size    image.Point

	id int

	inputChan chan PaneMsger

	visualer Visualer
	layouter Layouter
	actioner Actioner

	windower Windower
}

func (pn *standardPane) MinMax() (image.Point, image.Point) {
	return pn.minSize, pn.maxSize
}

func (pn *standardPane) Id() int {
	return pn.id
}

func (pn *standardPane) OtherPaneMsg(msg PaneMsger) {
}

func (pn *standardPane) CreatePaneMsgChan() int {
	if pn.inputChan != nil {
		return 1
	}

	pn.inputChan = make(chan PaneMsger, 0)

	return 0
}

func (pn *standardPane) PaneMsgChan() chan PaneMsger {
	return pn.inputChan
}

func (pn *standardPane) Visualer() Visualer {
	return pn.visualer
}

func (pn *standardPane) SetVisualer(vs Visualer) {
	pn.visualer = vs
}

func (pn *standardPane) Layouter() Layouter {
	return pn.layouter
}

func (pn *standardPane) SetLayouter(ly Layouter) {
	pn.layouter = ly
}

func (pn *standardPane) Actioner() Actioner {
	return pn.actioner
}

func (pn *standardPane) SetActioner(ac Actioner) {
	pn.actioner = ac
}

func (pn *standardPane) Windower() Windower {
	return pn.windower
}

func (pn *standardPane) SetWindower(wn Windower) {
	pn.windower = wn
}

func (pn *standardPane) SetSize(rs ResizeMsg) {

}

func (pn *standardPane) Size() image.Point {
	return pn.size
}

func NewStandardPane() *standardPane {
	pn := new(standardPane)
	pn.minSize = image.Point{0, 0}
	pn.maxSize = image.Point{0, 0}
	pn.id = SwtkId()
	pn.inputChan = make(chan PaneMsger)
	return pn
}

func NewStandardPaneActor() *standardPane {
	pn := NewStandardPane()
	go PaneActor(pn)
	return pn
}
