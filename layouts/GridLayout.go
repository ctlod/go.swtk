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
	iChan	chan swtk.LayoutMsger
}

func NewGridLayout() *GridLayout {
	pn := new(GridLayout)
	pn.sizes = make(map[int]*image.Point)
	pn.paneCell = make(map[int]*image.Point)
	pn.gridCell = make(map[image.Point]swtk.Pane)
	pn.gridAlign = make(map[image.Point]swtk.Alignmenter)
	pn.gridSize = image.Point{0, 0}
	pn.iChan = make(chan swtk.LayoutMsger)
	go swtk.LayoutActor(pn)
	return pn
}

func (gl *GridLayout) LayoutMsgChan() chan swtk.LayoutMsger {
	return gl.iChan
}

func (gl *GridLayout) OtherLayoutMsg(msg swtk.LayoutMsger) {
}

func (gl *GridLayout) Pane() swtk.Pane {
	return gl.thePane
}

func (gl *GridLayout) SetPane(pane swtk.Pane) {
	gl.thePane = pane
}

func (pn *GridLayout) MapResize(re swtk.ResizeMsg) {
	log.Println("Mapping Resize for Pane ", pn.thePane.Id())
	pn.size = re.Size
	for _, child := range pn.children {
		log.Println("Mapping Resize to Pane ", child.Id())
		cell := pn.paneCell[child.Id()]

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
		child.PaneMsgChan() <- swtk.ResizeMsg{Size: size, View: p.Sub(origin)}
	}
}

func (gl *GridLayout) RemovePane(pn swtk.Pane) {
}

func (pn *GridLayout) MapClose() {
	for _, c := range pn.children {
		close(c.PaneMsgChan())
	}
}

func (gl *GridLayout) AddPane(pane swtk.Pane, x, y int) {
	if gl.paneCell[pane.Id()] != nil {
		log.Println("GridLayout - Pane already exists: ", pane.Id(), x, y)
		return
	}

	if gl.gridSize.X <= x {
		gl.gridSize.X = x + 1
	}

	if gl.gridSize.Y <= y {
		gl.gridSize.Y = y + 1
	}

	gl.children = append(gl.children, pane)
	p := image.Point{x, y}
	gl.paneCell[pane.Id()] = &p
	gl.gridCell[p] = pane
	gl.gridAlign[p] = swtk.AlignCenter
	gl.thePane.Renderer().RegisterPane(pane, gl.thePane)
	pane.PaneMsgChan() <- swtk.SetRendererMsg{Renderer: gl.thePane.Renderer()}
	
}

func (pn *GridLayout) SetAlignment(pane swtk.Pane, a swtk.Alignmenter) {
	c := pn.paneCell[pane.Id()]
	pn.gridAlign[*c] = a
}
