package core

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
)

const FINDER_PATTERN_SIZE = 7
const QUIET_ZONE_SIZE = 4

var BLACK = color.RGBA{0, 0, 0, 255}
var WHITE = color.RGBA{255, 255, 255, 255}

type QRRenderer struct {
	data     []byte   // data
	scale    int      // 1 is means literal qr size
	version  int      // currently only 1
	ECLevel  string   // e.g. "L", "M", "Q", "H"
	mask     int      // 0-7 inclusive
	matrix   [][]byte // raw size matrix we are goona scale this with scale
	reserved [][]byte
	img      *image.RGBA
}

func (r *QRRenderer) SetConfig(data []byte, scale int, version int, mask int, ECLevel string) {
	r.scale = scale
	r.version = version
	r.mask = mask
	r.ECLevel = ECLevel
	r.data = data

	qrSize := r.getQrSize()

	r.img = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{(qrSize * scale) + (QUIET_ZONE_SIZE * scale * 2), (qrSize * scale) + (QUIET_ZONE_SIZE * scale * 2)}})
	r.matrix = make([][]byte, qrSize)
	r.reserved = make([][]byte, qrSize)
	for i := range r.matrix {
		r.matrix[i] = make([]byte, qrSize)
		r.reserved[i] = make([]byte, qrSize)
		for j := range r.matrix[i] {
			r.matrix[i][j] = 3   // 3 MEANS UNASSIGNED
			r.reserved[i][j] = 0 // 0 MEANS NOT RESERVED
		}
	}

	draw.Draw(r.img, r.img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
}

func (r *QRRenderer) getQrSize() int {
	return r.version*4 + 17
}

func (r *QRRenderer) getFinderPos() [][2]int {
	size := r.getQrSize()

	return [][2]int{
		{0, 0},
		{size - FINDER_PATTERN_SIZE, 0},
		{0, size - FINDER_PATTERN_SIZE},
	}
}

func (ren *QRRenderer) SetAlignments() {
	if ren.version == 1 {
		return
	}

	alignmentValues, exists := GetAlignmentValues(ren.version)
	posCount := len(alignmentValues)

	if !exists {
		panic("Invalid QR version for alignment positions")
	}

	for i := 0; i < int(posCount); i++ {
		for j := 0; j < int(posCount); j++ {
			if (i == 0 && j == 0) ||
				(i == 0 && j == int(posCount)-1) ||
				(i == int(posCount)-1 && j == 0) {
				continue // Skip alignments that overlap with finder patterns
			}

			row := alignmentValues[i]
			col := alignmentValues[j]

			for r := -2; r <= 2; r++ {
				for c := -2; c <= 2; c++ {
					if r == -2 || r == 2 || c == -2 || c == 2 ||
						(r == 0 && c == 0) {
						ren.matrix[row+r][col+c] = 1
						ren.reserved[row+r][col+c] = 1
					} else {
						ren.matrix[row+r][col+c] = 0
						ren.reserved[row+r][col+c] = 1
					}
				}
			}
		}
	}

}

func (ren *QRRenderer) SetFinderPattern() {
	finderPositions := ren.getFinderPos()
	size := ren.getQrSize()
	for _, pos := range finderPositions {
		row := pos[0] // 0 14 0
		col := pos[1] // 0 0 14

		for r := -1; r <= 7; r++ {
			if row+r <= -1 || size <= row+r {
				continue
			}
			for c := -1; c <= 7; c++ {
				if col+c <= -1 || size <= col+c {
					continue
				}

				if (r >= 0 && r <= 6 && (c == 0 || c == 6)) ||
					(c >= 0 && c <= 6 && (r == 0 || r == 6)) ||
					(r >= 2 && r <= 4 && c >= 2 && c <= 4) {
					ren.matrix[row+r][col+c] = 1
					ren.reserved[row+r][col+c] = 1
				} else {
					ren.matrix[row+r][col+c] = 0
					ren.reserved[row+r][col+c] = 1
				}
			}
		}
	}
}

func (r *QRRenderer) SetTimingPattern() {
	size := len(r.matrix)
	for i := 8; i < size-8; i++ {
		value := byte((i + 1) % 2)

		r.matrix[i][6] = value
		r.matrix[6][i] = value
		r.reserved[i][6] = 1
		r.reserved[6][i] = 1
	}
}

