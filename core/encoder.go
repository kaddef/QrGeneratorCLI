package core

import (
	"fmt"
	"strconv"
)

type RSEncoder struct { // Reed-Solomon Encoder
	PlainTextData    string
	PlainByteArray   []byte
	Encoding         string
	Length           int
	BinaryData       string
	DataByteArray    []byte
	EncodedByteArray []byte
	QrVersion        int
	ECLevel          string
}

func InitEncoder(version int, ecLevel string) *RSEncoder {
	return &RSEncoder{
		QrVersion: version,
		ECLevel:   ecLevel,
	}
}

func (r *RSEncoder) SetPlainMessage(msg string) {
	r.PlainTextData = msg
	r.Length = len(msg)
	r.Encoding = "Byte" // Numeric Alphanumeric Byte Kanji
	r.PlainByteArray = []byte(msg)
}

func (r *RSEncoder) Debug() {
	fmt.Print("Plain text: ")
	fmt.Println(r.PlainTextData)

	fmt.Print("PlainByteArray: ")
	fmt.Println(r.PlainByteArray)

	fmt.Print("Encoding: ")
	fmt.Print(r.Encoding)
	fmt.Println(" (", r.getEncodingBits(), ")")

	fmt.Print("Length: ")
	fmt.Print(r.Length)
	fmt.Println(" (", fmt.Sprintf("%08b", r.Length), ")")

	fmt.Print("BinaryData: ")
	fmt.Println(r.BinaryData)

	fmt.Print("DataByteArray: ")
	fmt.Println(r.DataByteArray)
}

func (r *RSEncoder) getEncodingBits() string {
	switch r.Encoding {
	case "Numeric":
		return "0001"
	case "Alphanumeric":
		return "0010"
	case "Byte":
		return "0100"
	case "Kanji":
		return "1000"
	default:
		panic("Invalid Encoding Type") // maybe return error
	}
}

func (r *RSEncoder) binaryToByte() {
	for i := 0; i < len(r.BinaryData); i += 8 {
		byteVal, _ := strconv.ParseUint(r.BinaryData[i:i+8], 2, 8)
		r.DataByteArray = append(r.DataByteArray, byte(byteVal))
	}
}

func (r *RSEncoder) CreateData() {
	encodingBits := r.getEncodingBits()
	var lengthBits string
	if r.QrVersion >= 1 && r.QrVersion <= 9 {
		lengthBits = fmt.Sprintf("%08b", r.Length)
	} else if r.QrVersion >= 10 && r.QrVersion <= 26 {
		lengthBits = fmt.Sprintf("%016b", r.Length)
	} else if r.QrVersion >= 27 && r.QrVersion <= 40 {
		lengthBits = fmt.Sprintf("%016b", r.Length)
	} else {
		panic("Invalid QR version for encoding")
	}

	r.BinaryData += encodingBits
	r.BinaryData += lengthBits

	for _, b := range r.PlainByteArray {
		r.BinaryData += fmt.Sprintf("%08b", b)
	}

	totalBitCount := GetTotalCodewordsCount(r.QrVersion)
	ecBitCount := GetECCodewordsCount(r.QrVersion, r.ECLevel)
	maxBits := totalBitCount*8 - ecBitCount*8
	remaining := maxBits - len(r.BinaryData)

	// Terminator Bits
	if remaining >= 4 {
		r.BinaryData += "0000"
		remaining -= 4
	} else {
		for range remaining {
			r.BinaryData += "0"
			remaining--
		}
	}

	// Padding to Multiple of 8, This Only Happens in Numeric or Alphanumeric not Byte
	if len(r.BinaryData)%8 != 0 {
		padLength := 8 - (len(r.BinaryData) % 8)
		for range padLength {
			r.BinaryData += "0"
			remaining--
		}
	}

	// After these step if still place left pad with this data
	extraPad := []string{"11101100", "00010001"}
	padIndex := 0
	for remaining >= 8 {
		r.BinaryData += extraPad[padIndex]
		padIndex = (padIndex + 1) % 2
		remaining -= 8
	}

	// THIS SHOULDNT HAPPEN IF WE HIT HERE SOMETHING IS WRONG
	if remaining > 0 {
		r.BinaryData += fmt.Sprintf("%0*s", remaining, "")
		panic("WE HIT A MINE")
	}

	r.binaryToByte()
}

func (r *RSEncoder) Encode() []byte {
	generator := GenerateECPolynomial(7) // 7 HARDCODED
	paddedData := make([]byte, len(r.DataByteArray)+7)
	copy(paddedData, r.DataByteArray)
	remainder := PolyMod(paddedData, generator)
	copy(paddedData[len(paddedData)-7:], remainder)
	r.EncodedByteArray = paddedData
	return r.EncodedByteArray
}
