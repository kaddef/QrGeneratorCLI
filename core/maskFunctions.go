package core

func Mask0(row, col int) bool { return (row+col)%2 == 0 }

func Mask1(row, col int) bool { return row%2 == 0 }

func Mask2(row, col int) bool { return col%3 == 0 }

func Mask3(row, col int) bool { return (row+col)%3 == 0 }

func Mask4(row, col int) bool { return (row/2+col/3)%2 == 0 }

func Mask5(row, col int) bool { return ((row*col)%2)+((row*col)%3) == 0 }

func Mask6(row, col int) bool { return (((row*col)%2)+((row*col)%3))%2 == 0 }

func Mask7(row, col int) bool { return (((row+col)%2)+((row*col)%3))%2 == 0 }

var maskPatterns = []func(row, col int) bool{
	Mask0, Mask1, Mask2, Mask3, Mask4, Mask5, Mask6, Mask7,
}

func MaskEvaluation(matrix [][]byte) {

}

func MaskEval1(matrix [][]byte) int {
	rows := len(matrix)
	cols := len(matrix[0])
	penalty := 0

	// Check rows
	for i := 0; i < rows; i++ {
		runColor := matrix[i][0]
		runLength := 1
		for j := 1; j < cols; j++ {
			if matrix[i][j] == runColor {
				runLength++
			} else {
				if runLength >= 5 {
					penalty += 3 + (runLength - 5)
				}
				runColor = matrix[i][j]
				runLength = 1
			}
		}
		if runLength >= 5 {
			penalty += 3 + (runLength - 5)
		}
	}

	// Check columns
	for j := 0; j < cols; j++ {
		runColor := matrix[0][j]
		runLength := 1
		for i := 1; i < rows; i++ {
			if matrix[i][j] == runColor {
				runLength++
			} else {
				if runLength >= 5 {
					penalty += 3 + (runLength - 5)
				}
				runColor = matrix[i][j]
				runLength = 1
			}
		}
		if runLength >= 5 {
			penalty += 3 + (runLength - 5)
		}
	}

	return penalty
}

func MaskEval2(matrix [][]byte) int {
	rows := len(matrix)
	cols := len(matrix[0])
	penalty := 0

	for i := 0; i < rows-1; i++ {
		for j := 0; j < cols-1; j++ {
			current := matrix[i][j]
			if matrix[i][j+1] == current &&
				matrix[i+1][j] == current &&
				matrix[i+1][j+1] == current {
				penalty += 3
			}
		}
	}

	return penalty
}

func MaskEval3(matrix [][]byte) int {
	rows := len(matrix)
	cols := len(matrix[0])
	penalty := 0

	// Not sure of these patterns
	pattern1 := []byte{1, 0, 1, 1, 1, 0, 1}
	pattern2 := []byte{0, 1, 0, 0, 0, 1, 0}

	// Check rows
	for i := 0; i < rows; i++ {
		for j := 0; j <= cols-7; j++ {
			segment := matrix[i][j : j+7]
			if matchesPattern(segment, pattern1) || matchesPattern(segment, pattern2) {
				penalty += 40
			}
		}
	}

	// Check columns
	for j := 0; j < cols; j++ {
		for i := 0; i <= rows-7; i++ {
			segment := make([]byte, 7)
			for k := 0; k < 7; k++ {
				segment[k] = matrix[i+k][j]
			}
			if matchesPattern(segment, pattern1) || matchesPattern(segment, pattern2) {
				penalty += 40
			}
		}
	}

	return penalty
}

func MaskEval4(matrix [][]byte) int {
	rows := len(matrix)
	cols := len(matrix[0])
	totalModules := rows * cols
	darkModules := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if matrix[i][j] == 1 {
				darkModules++
			}
		}
	}

	percentage := (darkModules * 100) / totalModules
	diff := abs(percentage - 50)
	penalty := (diff / 5) * 10

	return penalty
}

func matchesPattern(segment, pattern []byte) bool {
	for i := 0; i < 7; i++ {
		if segment[i] != pattern[i] {
			return false
		}
	}
	return true
}

// Integer abs
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
