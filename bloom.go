package bloom

const bits = 64

// Hash is a function that takes any input
// and should return a uniformly randomly
// distributed 'uint'.
type Hash func(interface{}) uint

// Filter is the Bloom Filter interface of allowable actions.
type Filter interface {

	// Contains returns true if it's probable that
	// this value was added to the Bloom Filter.
	Contains(interface{}) bool

	// Add puts a new value into the Bloom Filter.
	Add(interface{})
}

type filter struct {
	bits   []uint64
	m      uint
	hashes []Hash
}

func (f *filter) Contains(d interface{}) bool {
	for _, bh := range f.hashes {
		realPos := bh(d)
		index := realPos / bits
		offset := realPos % bits

		if res := f.bits[index] & (uint64(1) << offset); res == 0 {
			return false
		}
	}

	return true
}

func (f *filter) Add(d interface{}) {
	for _, bh := range f.hashes {
		realPos := bh(d)
		index := realPos / bits
		offset := realPos % bits

		f.bits[index] |= (uint64(1) << offset)
	}
}

// New returns a new Bloom Filter of size 'm' bits and with the
// specified number of hash functions.
func New(m uint, k ...Hash) Filter {
	for _, f := range k {
		if f == nil {
			panic("can't pass nil as a hash function")
		}
	}

	return &filter{
		bits:   make([]uint64, (m/bits)+1),
		m:      m,
		hashes: k,
	}
}
