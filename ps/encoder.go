package ps

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func utf16LeEncode(data []byte) []byte {
	result, _, err := transform.Bytes(
		unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(),
		data,
	)
	if err != nil {
		panic(err)
	}
	return result
}

func cp866Decode(data []byte) []byte {
	result, err := charmap.CodePage866.NewDecoder().Bytes(data)
	if err != nil {
		panic(err)
	}
	return result
}
