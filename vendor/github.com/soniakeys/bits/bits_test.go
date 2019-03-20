// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package bits_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/soniakeys/bits"
)

func ExampleNew() {
	b := bits.New(80)
	fmt.Printf("%#v\n", b)
	// Output:
	// bits.Bits{Num:80, Bits:[]uint64{0x0, 0x0}}
}

func ExampleNewGivens() {
	b := bits.NewGivens(0, 63, 65, 2)
	fmt.Println(b.Slice())
	// Output:
	// [0 2 63 65]
}

func ExampleBits_AllOnes() {
	b := bits.New(5)
	b.SetAll()
	fmt.Println(b.AllOnes())
	b.SetBit(2, 0)
	fmt.Println(b.AllOnes())
	// Output:
	// true
	// false
}

func ExampleBits_AllZeros() {
	b := bits.New(5)
	fmt.Println(b.AllZeros())
	b.SetBit(2, 1)
	fmt.Println(b.AllZeros())
	// Output:
	// true
	// false
}

func ExampleBits_And() {
	x := bits.NewGivens(3, 5, 6)
	y := bits.NewGivens(4, 5, 6)
	x.And(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [5 6]
}

func ExampleBits_AndNot() {
	x := bits.NewGivens(3, 5, 6)
	y := bits.NewGivens(4, 5, 6)
	x.AndNot(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [3]
}

func ExampleBits_Bit() {
	b := bits.NewGivens(0, 63, 65, 2)
	for _, n := range []int{0, 1, 2, 63, 64, 65} {
		fmt.Printf("bit %d: %d\n", n, b.Bit(n))
	}
	// Output:
	// bit 0: 1
	// bit 1: 0
	// bit 2: 1
	// bit 63: 1
	// bit 64: 0
	// bit 65: 1
}

func ExampleBits_ClearAll() {
	x := bits.NewGivens(3, 5, 6)
	fmt.Println(x)
	x.ClearAll()
	fmt.Println(x)
	// Output:
	// 1101000
	// 0000000
}

func ExampleBits_ClearBits() {
	b := bits.NewGivens(0, 2, 63, 65)
	b.ClearBits(0, 63)
	fmt.Println(b.Slice())
	// Output:
	// [2 65]
}

func ExampleBits_Equal() {
	a := bits.NewGivens(1, 3)
	b := bits.NewGivens(3, 1)
	// Bits values with same bit numbers set compare equal
	fmt.Println(a, b, a.Equal(b))
	// Output:
	// 1010 1010 true
}

func ExampleBits_IterateOnes() {
	b := bits.NewGivens(0, 63, 65, 2)
	b.IterateOnes(func(n int) bool {
		fmt.Print(n, " ")
		return true
	})
	fmt.Println()
	// Output:
	// 0 2 63 65
}

func ExampleBits_IterateZeros() {
	b := bits.NewGivens(0, 63, 65, 2)
	b.Not(b)
	fmt.Println(b)
	b.IterateZeros(func(n int) bool {
		fmt.Print(n, " ")
		return true
	})
	fmt.Println()
	// Output:
	// 010111111111111111111111111111111111111111111111111111111111111010
	// 0 2 63 65
}

func ExampleBits_Not() {
	x := bits.NewGivens(3, 5, 6)
	x.Not(x)
	fmt.Println(x.Slice())
	// Output:
	// [0 1 2 4]
}

func ExampleBits_OneFrom() {
	b := bits.NewGivens(0, 63, 65, 2)
	for n := 0; ; n++ {
		n = b.OneFrom(n)
		if n < 0 {
			break
		}
		fmt.Print(n, " ")
	}
	fmt.Println()
	// Output:
	// 0 2 63 65
}

func ExampleBits_OneFrom_sieve() {
	q := 8
	n := q * q
	b := bits.New(n)
	b.SetAll()
	b.ClearBits(0, 1)
	for p := 2; p < q; p = b.OneFrom(p + 1) {
		for c := p + p; c < n; c += p {
			b.SetBit(c, 0)
		}
	}
	fmt.Println(b.Slice())
	// Output:
	// [2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61]
}

func ExampleBits_Or() {
	x := bits.NewGivens(3, 5, 6)
	y := bits.NewGivens(4, 5, 6)
	x.Or(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [3 4 5 6]
}

func ExampleBits_OnesCount() {
	b := bits.NewGivens(0, 2, 128)
	fmt.Println(b.OnesCount())
	// Output:
	// 3
}

func ExampleBits_Set() {
	x := bits.NewGivens(0, 2)
	var z bits.Bits
	z.Set(x)
	fmt.Println(z.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_SetAll() {
	b := bits.New(5)
	b.SetAll()
	fmt.Println(b.Slice())
	// Output:
	// [0 1 2 3 4]
}

func ExampleBits_SetBit() {
	b := bits.New(5)
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	fmt.Println(b.Slice())
	// Output:
	// [0 2]
}

func ExampleBits_SetBits() {
	b := bits.NewGivens(2, 65)
	b.SetBits(0, 63)
	fmt.Println(b.Slice())
	// Output:
	// [0 2 63 65]
}

func ExampleBits_Single() {
	x := bits.NewGivens(0, 2)
	y := bits.NewGivens(129)
	var z bits.Bits

	fmt.Println(x.OnesCount(), "bits, single =", x.Single())
	fmt.Println(y.OnesCount(), "bit,  single =", y.Single())
	fmt.Println(z.OnesCount(), "bits, single =", z.Single())
	// Output:
	// 2 bits, single = false
	// 1 bit,  single = true
	// 0 bits, single = false
}

func ExampleBits_Slice() {
	b := bits.NewGivens(0, 63, 65, 2)
	fmt.Println(b.Slice())
	// Output:
	// [0 2 63 65]
}

func ExampleBits_String() {
	b := bits.New(66)
	b.SetBits(0, 2, 63, 64)
	fmt.Println("bit 65                                                       bit 0")
	fmt.Println("|                                                                |")
	fmt.Println("v                                                                v")
	fmt.Println(b.String())
	fmt.Println(b)
	// Output:
	// bit 65                                                       bit 0
	// |                                                                |
	// v                                                                v
	// 011000000000000000000000000000000000000000000000000000000000000101
	// 011000000000000000000000000000000000000000000000000000000000000101
}

func ExampleBits_Xor() {
	x := bits.NewGivens(3, 5, 6)
	y := bits.NewGivens(4, 5, 6)
	x.Xor(x, y)
	fmt.Println(x.Slice())
	// Output:
	// [3 4]
}

func ExampleBits_ZeroFrom() {
	b := bits.NewGivens(0, 63, 65, 2)
	b.Not(b)
	fmt.Println(b)
	for n := 0; ; n++ {
		n = b.ZeroFrom(n)
		if n < 0 {
			break
		}
		fmt.Print(n, " ")
	}
	fmt.Println()
	// Output:
	// 010111111111111111111111111111111111111111111111111111111111111010
	// 0 2 63 65
}

// Tests probe some boundary conditions and push coverage to 100%

func TestNew(t *testing.T) {
	// test that proper length slice is allocated
	for _, tc := range []struct{ n, l int }{
		{0, 0},
		{1, 1},
		{64, 1},
		{65, 2},
	} {
		b := bits.New(tc.n)
		if len(b.Bits) != tc.l {
			t.Fatal("len(b.Bits) = ", len(b.Bits), " want ", tc.l)
		}
	}
	// test that negative bit number panics
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	bits.New(-1)
}

func TestNewGivens(t *testing.T) {
	// test negative bit number panics
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	bits.NewGivens(0, -1, 3)
}

func TestAllOnes(t *testing.T) {
	// exercise early return
	b := bits.NewGivens(63, 64)
	b.Not(b)
	if b.AllOnes() {
		t.Fatal("real problem")
	}
}

func TestAllZeros(t *testing.T) {
	// exercise early return
	b := bits.NewGivens(63, 64)
	if b.AllZeros() {
		t.Fatal("real problem")
	}
}

func TestAnd(t *testing.T) {
	// test allocate z if Num is wrong size
	x := bits.New(1)
	z := bits.Bits{}
	z.And(x, x)
	if z.Num != 1 || len(z.Bits) != 1 {
		t.Fatal("z not allocated to size of args")
	}
	// test different Nums panic
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	z.And(z, bits.Bits{})
}

func TestAndNot(t *testing.T) {
	// test allocate z if Num is wrong size
	x := bits.New(1)
	z := bits.Bits{}
	z.AndNot(x, x)
	if z.Num != 1 || len(z.Bits) != 1 {
		t.Fatal("z not allocated to size of args")
	}
	// test different Nums panic
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	z.AndNot(z, bits.Bits{})
}

func TestBit(t *testing.T) {
	b := bits.New(4)
	// test negative bit number panics
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("panic expected")
			}
		}()
		b.Bit(-1)
	}()
	// test number out of range number panics
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("panic expected")
			}
		}()
		b.Bit(4)
	}()
}

