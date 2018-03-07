package processor

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/codenaut/barcoder/images"
	"github.com/fogleman/gg"
	"gopkg.in/urfave/cli.v1"
)

type internal struct {
	config PageConfig
	height int
	width  int
	canvas *gg.Context
}

type PageConfig struct {
	Image  []ImageFileConfig
	Text   []TextConfig
	Height int
	Width  int
}
type ImageFileConfig struct {
	File       string
	FileInput  int
	Properties ImageConfig
}
type TextConfig struct {
	Input      int
	Value      string
	Font       string
	FontSize   float64
	Properties ImageConfig
}

type ImageConfig struct {
	Size     []int
	Position []int
	AutoTrim bool
	Border   int
	Rotate   float64
	Center   bool
	Top      bool
	Bottom   bool
	Left     bool
	Right    bool
}

func New(config PageConfig) *internal {

	obj := &internal{config: config, width: config.Width, height: config.Height}

	if config.Height != 0 && config.Width != 0 {
		obj.canvas = gg.NewContext(obj.width, obj.height)
	}
	return obj
}

func scale(factor float64, size []int) (int, int) {
	if len(size) > 1 {
		x := float64(size[0])
		y := float64(size[1])
		return int(factor * x), int(factor * y)
	}

	return 0, 0
}

func (t *internal) processText(value string, txt TextConfig, args cli.Args) error {
	str := value
	if str == "" {
		str = txt.Value
	}
	if str == "" {
		str = args.Get(txt.Input)
	}
	ctx := gg.NewContext(t.width, t.height)
	font := txt.Font
	fontSize := txt.FontSize
	if font == "" {
		font = "/Library/Fonts/Verdana.ttf"
	}
	if fontSize == 0 {
		fontSize = 72
	}
	if err := ctx.LoadFontFace(font, fontSize); err != nil {
		return err
	}
	w, h := ctx.MeasureString(str)

	strCtx := gg.NewContext(int(w*1.1), int(h*1.4))
	if err := strCtx.LoadFontFace(font, fontSize); err != nil {
		return err
	}
	strCtx.SetColor(image.White)
	strCtx.Clear()
	strCtx.SetColor(image.Black)
	x, y := 0, 0
	strCtx.DrawString(str, float64(x), float64(y)+h)
	img := strCtx.Image()

	t.insertImage(img, txt.Properties)
	return nil
}

func (t *internal) Process(output io.Writer, args cli.Args) error {

	for _, img := range t.config.Image {
		filename := img.File
		if filename == "" {
			filename = args.Get(img.FileInput)
		}
		if i, err := images.OpenPng(filename); err != nil {
			return err
		} else {
			flat := images.FlattenImage(i)
			if err := t.insertImage(flat, img.Properties); err != nil {
				return err
			}

		}
	}
	for _, txt := range t.config.Text {
		if err := t.processText("", txt, args); err != nil {
			return err
		}
	}

	if t.canvas != nil {
		if err := png.Encode(output, t.canvas.Image()); err != nil {
			return err
		}
	}
	return nil
}

func (t *internal) insertImage(img image.Image, imageConf ImageConfig) error {
	backgroundColour := uint16(65535)
	border := imageConf.Border

	size := img.Bounds().Size()
	maxX := 0
	maxY := 0
	minX := size.X
	minY := size.Y
	if imageConf.AutoTrim {
		for x := 0; x < size.X; x++ {
			for y := 0; y < size.Y; y++ {
				p := img.At(x, y)
				colour := color.Gray16Model.Convert(p).(color.Gray16).Y
				if colour != backgroundColour {
					if maxX < x {
						maxX = x
					}
					if maxY < y {
						maxY = y
					}
					if minX > x {
						minX = x
					}
					if minY > y {
						minY = y
					}
				}
			}
		}
	} else {
		maxX = size.X
		minX = 0
		maxY = size.Y
		minY = 0
	}

	height := maxY - minY + (border * 2)
	width := maxX - minX + (border * 2)
	if t.canvas == nil {
		t.canvas = gg.NewContext(width, height)
	} else if height > t.canvas.Height() || width > t.canvas.Width() {
		oldImage := t.canvas.Image()
		t.canvas = gg.NewContext(width, height)
		t.canvas.DrawImage(oldImage, 0, 0)
	}
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			p := img.At(x, y)
			t.canvas.SetColor(p)
			x1 := border + x - minX
			y1 := border + y - minY
			t.canvas.SetPixel(x1, y1)
		}
	}

	return nil

}
