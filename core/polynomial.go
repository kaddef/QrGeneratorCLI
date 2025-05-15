package core

func PolyMul(p1, p2 []byte) []byte {
	coeff := make([]byte, len(p1)+len(p2)-1)

	for i := 0; i < len(p1); i++ {
		for j := 0; j < len(p2); j++ {
			coeff[i+j] ^= Mul(p1[i], p2[j])
		}
	}
	return coeff
}

func PolyMod(divident, divisor []byte) []byte {
	result := make([]byte, len(divident))
	copy(result, divident)

	for len(result)-len(divisor) >= 0 {
		coeff := result[0]

		for i := 0; i < len(divisor); i++ {
			result[i] ^= Mul(divisor[i], coeff)
		}

		offset := 0
		for offset < len(result) && result[offset] == 0 {
			offset++
		}
		result = result[offset:]
	}
	return result
}

func GenerateECPolynomial(degree int) []byte {
	poly := []byte{1}

	for i := 0; i < degree; i++ {
		term := []byte{1, Exp(byte(i))}
		poly = PolyMul(poly, term)
	}

	return poly
}
