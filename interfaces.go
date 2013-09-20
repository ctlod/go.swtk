package swtk

import "image"
import "log"

type Pane interface {
	MinMax() (image.Point, image.Point)
	Size() image.Point
	Id() int

	SetSize(rs ResizeMsg)

	PaneMsgChan() chan PaneMsger
	CreatePaneMsgChan() int
	OtherPaneMsg(msg PaneMsger)

	Visualer() Visualer
	SetVisualer(vs Visualer)

	Actioner() Actioner
	SetActioner(ac Actioner)

	Layouter() Layouter
	SetLayouter(ly Layouter)

	Renderer() Renderer
	SetRenderer(rd Renderer)
}

type PaneMsger interface {
	PaneMsg()
}

type VisualMsger interface {
	VisualMsg()
}

type Visualer interface {
	Draw()
	ResizeCanvas(re ResizeMsg)

	OtherVisualMsg(msg VisualMsger)

	VisualMsgChan() chan VisualMsger
	CreateVisualMsgChan() int

	Pane() Pane
	SetPane(p Pane)
}

type Alignmenter interface {
	Alignment() alignment
}

type Actioner interface {
	HandleMouseState(ms MouseState)
}

type Layouter interface {
	MapClose()
	MapResize(re ResizeMsg)

	AddPane(pane Pane, x, y int)
	RemovePane(pane Pane)
}

type Renderer interface {
	RegisterPane(pane Pane, parentPane Pane)
	SetAspect(pi PaneImage)
	SetLocation(pc PaneCoords)
	SetBasePane(pane Pane)
	Run()
	BackEndRun()
}

func VisualActor(vs Visualer) {
	for {
		select {
		case msg, inChanOk := <-vs.VisualMsgChan():
			if !inChanOk {
				return
			}
			switch msg := msg.(type) {
			case ResizeMsg:
				log.Println("Drawing ", vs.Pane().Id())
				vs.ResizeCanvas(msg)
				vs.Draw()
			case SetPaneMsg:
				vs.SetPane(msg.Pane)
			default:
				vs.OtherVisualMsg(msg)
			}
		}
	}
}

func PaneActor(pn Pane) {
	log.Println("Starting ", pn.Id())
	for {
		select {
		case msg, inChanOk := <-pn.PaneMsgChan():
			if !inChanOk {
				if pn.Layouter() != nil {
					pn.Layouter().MapClose()
				}
				if pn.Visualer() != nil {
					close(pn.Visualer().VisualMsgChan())
				}
				if pn.Actioner() != nil {
				}
				return
			}
			log.Println("Received Msg for Pane ", pn.Id())
			switch msg := msg.(type) {
			case ResizeMsg:
				log.Println("Resize ", pn.Id())
				pn.SetSize(msg)
				if pn.Layouter() != nil {
					pn.Layouter().MapResize(msg)
				}
				if pn.Visualer() != nil {
					pn.Visualer().VisualMsgChan() <- msg
				}
			case SetRendererMsg:
				if pn.Renderer() != nil {
					if pn.Renderer() != msg.Renderer {
						//changing ?
					}
				} else {
					pn.SetRenderer(msg.Renderer)
				}
			case SetLayouterMsg:
				if pn.Layouter() != nil {
					if pn.Layouter() != msg.Layouter {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetLayouter(msg.Layouter)
				}
			case SetActionerMsg:
				if pn.Actioner() != nil {
					if pn.Actioner() != msg.Actioner {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetActioner(msg.Actioner)
				}
			case SetVisualerMsg:
				if pn.Visualer() != nil {
					if pn.Visualer() != msg.Visualer {
						//changing ?
					}
				} else {
					log.Println("Adding visual to ", pn.Id())
					pn.SetVisualer(msg.Visualer)
					msg.Visualer.VisualMsgChan() <- SetPaneMsg{Pane: pn}
				}
			default:
				pn.OtherPaneMsg(msg)
			}
		}
	}
	log.Println("Stopping ", pn.Id())
}

func SwtkId() int {
	i := <-idChan
	return i
}
