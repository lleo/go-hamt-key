package key

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type HashVal30 uint32

const BitsPerLevel30 uint = 5

const MaxDepth30 uint = 5

func indexMask(depth uint) HashVal30 {
	return HashVal30((1<<BitsPerLevel30)-1) << (depth * BitsPerLevel30)
}

// Index() will return a 5bit (aka BitsPerLevel30) value 'depth' number
// of 5bits from the beginning of the HashVal30 (aka uint32) h30 value.
func (h30 HashVal30) Index(depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint((h30 & idxMask) >> (depth * BitsPerLevel30))
	return idx
}

func hashPathMask(depth uint) HashVal30 {
	return HashVal30(1<<(depth*BitsPerLevel30)) - 1
}

// HashPathString() returns a string representation of the index path of
// a HashVal30 30 bit value; that is depth number of zero padded numbers between
// "00" and "63" separated by '/' characters and a leading '/'. If the depth
// parameter is 0 then the method will simply return a solitary "/".
// Warning: It will panic() if depth > MaxDepth30.
// Example: "/00/24/46/17" for depth=4 of a hash30 value represented
//       by "/00/24/46/17/34/08".
func (h30 HashVal30) HashPathString(depth uint) string {
	if depth > MaxDepth30 {
		panic(fmt.Sprintf("HashPathString: depth,%d > MaxDepth30,%d\n", depth, MaxDepth30))
	}

	if depth == 0 {
		return "/"
	}

	// Remember we want to include the indexes from [0, depth] (hence including depth)
	// So strs has to be depth+1 in size, and the for loop has to include i=depth.

	var strs = make([]string, depth+1)

	for d := uint(0); d <= depth; d++ {
		var idx = h30.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// String() returns a string representation of the h30 HashVal30 value. This
// is MaxDepth30+1(6) two digit numbers (zero padded) between "00" and "31"
// seperated by '/' characters and given a leading '/'.
// Example: "/08/14/28/20/00/31"
func (h30 HashVal30) String() string {
	return h30.HashPathString(MaxDepth30)
}

func StringToHashVal30(string) HashVal30 {
	if !strings.HasPrefix(s, "/") {
		panic(errors.New("does not start with '/'"))
	}
	var s0 = s[1:]
	var as = strings.Split(s0, "/")

	var h30 HashVal30 = 0
	for i, s1 := range as {
		var ui, err = strconv.ParseUint(s1, 10, int(Nbits))
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("strconv.ParseUint(%q, %d, %d) failed", s1, 10, Nbits)))
		}
		h30 |= HashVal30(ui << (uint(i) * BitsPerLevel30))
		//fmt.Printf("%d: h30 = %q %2d %#02x %05b\n", i, s1, ui, ui, ui)
	}

	return h30
}
