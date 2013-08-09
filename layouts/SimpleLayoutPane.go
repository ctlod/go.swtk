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

func (pn *SimpleLayoutPane) HandleResizeEvent(re image.Point) {
	pn.size = re
	for _, child := range pn.children {
		c := pn.coords[child]
		p := image.Point{re.X - c.X, re.Y - c.Y}
		min, max := child.MinMax()

		if min.X > p.X {
			p.X = 0
		} else if max.X != 0 && max.X < p.X {
			p.X = max.X
		}

		if min.Y > p.Y {
			p.Y = 0
		} else if max.Y != 0 && max.Y < p.Y {
			p.Y = max.Y
		}

		pn.renderer.RefreshLocationChannel() <- swtk.PaneCoords{child, *c}
		child.ResizeChannel() <- p
	}
}

func (pn *SimpleLayoutPane) HandleCloseEvent(ce int) {
	for _, c := range pn.children {
		c.CloseChannel() <- ce
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

	if pane.LayoutPane() != nil {
		pane.LayoutPane().RegisterRenderer(pn.renderer)
	}

	if pane.DisplayPane() != nil {
		pane.DisplayPane().SetRenderChannel(pn.renderer.UpdateNotifyChannel())
	}

	go pane.EventHandler()
}
