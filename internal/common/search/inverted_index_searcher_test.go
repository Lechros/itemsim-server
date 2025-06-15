package search

import (
	"reflect"
	"strings"
	"testing"
)

func Test_invertedIndexSearcher_Search_Count(t *testing.T) {
	type args struct {
		query string
		size  int
	}
	type testCase[T any] struct {
		name      string
		args      args
		wantCount int
	}
	tests := []testCase[string]{
		{"반환 결과는 size를 넘지 않는다.",
			args{"ㅅ", 2}, 2},
		{"일치 개수가 size보다 작을 경우 해당 개수만큼 반환된다.",
			args{"ㅅ", 100}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_invertedIndexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, strings.Compare, nil); !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("len(Search()) = %v, wantCount %v", len(got), tt.wantCount)
			}
		})
	}
}

func Test_invertedIndexSearcher_Search_Sort(t *testing.T) {
	type args struct {
		query string
		size  int
	}
	type testCase[T any] struct {
		name          string
		args          args
		wantFirstText string
	}
	tests := []testCase[string]{
		{"매치의 첫 글자 위치가 빠른 순서대로 결과가 정렬되어 있다.",
			args{"ㅇ", 5}, "양"},
		{"매치의 첫 글자 위치가 동일할 경우 cmp로 정렬된 상태의 결과가 반환된다.",
			args{"서명", 5}, "서명"},
	}
	cmp := func(a, b string) int {
		return -strings.Compare(a, b)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_invertedIndexSearcher()
			got := s.Search(tt.args.query, tt.args.size, cmp, nil)
			text := ""
			if len(got) > 0 {
				text = got[0].Text
			}
			if !reflect.DeepEqual(text, tt.wantFirstText) {
				t.Errorf("Search() = %v, wantFirstText %v", got, tt.wantFirstText)
			}
		})
	}
}

func Test_invertedIndexSearcher_Search_Highlight(t *testing.T) {
	type args struct {
		query string
		size  int
	}
	type testCase[T any] struct {
		name               string
		args               args
		wantFirstHighlight string
	}
	tests := []testCase[string]{
		{"하이라이트에서 부분 매칭된 문자도 1로 표시된다.",
			args{"아방", 5}, "10110"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_invertedIndexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, strings.Compare, nil); !reflect.DeepEqual(got[0].Highlight, tt.wantFirstHighlight) {
				t.Errorf("Search() = %v, wantFirstHighlight %v", got, tt.wantFirstHighlight)
			}
		})
	}
}

func Test_invertedIndexSearcher_Search_Hangul(t *testing.T) {
	type args struct {
		query string
		size  int
	}
	type testCase[T any] struct {
		name          string
		args          args
		wantFirstText string
	}
	tests := []testCase[string]{
		{"복합 모음 결과를 포함한다.",
			args{"크리우", 5}, "크리스탈 웬투스"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_invertedIndexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, strings.Compare, nil); !reflect.DeepEqual(got[0].Text, tt.wantFirstText) {
				t.Errorf("Search()[0] = %v, wantFirstText %v", got[0], tt.wantFirstText)
			}
		})
	}
}

func fixture_invertedIndexSearcher() Searcher[string] {
	searcher := NewInvertedIndexSearcher[string](100)
	for _, value := range [...]string{
		"서명", "청소기", "학생", "아르바이트", "손", "삼겹살", "주머니", "파일럿", "양", "김연아", "크리스탈 웬투스", "다크 티베리안",
	} {
		searcher.Add(value, value)
	}
	return searcher
}
