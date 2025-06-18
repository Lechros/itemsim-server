package search

import (
	"container/heap"
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

type itemStart[T any] struct {
	index int
	start int
}

// invertedIndexSearcher 아이템 목록에서 regexp를 사용해 검색합니다.
type invertedIndexSearcher[T any] struct {
	items   []T
	texts   []string                    // 각 아이템의 text 목록
	iindex  map[rune][]extendedItemInfo // Full 문자 inverted index
	piindex map[rune][]extendedItemInfo // 부분 inverted index, 한글 글자를 순차적으로 분리해서 저장
}

func NewInvertedIndexSearcher[T any](cap int) Searcher[T] {
	return &invertedIndexSearcher[T]{
		items:   make([]T, 0, cap),
		texts:   make([]string, 0, cap),
		iindex:  make(map[rune][]extendedItemInfo),
		piindex: make(map[rune][]extendedItemInfo),
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

func (s *invertedIndexSearcher[T]) Search(query string, size int, cmp ItemCmp[T], filter ItemFilter[T]) []SearchResult[T] {
	set := newItemInfoSet()
	for i, r := range query {
		if set.IsEmpty() {
			return []SearchResult[T]{}
		}
		if i+utf8.RuneLen(r) == len(query) && isHangul(r) { // 마지막 글자가 한글일 경우
			if !hasBatchim(r) { // 받침이 없다면 단순 부분 매칭,
				set.Intersect(s.piindex[r])
			} else { // 받침이 있다면 케이스를 나눠야 한다.
				choseongIndex, jungseongIndex, jongseongIndex := disassembleIntoIndexes(r)
				firstJongseong, secondJongseong := splitConsonant(jongseongs[jongseongIndex])
				if secondJongseong >= 0 { // 이중 받침인 경우 (쌍자음 X)
					// 해당 글자도 매칭하고, 분리한 매칭도 수행해야 함. ex) 앎 -> 앎 / 알+ㅁ / 아+ㄹ+ㅁ
					set1 := set.Intersection(s.iindex[r])
					set2 := set.Intersection(s.iindex[assembleFromIndexes(choseongIndex, jungseongIndex, getJongseongIndex(firstJongseong))])
					set2.Intersect(s.piindex[secondJongseong])
					set3 := set
					set3.Intersect(s.iindex[assembleFromIndexes(choseongIndex, jungseongIndex, 0)])
					set3.Intersect(s.piindex[firstJongseong])
					set3.Intersect(s.piindex[secondJongseong])
					set = Union(set1, set2, set3)
				} else {
					// 해당 글자도 매칭하고, 받침을 초성으로 따로 매칭해야 함. ex) 한 -> 한 / 하+ㄴ
					set1 := set.Intersection(s.iindex[r])
					set2 := set
					set2.Intersect(s.iindex[assembleFromIndexes(choseongIndex, jungseongIndex, 0)])
					set2.Intersect(s.piindex[firstJongseong])
					set = Union(set1, set2)
				}
			}
		} else { // 마지막 글자가 아닌 문자는 초성 매칭이거나, 완전 매칭이어야 함
			if isConsonant(r) {
				firstConsonant, secondConsonant := splitConsonant(r)
				set.Intersect(s.piindex[firstConsonant])
				if secondConsonant >= 0 {
					set.Intersect(s.piindex[secondConsonant])
				}
			} else {
				// iindex에 매칭하여 부분 매칭 방지
				r = unicode.ToLower(r)
				set.Intersect(s.iindex[r])
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

	result := make([]SearchResult[T], size)
	for i, info := range h.items {
		if i == size {
			break
		}
		result[i] = SearchResult[T]{
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
	infos, exists := s.iindex[r]
	if !exists {
		infos = make([]extendedItemInfo, 0)
	}
	if len(infos) == 0 {
		infos = append(infos, extendedItemInfo{index, []int{position}})
	} else if infos[len(infos)-1].index == index {
		infos[len(infos)-1].positions = append(infos[len(infos)-1].positions, position)
	} else {
		infos = append(infos, extendedItemInfo{index, []int{position}})
	}
	s.iindex[r] = infos
}

func (s *invertedIndexSearcher[T]) addPartialRuneToIndex(r rune, index int, position int) {
	infos, exists := s.piindex[r]
	if !exists {
		infos = make([]extendedItemInfo, 0)
	}
	if len(infos) == 0 {
		infos = append(infos, extendedItemInfo{index, []int{position}})
	} else if infos[len(infos)-1].index == index {
		infos[len(infos)-1].positions = append(infos[len(infos)-1].positions, position)
	} else {
		infos = append(infos, extendedItemInfo{index, []int{position}})
	}
	s.piindex[r] = infos
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

func isHangul(r rune) bool {
	return '가' <= r && r <= '힣'
}

func isConsonant(r rune) bool {
	return 'ㄱ' <= r && r < 'ㅎ'
}

func hasBatchim(hangul rune) bool {
	_, _, jongseongIndex := disassembleIntoIndexes(hangul)
	return jongseongIndex > 0
}

func disassembleIntoIndexes(hangul rune) (int, int, int) {
	hangul -= '가'
	return int(hangul / 28 / 21), int(hangul / 28 % 21), int(hangul % 28)
}

func assembleFromIndexes(choseongIndex int, jungseongIndex int, jongseongIndex int) rune {
	return rune('가' + (choseongIndex*21+jungseongIndex)*28 + jongseongIndex)
}

func getFirstVowelPart(jungseong rune) rune {
	switch jungseong {
	case 'ㅘ':
		return 'ㅗ'
	case 'ㅙ':
		return 'ㅗ'
	case 'ㅚ':
		return 'ㅗ'
	case 'ㅝ':
		return 'ㅜ'
	case 'ㅞ':
		return 'ㅜ'
	case 'ㅟ':
		return 'ㅜ'
	case 'ㅢ':
		return 'ㅡ'
	default:
		return -1
	}
}

func getFirstConsonantPart(jongseong rune) rune {
	switch jongseong {
	case 'ㄳ':
		return 'ㄱ'
	case 'ㄵ':
		return 'ㄴ'
	case 'ㄶ':
		return 'ㄴ'
	case 'ㄺ':
		return 'ㄹ'
	case 'ㄻ':
		return 'ㄹ'
	case 'ㄼ':
		return 'ㄹ'
	case 'ㄽ':
		return 'ㄹ'
	case 'ㄾ':
		return 'ㄹ'
	case 'ㄿ':
		return 'ㄹ'
	case 'ㅀ':
		return 'ㄹ'
	case 'ㅄ':
		return 'ㅂ'
	default:
		return -1
	}
}

func splitConsonant(consonant rune) (rune, rune) {
	switch consonant {
	case 'ㄳ':
		return 'ㄱ', 'ㅅ'
	case 'ㄵ':
		return 'ㄴ', 'ㅈ'
	case 'ㄶ':
		return 'ㄴ', 'ㅎ'
	case 'ㄺ':
		return 'ㄹ', 'ㄱ'
	case 'ㄻ':
		return 'ㄹ', 'ㅁ'
	case 'ㄼ':
		return 'ㄹ', 'ㅂ'
	case 'ㄽ':
		return 'ㄹ', 'ㅅ'
	case 'ㄾ':
		return 'ㄹ', 'ㅌ'
	case 'ㄿ':
		return 'ㄹ', 'ㅍ'
	case 'ㅀ':
		return 'ㄹ', 'ㅎ'
	case 'ㅄ':
		return 'ㅂ', 'ㅅ'
	default:
		return consonant, -1
	}
}

func getChoseongIndex(choseong rune) int {
	switch choseong {
	case 'ㄱ':
		return 0
	case 'ㄲ':
		return 1
	case 'ㄴ':
		return 2
	case 'ㄷ':
		return 3
	case 'ㄸ':
		return 4
	case 'ㄹ':
		return 5
	case 'ㅁ':
		return 6
	case 'ㅂ':
		return 7
	case 'ㅃ':
		return 8
	case 'ㅅ':
		return 9
	case 'ㅆ':
		return 10
	case 'ㅇ':
		return 11
	case 'ㅈ':
		return 12
	case 'ㅉ':
		return 13
	case 'ㅊ':
		return 14
	case 'ㅋ':
		return 15
	case 'ㅌ':
		return 16
	case 'ㅍ':
		return 17
	case 'ㅎ':
		return 18
	default:
		return -1
	}
}

func getJungseongIndex(jungseong rune) int {
	switch jungseong {
	case 'ㅏ':
		return 0
	case 'ㅐ':
		return 1
	case 'ㅑ':
		return 2
	case 'ㅒ':
		return 3
	case 'ㅓ':
		return 4
	case 'ㅔ':
		return 5
	case 'ㅕ':
		return 6
	case 'ㅖ':
		return 7
	case 'ㅗ':
		return 8
	case 'ㅘ':
		return 9
	case 'ㅙ':
		return 10
	case 'ㅚ':
		return 11
	case 'ㅛ':
		return 12
	case 'ㅜ':
		return 13
	case 'ㅝ':
		return 14
	case 'ㅞ':
		return 15
	case 'ㅟ':
		return 16
	case 'ㅠ':
		return 17
	case 'ㅡ':
		return 18
	case 'ㅢ':
		return 19
	case 'ㅣ':
		return 20
	default:
		return -1
	}
}

func getJongseongIndex(jongseong rune) int {
	switch jongseong {
	case -1:
		return 0
	case 'ㄱ':
		return 1
	case 'ㄲ':
		return 2
	case 'ㄳ':
		return 3
	case 'ㄴ':
		return 4
	case 'ㄵ':
		return 5
	case 'ㄶ':
		return 6
	case 'ㄷ':
		return 7
	case 'ㄹ':
		return 8
	case 'ㄺ':
		return 9
	case 'ㄻ':
		return 10
	case 'ㄼ':
		return 11
	case 'ㄽ':
		return 12
	case 'ㄾ':
		return 13
	case 'ㄿ':
		return 14
	case 'ㅀ':
		return 15
	case 'ㅁ':
		return 16
	case 'ㅂ':
		return 17
	case 'ㅄ':
		return 18
	case 'ㅅ':
		return 19
	case 'ㅆ':
		return 20
	case 'ㅇ':
		return 21
	case 'ㅈ':
		return 22
	case 'ㅊ':
		return 23
	case 'ㅋ':
		return 24
	case 'ㅌ':
		return 25
	case 'ㅍ':
		return 26
	case 'ㅎ':
		return 27
	default:
		return -1
	}
}

var choseongs = [...]rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
var jungseongs = [...]rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ'}
var jongseongs = [...]rune{0, 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
