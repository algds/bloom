package bloom

import "testing"

// nolint
var testcases = []struct {
	data   []interface{}
	m      uint
	hashes []Hash
}{
	{
		[]interface{}{"foo", "bar", "baz", "fizz", "buzz"},
		100,
		[]Hash{
			func(d interface{}) uint {
				s := d.(string)

				return uint(s[0])
			},
			func(d interface{}) uint {
				s := d.(string)

				return uint(s[1])
			},
			func(d interface{}) uint {
				s := d.(string)

				return uint(s[2])
			},
		},
	},
}

func TestBloomFilter(t *testing.T) {
	t.Parallel()

	for _, tc := range testcases {
		blf := New(tc.m, tc.hashes...)

		for _, data := range tc.data {
			blf.Add(data)

			if !blf.Contains(data) {
				t.Errorf("Newly added element is not found")
			}
		}
	}
}

func TestBloomFilterPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected a panic for a nil hash function")
		}
	}()

	New(100, nil)
}

func TestBloomFilterNotFound(t *testing.T) {
	t.Parallel()

	f1 := func(d interface{}) uint {
		s := d.(string)

		return uint(s[0])
	}
	f2 := func(d interface{}) uint {
		s := d.(string)

		return uint(s[1])
	}
	f3 := func(d interface{}) uint {
		s := d.(string)

		return uint(s[2])
	}

	blf := New(100, f1, f2, f3)

	blf.Add("foo")
	blf.Add("bar")
	blf.Add("baz")

	if blf.Contains("xyz") {
		t.Errorf("this value is so different it can't exist given the hashes above")
	}
}

func BenchmarkBloom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range testcases {
			blf := New(tc.m, tc.hashes...)

			for _, data := range tc.data {
				blf.Add(data)
				_ = blf.Contains
			}
		}
	}
}