func TestEqual(t *testing.T) {
	// test empty Bits
	var a bits.Bits
	if !a.Equal(a) {
		t.Fatal("empty")
	}
	// test unequal first word
	a = bits.NewGivens(200)
	if a.Equal(bits.NewGivens(0, 200)) {
		t.Fatal("nope 0")
	}
	// test unequal last word
	a = bits.NewGivens(200)
	if a.Equal(bits.NewGivens(199, 200)) {
		t.Fatal("nope 199")
	}
	// test unqual Nums
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	a.Equal(bits.Bits{})
}

func TestIterateOnes(t *testing.T) {
	// test visitor abort
	b := bits.NewGivens(20)
	if b.IterateOnes(func(int) bool { return false }) {
		t.Fatal("but, but")
	}
	// test 1 after Num
	b.Num = 10
	b.IterateOnes(func(int) bool {
		t.Fatal("just no")
		return false
	})
}

func TestIterateZeros(t *testing.T) {
	// test visitor abort
	b := bits.NewGivens(20)
	b.Not(b)
	if b.IterateZeros(func(int) bool { return false }) {
		t.Fatal("but, but")
	}
	// test 0 after Num
	b.Num = 10
	b.IterateZeros(func(int) bool {
		t.Fatal("just no")
		return false
	})
}

