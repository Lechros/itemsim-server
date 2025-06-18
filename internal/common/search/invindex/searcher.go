package invindex

import (
	"container/heap"
	"itemsim-server/internal/common/search"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// itemInfo 아이템의 items 내 인덱스와 매칭된 문자의 item text 내 position을 저장합니다.
type itemInfo struct {
	index    int // items 내 인덱스
	position int // 매칭된 문자의 item text 내 position
}

// extendedItemInfo 아이템의 items 내 인덱스와 매칭된 모든 문자의 item text 내 위치들을 저장합니다.
type extendedItemInfo struct {
	index     int   // items 내 인덱스
	positions []int // 매칭된 모든 문자의 item text 내 position 목록
}

// invertedIndexSearcher 아이템 목록에서 regexp를 사용해 검색합니다.
type invertedIndexSearcher[T any] struct {
	items           []T
	texts           []string            // 각 아이템의 text 목록
	invIndex        map[rune][]itemInfo // Full 문자 inverted index
	partialInvIndex map[rune][]itemInfo // 부분 inverted index, 한글 글자를 순차적으로 분리해서 저장
}

func NewSearcher[T any](cap int) search.Searcher[T] {
	return &invertedIndexSearcher[T]{
		items:           make([]T, 0, cap),
		texts:           make([]string, 0, cap),
		invIndex:        make(map[rune][]itemInfo),
		partialInvIndex: make(map[rune][]itemInfo),
	}
}

func (s *invertedIndexSearcher[T]) Add(item T, text string) {
	index := len(s.items)
	s.items = append(s.items, item)
	s.texts = append(s.texts, text)
	for pos, r := range text {
		if isHangul(r) {
			for _, partial := range getHangulRunesForIIndex(r) {
				s.addPartialRuneToIndex(partial, index, pos)
			}
			s.addFullRuneToIndex(r, index, pos)
		} else if isConsonant(r) { // 검색 시 편의성을 위해 초성을 PartialRune 에도 저장
			s.addPartialRuneToIndex(r, index, pos)
			s.addFullRuneToIndex(r, index, pos)
		} else {
			r = unicode.ToLower(r)
			s.addFullRuneToIndex(r, index, pos)
		}
	}
}

func (s *invertedIndexSearcher[T]) Search(query string, size int, cmp search.ItemCmp[T], filter search.ItemFilter[T]) []search.Result[T] {
	set := newItemInfoSet()
	for i, r := range query {
		if set.IsEmpty() {
			return []search.Result[T]{}
		}
		if i+utf8.RuneLen(r) == len(query) && isHangul(r) { // 마지막 글자가 한글일 경우
			if !hasBatchim(r) { // 받침이 없다면 단순 부분 매칭,
				set.Intersect(s.partialInvIndex[r])
			} else { // 받침이 있다면 케이스를 나눠야 한다.
				choseongIndex, jungseongIndex, jongseongIndex := disassembleIntoIndexes(r)
				firstJongseong, secondJongseong := splitConsonant(jongseongs[jongseongIndex])
				if secondJongseong >= 0 { // 이중 받침인 경우 (쌍자음 X)
					// 해당 글자도 매칭하고, 분리한 매칭도 수행해야 함. ex) 앎 -> 앎 / 알+ㅁ / 아+ㄹ+ㅁ
					set1 := set.Intersection(s.invIndex[r])
					set2 := set.Intersection(s.invIndex[assembleFromIndexes(choseongIndex, jungseongIndex, getJongseongIndex(firstJongseong))])
					set2.Intersect(s.partialInvIndex[secondJongseong])
					set3 := set
					set3.Intersect(s.invIndex[assembleFromIndexes(choseongIndex, jungseongIndex, 0)])
					set3.Intersect(s.partialInvIndex[firstJongseong])
					set3.Intersect(s.partialInvIndex[secondJongseong])
					set = Union(set1, set2, set3)
				} else {
					// 해당 글자도 매칭하고, 받침을 초성으로 따로 매칭해야 함. ex) 한 -> 한 / 하+ㄴ
					set1 := set.Intersection(s.invIndex[r])
					set2 := set
					set2.Intersect(s.invIndex[assembleFromIndexes(choseongIndex, jungseongIndex, 0)])
					set2.Intersect(s.partialInvIndex[firstJongseong])
					set = Union(set1, set2)
				}
			}
		} else { // 마지막 글자가 아닌 문자는 초성 매칭이거나, 완전 매칭이어야 함
			if isConsonant(r) {
				firstConsonant, secondConsonant := splitConsonant(r)
				set.Intersect(s.partialInvIndex[firstConsonant])
				if secondConsonant >= 0 {
					set.Intersect(s.partialInvIndex[secondConsonant])
				}
			} else {
				// invIndex 에 매칭하여 부분 매칭 방지
				r = unicode.ToLower(r)
				set.Intersect(s.invIndex[r])
			}
		}
	}

	infos := set.Infos()

	if filter != nil {
		i := 0
		for _, info := range infos {
			if filter(s.items[info.index]) {
				infos[i] = info
				i++
			}
		}
		infos = infos[:i]
	}

	size = min(len(infos), size)

	itemCmp := func(a, b extendedItemInfo) int {
		// 첫 매치 위치가 빠른 결과 우선
		if a.positions[0] != b.positions[0] {
			return a.positions[0] - b.positions[0]
		}
		// 이후 입력 정렬 기준으로 정렬
		return cmp(s.items[a.index], s.items[b.index])
	}

	h := &Heap{items: make([]extendedItemInfo, 0, size), cmp: itemCmp}
	heap.Init(h)
	for _, info := range infos {
		if h.Len() < size {
			heap.Push(h, info)
		} else if itemCmp(info, h.items[0]) < 0 { // 새 원소가 루트(최대값)보다 작으면 교체
			h.items[0] = info
			heap.Fix(h, 0)
		}
	}
	slices.SortFunc(h.items, itemCmp)

	result := make([]search.Result[T], size)
	for i, info := range h.items {
		if i == size {
			break
		}
		result[i] = search.Result[T]{
			Item:      s.items[info.index],
			Text:      s.texts[info.index],
			Highlight: getHighlightByInfo(s.texts[info.index], info.positions),
		}
	}
	return result
}

type Heap struct {
	items []extendedItemInfo
	cmp   func(a, b extendedItemInfo) int
}

func (h *Heap) Len() int           { return len(h.items) }
func (h *Heap) Less(i, j int) bool { return h.cmp(h.items[i], h.items[j]) > 0 }
func (h *Heap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *Heap) Push(x any)         { h.items = append(h.items, x.(extendedItemInfo)) }
func (h *Heap) Pop() any {
	n := len(h.items)
	val := h.items[n-1]
	h.items = h.items[:n-1]
	return val
}

func (s *invertedIndexSearcher[T]) addFullRuneToIndex(r rune, index int, position int) {
	info, exists := s.invIndex[r]
	if !exists {
		info = make([]itemInfo, 0)
	}
	info = append(info, itemInfo{index, position})
	s.invIndex[r] = info
}

func (s *invertedIndexSearcher[T]) addPartialRuneToIndex(r rune, index int, position int) {
	info, exists := s.partialInvIndex[r]
	if !exists {
		info = make([]itemInfo, 0)
	}
	info = append(info, itemInfo{index, position})
	s.partialInvIndex[r] = info
}

func getHangulRunesForIIndex(hangul rune) []rune {
	choseongIndex, jungseongIndex, jongseongIndex := disassembleIntoIndexes(hangul)
	result := make([]rune, 0, 5)
	result = append(result, choseongs[choseongIndex])
	jungseongPart := getFirstVowelPart(jungseongs[jungseongIndex])
	if jungseongPart >= 0 {
		result = append(result, assembleFromIndexes(choseongIndex, getJungseongIndex(jungseongPart), 0))
	}
	result = append(result, assembleFromIndexes(choseongIndex, jungseongIndex, 0))
	if jongseongIndex > 0 {
		jongseongPart := getFirstConsonantPart(jongseongs[jongseongIndex])
		if jongseongPart >= 0 {
			result = append(result, assembleFromIndexes(choseongIndex, jungseongIndex, getJongseongIndex(jongseongPart)))
		}
		result = append(result, assembleFromIndexes(choseongIndex, jungseongIndex, jongseongIndex))
	}
	return result
}

func getHighlightByInfo(text string, positions []int) string {
	builder := strings.Builder{}
	builder.Grow(len(text)) // Generous amount of buffer is faster than utf8.RuneCountInString, or reallocation

	i := 0
	for pos := range text {
		if i == len(positions) || positions[i] > pos {
			builder.WriteRune('0')
		} else if positions[i] == pos {
			builder.WriteRune('1')
			i++
		} else {
			panic("invalid positions")
		}
	}
	return builder.String()
}
