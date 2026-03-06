package internal

type GoGenerator struct {
	Config Config
}

func NewGoGenerator(cfg Config) *GoGenerator {
	return &GoGenerator{Config: cfg}
}
