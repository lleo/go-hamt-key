package key

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type HashVal60 uint64

const BitsPerLevel60 uint = 6

const MaxDepth60 uint = 9

func indexMask(depth uint) HashVal60 {
	return HashVal60((1<<BitsPerLevel60)-1) << (depth * BitsPerLevel60)
}

func (h60 HashVal60) Index(depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint((h60 & idxMask) >> (depth * BitsPerLevel60))
	return idx
}

func hashPathMask(depth uint) HashVal60 {
	return HashVal60(1<<(depth*BitsPerLevel60)) - 1
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

	var strs = make([]string, depth+1)

	for d := uint(0); d <= depth; d++ {
		var idx = h60.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// String() returns a string representation of the h60 HashVal60 value. This
// is MaxDepth60+1(10) two digit numbers (zero padded) between "00" and "63"
// seperated by '/' characters and given a leading '/'.
// Example: "/08/14/28/20/00/31/56/01/24/63"
func (h60 HashVal60) String() string {
	return h60.HashPathString(MaxDepth60)
}

func StringToHashVal60(string) HashVal60 {
	if !strings.HasPrefix(s, "/") {
		panic(errors.New("does not start with '/'"))
	}
	var s0 = s[1:]
	var as = strings.Split(s0, "/")

	var h60 HashVal60 = 0
	for i, s1 := range as {
		var ui, err = strconv.ParseUint(s1, 10, int(Nbits))
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("strconv.ParseUint(%q, %d, %d) failed", s1, 10, Nbits)))
		}
		h60 |= HashVal60(ui << (uint(i) * BitsPerLevel60))
		//fmt.Printf("%d: h60 = %q %2d %#02x %05b\n", i, s1, ui, ui, ui)
	}

	return h60
}
