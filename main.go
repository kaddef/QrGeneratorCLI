package main

import (
	// "flag"
	// "fmt"
	// "os"
	"qrGenerator/core"
)

func main() {
	// version := flag.Int("version", 1, "QR version (1-40)")
	// mask := flag.Int("mask", -1, "Mask pattern (0-7), -1 for automatic")
	// scale := flag.Int("scale", 4, "Scale of the QR image")
	// ecLevel := flag.String("ecLevel", "L", "Error correction level: L, M, Q, or H")

	// // Parse flags
	// flag.Parse()

	// // Positional argument: message
	// if flag.NArg() == 0 {
	// 	fmt.Println("Error: message is required.")
	// 	fmt.Println("Usage: qrgen [options] <message>")
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	// message := flag.Arg(0)

	// fmt.Println("Generating QR code with:")
	// fmt.Println("Message:", message)
	// fmt.Println("Version:", *version)
	// fmt.Println("Mask:", *mask)
	// fmt.Println("Scale:", *scale)
	// fmt.Println("Error Correction Level:", *ecLevel)

	core.InitTables()

	e := core.InitEncoder(1, "L")
	e.SetPlainMessage("Türkiye")
	e.CreateData()
	encodedData := e.Encode()
	e.Debug()

	r := core.QRRenderer{}
	r.SetConfig(encodedData, 1, 1, 2, "L")
	r.SetFinderPattern()
	r.SetTimingPattern()
	r.SetFormatInfo()
	r.SetDarkModule()
	r.SetAlignments()  // Works After version 2
	r.SetVersionInfo() // Works After version 7
	r.SetData()
	r.ApplyMask()
	r.Save()
}

// TODO Add Other modes than byte Numeric Alphanumeric Kanji
// TODO Implement Mask Patterns
// TODO Implement reserved matrix. The reserved matrix prevents the mask from being applied to the alignment patterns.
// TODO Get rid of lookup tables
