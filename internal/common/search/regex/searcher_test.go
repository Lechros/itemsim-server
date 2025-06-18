package regex

import (
	"fmt"
	"itemsim-server/internal/common/search"
	"reflect"
	"strings"
	"testing"
)

func Test_regexSearcher_Search_Count(t *testing.T) {
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
			args{"서명", 5}, 5},
		{"일치 개수가 size보다 작을 경우 해당 개수만큼 반환된다.",
			args{"서명", 100}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_regexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, strings.Compare, nil); !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("len(Search()) = %v, wantCount %v", len(got), tt.wantCount)
			}
		})
	}
}

func Test_regexSearcher_Search_Sort(t *testing.T) {
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
			args{"0", 5}, "양0"},
		{"매치의 첫 글자 위치가 동일할 경우 cmp로 정렬된 상태의 결과가 반환된다.",
			args{"서명", 5}, "서명9"},
	}
	cmp := func(a, b string) int {
		return -strings.Compare(a, b)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_regexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, cmp, nil); !reflect.DeepEqual(got[0].Text, tt.wantFirstText) {
				t.Errorf("Search() = %v, wantFirstText %v", got, tt.wantFirstText)
			}
		})
	}
}

func Test_regexSearcher_Search_Highlight(t *testing.T) {
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
			args{"아방", 5}, "101100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fixture_regexSearcher()
			if got := s.Search(tt.args.query, tt.args.size, strings.Compare, nil); !reflect.DeepEqual(got[0].Highlight, tt.wantFirstHighlight) {
				t.Errorf("Search() = %v, wantFirstHighlight %v", got, tt.wantFirstHighlight)
			}
		})
	}
}

func fixture_regexSearcher() search.Searcher[string] {
	searcher := NewSearcher[string](100)
	for _, prefix := range [...]string{
		"서명", "청소기", "학생", "아르바이트", "손", "삼겹살", "주머니", "파일럿", "양", "김연아",
	} {
		for i := range 10 {
			value := fmt.Sprintf("%s%d", prefix, i)
			searcher.Add(value, value)
		}
	}
	return searcher
}
