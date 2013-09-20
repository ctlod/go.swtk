package layouts

import "image"
import "github.com/ctlod/go.swtk"
import "log"

type GridLayout struct {
	thePane   swtk.Pane
	children  []swtk.Pane
	coords    map[int]*image.Point
	sizes     map[int]*image.Point
	paneCell  map[int]*image.Point
	gridCell  map[image.Point]swtk.Pane
	gridAlign map[image.Point]swtk.Alignmenter
	gridSize  image.Point
	size      image.Point
}

func NewGridLayout(pane swtk.Pane) *GridLayout {
	pn := new(GridLayout)
	pn.sizes = make(map[int]*image.Point)
	pn.paneCell = make(map[int]*image.Point)
	pn.gridCell = make(map[image.Point]swtk.Pane)
	pn.gridAlign = make(map[image.Point]swtk.Alignmenter)
	pn.gridSize = image.Point{0, 0}
	pn.thePane = pane
	return pn
}

func (pn *GridLayout) HandleResizeEvent(re swtk.ResizeEvent) {
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
		if size.X < gc_w {
			if a == swtk.AlignCenter || a == swtk.AlignTop || a == swtk.AlignBottom {
				origin.X = origin.X + (gc_w-size.X)/2
			} else if a == swtk.AlignCenterRight || a == swtk.AlignTopRight || a == swtk.AlignBottomRight {
				origin.X = origin.X + (gc_w - size.X)
			}
		}

		//Align Vertically
		if size.Y < gc_h {
			if a == swtk.AlignCenter || a == swtk.AlignCenterLeft || a == swtk.AlignCenterRight {
				origin.Y = origin.Y + (gc_h-size.Y)/2
			} else if a == swtk.AlignBottom || a == swtk.AlignBottomLeft || a == swtk.AlignBottomRight {
				origin.Y = origin.Y + (gc_h - size.Y)
			}
		}

		//Work out view to tell child
		p := image.Rect(0, 0, size.X, size.Y).Add(origin)
		p = p.Intersect(re.View)

		pn.thePane.Renderer().SetLocation(swtk.PaneCoords{child, p.Min, image.Point{p.Dx(), p.Dy()}})
		child.SetSize(swtk.ResizeEvent{size, p.Sub(origin)})
	}
}

func (pn *GridLayout) HandleCloseEvent() {
	for _, c := range pn.children {
		c.Close()
	}
}

func (pn *GridLayout) AddPane(pane swtk.Pane, x, y int) {
	if pn.paneCell[pane] != nil {
		log.Println("GridLayout - Pane already exists: ", pane, x, y)
		return
	}

	if pn.gridSize.X <= x {
		pn.gridSize.X = x + 1
	}

	if pn.gridSize.Y <= y {
		pn.gridSize.Y = y + 1
	}

	pn.children = append(pn.children, pane)
	p := image.Point{x, y}
	pn.paneCell[pane] = &p
	pn.gridCell[p] = pane
	pn.gridAlign[p] = swtk.AlignCenter
}

func (pn *GridLayout) StartPanes() {
	for _, child := range pn.children {
		child.SetRenderer(pn.thePane.Renderer())
		pn.thePane.Renderer().RegisterPane(child, pn.thePane)
		go child.PaneHandler()
	}
}

func (pn *GridLayout) SetAlignment(p swtk.Pane, a swtk.Alignmenter) {
	c := pn.paneCell[p]
	pn.gridAlign[*c] = a
}