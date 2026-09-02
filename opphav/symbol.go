package opphav

type Symbol interface {
	Name() string
	Index() uint8
	IsZero() bool

	isSymbol()
}

type Sign interface {
	Symbol
	Observable

	isSign()
}

type Dimension interface {
	Symbol

	isDimension()
}

type Category interface {
	Dimension

	isCategory()
}

type Spectrum interface {
	Dimension

	isSpectrum()
}

type entry struct {
	name  string
	index uint8
}

func (e entry) Name() string {
	return e.name
}

func (e entry) Index() uint8 {
	return e.index
}

func (e entry) String() string {
	return e.name
}

func (e entry) IsZero() bool {
	return e.name == ""
}

type SignEntry struct {
	entry
}

func (SignEntry) isSymbol() {}
func (SignEntry) isSign()   {}

type CategoryEntry struct {
	entry
}

func (CategoryEntry) isSymbol()    {}
func (CategoryEntry) isDimension() {}
func (CategoryEntry) isCategory()  {}

type SpectrumEntry struct {
	entry
}

func (SpectrumEntry) isSymbol()    {}
func (SpectrumEntry) isDimension() {}
func (SpectrumEntry) isSpectrum()  {}
