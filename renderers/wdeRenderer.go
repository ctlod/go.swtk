package renderers

import "log"
import "image"
import "github.com/skelterjohn/go.wde"
import _ "github.com/skelterjohn/go.wde/init"
import "image/draw"
import "image/color"
import "github.com/ctlod/go.swtk"

const (
	wdeDebug = true
)

type paneNode struct {
	children []*paneNode
	parent   *paneNode
	im       draw.Image
	x, y     int
	dx, dy   int
	pane     swtk.Pane
}

type wdeRenderer struct{}

type wdeWindow struct {
	swtk.Pane
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

func (wr *wdeRenderer) NewWindowActor(title string, bg color.Color, x int, y int) *wdeWindow {
	w := new(wdeWindow)
	w.Pane = swtk.NewStandardPane()

	w.wdeWindow, _ = wde.NewWindow(x, y)
	w.wdeWindow.SetTitle(title)
	w.wdeWindow.Show()

	w.paneMap = make(map[int]*paneNode)
	w.panesCom = make(chan swtk.PaneImage, 100)
	w.panesCoords = make(chan swtk.PaneCoords, 100)

	w.minSize = image.Point{0, 0}
	w.maxSize = image.Point{0, 0}

	w.bgImage = image.NewUniform(bg)

	w.mouseState = swtk.MouseState{int8(-1), int16(-1), int16(-1)}
	w.pointerPanes = make(map[image.Point]swtk.Pane)

	w.SetWindower(w)
	if wdeDebug {
		log.Println("WindowPane: ", w.Id())
	}
	w.paneMap[w.Id()] = new(paneNode)
	w.paneMap[w.Id()].pane = w
	w.renderList = buildRenderList(w.paneMap[w.Id()], nil)
	if wdeDebug {
		log.Println(w.renderList, w.renderList[0].pane.Id())
	}

	go w.Run()

	return w
}

func NewWdeRenderer() *wdeRenderer {
	wr := new(wdeRenderer)
	return wr
}

func (wr *wdeRenderer) BackEndRun() {
	wde.Run()
}

func (wr *wdeWindow) Run() {
	wr.handleWindowResize()
	for {
		select {
		case e, ok := <-wr.wdeWindow.EventChan():
			if !ok {
				return
			} else {
				switch e := e.(type) {
				case wde.CloseEvent:
					if wr.Layouter() != nil {
						close(wr.Layouter().LayoutMsgChan())
					}
					wr.wdeWindow.Close()
					wde.Stop()
					return
				case wde.ResizeEvent:
					x, y := wr.wdeWindow.Size()
					if (x != 0) && (y != 0) {
						wr.handleWindowResize()
						p := image.Point{x, y}
						wr.refreshLocation(wr, image.ZP, p)
						b := image.Rectangle{image.ZP, p}
						msg := swtk.ResizeMsg{Size: p, View: b}
						wr.SetSize(msg)
						if wr.Layouter() != nil {
							wr.Layouter().LayoutMsgChan() <- msg
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
					if wdeDebug {
						log.Println(e.Key, e.Glyph, e.Chord)
					}
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
		case msg := <-wr.PaneMsgChan():
			switch msg := msg.(type) {
			case swtk.ResizeMsg:
				wr.SetSize(msg)
				if wr.Layouter() != nil {
					wr.Layouter().LayoutMsgChan() <- msg
				}
			case swtk.SetLayouterMsg:
				if wr.Layouter() == nil {
					wr.SetLayouter(msg.Layouter)
					msg.Layouter.LayoutMsgChan() <- swtk.SetPaneMsg{Pane: wr}
				}
			}
		}
	}
}

func (wr *wdeWindow) wdeHandleMouseState(ms swtk.MouseState) {
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

func (wr *wdeWindow) findNode(x, y int) *paneNode {
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

func (wr *wdeWindow) RegisterPane(pane swtk.Pane, parentPane swtk.Pane) {
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
	wr.renderList = buildRenderList(wr.paneMap[wr.Id()], nil)
}

func (wr *wdeWindow) refreshLocation(pane swtk.Pane, point image.Point, sizes image.Point) {
	log.Println("wdeRendere refreshlocation", pane.Id(), point, sizes)
	node := wr.paneMap[pane.Id()]
	if node.parent != nil {
		node.x = node.parent.x + point.X
		node.y = node.parent.y + point.Y
	}
	node.dx = sizes.X
	node.dy = sizes.Y
}

func (wr *wdeWindow) refreshBuffer(pane swtk.Pane, im draw.Image) {
	log.Println("wdeRendere refreshbuffer", pane, &im)
	node := wr.paneMap[pane.Id()]
	node.im = im
}

func (wr *wdeWindow) SetAspect(im swtk.PaneImage) {
	wr.panesCom <- im
}

func (wr *wdeWindow) SetLocation(pc swtk.PaneCoords) {
	wr.panesCoords <- pc
}

func (wr *wdeWindow) render(section image.Rectangle) {
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

func (wr *wdeWindow) handleWindowResize() {
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
