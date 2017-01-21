// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package bitset

import (
	"testing"
	"testing/quick"

	"github.com/tmthrgd/go-byte-test"
)

func TestSet(t *testing.T) {
	b := make(Bitset, 10)

	b.Set(50)

	if !bytetest.Test(b[:6], 0) || b[6] != 0x04 || !bytetest.Test(b[7:], 0) {
		t.Error("Set failed")
	}

	b.Set(60)

	if !bytetest.Test(b[:6], 0) || b[6] != 0x04 || b[7] != 0x10 || !bytetest.Test(b[8:], 0) {
		t.Error("Set failed")
	}
}

func TestClear(t *testing.T) {
	b := make(Bitset, 10)

	b.SetAll()

	b.Clear(50)

	if !bytetest.Test(b[:6], 0xff) || b[6] != 0xfb || !bytetest.Test(b[7:], 0xff) {
		t.Error("Clear failed")
	}
}

func TestInvert(t *testing.T) {
	b := make(Bitset, 10)

	b.Invert(50)

	if !bytetest.Test(b[:6], 0) || b[6] != 0x04 || !bytetest.Test(b[7:], 0) {
		t.Error("Invert failed")
	}

	b.Invert(50)

	if !bytetest.Test(b, 0) {
		t.Error("Invert failed")
	}
}

func TestSetRange(t *testing.T) {
	if err := quick.CheckEqual(func(size, start, end uint) []byte {
		b := New(size)

		for i := start; i < end; i++ {
			b.Set(i)
		}

		return b
	}, func(size, start, end uint) []byte {
		b := New(size)
		b.SetRange(start, end)
		return b
	}, &quick.Config{
		Values:        rangeTestValues,
		MaxCountScale: 100,
	}); err != nil {
		t.Error(err)
	}
}

func TestClearRange(t *testing.T) {
	if err := quick.CheckEqual(func(size, start, end uint) []byte {
		b := New(size)
		b.SetAll()

		for i := start; i < end; i++ {
			b.Clear(i)
		}

		return b
	}, func(size, start, end uint) []byte {
		b := New(size)
		b.SetAll()
		b.ClearRange(start, end)
		return b
	}, &quick.Config{
		Values:        rangeTestValues,
		MaxCountScale: 100,
	}); err != nil {
		t.Error(err)
	}
}

func TestInvertRange(t *testing.T) {
	if err := quick.CheckEqual(func(size, start, end uint) []byte {
		b := New(size)

		for i := start; i < end; i++ {
			b.Invert(i)
		}

		return b
	}, func(size, start, end uint) []byte {
		b := New(size)
		b.InvertRange(start, end)
		return b
	}, &quick.Config{
		Values:        rangeTestValues,
		MaxCountScale: 100,
	}); err != nil {
		t.Error(err)
	}
}

func BenchmarkSet(b *testing.B) {
	bs := make(Bitset, 10)

	for i := 0; i < b.N; i++ {
		bs.Set(50)
	}
}

func BenchmarkClear(b *testing.B) {
	bs := make(Bitset, 10)

	for i := 0; i < b.N; i++ {
		bs.Clear(50)
	}
}

func BenchmarkSetRange(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			bs := make(Bitset, size.l)
			l := bs.Len()

			if size.l > 1024 {
				b.ResetTimer()
			}

			for i := 0; i < b.N; i++ {
				bs.SetRange(1, l-1)
			}
		})
	}
}

func BenchmarkClearRange(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			bs := make(Bitset, size.l)
			l := bs.Len()

			if size.l > 1024 {
				b.ResetTimer()
			}

			for i := 0; i < b.N; i++ {
				bs.ClearRange(1, l-1)
			}
		})
	}
}