func (r *QRRenderer) SetFormatInfo() {
	data, exists := GetFormatValue(r.ECLevel, r.mask)
	if !exists {
		panic("Wrong EC or Mask")
	}
	binaryData := fmt.Sprintf("%015b", data)

	for i := 0; i <= 14; i++ {
		binary := binaryData[i] - 48

		if i <= 5 {
			r.matrix[i][8] = binary
			r.reserved[i][8] = 1
		} else if i == 6 {
			r.matrix[i+1][8] = binary
			r.reserved[i+1][8] = 1
		} else {
			rowIndex := r.getQrSize() - 8 + (i - 7)
			r.matrix[rowIndex][8] = binary
			r.reserved[rowIndex][8] = 1
		}

		if i < 7 {
			r.matrix[8][r.getQrSize()-i-1] = binary
			r.reserved[8][r.getQrSize()-i-1] = 1
		} else if i < 9 {
			r.matrix[8][15-i-1+1] = binary
			r.reserved[8][15-i-1+1] = 1
		} else {
			r.matrix[8][15-i-1] = binary
			r.reserved[8][15-i-1] = 1
		}
	}
}

func (r *QRRenderer) SetVersionInfo() {
	if r.version < 7 {
		return // Version info is only applicable for versions 7 and above
	}
	size := r.getQrSize()
	binaryData, exists := GetVersionBits(r.version)
	if !exists {
		panic("Invalid QR version for version info")
	}
	counter := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			bit := byte(binaryData[17-counter] - '0') // subtracts int32 values
			r.matrix[i][size-11+j] = bit              //BottomLeft
			r.matrix[size-11+j][i] = bit              //TopRight
			r.reserved[i][size-11+j] = 1
			r.reserved[size-11+j][i] = 1
			counter++
		}
	}
}

func (r *QRRenderer) SetData() {
	goingUp := true
	binary := ""
	for _, b := range r.data {
		binary += fmt.Sprintf("%08b", b)
	}
	binaryIndex := 0
	binary += strings.Repeat("0", GetRemainderBitCount(r.version))

	for i := r.getQrSize() - 1; i > 0; i -= 2 {
		if i == 6 {
			i--
		}

		if goingUp {
			for j := r.getQrSize() - 1; j >= 0; j-- {
				if r.matrix[i][j] == 3 {
					val, _ := strconv.ParseUint(string(binary[binaryIndex]), 2, 8)
					r.matrix[i][j] = byte(val)
					binaryIndex++
				}
				if r.matrix[i-1][j] == 3 {
					val, _ := strconv.ParseUint(string(binary[binaryIndex]), 2, 8)
					r.matrix[i-1][j] = byte(val)
					binaryIndex++
				}
			}
		} else {
			for j := 0; j < r.getQrSize(); j++ {
				if r.matrix[i][j] == 3 {
					val, _ := strconv.ParseUint(string(binary[binaryIndex]), 2, 8)
					r.matrix[i][j] = byte(val)
					binaryIndex++
				}
				if r.matrix[i-1][j] == 3 {
					val, _ := strconv.ParseUint(string(binary[binaryIndex]), 2, 8)
					r.matrix[i-1][j] = byte(val)
					binaryIndex++
				}
			}
		}
		goingUp = !goingUp
	}
}

func (r *QRRenderer) SetDarkModule() {
	r.matrix[8][(4*r.version)+9] = 1
	r.reserved[8][(4*r.version)+9] = 1
}

func (r *QRRenderer) ApplyMask() {
	maskFunc := maskPatterns[2]

	size := r.getQrSize()
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if r.reserved[row][col] == 1 { // Reserved bit
				continue
			}
			if maskFunc(row, col) {
				r.matrix[col][row] ^= 1 // flips bit
			}
		}
	}
}

func (r *QRRenderer) Save() error {
	fmt.Println(MaskEval1(r.matrix))
	for i := 0; i < len(r.matrix); i++ {
		for j := 0; j < len(r.matrix[0]); j++ {
			if r.matrix[i][j] == 1 {
				// r.img.SetRGBA(i, j, color.RGBA{0, 0, 0, 255})
				draw.Draw(r.img, image.Rect((QUIET_ZONE_SIZE*r.scale)+(i*r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale), (QUIET_ZONE_SIZE*r.scale)+(i*r.scale+r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale+r.scale)), &image.Uniform{BLACK}, image.Point{}, draw.Src)
			} else if r.matrix[i][j] == 0 {
				// r.img.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
				draw.Draw(r.img, image.Rect((QUIET_ZONE_SIZE*r.scale)+(i*r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale), (QUIET_ZONE_SIZE*r.scale)+(i*r.scale+r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale+r.scale)), &image.Uniform{WHITE}, image.Point{}, draw.Src)
			} else if r.matrix[i][j] == 4 {
				// 4 IS USED FOR DEBUGGING
				draw.Draw(r.img, image.Rect((QUIET_ZONE_SIZE*r.scale)+(i*r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale), (QUIET_ZONE_SIZE*r.scale)+(i*r.scale+r.scale), (QUIET_ZONE_SIZE*r.scale)+(j*r.scale+r.scale)), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)
			} else { // 3 UNASSIGNED
				continue
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
