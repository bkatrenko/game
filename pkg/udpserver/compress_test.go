package udpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testData is just a random JSON example
const (
	testData = `{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}
`

	expectedCompressedSizePercent = 60
)

func TestCompress(t *testing.T) {
	compressor := NewCompressor()

	compressedData := compressor.Compress([]byte(testData))
	compressedSizePercent := len(compressedData) / (len(testData) / 100)

	assert.Less(t, compressedSizePercent, expectedCompressedSizePercent)
}

func TestDecompress(t *testing.T) {
	compressor := NewCompressor()

	compressedData := compressor.Compress([]byte(testData))
	decompressed, err := compressor.Decompress(compressedData)

	assert.Nil(t, err)
	assert.Equal(t, testData, string(decompressed))
}
