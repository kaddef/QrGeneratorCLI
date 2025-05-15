package qrmeta

type QRVersionInfo struct {
	Size                 int // matrix size (width = height)
	TotalCodewords       int
	DataCodewords        int
	ErrorCorrectionWords int
}

var Versions = map[int]QRVersionInfo{
	1: {
		Size:                 21,
		TotalCodewords:       26,
		DataCodewords:        19,
		ErrorCorrectionWords: 7,
	},
	// Future:
	// 2: {Size: 25, TotalCodewords: 44, DataCodewords: 34, ErrorCorrectionWords: 10},
}
