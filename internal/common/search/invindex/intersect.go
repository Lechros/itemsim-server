package invindex

import (
	"github.com/achille-roussel/kway-go"
	"iter"
)

type indexPositionSet struct {
	data []indexPositions // 실제로 Intersect 가 수행된 이후에는 data 에 저장된다.
	temp []indexPosition  // 첫 Intersect 인 경우 other 를 그대로 가져온다.
}

func newIndexPositionSet() *indexPositionSet {
	return &indexPositionSet{nil, nil}
}

func (s *indexPositionSet) IsEmpty() bool {
	if s.data != nil {
		return len(s.data) == 0
	}
	if s.temp != nil {
		return len(s.temp) == 0
	}
	return false
}

func (s *indexPositionSet) Infos() []indexPositions {
	if s.data != nil {
		return s.data
	} else if s.temp != nil {
		result := make([]indexPositions, 0, len(s.temp))
		lastIndex := -1
		for _, item := range s.temp {
			if item.index > lastIndex {
				result = append(result, indexPositions{
					index:     item.index,
					positions: []int{item.position},
				})
				lastIndex = item.index
			}
		}
		return result
	} else {
		return []indexPositions{}
	}
}

func (s *indexPositionSet) Intersect(other []indexPosition) {
	if other == nil {
		s.data = []indexPositions{}
	} else if s.data != nil {
		s.data = intersectExisting(s.data, s.data, other)
	} else if s.temp != nil {
		s.data = intersectNew(s.temp, other)
		s.temp = nil
	} else {
		s.temp = other
	}
}

func (s *indexPositionSet) Intersection(other []indexPosition) *indexPositionSet {
	result := newIndexPositionSet()
	if other == nil {
		result.data = []indexPositions{}
	} else if s.data != nil {
		result.data = intersectExisting(make([]indexPositions, min(len(s.data), len(other))), s.data, other)
	} else if s.temp != nil {
		result.data = intersectNew(s.temp, other)
	} else {
		result.temp = other
	}
	return result
}

func Union(sets ...*indexPositionSet) *indexPositionSet {
	sequence := func(sets *indexPositionSet) iter.Seq2[indexPositions, error] {
		return func(yield func(indexPositions, error) bool) {
			for _, item := range sets.Infos() {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
	cmp := func(a, b indexPositions) int {
		if a.index != b.index {
			return a.index - b.index
		}
		return a.positions[len(a.positions)-1] - b.positions[len(b.positions)-1]
	}
	seqs := make([]iter.Seq2[indexPositions, error], len(sets))
	for i, item := range sets {
		seqs[i] = sequence(item)
	}
	result := make([]indexPositions, 0)
	for e := range kway.MergeFunc(cmp, seqs...) {
		if len(result) == 0 || e.index != result[len(result)-1].index {
			result = append(result, e)
		}
	}
	return &indexPositionSet{result, nil}
}

func intersectNew(cur, next []indexPosition) []indexPositions {
	if len(cur) == 0 {
		return []indexPositions{}
	}
	if len(next) == 0 {
		return []indexPositions{}
	}
	buf := make([]indexPositions, min(len(cur), len(next)))
	bi := 0
	ci := 0
	ni := 0
	for ci < len(cur) && ni < len(next) {
		if cur[ci].index == next[ni].index {
			index := cur[ci].index
			position := cur[ci].position
			for ni < len(next) && next[ni].index == index && position >= next[ni].position {
				ni++
			}
			if ni < len(next) && next[ni].index == index && position < next[ni].position {
				buf[bi] = indexPositions{
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

func intersectExisting(buf []indexPositions, cur []indexPositions, next []indexPosition) []indexPositions {
	if len(cur) == 0 {
		return []indexPositions{}
	}
	if len(next) == 0 {
		return []indexPositions{}
	}
	bi := 0
	ci := 0
	ni := 0
	for ci < len(cur) && ni < len(next) {
		if cur[ci].index == next[ni].index {
			index := cur[ci].index
			position := cur[ci].positions[len(cur[ci].positions)-1]
			for ni < len(next) && next[ni].index == index && position >= next[ni].position {
				ni++
			}
			if ni < len(next) && next[ni].index == index && position < next[ni].position {
				buf[bi] = indexPositions{
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
