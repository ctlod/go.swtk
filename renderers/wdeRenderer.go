package renderers

import "fmt"
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

	return wr
}

func (wr *wdeRenderer) BackEndRun() {
	wde.Run()
}

func (wr *wdeRenderer) Run() {
	for {
		select {
		case e, ok := <-wr.wdeWindow.EventChan():
			if !ok {
				return
			} else {
				switch e := e.(type) {
				case wde.CloseEvent:
					wr.basePane.CloseChannel() <- 1
					wr.wdeWindow.Close()
					wde.Stop()
					return
				case wde.ResizeEvent:
					x, y := wr.wdeWindow.Size()
					if (x != 0) && (y != 0) {
						re := image.Point{wr.wdeWindow.Screen().Bounds().Dx(), wr.wdeWindow.Screen().Bounds().Dy()}
						wr.handleWindowResize()
						wr.basePane.ResizeChannel() <- re
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
					fmt.Println(e.Key, e.Glyph, e.Chord)
				default:
				}
			}
		case pn := <-wr.panesCom:
			wr.refreshBuffer(pn.Pane, pn.Image)
			cond := true
			//count := cap(wr.panesCom)
			for cond {
				select {
				case pn2 := <-wr.panesCom:
					wr.refreshBuffer(pn2.Pane, pn2.Image)
					/*count--
					if count < 0 {
						cond = false
					}*/
				default:
					cond = false
				}
			}
			wr.render()
		case pc := <-wr.panesCoords:
			wr.refreshLocation(pc.Pane, pc.Coords)
			cond := true
			//count := cap(wr.panesCoords)
			for cond {
				select {
				case pc2 := <-wr.panesCoords:
					wr.refreshLocation(pc2.Pane, pc2.Coords)
					/*count--
					if count < 0 {
						cond = false
					}*/
				default:
					cond = false
				}
			}
			wr.render()
		}
	}
}

func (wr *wdeRenderer) wdeHandleMouseState(ms swtk.MouseState) {
	if ms.B < 0 && ms.X < 0 {
		//exit event
		if wr.mousePane != nil {
			wr.mousePane.MouseStateChannel() <- ms
		}
		wr.mousePane = nil
	} else {
		currentNode := wr.renderList[0]
		nextNode := currentNode
		for {
			nbChildren := len(currentNode.children) - 1
			if nbChildren >= 0 {
				for nbChildren >= 0 {
					child := currentNode.children[nbChildren]
					size := image.Point{0, 0}
					if child.im != nil {
						size = child.im.Bounds().Size()
					}
					if int(ms.X)-child.x >= 0 && int(ms.X)-child.x < size.X && int(ms.Y)-child.y >= 0 && int(ms.Y)-child.y < size.Y {
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
				if wr.mousePane != currentNode.pane && wr.mousePane != nil {
					wr.mousePane.MouseStateChannel() <- swtk.MouseState{int8(-1), int16(-1), int16(-1)}
					currentNode.pane.MouseStateChannel() <- swtk.MouseState{int8(-1), ms.X - int16(currentNode.x), ms.Y - int16(currentNode.y)}
				}
				currentNode.pane.MouseStateChannel() <- swtk.MouseState{ms.B, ms.X - int16(currentNode.x), ms.Y - int16(currentNode.y)}
				wr.mousePane = currentNode.pane
				break
			}
		}
	}
	wr.mouseState = ms

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
		pane.DisplayPane().SetRenderChannel(wr.panesCom)
	}
	wr.renderList = buildRenderList(wr.paneMap[wr.basePane], nil)
	go wr.basePane.EventHandler()
}

func (wr *wdeRenderer) refreshLocation(pane swtk.Pane, point image.Point) {
	node := wr.paneMap[pane]
	node.x = node.parent.x + point.X
	node.y = node.parent.y + point.Y
}

func (wr *wdeRenderer) refreshBuffer(pane swtk.Pane, im image.Image) {
	node := wr.paneMap[pane]
	if im == nil {
		node.im = nil
	} else {	
		if node.im == nil || node.im.Bounds() != im.Bounds() {
			node.im = image.NewRGBA(im.Bounds())
		}
		draw.Draw(node.im, node.im.Bounds(), im, im.Bounds().Min, draw.Src)
	}
}

func (wr *wdeRenderer) UpdateNotifyChannel() chan swtk.PaneImage {
	return wr.panesCom
}

func (wr *wdeRenderer) RefreshLocationChannel() chan swtk.PaneCoords {
	return wr.panesCoords
}

func (wr *wdeRenderer) render() {
	draw.Draw(wr.wdeWindow.Screen(), wr.wdeWindow.Screen().Bounds(), wr.bgImage, image.Point{0, 0}, draw.Src)
	for _, src := range wr.renderList {
		if src.im != nil {
			dp := image.Point{src.x, src.y}
			r := image.Rectangle{dp, dp.Add(src.im.Bounds().Size())}
			draw.Draw(wr.wdeWindow.Screen(), r, src.im, src.im.Bounds().Min, draw.Over)
		}
	}
	wr.wdeWindow.FlushImage(wr.wdeWindow.Screen().Bounds())
}

func (wr *wdeRenderer) handleWindowResize() {
	for _, src := range wr.renderList {
		if src.im != nil {
			src.im = nil
			src.x = 0
			src.y = 0
		}
	}
}

func buildRenderList(pn *paneNode, list []*paneNode) []*paneNode {
	list = append(list, pn)
	for _, c := range pn.children {
		list = buildRenderList(c, list)
	}
	return list
}