func TestNot(t *testing.T) {
	// test allocation
	var z bits.Bits
	z.Not(bits.New(5))
	if z.Num != 5 {
		t.Fatal("z5")
	}
}

func TestOneFrom(t *testing.T) {
	// test 1 bit >= Num
	b := bits.NewGivens(15)
	b.Num = 8
	if b.OneFrom(0) != -1 {
		t.Fatal("iterate past Num")
	}
	// test all words zeros
	b = bits.New(100)
	if b.OneFrom(0) != -1 {
		t.Fatal("no 100")
	}
	// test first word 0
	b.SetBit(99, 1)
	if b.OneFrom(0) != 99 {
		t.Fatal("no 99")
	}
	// test 1 after Num
	b.Num = 90
	if b.OneFrom(0) != -1 {
		t.Fatal("no 90")
	}
}

func TestOr(t *testing.T) {
	// test allocate z if Num is wrong size
	x := bits.New(1)
	z := bits.Bits{}
	z.Or(x, x)
	if z.Num != 1 || len(z.Bits) != 1 {
		t.Fatal("z not allocated to size of args")
	}
	// test different Nums panic
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	z.Or(z, bits.Bits{})
}

func TestSetBit(t *testing.T) {
	// primitive test, independent of other methods
	b := bits.New(80)
	for _, tc := range []struct {
		pos, x int
		bits   []uint64
	}{
		{0, 1, []uint64{1, 0}},
		{2, 1, []uint64{5, 0}},
		{63, 1, []uint64{0x8000000000000005, 0}},
		{65, 1, []uint64{0x8000000000000005, 2}},
		{63, 0, []uint64{5, 2}},
		{0, 0, []uint64{4, 2}},
	} {
		b.SetBit(tc.pos, tc.x)
		if !reflect.DeepEqual(b.Bits, tc.bits) {
			t.Fatal("got ", b.Bits, " want ", tc.bits)
		}
	}
	// test set out of range
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("panic expected")
			}
		}()
		b.SetBit(-1, 1)
	}()
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("panic expected")
			}
		}()
		b.SetBit(80, 1)
	}()
}

func TestSetBits(t *testing.T) {
	b := bits.New(10)
	// test set out of range
	defer func() {
		if recover() == nil {
			t.Fatal("panic expected")
		}
	}()
	b.SetBits(1, 10)
}

func TestSingle(t *testing.T) {
	// test early return
	b := bits.NewGivens(7, 8, 78)
	if b.Single() {
		t.Fatal(78)
	}
}

func TestSlice(t *testing.T) {
	// test zero word
	b := bits.New(5)
	if len(b.Slice()) != 0 {
		panic("of nothing")
	}
}

func TestString(t *testing.T) {
	// test Num == 0
	var b bits.Bits
	if s := b.String(); s > "" {
		t.Fatal(s)
	}
}

func TestXor(t *testing.T) {
	// test allocate z if Num is wrong size
	x := bits.New(1)
	z := bits.Bits{}
	z.Xor(x, x)
	if z.Num != 1 || len(z.Bits) != 1 {
		t.Fatal("z not allocated to size of args")
	}
	// test different Nums panic
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("panic expected")
		}
	}()
	z.Xor(z, bits.Bits{})
}

func TestZeroFrom(t *testing.T) {
	// test 0 bit >= Num
	b := bits.NewGivens(15)
	b.Not(b)
	b.Num = 8
	if b.ZeroFrom(0) != -1 {
		t.Fatal("iterate past Num")
	}
	// test all words ones
	b = bits.New(100)
	b.SetAll()
	if b.ZeroFrom(0) != -1 {
		t.Fatal("no 100")
	}
	// test first word 1s
	b.SetBit(99, 0)
	if b.ZeroFrom(0) != 99 {
		t.Fatal("no 99")
	}
	// test 0 after Num
	b.Num = 90
	if b.ZeroFrom(0) != -1 {
		t.Fatal("no 90")
	}
}
