package renderers

import "log"
import "image"
import "github.com/skelterjohn/go.wde"
import _ "github.com/skelterjohn/go.wde/init"
import "image/draw"
import "image/color"
import "github.com/ctlod/go.swtk"

type paneNode struct {
	children []*paneNode
	parent   *paneNode
	im       draw.Image
	x, y     int
	dx, dy   int
	pane     swtk.Pane
}

type wdeRenderer struct {
	basePane    swtk.Pane
	wdeWindow   wde.Window
	renderList  []*paneNode
	paneMap     map[swtk.Pane]*paneNode
	panesCom    chan swtk.PaneImage
	panesCoords chan swtk.PaneCoords
	minSize     image.Point
	maxSize     image.Point
	bgImage     image.Image
	mouseState  swtk.MouseState
	mousePane   swtk.Pane
	pointerPanes map[image.Point]swtk.Pane
}

func NewWdeRenderer(title string, bg color.Color) *wdeRenderer {
	wr := new(wdeRenderer)

	wr.wdeWindow, _ = wde.NewWindow(400, 300)
	wr.wdeWindow.Show()
	wr.wdeWindow.SetTitle(title)

	wr.paneMap = make(map[swtk.Pane]*paneNode)

	wr.panesCom = make(chan swtk.PaneImage, 100)
	wr.panesCoords = make(chan swtk.PaneCoords, 100)

	wr.minSize = image.Point{0, 0}
	wr.maxSize = image.Point{0, 0}

	wr.bgImage = image.NewUniform(bg)

	wr.mouseState = swtk.MouseState{int8(-1), int16(-1), int16(-1)}

	wr.pointerPanes = make(map[image.Point]swtk.Pane)

	return wr
}

func (wr *wdeRenderer) BackEndRun() {
	wde.Run()
}

func (wr *wdeRenderer) Run() {
	wr.handleWindowResize()
	for {
		select {
		case e, ok := <-wr.wdeWindow.EventChan():
			if !ok {
				return
			} else {
				switch e := e.(type) {
				case wde.CloseEvent:
					wr.basePane.Close()
					wr.wdeWindow.Close()
					wde.Stop()
					return
				case wde.ResizeEvent:
					x, y := wr.wdeWindow.Size()
					if (x != 0) && (y != 0) {
						b := wr.wdeWindow.Screen().Bounds()
						p := image.Point{b.Dx(), b.Dy()}
						wr.handleWindowResize()
						wr.basePane.SetSize(swtk.ResizeEvent{p,b})
					}
				case wde.MouseDownEvent:
					me := swtk.MouseState{int8(e.Which) | wr.mouseState.B, int16(e.Where.X), int16(e.Where.Y)}
					wr.wdeHandleMouseState(me)
				case wde.MouseUpEvent:
					me := swtk.MouseState{wr.mouseState.B &^ int8(e.Which), int16(e.Where.X), int16(e.Where.Y)}
					wr.wdeHandleMouseState(me)
				case wde.MouseDraggedEvent:
					me := swtk.MouseState{int8(e.Which), int16(e.Where.X), int16(e.Where.Y)}
					wr.wdeHandleMouseState(me)
				case wde.MouseMovedEvent:
					me := swtk.MouseState{int8(0), int16(e.Where.X), int16(e.Where.Y)}
					wr.wdeHandleMouseState(me)
				case wde.MouseEnteredEvent:
					me := swtk.MouseState{int8(-1), int16(e.Where.X), int16(e.Where.Y)}
					wr.wdeHandleMouseState(me)
				case wde.MouseExitedEvent:
					me := swtk.MouseState{int8(-1), int16(-1), int16(-1)}
					wr.wdeHandleMouseState(me)
				case wde.KeyDownEvent:
				case wde.KeyUpEvent:
				case wde.KeyTypedEvent:
					log.Println(e.Key, e.Glyph, e.Chord)
				default:
				}
			}
		case pn := <-wr.panesCom:
			wr.refreshBuffer(pn.Pane, pn.Image)
			p := wr.paneMap[pn.Pane]
			r := pn.Image.Bounds().Add(image.Point{p.x, p.y})
			rr := r
			panesCom:for {
				select {
				case pn = <-wr.panesCom:
					wr.refreshBuffer(pn.Pane, pn.Image)
					p = wr.paneMap[pn.Pane]
					r = pn.Image.Bounds().Add(image.Point{p.x, p.y})
					rr = rr.Union(r)
				default:
					break panesCom
				}
			}
			wr.render(rr)
		case pc := <-wr.panesCoords:
			wr.refreshLocation(pc.Pane, pc.Coords, pc.Size)
			panesCoords:for {
				select {
				case pc = <-wr.panesCoords:
					wr.refreshLocation(pc.Pane, pc.Coords, pc.Size)
				default:
					break panesCoords
				}
			}
		}
	}
}

