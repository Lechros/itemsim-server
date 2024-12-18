package util

import (
	"github.com/Lechros/hangul_regexp"
	"strconv"
	"unicode/utf8"
)

type HangulInvertedIndex[T any] struct {
	items []T
	// key: rune, value: slice of idx of data
	index map[rune][]int
}

func (ii *HangulInvertedIndex[T]) Init() {
	ii.items = make([]T, 0)
	ii.index = make(map[rune][]int)
}

func (ii *HangulInvertedIndex[T]) Add(item T, key string) {
	itemIdx := len(ii.items)
	ii.items = append(ii.items, item)
	visited := make(map[rune]struct{}, len(key))
	for _, ch := range key {
		if _, ok := visited[ch]; !ok {
			ii.index[ch] = append(ii.index[ch], itemIdx)
			visited[ch] = struct{}{}
		}
	}
}

func (ii *HangulInvertedIndex[T]) FindAll(key string) []T {
	if len(key) <= 1 {
		panic("key should have length > 1 but was " + strconv.Itoa(len(key)))
	}

	var idxList []int = nil
	for i, ch := range key {
		if i+utf8.RuneLen(ch) == len(key) {
			break
		}
		if hangul_regexp.CanBeChoseongOrJongseong(ch) {
			continue
		}
		if idxList == nil {
			idxList = make([]int, len(ii.index[ch]))
			copy(idxList, ii.index[ch])
		} else {
			idxList = intersect(idxList, ii.index[ch])
		}
	}
	if idxList == nil {
		return nil
	}
	result := make([]T, len(idxList))
	for i, idx := range idxList {
		result[i] = ii.items[idx]
	}
	return result
}

func intersect(dst, other []int) []int {
	count := 0
	i, j := 0, 0
	for i < len(dst) && j < len(other) {
		if dst[i] < other[j] {
			i++
		} else if dst[i] > other[j] {
			j++
		} else {
			dst[count] = dst[i]
			count++
			i++
			j++
		}
	}
	return dst[:count]
}
