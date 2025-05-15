package core

// Sourced https://github.com/soldair/node-qrcode/blob/master/lib/core/galois-field.js

var expTable [512]byte
var logTable [256]byte

func InitTables() {
	// Generator α = 2
	x := 1
	for i := 0; i < 255; i++ {
		expTable[i] = byte(x)
		logTable[x] = byte(i)
		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11D // primitive polynomial
		}
	}

	// Duplicate for easy overflow handling in expTable
	for i := 256; i < 512; i++ {
		expTable[i] = expTable[i-256]
	}
}

func Log(n byte) byte {
	return logTable[n]
}

func Exp(n byte) byte {
	return expTable[n]
}

func Mul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return expTable[int(logTable[a])+int(logTable[b])]
}

// func Add(a, b byte) byte {
// 	return a ^ b // same as subtract
// }

// func Sub(a, b byte) byte {
// 	return a ^ b
// }

// func Mul(a, b byte) byte {
// 	if a == 0 || b == 0 {
// 		return 0
// 	}
// 	return expTable[int(logTable[a])+int(logTable[b])]
// }

// func Div(a, b byte) byte {
// 	if b == 0 {
// 		panic("divide by zero")
// 	}
// 	if a == 0 {
// 		return 0
// 	}
// 	return expTable[int(logTable[a])+255-int(logTable[b])]
// }
