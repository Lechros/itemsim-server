package repository

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSearchGearByName(t *testing.T) {
	type args struct {
		search string
	}
	tests := []struct {
		args    args
		wantLen int
	}{
		{args{"1"}, 214},
		{args{"a"}, 18},
		{args{"ㄱ"}, 2919},
		{args{"ㄻ"}, 1065},
		{args{"ㅇ"}, 6995},
		{args{"ㅇ ㅅ ㅇ"}, 170},
		{args{"이"}, 3207},
		{args{"앱솔"}, 107},
		{args{"에텔"}, 36},
		{args{"이 이"}, 375},
	}
	name := "len(SearchGearByName(%s))==%d"
	sizeArg := 9999
	for _, tt := range tests {
		t.Run(fmt.Sprintf(name, tt.args.search, tt.wantLen), func(t *testing.T) {
			if got, _ := SearchGearByName(tt.args.search, sizeArg); !reflect.DeepEqual(len(got), tt.wantLen) {
				t.Errorf("SearchGearByName() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func BenchmarkSearchGearByName(b *testing.B) {
	searches := []string{
		"a",
		"ㅇ",
		"ㅇ ㅅ ㅇ",
		"이",
		"마깃안",
		"앱솔",
		"젠",
		"제네",
		"에광",
		"ㅇㅋㅇㅅㅇㄷ ㅇ",
		"아케인셰이드 아처",
	}
	for _, search := range searches {
		b.Run(search, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				SearchGearByName("ㅇ", 9999)
			}
		})
	}
}
