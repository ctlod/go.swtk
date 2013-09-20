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
	basePane      swtk.Pane
	wdeWindow     wde.Window
	renderList    []*paneNode
	paneMap       map[int]*paneNode
	panesCom      chan swtk.PaneImage
	panesCoords   chan swtk.PaneCoords
	minSize       image.Point
	maxSize       image.Point
	bgImage       image.Image
	mouseState    swtk.MouseState
	mouseActioner swtk.Actioner
	pointerPanes  map[image.Point]swtk.Pane
}

func NewWdeRenderer(title string, bg color.Color) *wdeRenderer {
	wr := new(wdeRenderer)

	wr.wdeWindow, _ = wde.NewWindow(400, 300)
	wr.wdeWindow.Show()
	wr.wdeWindow.SetTitle(title)

	wr.paneMap = make(map[int]*paneNode)

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
					if wr.basePane != nil {
						close(wr.basePane.PaneMsgChan())
					}
					wr.wdeWindow.Close()
					wde.Stop()
					return
				case wde.ResizeEvent:
					x, y := wr.wdeWindow.Size()
					if (x != 0) && (y != 0) {
						wr.handleWindowResize()
						if wr.basePane != nil {
							p := image.Point{x, y}
							wr.refreshLocation(wr.basePane, image.ZP, p)
							b := image.Rectangle{image.ZP, p}
							wr.basePane.PaneMsgChan() <- swtk.ResizeMsg{Size: p, View: b}
						}
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
			p := wr.paneMap[pn.Pane.Id()]
			r := pn.Image.Bounds().Add(image.Point{p.x, p.y})
			rr := r
		panesCom:
			for {
				select {
				case pn = <-wr.panesCom:
					wr.refreshBuffer(pn.Pane, pn.Image)
					p = wr.paneMap[pn.Pane.Id()]
					r = pn.Image.Bounds().Add(image.Point{p.x, p.y})
					rr = rr.Union(r)
				default:
					break panesCom
				}
			}
			wr.render(rr)
		case pc := <-wr.panesCoords:
			wr.refreshLocation(pc.Pane, pc.Coords, pc.Size)
		panesCoords:
			for {
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
		if wr.mouseActioner != nil {
			wr.mouseActioner.HandleMouseState(ms)
		}
		wr.mouseActioner = nil
	} else {
		targetNode := wr.findNode(int(ms.X), int(ms.Y))
		if targetNode != nil && targetNode.pane.Actioner() != nil {
			newActioner := targetNode.pane.Actioner()
			if wr.mouseActioner != nil && wr.mouseActioner != newActioner {
				//Pointer Exit Event (ie, -Id, -1, -1 coordinates)
				wr.mouseActioner.HandleMouseState(swtk.MouseState{int8(-1), int16(-1), int16(-1)})
				//Pointer Enter Event (Doesn't need to exist...)
				newActioner.HandleMouseState(swtk.MouseState{int8(-1), ms.X - int16(targetNode.x), ms.Y - int16(targetNode.y)})
			}
			//make Pointer Event
			newActioner.HandleMouseState(swtk.MouseState{ms.B, ms.X - int16(targetNode.x), ms.Y - int16(targetNode.y)})
			//make contact Event
			wr.mouseActioner = newActioner
		}
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
	if len(wr.renderList) == 0 {
		return nil
	}
	currentNode := wr.renderList[0]
	nextNode := currentNode
	for {
		nbChildren := len(currentNode.children) - 1
		if nbChildren >= 0 {
			for nbChildren >= 0 {
				child := currentNode.children[nbChildren]
				if x-child.x >= 0 && x-child.x < child.dx && y-child.y >= 0 && y-child.y < child.dy {
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
	if wr.paneMap[pane.Id()] != nil {
		return
	}
	if parentPane == nil {
		return
	}
	if wr.paneMap[parentPane.Id()] == nil {
		return
	}
	wr.paneMap[pane.Id()] = new(paneNode)
	wr.paneMap[pane.Id()].parent = wr.paneMap[parentPane.Id()]
	wr.paneMap[parentPane.Id()].children = append(wr.paneMap[parentPane.Id()].children, wr.paneMap[pane.Id()])
	wr.paneMap[pane.Id()].pane = pane
	wr.renderList = buildRenderList(wr.paneMap[wr.basePane.Id()], nil)
}

func (wr *wdeRenderer) SetBasePane(pane swtk.Pane) {
	if wr.basePane != nil {
		log.Println("wdeRenderer.go - SetBasePane: Already Set!")
	}
	wr.basePane = pane
	pane.PaneMsgChan() <- swtk.SetRendererMsg{Renderer: wr}
	log.Println("BasePane: ", wr.basePane.Id())
	wr.paneMap[pane.Id()] = new(paneNode)
	wr.paneMap[pane.Id()].pane = pane
	wr.renderList = buildRenderList(wr.paneMap[wr.basePane.Id()], nil)
	log.Println(wr.renderList, wr.renderList[0].pane.Id())
}

func (wr *wdeRenderer) refreshLocation(pane swtk.Pane, point image.Point, sizes image.Point) {
	log.Println("wdeRendere refreshlocation", pane.Id(), point, sizes)
	node := wr.paneMap[pane.Id()]
	if node.parent != nil {
		node.x = node.parent.x + point.X
		node.y = node.parent.y + point.Y
	}
	node.dx = sizes.X
	node.dy = sizes.Y
}

func (wr *wdeRenderer) refreshBuffer(pane swtk.Pane, im draw.Image) {
	log.Println("wdeRendere refreshbuffer", pane, &im)
	node := wr.paneMap[pane.Id()]
	node.im = im
}

func (wr *wdeRenderer) SetAspect(im swtk.PaneImage) {
	wr.panesCom <- im
}

func (wr *wdeRenderer) SetLocation(pc swtk.PaneCoords) {
	wr.panesCoords <- pc
}

func (wr *wdeRenderer) render(section image.Rectangle) {
	log.Println("wdeRenderer render", section)
	draw.Draw(wr.wdeWindow.Screen(), section, wr.bgImage, image.ZP.Sub(section.Min), draw.Src)
	for _, src := range wr.renderList {
		log.Println("wdeRenderer render src", src)
		if src.im != nil && src.dx > 0 && src.dy > 0 {
			log.Println("wdeRenderer render src drawing")
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
