package search

import (
	"github.com/achille-roussel/kway-go"
	"iter"
)

type extendedItemInfo struct {
	index     int
	positions []int
}

type ItemInfoSet struct {
	data []extendedItemInfo // 실제로 Intersect 가 수행된 이후에는 data 에 저장된다.
	temp []itemInfo         // 첫 Intersect 인 경우 other 를 그대로 가져온다.
}

func newItemInfoSet() *ItemInfoSet {
	return &ItemInfoSet{nil, nil}
}

func (s *ItemInfoSet) IsEmpty() bool {
	if s.data != nil {
		return len(s.data) == 0
	}
	if s.temp != nil {
		return len(s.temp) == 0
	}
	return false
}

func (s *ItemInfoSet) Infos() []extendedItemInfo {
	if s.data != nil {
		return s.data
	} else if s.temp != nil {
		result := make([]extendedItemInfo, 0, len(s.temp))
		lastIndex := -1
		for _, item := range s.temp {
			if item.index > lastIndex {
				result = append(result, extendedItemInfo{
					index:     item.index,
					positions: []int{item.position},
				})
				lastIndex = item.index
			}
		}
		return result
	} else {
		return []extendedItemInfo{}
	}
}

func (s *ItemInfoSet) Intersect(other []itemInfo) {
	if other == nil {
		s.data = []extendedItemInfo{}
	} else if s.data != nil {
		s.data = intersectExisting(s.data, s.data, other)
	} else if s.temp != nil {
		s.data = intersectNew(s.temp, other)
		s.temp = nil
	} else {
		s.temp = other
	}
}

func (s *ItemInfoSet) Intersection(other []itemInfo) *ItemInfoSet {
	result := newItemInfoSet()
	if other == nil {
		result.data = []extendedItemInfo{}
	} else if s.data != nil {
		result.data = intersectExisting(make([]extendedItemInfo, min(len(s.data), len(other))), s.data, other)
	} else if s.temp != nil {
		result.data = intersectNew(s.temp, other)
	} else {
		result.temp = other
	}
	return result
}

func Union(sets ...*ItemInfoSet) *ItemInfoSet {
	sequence := func(sets *ItemInfoSet) iter.Seq2[extendedItemInfo, error] {
		return func(yield func(extendedItemInfo, error) bool) {
			for _, item := range sets.Infos() {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
	cmp := func(a, b extendedItemInfo) int {
		if a.index != b.index {
			return a.index - b.index
		}
		return a.positions[len(a.positions)-1] - b.positions[len(b.positions)-1]
	}
	seqs := make([]iter.Seq2[extendedItemInfo, error], len(sets))
	for i, item := range sets {
		seqs[i] = sequence(item)
	}
	result := make([]extendedItemInfo, 0)
	for e := range kway.MergeFunc(cmp, seqs...) {
		if len(result) == 0 || e.index != result[len(result)-1].index {
			result = append(result, e)
		}
	}
	return &ItemInfoSet{result, nil}
}

func intersectNew(cur, next []itemInfo) []extendedItemInfo {
	if len(cur) == 0 {
		return []extendedItemInfo{}
	}
	if len(next) == 0 {
		return []extendedItemInfo{}
	}
	buf := make([]extendedItemInfo, min(len(cur), len(next)))
	bi := 0
	ci := 0
	ni := 0
	for ci < len(cur) && ni < len(next) {
		if cur[ci].index == next[ni].index {
			index := cur[ci].index
			for ni < len(next) && next[ni].index == index && cur[ci].position >= next[ni].position {
				ni++
			}
			if ni < len(next) && next[ni].index == index && cur[ci].position < next[ni].position {
				buf[bi] = extendedItemInfo{
					index:     index,
					positions: []int{cur[ci].position, next[ni].position},
				}
				bi++
			}
			for ci < len(cur) && cur[ci].index == index {
				ci++
			}
			for ni < len(next) && next[ni].index == index {
				ni++
			}
		} else if cur[ci].index < next[ni].index {
			ci++
		} else {
			ni++
		}
	}
	return buf[:bi]
}

func intersectExisting(buf []extendedItemInfo, cur []extendedItemInfo, next []itemInfo) []extendedItemInfo {
	if len(cur) == 0 {
		return []extendedItemInfo{}
	}
	if len(next) == 0 {
		return []extendedItemInfo{}
	}
	bi := 0
	ci := 0
	ni := 0
	for ci < len(cur) && ni < len(next) {
		if cur[ci].index == next[ni].index {
			index := cur[ci].index
			for ni < len(next) && next[ni].index == index && cur[ci].positions[len(cur[ci].positions)-1] >= next[ni].position {
				ni++
			}
			if ni < len(next) && next[ni].index == index && cur[ci].positions[len(cur[ci].positions)-1] < next[ni].position {
				buf[bi] = extendedItemInfo{
					index:     index,
					positions: append(cur[ci].positions, next[ni].position),
				}
				bi++
			}
			for ci < len(cur) && cur[ci].index == index {
				ci++
			}
			for ni < len(next) && next[ni].index == index {
				ni++
			}
		} else if cur[ci].index < next[ni].index {
			ci++
		} else {
			ni++
		}
	}
	return buf[:bi]
}
