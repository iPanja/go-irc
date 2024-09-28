package parser

type Parser interface {
	ExtractData(buf []byte) map[string]string
}
