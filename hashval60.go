package key

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// HashVal60 stores 60 bits of a hash value.
type HashVal60 uint64

// BitsPerLevel60 is the number of bits per depth level of the HashVal60.
const BitsPerLevel60 uint = 6

// MaxDepth60 represents the maximum depth of the HashVal60.
const MaxDepth60 uint = 9

func indexMask60(depth uint) HashVal60 {
	return HashVal60((1<<BitsPerLevel60)-1) << (depth * BitsPerLevel60)
}

// Index() will return a 6 bit (aka BitsPerLevel60) value 'depth' number
// of 6 bits from the beginning of the HashVal60 (aka uint32) h60 value.
func (h60 HashVal60) Index(depth uint) uint {
	var idxMask = indexMask60(depth)
	var idx = uint((h60 & idxMask) >> (depth * BitsPerLevel60))
	return idx
}

// HashPathMask60() returns the mask for a 60 bit HashPath value.
func HashPathMask60(depth uint) HashVal60 {
	//return HashVal60(1<<(depth*BitsPerLevel60)) - 1
	return HashVal60(1<<((depth+1)*BitsPerLevel60)) - 1
}

// BuildHashPath() method adds a idx at depth level of the hashPath.
// Given a hashPath = "/11/22/33" and you call hashPath.BuildHashPath(44, 3)
// the method will return hashPath "/11/22/33/44". hashPath is shown here
// in the string representation, but the real value is HashVal60 (aka uint32).
func (hashPath HashVal60) BuildHashPath(idx, depth uint) HashVal60 {
	//var mask = HashPathMask60(depth-1)
	var mask HashVal60 = (1 << (depth * BitsPerLevel60)) - 1
	var hp = hashPath & mask

	return hp | HashVal60(idx<<(depth*BitsPerLevel60))
}

// HashPathString() returns a string representation of the index path of
// a HashVal60 60 bit value; that is depth number of zero padded numbers between
// "00" and "63" separated by "/" characters and a leading '/'. If the depth
// parameter is 0 then the method will simply return a solitary "/".
// Warning: It will panic() if depth > MaxDepth60.
// Example: "/00/24/46/17/34/08/54" for depth=7 of a hash60 value represented
//       by "/00/24/46/17/34/08/54/28/59/51".
func (h60 HashVal60) HashPathString(depth uint) string {
	if depth > MaxDepth60 {
		panic(fmt.Sprintf("HashPathString: depth,%d > MaxDepth60,%d\n", depth, MaxDepth60))
	}

	if depth == 0 {
		return "/"
	}

	// Remember we want to include the indexes from [0, depth] (hence including depth)
	// So strs has to be depth+1 in size, and the for loop has to include i=depth.

	var strs = make([]string, depth)

	for d := uint(0); d < depth; d++ {
		var idx = h60.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// Return the HashVal60 as a 60 bit bit string separated into groups of 6 bits
// (aka BitsPerLevel60).
func (h60 HashVal60) BitString() string {
	var s = make([]string, MaxDepth60+1)
	for i := uint(0); i <= MaxDepth60; i++ {
		s[MaxDepth60-i] += fmt.Sprintf("%06b", h60.Index(i))
	}
	return "00 " + strings.Join(s, " ")
}

// String() returns a string representation of the h60 HashVal60 value. This
// is MaxDepth60+1(10) two digit numbers (zero padded) between "00" and "63"
// seperated by '/' characters and given a leading '/'.
// Example: "/08/14/28/20/00/31/56/01/24/63"
func (h60 HashVal60) String() string {
	return h60.HashPathString(MaxDepth60)
}

// ParseHashVal60() parses a string with a leading '/' and MaxDepth60+1 number
// of two digit numbers zero padded between "00" and "63" joined by '/' characters.
// Example: var h60 key.HashVal60 = key.ParseHashVal60("/00/01/02/0\/03/04/05/06/07/08/09")
func ParseHashVal60(s string) HashVal60 {
	if !strings.HasPrefix(s, "/") {
		panic(errors.New("does not start with '/'"))
	}
	var s0 = s[1:]
	var as = strings.Split(s0, "/")

	var h60 HashVal60 = 0
	for i, s1 := range as {
		var ui, err = strconv.ParseUint(s1, 10, int(BitsPerLevel60))
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("strconv.ParseUint(%q, %d, %d) failed", s1, 10, BitsPerLevel60)))
		}
		h60 |= HashVal60(ui << (uint(i) * BitsPerLevel60))
		//fmt.Printf("%d: h60 = %q %2d %#02x %05b\n", i, s1, ui, ui, ui)
	}

	return h60
}
