package swtk

import "flag"
import "io/ioutil"
import "code.google.com/p/freetype-go/freetype"
import "code.google.com/p/freetype-go/freetype/truetype"
import "log"
import "image"
import "image/draw"

var (
	Dpi         float64
	Fontfile    string
	FontSize    float64
	LineSpacing float64
	Font        *truetype.Font
)

func init() {
	flag.Float64Var(&Dpi, "dpi", 72, "screen resolution in Dots Per Inch")
	flag.StringVar(&Fontfile, "fontfile", "arial.ttf", "filename of the ttf font")
	flag.Float64Var(&FontSize, "size", 12, "font size in points")
	flag.Float64Var(&LineSpacing, "spacing", 1.5, "line spacing (e.g. 2 means double spaced)")

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(Fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	font1, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	Font = font1
}

type PaneCoords struct {
	Pane   Pane
	Coords image.Point
	Size   image.Point
}

type PaneImage struct {
	Pane  Pane
	Image draw.Image
}

type alignment int

func (a alignment) Alignment() alignment {
	return a
}

const (
	AlignCenter alignment = iota
	AlignCenterLeft
	AlignCenterRight
	AlignTop
	AlignTopLeft
	AlignTopRight
	AlignBottom
	AlignBottomLeft
	AlignBottomRight
)

type MouseState struct {
	B    int8
	X, Y int16
}

//This is where on screen the pointer is
//There may be multiple pointers from multiple Devices
type PointerState struct {
	Device int
	Id     int
	X, Y   int
}

//This is where on screen has been 'touched'
// ie: with finger, or mouse button down
//There will certainly be multiple contacts from multiple Devices
type ContactState struct {
	Device int
	Id     int
	X, Y   int
}
