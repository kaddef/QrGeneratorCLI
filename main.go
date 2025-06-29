package main

import (
	"flag"
	"qrGenerator/core"
)

func main() {
	version := flag.Int("V", -1, "QR version (1-40)")
	mask := flag.Int("M", -1, "Mask pattern (0-7)")
	scale := flag.Int("S", -1, "Scale of the QR image")
	ecLevel := flag.String("EC", "", "Error correction level: L, M, Q, or H")
	message := flag.String("MSG", "", "Message to encode in the QR code")

	// Parse flags
	flag.Parse()

	if *message == "" {
		panic("Empty Message")
	}
	if *ecLevel == "" {
		*ecLevel = "L" // Default to Low error correction level
	}
	if *scale == -1 {
		*scale = 1 // Determine this according to the size of qr code
	}
	if *mask == -1 {
		*mask = 2 // Mask is gonna be determined dynamically but we can take for more customization
	}
	if *version == -1 {
		*version = 3 // Default version, can be determined dynamically based on message length and error correction level
	}

	core.InitTables()

	e := core.InitEncoder(*version, *ecLevel) //e := core.InitEncoder(3, "H")
	e.SetPlainMessage(*message)               //e.SetPlainMessage("deneme")
	e.CreateData()
	encodedData := e.Encode()
	e.Debug()

	r := core.QRRenderer{}
	r.SetConfig(encodedData, *scale, *version, *mask, *ecLevel) //r.SetConfig(encodedData, 1, 3, 2, "H")
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

// DONE: Implement mask patterns
// DONE: Implement reserved matrix — it prevents the mask from being applied to static patterns
// DONE: Divide data into blocks and apply Reed-Solomon encoding each block
// DONE: In ApplyMask, use only the reserved matrix
// TODO: Add support for modes other than Byte (Numeric, Alphanumeric, Kanji)
// TODO: Eliminate lookup tables calculate them dynamically
// TODO: Add ECI (Extended Channel Interpretation) compatibility
// TODO: Dynamically determine the optimal mask pattern
// TODO: Dynamically determine the version based on the message length and error correction level

// go run main.go -EC L -S 5 -V 2 -MSG "Türkiye"
