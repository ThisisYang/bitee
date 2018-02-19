package bitee

import (
	"fmt"
	"strings"
)

/*
BitArray implementation

Practice Golang and implement bit array.

Bit array is a sequence of bit.
It implemented methods including Size, SetBit, UnSetBit, ToValues and ToString.
This implementation use byte slice underneath.
Each byte is 8 bits.
The highest bit is stored on right most block of byte.
Within the byte, the highest bit is stored on right most bit as well
Postition of the bit starts from 1

example: when set bit 19
blocks = 3, set position 2 ( 2*8 + 2+1),
result: 00000000 00000000 00200000
*/

var (
	bitsPerUnit   = 8
	errSetZero    = fmt.Errorf("has to be larger then 0")
	errOutOfIndex = fmt.Errorf("out of index")
	errUnset      = fmt.Errorf("bit was not set")
	errSet        = fmt.Errorf("bit was set")
)

// BitArray is unsigned 64 bit int, up to 18446744073709551615
type BitArray struct {
	length int
	arr    []byte
}

// Size return the size of the bit array
func (b *BitArray) Size() int {
	return b.length
}

// calPosition c
// highest bit stored on right most bit
// if blocks = 3, set position 2 which means set 19 (2*8 + 2 + 1)
// result: 00000000 00000000 00200000
func calPosition(p int) (uint, uint) {
	// p = p-1 first since index start from 0, instead of 1
	p--
	return uint(p / bitsPerUnit), uint(p % bitsPerUnit)
}

// New create a new bitarray which can have 256 bits
func New(length int) *BitArray {

	blocks := length / bitsPerUnit
	remainer := length % bitsPerUnit

	if remainer > 0 {
		blocks++
	}

	return &BitArray{
		length: length,
		arr:    make([]byte, blocks, blocks),
	}

}

// SetBit will set the bit
func (b *BitArray) SetBit(n int) error {

	if err := b.validPosition(n); err != nil {
		return err
	}
	block, position := calPosition(n)

	if isSet(b.arr[block], position) == true {
		return errSet
	}
	b.arr[block] = b.arr[block] | (0x80 >> position)

	// or use xor logic
	// b.arr[block] = b.arr[block] ^ (1 << position)
	return nil
}

// UnSetBit unset a position. Return error if was not set or out of range
func (b *BitArray) UnSetBit(n int) error {
	if err := b.validPosition(n); err != nil {
		return err
	}
	block, position := calPosition(n)

	if isSet(b.arr[block], position) == false {
		return errUnset
	}

	// convert the byte to int, then minus a int that bitwise moved n poistion
	b.arr[block] = byte(int(b.arr[block]) - int(0x80>>position))

	/* or use xor logic
		   xor table:
	        a   b  a^b
	        0   0   0
	        1   0   1
	        0   1   1
	        1   1   0
	*/
	// b.arr[block] = b.arr[block] ^ (0x80 >> position)

	return nil
}

// ToValue returns an array of int represent the bitArray
// Most left bit is highest and right most is lowest bit
func (b *BitArray) ToValue() ([]int, error) {
	// set length = 0, so initially contains 0 elements.
	// otherwise, it will filled 0s then append after 0s
	vals := make([]int, 0, b.length)
	/*
		c := 0
		for _, bytes := range b.arr {
			for j := 0; j < 8; j++ {
				c++
				if isSet(bytes, uint(j)) {
					vals = append(vals, c)
				}
			}
		}
	*/
	for i := 1; i <= b.length; i++ {
		block, p := calPosition(i)
		if isSet(b.arr[block], p) {
			vals = append(vals, i)
		}
	}

	return vals, nil
}

func (b *BitArray) validPosition(p int) error {
	if p < 1 {
		return errSetZero
	}
	if p > b.length {
		return errOutOfIndex
	}
	return nil
}

// ToString convert the bitArray to string, for example: "0100010"
func (b *BitArray) ToString() string {
	vals := make([]string, 0, b.length*2)
	for _, bt := range b.arr {
		// %b will only return base 2 value.
		// 0x00 will return 0 instead of 00000000. need to prepend 0s
		toBinary := fmt.Sprintf("%b", bt)
		if len(toBinary) < 8 {
			vals = append(vals, strings.Repeat("0", 8-len(toBinary)))
		}
		vals = append(vals, toBinary)
	}
	s := strings.Join(vals, "")
	return s
}

// IsSet is exposble method, check if the Nth bit is set.
// Return bool, and error in case OutOfRange or negative position
func (b *BitArray) IsSet(p int) (bool, error) {
	if err := b.validPosition(p); err != nil {
		return false, err
	}
	block, position := calPosition(p)
	return isSet(b.arr[block], position), nil
}

// isSet return True if bit is set in the byte
// p has range 0 - 7
func isSet(b byte, p uint) bool {
	// use $ logic
	val := b & (0x80 >> p)
	if val > 0 {
		return true
	}
	return false
}

func byteEqual(b1, b2 byte) bool {
	return b1^b2 == 0
}

func sliceEqual(s1, s2 []byte) bool {
	// they should have same length
	// XOR operation
	if len(s1) != len(s2) {
		return false
	}
	for n, s1Byte := range s1 {
		s2Byte := s2[n]
		if byteEqual(s1Byte, s2Byte) == false {
			return false
		}
	}
	return true
}
