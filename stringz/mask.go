package stringz

// MaskPrefix masks the prefix of a string with a mask.
//
// e.g.
//
//	MaskPrefix("ABCDEFGH", "*", 3) // returns "*****FGH"
//	MaskPrefix("ABCDEFGH", "*", 4) // returns "****EFGH"
func MaskPrefix(s, mask string, unmaskLen int) (masked string) {
	for i, r := range s {
		if len(s)-unmaskLen <= i {
			masked += string(r)
			continue
		}
		masked += mask
	}
	return masked
}

// MaskSuffix masks the suffix of a string with a mask.
//
// e.g.
//
//	MaskSuffix("ABCDEFGH", "*", 3) // returns "ABC*****"
//	MaskSuffix("ABCDEFGH", "*", 4) // returns "ABCD****"
func MaskSuffix(s, mask string, unmaskLen int) (masked string) {
	for i, r := range s {
		if unmaskLen > i {
			masked += string(r)
			continue
		}
		masked += mask
	}
	return masked
}
