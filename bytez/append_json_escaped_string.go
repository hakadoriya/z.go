package bytez

//nolint:cyclop
func AppendJSONEscapedString(dst []byte, s string) []byte {
	for i := range len(s) {
		if s[i] != '"' && s[i] != '\\' && s[i] > 0x1F {
			dst = append(dst, s[i])

			continue
		}

		// cf. https://tools.ietf.org/html/rfc8259#section-7
		// ... MUST be escaped: quotation mark, reverse solidus, and the control characters (U+0000 through U+001F).
		switch s[i] {
		case '"', '\\':
			dst = append(dst, '\\', s[i])
		case '\b' /* 0x08 */ :
			dst = append(dst, '\\', 'b')
		case '\f' /* 0x0C */ :
			dst = append(dst, '\\', 'f')
		case '\n' /* 0x0A */ :
			dst = append(dst, '\\', 'n')
		case '\r' /* 0x0D */ :
			dst = append(dst, '\\', 'r')
		case '\t' /* 0x09 */ :
			dst = append(dst, '\\', 't')
		default:
			const hexTable string = "0123456789abcdef"
			// cf. https://github.com/golang/go/blob/70deaa33ebd91944484526ab368fa19c499ff29f/src/encoding/hex/hex.go#L28-L29
			dst = append(dst, '\\', 'u', '0', '0', hexTable[s[i]>>4], hexTable[s[i]&0x0f])
		}
	}

	return dst
}
