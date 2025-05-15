package core

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const FINDER_PATTERN_SIZE = 7

var BLACK = color.RGBA{0, 0, 0, 255}
var WHITE = color.RGBA{255, 255, 255, 255}

type QRRenderer struct {
	// data    []byte // data
	scale   int    // 1 is means literal qr size
	version int    // currently only 1
	ECLevel string // e.g. "L", "M", "Q", "H"
	mask    int    // 0-7 inclusive
	matrix  [][]byte
	img     *image.RGBA
}

func (r *QRRenderer) SetConfig(scale int, version int, mask int, ECLevel string) {
	r.scale = scale
	r.version = 1
	r.mask = mask
	r.ECLevel = ECLevel

	r.img = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{21, 21}})
	r.matrix = make([][]byte, 21)
	for i := range r.matrix {
		r.matrix[i] = make([]byte, 21)
	}

	gray := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	draw.Draw(r.img, r.img.Bounds(), &image.Uniform{C: gray}, image.Point{}, draw.Src)
}

func (r *QRRenderer) getFinderPos() [][2]int {
	size := 21 // Get qr size from version here

	return [][2]int{
		{0, 0},
		{size - FINDER_PATTERN_SIZE, 0},
		{0, size - FINDER_PATTERN_SIZE},
	}
}

func (r *QRRenderer) getQRSize() int {
	return 21 + (r.version-1)*4
}

func (ren *QRRenderer) SetFinderPattern(x, y, scale, orientation int) {
	finderPositions := ren.getFinderPos()
	fmt.Println(finderPositions)
	for _, pos := range finderPositions {
		row := pos[0] // 0 14 0
		col := pos[1] // 0 0 14

		for r := -1; r <= 7; r++ {
			if row+r <= -1 || 21 <= row+r {
				continue
			}
			for c := -1; c <= 7; c++ {
				if col+c <= -1 || 21 <= col+c {
					continue
				}

				if (r >= 0 && r <= 6 && (c == 0 || c == 6)) ||
					(c >= 0 && c <= 6 && (r == 0 || r == 6)) ||
					(r >= 2 && r <= 4 && c >= 2 && c <= 4) {
					ren.matrix[row+r][col+c] = 1
				} else {
					ren.matrix[row+r][col+c] = 0
				}
			}
		}
	}
}

func (r *QRRenderer) SetTimingPattern() {
	for i := 8; i <= len(r.matrix)-8; i += 2 {
		r.matrix[i][6] = 1
		r.matrix[6][i] = 1
	}
}

func (r *QRRenderer) SetFormatInfo() {
	data, exists := GetFormatValue(r.ECLevel, r.mask)
	if !exists {
		panic("Wrong EC or Mask")
	}
	binaryData := fmt.Sprintf("%08b", data)
	fmt.Println(binaryData)

	for i := range 15 {
		binary := binaryData[i] - 48

		if i < 6 {
			r.matrix[i][8] = binary
		} else if i < 8 {
			r.matrix[i][8] = binary
		} else {
			r.matrix[r.getQRSize()-15+i][8] = binary
		}

		if i < 8 {
			r.matrix[8][r.getQRSize()-i-1] = binary
		} else if i < 9 {
			r.matrix[8][15-i-1+1] = binary
		} else {
			r.matrix[8][15-i-1] = binary
		}
	}
}

func (r *QRRenderer) SetDarkModule() {
	r.matrix[8][(4*r.version)+9] = 1
}

func (r *QRRenderer) Save() error {
	for rowIndex, row := range r.matrix {
		for colIndex := range row {
			if r.matrix[rowIndex][colIndex] == 1 {
				r.img.SetRGBA(rowIndex, colIndex, color.RGBA{0, 0, 0, 255})
			} else {
				continue
				r.img.SetRGBA(rowIndex, colIndex, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	file, err := os.Create("output.png")
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, r.img)
}
