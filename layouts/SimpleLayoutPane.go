package layouts

import "image"
import "github.com/ctlod/go.swtk"

type SimpleLayoutPane struct {
	thePane  swtk.Pane
	renderer swtk.Renderer
	children []swtk.Pane
	coords   map[swtk.Pane]*image.Point
	sizes    map[swtk.Pane]*image.Point
	size     image.Point
}

func NewSimpleLayoutPane() *SimpleLayoutPane {
	pn := new(SimpleLayoutPane)
	pn.coords = make(map[swtk.Pane]*image.Point)
	pn.sizes = make(map[swtk.Pane]*image.Point)
	return pn
}

func (pn *SimpleLayoutPane) RegisterRenderer(wr swtk.Renderer) {
	pn.renderer = wr
}

func (pn *SimpleLayoutPane) HandleResizeEvent(re swtk.ResizeEvent) {
	pn.size = re.Size
	for _, child := range pn.children {
		c := pn.coords[child]
		origin := *c
		size := image.Point{re.Size.X - origin.X, re.Size.Y - origin.Y}
		min, max := child.MinMax()

		if min.X > size.X {
			size.X = 0
		} else if max.X != 0 && max.X < size.X {
			size.X = max.X
		}

		if min.Y > size.Y {
			size.Y = 0
		} else if max.Y != 0 && max.Y < size.Y {
			size.Y = max.Y
		}

		//Work out view to tell child
		p := image.Rect(0, 0, size.X, size.Y).Add(origin)
		p = p.Intersect(re.View)

		pn.renderer.SetLocation(swtk.PaneCoords{child, p.Min, size})
		child.SetSize(swtk.ResizeEvent{size, p.Sub(origin)})
	}
}

func (pn *SimpleLayoutPane) HandleCloseEvent(ce int) {
	for _, c := range pn.children {
		c.Close()
	}
}

func (lp *SimpleLayoutPane) SetPane(pn swtk.Pane) {
	lp.thePane = pn
}

func (pn *SimpleLayoutPane) AddPane(pane swtk.Pane, x, y int) {
	if pn.coords[pane] != nil {
		return
	}

	pn.children = append(pn.children, pane)
	pn.coords[pane] = &image.Point{x, y}
	pn.renderer.RegisterPane(pane, pn.thePane)

	go pane.PaneHandler()
}
