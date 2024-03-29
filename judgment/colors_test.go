package judgment

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func hex(s string) color.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("hex: " + err.Error())
	}
	return c
}

func TestCreatePalette(t *testing.T) {
	type args struct {
		amountOfColors int
	}
	tests := []struct {
		name                   string
		args                   args
		want                   color.Palette
		expectedAmountOfColors int
	}{
		{
			name: "Palette of -1 is empty",
			args: args{
				amountOfColors: -1,
			},
			want: []color.Color{},
		},
		{
			name: "Palette of 0 is empty",
			args: args{
				amountOfColors: 0,
			},
			want: []color.Color{},
		},
		{
			name: "Palette of 1",
			args: args{
				amountOfColors: 1,
			},
			want: []color.Color{
				hex("#00a249"),
			},
		},
		{
			name: "Palette of 2",
			args: args{
				amountOfColors: 2,
			},
			want: []color.Color{
				hex("#df3222"),
				hex("#00a249"),
			},
		},
		{
			name: "Palette of 3",
			args: args{
				amountOfColors: 3,
			},
			want: []color.Color{
				hex("#df3222"),
				hex("#fab001"),
				hex("#00a249"),
			},
		},
		{
			name: "Palette of 4",
			args: args{
				amountOfColors: 4,
			},
			want: []color.Color{
				hex("#df3222"),
				hex("#fab001"),
				hex("#7bbd3e"),
				hex("#017a36"),
			},
		},
		{
			name: "Palette of 5",
			args: args{
				amountOfColors: 5,
			},
			want: []color.Color{
				hex("#df3222"),
				hex("#ed6f01"),
				hex("#fab001"),
				hex("#7bbd3e"),
				hex("#00a249"),
			},
		},
		{
			name: "Palette of 6",
			args: args{
				amountOfColors: 6,
			},
			expectedAmountOfColors: 6,
		},
		{
			name: "Palette of 7",
			args: args{
				amountOfColors: 7,
			},
			expectedAmountOfColors: 7,
		},
		{
			name: "Palette of 32",
			args: args{
				amountOfColors: 32,
			},
			expectedAmountOfColors: 32,
			want: []color.Color{
				hex("#df3222"),
				hex("#e3401d"),
				hex("#e64c18"),
				hex("#e95712"),
				hex("#eb630b"),
				hex("#ed6d02"),
				hex("#f07a00"),
				hex("#f38700"),
				hex("#f69300"),
				hex("#f8a000"),
				hex("#faac00"),
				hex("#f5b500"),
				hex("#ecbc00"),
				hex("#e2c300"),
				hex("#d7ca00"),
				hex("#cbd000"),
				hex("#bdd20c"),
				hex("#aece1d"),
				hex("#a0ca28"),
				hex("#92c531"),
				hex("#84c039"),
				hex("#76bc3e"),
				hex("#65b740"),
				hex("#54b142"),
				hex("#40ac44"),
				hex("#28a747"),
				hex("#00a148"),
				hex("#009944"),
				hex("#009141"),
				hex("#00893d"),
				hex("#008239"),
				hex("#017a36"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CreateDefaultPalette(tt.args.amountOfColors)
			print(DumpPaletteHexString(actual, ", ", "\"") + "\n")

			if tt.expectedAmountOfColors > 0 {
				assert.Equal(t, tt.expectedAmountOfColors, len(actual))
			}
			if nil != tt.want {
				for i, expectedColor := range tt.want {
					// our test values are not as precise as colorful's colors
					//assert.Equal(t, expectedColor, actual[i])
					// so we use equalish comparisons
					p := 300.0
					er, eg, eb, ea := expectedColor.RGBA()
					ar, ag, ab, aa := actual[i].RGBA()
					assert.InDelta(t, er, ar, p)
					assert.InDelta(t, eg, ag, p)
					assert.InDelta(t, eb, ab, p)
					assert.InDelta(t, ea, aa, p)
				}
			}
		})
	}
}
