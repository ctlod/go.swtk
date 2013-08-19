package layouts

import "image"
import "github.com/ctlod/go.swtk"

type GridLayoutPane struct {
	thePane  swtk.Pane
	renderer swtk.Renderer
	children []swtk.Pane
	coords   map[swtk.Pane]*image.Point
	sizes    map[swtk.Pane]*image.Point
	paneCell map[swtk.Pane]*image.Point
	gridCell map[image.Point]swtk.Pane
	gridSize image.Point
	size     image.Point
}

func NewGridLayoutPane() *GridLayoutPane {
	pn := new(GridLayoutPane)
	pn.sizes = make(map[swtk.Pane]*image.Point)
	pn.paneCell = make(map[swtk.Pane]*image.Point)
	pn.gridCell = make(map[image.Point]swtk.Pane)
	pn.gridSize = image.Point{0, 0}
	return pn
}

func (pn *GridLayoutPane) RegisterRenderer(wr swtk.Renderer) {
	pn.renderer = wr
}

func (pn *GridLayoutPane) HandleResizeEvent(re image.Point) {
	pn.size = re
	for _, child := range pn.children {
		cell := pn.paneCell[child]
		gc_w := re.X / pn.gridSize.X
		gc_h := re.Y / pn.gridSize.Y

		min, max := child.MinMax()
		c := image.Point{gc_w * cell.X, gc_h * cell.Y}

		p := image.Point{gc_w, gc_h}
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

		pn.renderer.SetLocation(swtk.PaneCoords{child, c, p})
		child.SetSize(p)
	}
}

func (pn *GridLayoutPane) HandleCloseEvent() {
	for _, c := range pn.children {
		c.Close()
	}
}

func (lp *GridLayoutPane) SetPane(pn swtk.Pane) {
	lp.thePane = pn
}

func (pn *GridLayoutPane) AddPane(pane swtk.Pane, x, y int) {
	if pn.paneCell[pane] != nil {
		return
	}

	if (pn.gridSize.X <= x) {
		pn.gridSize.X = x + 1
	}

	if (pn.gridSize.Y <= y) {
		pn.gridSize.Y = y + 1
	}

	pn.children = append(pn.children, pane)
	pn.paneCell[pane] = &image.Point{x, y}
	pn.gridCell[*pn.paneCell[pane]] = pane
	pn.renderer.RegisterPane(pane, pn.thePane)

	if pane.LayoutPane() != nil {
		pane.LayoutPane().RegisterRenderer(pn.renderer)
	}

	if pane.DisplayPane() != nil {
		pane.DisplayPane().SetRenderer(pn.renderer)
	}

	go pane.PaneHandler()
}
