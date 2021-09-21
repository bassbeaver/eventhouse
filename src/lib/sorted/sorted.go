package sorted

import (
	"sort"
)

func Uint64Insert(s []uint64, v uint64) []uint64 {
	p := sort.Search(len(s), func(j int) bool { return s[j] >= v })

	if len(s) == p {
		return append(s, v)
	}

	n := make([]uint64, len(s)+1)
	for i := 0; i < p; i++ {
		n[i] = s[i]
	}
	n[p] = v
	for i := p; i < len(s); i++ {
		n[1+i] = s[i]
	}

	return n
}
