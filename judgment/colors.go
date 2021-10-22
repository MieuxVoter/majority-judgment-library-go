package judgment

// Some code below is taken directly from colorful's doc examples.
// https://github.com/lucasb-eyer/go-colorful/blob/master/doc/gradientgen/gradientgen.go

// May be useful later:
// c, err := colorful.Hex(s)

import (
	"errors"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
)

// CreateDefaultPalette returns a Palette of amountOfColors colors.
// 7 colors we use, red to green:
// "#df3222", "#ed6f01", "#fab001", "#c5d300", "#7bbd3e", "#00a249", "#017a36"
// When requiring more than 7, we interpolate in HSV space.
// This tries to be fault-tolerant, and returns an empty palette upon trouble.
func CreateDefaultPalette(amountOfColors int) color.Palette {
	const Color0 = 0xdf3222
	const Color1 = 0xed6f01
	const Color2 = 0xfab001
	const Color3 = 0xc5d300
	const Color4 = 0x7bbd3e
	const Color5 = 0x00a249
	const Color6 = 0x017a36

	color0 := hexToRGB(Color0)
	color1 := hexToRGB(Color1)
	color2 := hexToRGB(Color2)
	color3 := hexToRGB(Color3)
	color4 := hexToRGB(Color4)
	color5 := hexToRGB(Color5)
	color6 := hexToRGB(Color6)

	if amountOfColors < 0 {
		amountOfColors = amountOfColors * -1
	}

	switch amountOfColors {
	case 0:
		return []color.Color{}
	case 1:
		return []color.Color{
			color5,
		}
	case 2:
		return []color.Color{
			color0,
			color5,
		}
	case 3:
		return []color.Color{
			color0,
			color2,
			color5,
		}
	case 4:
		return []color.Color{
			color0,
			color2,
			color4,
			color6,
		}
	case 5:
		return []color.Color{
			color0,
			color1,
			color2,
			color4,
			color5,
		}
	case 6:
		return []color.Color{
			color0,
			color1,
			color2,
			color4,
			color5,
			color6,
		}
	case 7:
		return []color.Color{
			color0,
			color1,
			color2,
			color3,
			color4,
			color5,
			color6,
		}
	default:
		palette, err := bakePalette(amountOfColors, []color.Color{
			color0,
			color1,
			color2,
			color3,
			color4,
			color5,
			color6,
		})
		if err != nil {
			//panic("CreateDefaultPalette: failed to bake: "+err.Error())
			return []color.Color{}
		}
		return palette
	}
}

// DumpPaletteHexString dumps the provided palette as a string
// Looks like: "#df3222", "#ed6f01", "#fab001", "#c5d300", "#7bbd3e", "#00a249", "#017a36"
func DumpPaletteHexString(palette color.Palette, separator string, quote string) string {
	out := ""

	for colorIndex, colorRgba := range palette {
		if colorIndex > 0 {
			out += separator
		}
		out += quote
		out += DumpColorHexString(colorRgba, "#", false)
		out += quote
	}

	return out
}

// DumpColorHexString outputs strings like #ff3399 or #ff3399ff with alpha
// Be mindful that PRECISION IS LOST because hex format has less bits
func DumpColorHexString(c color.Color, prefix string, withAlpha bool) string {
	out := prefix
	r, g, b, a := c.RGBA()
	out += fmt.Sprintf("%02x", r>>8)
	out += fmt.Sprintf("%02x", g>>8)
	out += fmt.Sprintf("%02x", b>>8)
	if withAlpha {
		out += fmt.Sprintf("%02x", a>>8)
	}
	return out
}

// there probably is a colorful way to do this
func hexToRGB(hexColor int) color.Color {
	rgba := color.RGBA{
		R: uint8((hexColor & 0xff0000) >> 16),
		G: uint8((hexColor & 0x00ff00) >> 8),
		B: uint8((hexColor & 0x0000ff) >> 0),
		A: 0xff,
	}
	c, success := colorful.MakeColor(rgba)
	if !success {
		panic("hexToRgb")
	}

	return c
}

// This table contains the "key" colors of the color gradient we want to generate.
// The position of each key has to live in the range [0,1]
type keyColor struct {
	Color    colorful.Color
	Position float64
}
type gradientTable []keyColor

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt gradientTable) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Position <= t && t <= c2.Position {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Position) / (c2.Position - c1.Position)
			return c1.Color.BlendHcl(c2.Color, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Color
}

func bakePalette(toLength int, keyColors color.Palette) (color.Palette, error) {
	if toLength < 2 {
		return nil, errors.New("bakePalette: the length of the palette must be > 1")
	}

	keyPoints := gradientTable{}
	paletteLen := len(keyColors)

	for colorIndex, colorObject := range keyColors {
		colorfulColor, success := colorful.MakeColor(colorObject)
		if !success {
			panic("Bad palette color: alpha channel is probably 0.")
		}
		keyPoints = append(keyPoints, keyColor{
			Color:    colorfulColor,
			Position: float64(colorIndex) / (float64(paletteLen) - 1.0),
		})
	}

	outPalette := make([]color.Color, 0, 7)

	for i := 0; i < toLength; i++ {
		c := keyPoints.getInterpolatedColorFor(float64(i) / (float64(toLength) - 1))
		outPalette = append(outPalette, c)
	}

	return outPalette, nil
}
