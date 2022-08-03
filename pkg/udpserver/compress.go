package udpserver

import (
	"github.com/klauspost/compress/zstd"
)

type (
	// compressor is an interface that responsible for compress/decompress data
	// for send and receive
	compressor interface {
		Compress(src []byte) []byte
		Decompress(src []byte) ([]byte, error)
	}

	compressorImpl struct {
		encoder *zstd.Encoder
		decoder *zstd.Decoder
	}
)

// NewCompressor is a constructor the the new compressor instance.
// It is panicking in case of any error while it is a necessary part of the application -
// and without compressor it is better to crash the application and fix the issue
func NewCompressor() compressor {
	// Create a writer that caches compressors.
	// For this operation type we supply a nil Reader.
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		panic(err)
	}
	// NewReader creates a new decoder.
	// A nil Reader can be provided in which case Reset can be used to start a decode.
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}

	return &compressorImpl{
		encoder: encoder,
		decoder: decoder,
	}
}

// Compress a buffer.
// If you have a destination buffer, the allocation in the call can also be eliminated.
func (c *compressorImpl) Compress(src []byte) []byte {
	return c.encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func (c *compressorImpl) Decompress(src []byte) ([]byte, error) {
	return c.decoder.DecodeAll(src, nil)
}
