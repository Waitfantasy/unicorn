package id

type GeneratorType int

type Data struct {
	sequence      uint64
	lastTimestamp uint64
}

const (
	MutexGeneratorType  GeneratorType = 1
	AtomicGeneratorType GeneratorType = 2
	SecondIdType                      = 0
	MilliSecondIdType                 = 1
)

type Generator interface {
	Make() (uint64, error)
	Extract(uuid uint64) (*MetaData)
}

type Factory interface {
	CreateGenerator(GeneratorType, *Meta) Generator
}

type GeneratorFactory struct {
}

func (f *GeneratorFactory) CreateGenerator(t GeneratorType, meta *Meta) Generator {
	switch t {
	case MutexGeneratorType:
		return NewMutexGenerator(meta)
	case AtomicGeneratorType:
		return NewAtomicGenerator(meta)
	default:
		return NewMutexGenerator(meta)
	}
}
