package main

import (
	"fmt"
	"qrGenerator/core"
	"strconv"
)

func main() {
	msg := "HELLO WORLD"
	byteArray := []byte(msg)
	// fmt.Println(byteArray)
	bitStream := "0100"
	bitStream += fmt.Sprintf("%08b", len(msg))
	for _, b := range byteArray {
		bitStream += fmt.Sprintf("%08b", b)
	}
	maxBits := 152
	remaining := maxBits - len(bitStream)
	if remaining >= 4 {
		bitStream += "0000"
	} else {
		for range remaining {
			bitStream += "0"
		}
	}
	if len(bitStream)%8 != 0 {
		fmt.Println("This should not happen")
		padding := 8 - (len(bitStream) % 8)
		for range padding {
			bitStream += "0"
		}
	}
	padBytes := []string{"11101100", "00010001"} //0xEC, 0x11
	padIndex := 0
	for len(bitStream) < maxBits {
		bitStream += padBytes[padIndex%2]
		padIndex++
	}
	bytes := []byte{}
	for i := 0; i < len(bitStream); i += 8 {
		byteVal, _ := strconv.ParseUint(bitStream[i:i+8], 2, 8)
		bytes = append(bytes, byte(byteVal))
	}
	fmt.Println(bytes)

	core.InitTables()
	// generator := []byte{87, 229, 146, 149, 238, 102, 21}
	genPoly := core.GenerateECPolynomial(7)

	paddedData := make([]byte, len(bytes)+7)

	copy(paddedData, bytes)
	fmt.Println(paddedData)

	remainder := core.PolyMod(paddedData, genPoly)
	fmt.Println(remainder)

	r := core.QRRenderer{}
	r.SetConfig(1, 1, 2, "L")
	r.SetFinderPattern(1, 1, 1, 1)
	r.SetTimingPattern()
	r.SetFormatInfo()
	r.SetDarkModule()
	r.SetData()
	r.ApplyMask()
	r.Save()
	fmt.Println(core.GenerateECPolynomial(7))

	r.SetFormatInfo()
}
