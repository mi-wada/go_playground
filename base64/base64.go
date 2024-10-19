package base64

const (
	StdPadding = '='
)

var (
	StdEncoding = NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	URLEncoding = NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
)

type Encoding struct {
	encoder []byte
	decoder map[byte]byte
	padding rune
}

func NewEncoding(encoder string) *Encoding {
	if len(encoder) != 64 {
		panic("`encoder` is not 64-bytes long")
	}

	decoder := make(map[byte]byte)
	for i, c := range encoder {
		decoder[byte(c)] = byte(i)
	}

	return &Encoding{encoder: []byte(encoder), decoder: decoder, padding: StdPadding}
}

func (e Encoding) WithPadding(padding rune) *Encoding {
	e.padding = padding

	return &e
}

func (e *Encoding) Encode(src []byte) []byte {
	// Base64 encoding handle 6bits as 1byte(8bits), so allocate len(src)*8/6 cap.
	res := make([]byte, 0, len(src)*8/6)

	for i := 0; i < len(src)*8; i += 6 {
		var sixBits byte
		for j := i; j < i+6; j++ {
			byteIdx := j / 8
			bitIdx := 7 - j%8

			var bit byte
			if byteIdx >= len(src) {
				bit = 0
			} else {
				bit = (src[byteIdx] & (1 << bitIdx)) >> bitIdx
			}
			sixBits += (bit << (5 - (j % 6)))
		}

		res = append(res, e.encoder[sixBits])
	}

	// Add padding
	for len(res)%4 != 0 {
		res = append(res, byte(e.padding))
	}

	return res
}

func (e *Encoding) Decode(src []byte) []byte {
	res := make([]byte, 0, len(src)*6/8)
	for i := 0; i < 6*len(src); i += 8 {
		var char byte
		meetsEq := false

		for j := i; j < i+8; j++ {
			byteIdx := j / 6
			bitIdx := 5 - j%6

			if src[byteIdx] == byte(e.padding) {
				meetsEq = true
				break
			}

			var bit byte
			bit = (e.decoder[src[byteIdx]] & (1 << bitIdx)) >> bitIdx
			char += (bit << (7 - (j % 8)))
		}

		if meetsEq {
			break
		}
		res = append(res, char)
	}
	return res
}
