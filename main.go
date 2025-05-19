package main

import (
	"qrGenerator/core"
)

func main() {
	core.InitTables()

	e := core.RSEncoder{}
	e.SetPlainMessage("deneme")
	e.CreateData()
	encodedData := e.Encode()
	// e.Debug()

	r := core.QRRenderer{}
	r.SetConfig(encodedData, 1, 1, 2, "L")
	r.SetFinderPattern(1, 1, 1, 1)
	r.SetTimingPattern()
	r.SetFormatInfo()
	r.SetDarkModule()
	r.SetData()
	r.ApplyMask()
	r.Save()
}
