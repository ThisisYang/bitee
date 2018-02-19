package bitee

import (
	"reflect"
	"testing"
)

func TestCalPostition(t *testing.T) {
	cases := []struct {
		testVal  int
		expBlock uint
		expBit   uint
	}{
		{testVal: 1, expBlock: 0, expBit: 0},
		{testVal: 8, expBlock: 0, expBit: 7},
		{testVal: 9, expBlock: 1, expBit: 0},
	}
	for _, c := range cases {
		block, bit := calPosition(c.testVal)
		if block != c.expBlock {
			t.Errorf("block expected: %v, got: %v", block, c.expBlock)
		}
		if bit != c.expBit {
			t.Errorf("bit expected: %v, got: %v", bit, c.expBit)
		}
	}
}

func TestSetBit(t *testing.T) {
	cases := []struct {
		name   string
		setBit int
		bitarr *BitArray
		expErr error
		expArr []byte
	}{
		{name: "case 1", setBit: 0, bitarr: New(16), expErr: errSetZero, expArr: []byte{0x00, 0x00}},
		{name: "case 2", setBit: 17, bitarr: New(16), expErr: errOutOfIndex, expArr: []byte{0x00, 0x00}},
		{name: "case 3", setBit: 16, bitarr: New(16), expErr: nil, expArr: []byte{0x00, 0x01}},
		{name: "case 4", setBit: 2, bitarr: New(16), expErr: nil, expArr: []byte{0x40, 0x00}},
		{name: "case 5", setBit: 8, bitarr: New(16), expErr: nil, expArr: []byte{0x01, 0x00}},
		{name: "case 6", setBit: 9, bitarr: New(17), expErr: nil, expArr: []byte{0x00, 0x80, 0x00}},
		{name: "case 7", setBit: 3, bitarr: &BitArray{length: 16, arr: []byte{0x20, 0x00}}, expErr: errSet, expArr: []byte{0x20, 0x00}},
	}

	for _, c := range cases {
		err := c.bitarr.SetBit(c.setBit)
		if err != c.expErr {
			t.Errorf("error case: %v\nexp err: %v\ngot err: %v", c.name, err, c.expErr)
		}
		if sliceEqual(c.expArr, c.bitarr.arr) == false {
			t.Errorf("error case: %v\nexp arr: %v\ngot arr: %v", c.name, c.expArr, c.bitarr.arr)
		}

	}
}

func TestUnSetBit(t *testing.T) {
	cases := []struct {
		name     string
		unSetBit int
		bitarr   *BitArray
		expErr   error
		expArr   []byte
	}{
		{name: "case 1", unSetBit: 0, bitarr: New(16), expErr: errSetZero, expArr: []byte{0x00, 0x00}},
		{name: "case 2", unSetBit: 17, bitarr: New(16), expErr: errOutOfIndex, expArr: []byte{0x00, 0x00}},
		{name: "case 3", unSetBit: 2, bitarr: New(16), expErr: errUnset, expArr: []byte{0x00, 0x00}},
		{name: "case 4", unSetBit: 3, bitarr: &BitArray{length: 17, arr: []byte{0x20, 0x00, 0x00}}, expErr: nil, expArr: []byte{0x00, 0x00, 0x00}},
		{name: "case 5", unSetBit: 8, bitarr: &BitArray{length: 16, arr: []byte{0x21, 0xFF}}, expErr: nil, expArr: []byte{0x20, 0xFF}},
	}

	for _, c := range cases {
		err := c.bitarr.UnSetBit(c.unSetBit)
		if err != c.expErr {
			t.Errorf("error case: %v\nexp err: %v\ngot err: %v", c.name, err, c.expErr)
		}
		if sliceEqual(c.expArr, c.bitarr.arr) == false {
			t.Errorf("error case: %v\nexp arr: %v\ngot arr: %v", c.name, c.expArr, c.bitarr.arr)
		}

	}
}

func TestToValue(t *testing.T) {
	cases := []struct {
		name   string
		bitarr *BitArray
		expArr []int
		expErr error
	}{
		{
			name: "case 1", bitarr: New(16),
			expArr: []int{}, expErr: nil,
		},
		{
			name: "case 2", bitarr: &BitArray{length: 24, arr: []byte{0x80, 0x30, 0x81}},
			expArr: []int{1, 11, 12, 17, 24}, expErr: nil,
		},
	}

	for _, c := range cases {
		gotarr, err := c.bitarr.ToValue()
		if err != c.expErr {
			t.Errorf("error case: %v\nexp: %v\ngot: %v", c.name, err, c.expErr)
		}
		if reflect.DeepEqual(c.expArr, gotarr) != true {
			t.Errorf("error case: %v\nexp: %v\ngot: %v", c.name, c.expArr, gotarr)
		}
	}
}

func TestToStringt(t *testing.T) {
	cases := []struct {
		name      string
		bitarr    *BitArray
		expString string
	}{
		{name: "case 1", bitarr: &BitArray{length: 16, arr: []byte{0x00, 0x00}}, expString: "0000000000000000"},
		{name: "case 2", bitarr: &BitArray{length: 16, arr: []byte{0xFF, 0xFF}}, expString: "1111111111111111"},
		{name: "case 3", bitarr: &BitArray{length: 16, arr: []byte{0xFF, 0x0F}}, expString: "1111111100001111"},
	}
	for _, c := range cases {
		gotString := c.bitarr.ToString()
		if gotString != c.expString {
			t.Errorf("error case: %v\nexp string: %v\ngot string: %v", c.name, c.expString, gotString)
		}
	}
}

func TestIsSet(t *testing.T) {
	cases := []struct {
		name   string
		b      byte
		p      uint
		result bool
	}{
		{name: "case 1", b: 0x00, p: 1, result: false},
		{name: "case 2", b: 0x04, p: 2, result: false},
		{name: "case 3", b: 0x01, p: 7, result: true},
		{name: "case 4", b: 0x01, p: 1, result: false},
		{name: "case 5", b: 0x00, p: 0, result: false},
	}
	for _, c := range cases {
		r := isSet(c.b, c.p)
		if r != c.result {
			t.Errorf("error case: %v\nexp: %v\ngot: %v", c.name, c.result, r)
		}
	}
}
