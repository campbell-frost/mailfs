package gmail

type Gmail struct {
	maxChunkSize int
}

func New(maxChunkSize int) *Gmail {
	return &Gmail{
		maxChunkSize: maxChunkSize,
	}
}

func (g *Gmail) MaxChunkSize() int {
	return g.maxChunkSize
}