func (wr *wdeRenderer) wdeHandleMouseState(ms swtk.MouseState) {
	//mouse event
	if ms.B < 0 && ms.X < 0 {
		//exit event
		if wr.mousePane != nil {
			wr.mousePane.SetMouseState(ms)
		}
		wr.mousePane = nil
	} else {
		targetNode := wr.findNode(int(ms.X), int(ms.Y))
		if wr.mousePane != targetNode.pane && wr.mousePane != nil {
			//Pointer Exit Event (ie, -Id, -1, -1 coordinates)
			wr.mousePane.SetMouseState(swtk.MouseState{int8(-1), int16(-1), int16(-1)})
			//Pointer Enter Event (Doesn't need to exist...)
			targetNode.pane.SetMouseState(swtk.MouseState{int8(-1), ms.X - int16(targetNode.x), ms.Y - int16(targetNode.y)})
		}
		//make Pointer Event
		targetNode.pane.SetMouseState(swtk.MouseState{ms.B, ms.X - int16(targetNode.x), ms.Y - int16(targetNode.y)})
		//make contact Event
		wr.mousePane = targetNode.pane
	}

	/*
	//handle pointer event from Mouse
	if ms.X < int16(0) {
		//create exit pointer
		if wr.pointerPanes[image.Point{0,0}] != nil {
			p := swtk.PointerState{0,0, -1, -1}
			log.Println(p)
			wr.pointerPanes[image.Point{0,0}] = nil
		}
	} else {
		targetNode := wr.findNode(int(ms.X), int(ms.Y))
		if wr.pointerPanes[image.Point{0,0}] != targetNode.pane && wr.pointerPanes[image.Point{0,0}] != nil {
			po := swtk.PointerState{0,0, -1, -1}
			log.Println(po)
		}
		pi := swtk.PointerState{0,0, int(ms.X) - targetNode.x, int(ms.Y) - targetNode.y}
		log.Println(pi)
		wr.pointerPanes[image.ZP] = targetNode.pane
	}
	*/
	wr.mouseState = ms
}

func (wr *wdeRenderer) findNode(x, y int) *paneNode {
	currentNode := wr.renderList[0]
	nextNode := currentNode
	for {
		nbChildren := len(currentNode.children) - 1
		if nbChildren >= 0 {
			for nbChildren >= 0 {
				child := currentNode.children[nbChildren]
				if x-child.x >= 0 && x - child.x < child.dx && y - child.y >= 0 && y - child.y < child.dy {
					nextNode = currentNode.children[nbChildren]
					nbChildren = -1
				} else {
					nbChildren = nbChildren - 1
				}
			}
		}

		if nextNode != currentNode {
			currentNode = nextNode
		} else {
			break
		}
	}
	return currentNode
}

func (wr *wdeRenderer) RegisterPane(pane swtk.Pane, parentPane swtk.Pane) {
	if wr.paneMap[pane] != nil {
		return
	}
	if parentPane == nil {
		return
	}
	if wr.paneMap[parentPane] == nil {
		return
	}

	wr.paneMap[pane] = new(paneNode)
	wr.paneMap[pane].parent = wr.paneMap[parentPane]
	wr.paneMap[parentPane].children = append(wr.paneMap[parentPane].children, wr.paneMap[pane])
	wr.paneMap[pane].pane = pane
	if pane.LayoutPane() != nil {
		pane.LayoutPane().RegisterRenderer(wr)
	}
	wr.renderList = buildRenderList(wr.paneMap[wr.basePane], nil)
}

func (wr *wdeRenderer) SetBasePane(pane swtk.Pane) {
	if wr.basePane != nil {
		return
	}
	wr.basePane = pane
	wr.paneMap[pane] = new(paneNode)
	wr.paneMap[pane].pane = pane
	if pane.LayoutPane() != nil {
		pane.LayoutPane().RegisterRenderer(wr)
	}
	if pane.DisplayPane() != nil {
		pane.DisplayPane().SetRenderer(wr)
	}
	wr.renderList = buildRenderList(wr.paneMap[wr.basePane], nil)
	go wr.basePane.PaneHandler()
}

func (wr *wdeRenderer) refreshLocation(pane swtk.Pane, point image.Point, sizes image.Point) {
	node := wr.paneMap[pane]
	node.x = node.parent.x + point.X
	node.y = node.parent.y + point.Y
	node.dx = sizes.X
	node.dy = sizes.Y
}

func (wr *wdeRenderer) refreshBuffer(pane swtk.Pane, im draw.Image) {
	node := wr.paneMap[pane]
	node.im = im
}

func (wr *wdeRenderer) SetAspect(im swtk.PaneImage) {
	wr.panesCom <-im
}

func (wr *wdeRenderer) SetLocation(pc swtk.PaneCoords) {
	wr.panesCoords <- pc
}

func (wr *wdeRenderer) render(section image.Rectangle) {
	draw.Draw(wr.wdeWindow.Screen(), section, wr.bgImage, image.ZP.Sub(section.Min), draw.Src)
	for _, src := range wr.renderList {
		if src.im != nil && src.dx > 0 && src.dy > 0 {
			orig := section.Min.Sub(image.Point{src.x, src.y})
			draw.Draw(wr.wdeWindow.Screen(), section, src.im, orig, draw.Over)
		}
	}
	wr.wdeWindow.FlushImage(section)
}

func (wr *wdeRenderer) handleWindowResize() {
	for _, src := range wr.renderList {
		src.im = nil
		src.x = 0
		src.y = 0
		src.dx = 0
		src.dy = 0
	}
	draw.Draw(wr.wdeWindow.Screen(), wr.wdeWindow.Screen().Bounds(), wr.bgImage, image.ZP, draw.Src)
	wr.wdeWindow.FlushImage(wr.wdeWindow.Screen().Bounds())
}

func buildRenderList(pn *paneNode, list []*paneNode) []*paneNode {
	list = append(list, pn)
	for _, c := range pn.children {
		list = buildRenderList(c, list)
	}
	return list
}
