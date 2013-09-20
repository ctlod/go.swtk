package swtk

import "image"

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

func VisualerActor(vs Visualer) {
	ok := vs.CreateVisualMsgChan()
	if ok != 0 {
		return
	}
	for {
		select {
		case msg, inChanOk := <- vs.VisualMsgChan():
			if !inChanOk {
				break
			}
			switch msg := msg.(type) {
			case ResizeMsg:
				vs.ResizeCanvas(msg)
				vs.Draw()
			default:
				vs.OtherVisualMsg(msg)
			}
		}
	}
}

func PaneHandler(pn Pane) {
	ok := pn.CreatePaneMsgChan()
	if ok != 0 {
		return
	}
	for {
		select {
		case msg, inChanOk := <- pn.PaneMsgChan():
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
			switch msg := msg.(type) {
			case ResizeMsg:
				pn.SetSize(msg)
				if pn.Layouter() != nil {
					pn.Layouter().MapResize(msg)
				}
				if pn.Visualer() != nil {
					pn.Visualer().VisualMsgChan() <- msg
				}
			case SetRendererMsg:
				if(pn.Renderer() != nil){
					if (pn.Renderer() != msg.Renderer) {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetRenderer(msg.Renderer)
				}
			case SetLayouterMsg:
				if(pn.Layouter() != nil){
					if (pn.Layouter() != msg.Layouter) {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetLayouter(msg.Layouter)
				}
			case SetActionerMsg:
				if(pn.Actioner() != nil){
					if (pn.Actioner() != msg.Actioner) {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetActioner(msg.Actioner)
				}
			case SetVisualerMsg:
				if(pn.Visualer() != nil){
					if (pn.Visualer() != msg.Visualer) {
						//changing ?
					} else {
						//adding same again ?
					}
				} else {
					pn.SetVisualer(msg.Visualer)
				}
			default:
				pn.OtherPaneMsg(msg)
			}
		}
	}
}

func SwtkId() int {
	i := <- idChan
	return i
}
