package swtk

import "flag"
import "io/ioutil"
import "code.google.com/p/freetype-go/freetype"
import "code.google.com/p/freetype-go/freetype/truetype"
import "log"

var (
	Dpi float64
	Fontfile string
	FontSize float64
	LineSpacing float64
	Font *truetype.Font
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

