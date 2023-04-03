package compress

type Compressor interface {
	Code() byte
	Compress(data []byte) ([]byte, error)
	Uncompress(data []byte) ([]byte, error)
}

type DoNothingCompressor struct {
}

func (d DoNothingCompressor) Code() byte {
	return 0
}

func (d DoNothingCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (d DoNothingCompressor) Uncompress(data []byte) ([]byte, error) {
	return data, nil
}
