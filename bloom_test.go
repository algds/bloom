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

func TestBloom(t *testing.T) {
	t.Parallel()

	for _, tc := range testcases {
		bl := New(tc.m, tc.hashes...)

		for _, data := range tc.data {
			bl.Add(data)

			if !bl.Contains(data) {
				t.Errorf("Newly added element is not found")
			}
		}
	}
}

func TestBloomPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected a panic for a nil hash function")
		}
	}()

	New(100, nil)
}

func BenchmarkBloom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range testcases {
			bl := New(tc.m, tc.hashes...)

			for _, data := range tc.data {
				bl.Add(data)
				_ = bl.Contains
			}
		}
	}
}
