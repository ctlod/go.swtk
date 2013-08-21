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
	gridAlign map[image.Point]swtk.Alignmenter
	gridSize image.Point
	size     image.Point
}

func NewGridLayoutPane() *GridLayoutPane {
	pn := new(GridLayoutPane)
	pn.sizes = make(map[swtk.Pane]*image.Point)
	pn.paneCell = make(map[swtk.Pane]*image.Point)
	pn.gridCell = make(map[image.Point]swtk.Pane)
	pn.gridAlign = make(map[image.Point]swtk.Alignmenter)
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

		a := pn.gridAlign[*cell]

		//Align horizontally
		if (p.X < gc_w) {
			if (a == swtk.AlignCenter || a== swtk.AlignTop || a == swtk.AlignBottom) {
				c.X = c.X + (gc_w - p.X) / 2
			} else if (a == swtk.AlignCenterRight || a== swtk.AlignTopRight || a == swtk.AlignBottomRight) {
				c.X = c.X + (gc_w - p.X)
			}
		}

		//Align Vertically
		if (p.Y < gc_h) {
			if (a == swtk.AlignCenter || a== swtk.AlignCenterLeft || a == swtk.AlignCenterRight) {
				c.Y = c.Y + (gc_h - p.Y) / 2
			} else if (a == swtk.AlignBottom || a== swtk.AlignBottomLeft || a == swtk.AlignBottomRight) {
				c.Y = c.Y + (gc_h - p.Y)
			}
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
	p := image.Point{x, y}
	pn.paneCell[pane] = &p
	pn.gridCell[p] = pane
	pn.gridAlign[p] = swtk.AlignCenter
	pn.renderer.RegisterPane(pane, pn.thePane)

	if pane.LayoutPane() != nil {
		pane.LayoutPane().RegisterRenderer(pn.renderer)
	}

	if pane.DisplayPane() != nil {
		pane.DisplayPane().SetRenderer(pn.renderer)
	}

	go pane.PaneHandler()
}

func (pn *GridLayoutPane) SetAlignment(p swtk.Pane, a swtk.Alignmenter) {
	c := pn.paneCell[p]
	pn.gridAlign[*c] = a
}