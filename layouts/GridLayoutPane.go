package layouts

import "image"
import "github.com/ctlod/go.swtk"
import "log"


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

func (pn *GridLayoutPane) HandleResizeEvent(re swtk.ResizeEvent) {
	pn.size = re.Size
	for _, child := range pn.children {
		cell := pn.paneCell[child]
		
		//work out width of current cell
		gc_w := re.Size.X / pn.gridSize.X
		gc_h := re.Size.Y / pn.gridSize.Y

		//work out origin of current cell
		origin := image.Point{gc_w * cell.X, gc_h * cell.Y}

		size := image.Point{gc_w, gc_h}
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

		a := pn.gridAlign[*cell]
		//Align horizontally
		if (size.X < gc_w) {
			if (a == swtk.AlignCenter || a== swtk.AlignTop || a == swtk.AlignBottom) {
				origin.X = origin.X + (gc_w - size.X) / 2
			} else if (a == swtk.AlignCenterRight || a== swtk.AlignTopRight || a == swtk.AlignBottomRight) {
				origin.X = origin.X + (gc_w - size.X)
			}
		}

		//Align Vertically
		if (size.Y < gc_h) {
			if (a == swtk.AlignCenter || a== swtk.AlignCenterLeft || a == swtk.AlignCenterRight) {
				origin.Y = origin.Y + (gc_h - size.Y) / 2
			} else if (a == swtk.AlignBottom || a== swtk.AlignBottomLeft || a == swtk.AlignBottomRight) {
				origin.Y = origin.Y + (gc_h - size.Y)
			}
		}

		//Work out view to tell child
		p := image.Rect(0, 0, size.X, size.Y).Add(origin)
		p = p.Intersect(re.View)

		pn.renderer.SetLocation(swtk.PaneCoords{child, p.Min, image.Point{p.Dx(), p.Dy()}})
		child.SetSize(swtk.ResizeEvent{size, p.Sub(origin)})
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
		log.Println("GridLayoutPane - Pane already exists: ", pane, x, y)
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