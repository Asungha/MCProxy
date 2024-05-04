package utils

type HexWrapper struct {
	Data  *[]byte
	Start int
	End   int
}

func NewHexWrapper(b *[]byte, start int, end int) *HexWrapper {
	return &HexWrapper{Data: b, Start: start, End: end}
}

func (h *HexWrapper) Get() []byte {
	if h.Start < 0 || h.End < 0 || h.Start > h.End {
		return []byte{}
	}
	if h.Start > len(*h.Data) || h.End > len(*h.Data) {
		return []byte{}
	}
	return (*h.Data)[h.Start:h.End]
}
