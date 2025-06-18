package invindex

import (
	"reflect"
	"testing"
)

func Test_ItemInfoSet_Intersect(t *testing.T) {
	type fields struct {
		data []extendedItemInfo
		temp []itemInfo
	}
	type args struct {
		other []itemInfo
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected []extendedItemInfo
	}{
		{
			name: "empty set should remain empty",
			fields: fields{
				data: []extendedItemInfo{},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 2}},
			},
			expected: []extendedItemInfo{},
		},
		{
			name: "nil other should result in empty set",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: nil,
			},
			expected: []extendedItemInfo{},
		},
		{
			name: "nil data with nil temp should store other in temp",
			fields: fields{
				data: nil,
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 2}},
			},
			expected: nil, // data remains nil, temp is set to other
		},
		{
			name: "nil data with existing temp should intersect temp with other",
			fields: fields{
				data: nil,
				temp: []itemInfo{{index: 1, position: 2}},
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			expected: []extendedItemInfo{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "no matching indices should result in empty set",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 2, position: 3}},
			},
			expected: []extendedItemInfo{},
		},
		{
			name: "matching indices but no valid positions should result in empty set",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{5}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			expected: []extendedItemInfo{},
		},
		{
			name: "matching indices with valid positions should be included",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			expected: []extendedItemInfo{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			expected: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			expected: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{4, 5}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ItemInfoSet{
				data: tt.fields.data,
				temp: tt.fields.temp,
			}
			s.Intersect(tt.args.other)

			// For the case where we expect temp to be set
			if tt.name == "nil data with nil temp should store other in temp" {
				if !reflect.DeepEqual(s.temp, tt.args.other) {
					t.Errorf("Intersect() temp = %v, want %v", s.temp, tt.args.other)
				}
				return
			}

			if !reflect.DeepEqual(s.data, tt.expected) {
				t.Errorf("Intersect() data = %v, want %v", s.data, tt.expected)
			}
		})
	}
}

func Test_ItemInfoSet_Intersection(t *testing.T) {
	type fields struct {
		data []extendedItemInfo
		temp []itemInfo
	}
	type args struct {
		other []itemInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ItemInfoSet
	}{
		{
			name: "basic intersection",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
				temp: nil,
			},
		},
		{
			name: "no matching indices",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 2, position: 3}},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{},
				temp: nil,
			},
		},
		{
			name: "matching indices but no valid positions",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{5}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{},
				temp: nil,
			},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 1}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{4, 5}}},
				temp: nil,
			},
		},
		{
			name: "nil data with nil temp should store other in temp",
			fields: fields{
				data: nil,
				temp: nil,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 2}},
			},
			want: &ItemInfoSet{
				data: nil,
				temp: []itemInfo{{index: 1, position: 2}},
			},
		},
		{
			name: "nil data with existing temp should intersect temp with other",
			fields: fields{
				data: nil,
				temp: []itemInfo{{index: 1, position: 2}},
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2, 3}}},
				temp: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ItemInfoSet{
				data: tt.fields.data,
				temp: tt.fields.temp,
			}
			if got := s.Intersection(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ItemInfoSet_IsEmpty(t *testing.T) {
	type fields struct {
		data []extendedItemInfo
		temp []itemInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil data and nil temp should not be empty",
			fields: fields{
				data: nil,
				temp: nil,
			},
			want: false,
		},
		{
			name: "empty data slice with nil temp should be empty",
			fields: fields{
				data: []extendedItemInfo{},
				temp: nil,
			},
			want: true,
		},
		{
			name: "non-empty data slice should not be empty",
			fields: fields{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			want: false,
		},
		{
			name: "nil data with empty temp slice should be empty",
			fields: fields{
				data: nil,
				temp: []itemInfo{},
			},
			want: true,
		},
		{
			name: "nil data with non-empty temp slice should not be empty",
			fields: fields{
				data: nil,
				temp: []itemInfo{{index: 1, position: 2}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ItemInfoSet{
				data: tt.fields.data,
				temp: tt.fields.temp,
			}
			if got := s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newItemInfoSet(t *testing.T) {
	tests := []struct {
		name string
		want *ItemInfoSet
	}{
		{
			name: "should create new empty set",
			want: &ItemInfoSet{
				data: nil,
				temp: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newItemInfoSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newItemInfoSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intersectNew(t *testing.T) {
	tests := []struct {
		name string
		cur  []itemInfo
		next []itemInfo
		want []extendedItemInfo
	}{
		{
			name: "empty current slice should return empty slice",
			cur:  []itemInfo{},
			next: []itemInfo{{index: 1, position: 2}},
			want: []extendedItemInfo{},
		},
		{
			name: "empty next slice should return empty slice",
			cur:  []itemInfo{{index: 1, position: 2}},
			next: []itemInfo{},
			want: []extendedItemInfo{},
		},
		{
			name: "no matching indices should return empty slice",
			cur:  []itemInfo{{index: 1, position: 2}},
			next: []itemInfo{{index: 2, position: 3}},
			want: []extendedItemInfo{},
		},
		{
			name: "matching indices but no valid positions should return empty slice",
			cur:  []itemInfo{{index: 1, position: 5}},
			next: []itemInfo{{index: 1, position: 3}},
			want: []extendedItemInfo{},
		},
		{
			name: "matching indices with valid positions should be included",
			cur:  []itemInfo{{index: 1, position: 2}},
			next: []itemInfo{{index: 1, position: 3}},
			want: []extendedItemInfo{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			cur:  []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}},
			next: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			want: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			cur:  []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}, {index: 3, position: 4}},
			next: []itemInfo{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			want: []extendedItemInfo{
				{index: 1, positions: []int{2, 3}},
				{index: 3, positions: []int{4, 5}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intersectNew(tt.cur, tt.next)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersectNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intersectExisting(t *testing.T) {
	tests := []struct {
		name string
		cur  []extendedItemInfo
		next []itemInfo
		want []extendedItemInfo
	}{
		{
			name: "empty current slice should return empty slice",
			cur:  []extendedItemInfo{},
			next: []itemInfo{{index: 1, position: 2}},
			want: []extendedItemInfo{},
		},
		{
			name: "empty next slice should return empty slice",
			cur:  []extendedItemInfo{{index: 1, positions: []int{2}}},
			next: []itemInfo{},
			want: []extendedItemInfo{},
		},
		{
			name: "no matching indices should return empty slice",
			cur:  []extendedItemInfo{{index: 1, positions: []int{2}}},
			next: []itemInfo{{index: 2, position: 3}},
			want: []extendedItemInfo{},
		},
		{
			name: "matching indices but no valid positions should return empty slice",
			cur:  []extendedItemInfo{{index: 1, positions: []int{5}}},
			next: []itemInfo{{index: 1, position: 3}},
			want: []extendedItemInfo{},
		},
		{
			name: "matching indices with valid positions should be included",
			cur:  []extendedItemInfo{{index: 1, positions: []int{2}}},
			next: []itemInfo{{index: 1, position: 3}},
			want: []extendedItemInfo{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			cur:  []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
			next: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			want: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			cur:  []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{3, 4}}},
			next: []itemInfo{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			want: []extendedItemInfo{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{3, 4, 5}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]extendedItemInfo, min(len(tt.cur), len(tt.next)))
			got := intersectExisting(buf, tt.cur, tt.next)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersectExisting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Union(t *testing.T) {
	tests := []struct {
		name string
		sets []*ItemInfoSet
		want *ItemInfoSet
	}{
		{
			name: "no sets should return empty set",
			sets: []*ItemInfoSet{},
			want: &ItemInfoSet{
				data: []extendedItemInfo{},
				temp: nil,
			},
		},
		{
			name: "single empty set should return empty set",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{},
				temp: nil,
			},
		},
		{
			name: "multiple empty sets should return empty set",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{}, temp: nil},
				{data: []extendedItemInfo{}, temp: nil},
				{data: []extendedItemInfo{}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{},
				temp: nil,
			},
		},
		{
			name: "single set with one item should return that item",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{{index: 1, positions: []int{2}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}},
				temp: nil,
			},
		},
		{
			name: "single set with multiple items should return all items",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
		},
		{
			name: "multiple sets with no overlapping indices",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{{index: 1, positions: []int{2}}}, temp: nil},
				{data: []extendedItemInfo{{index: 2, positions: []int{3}}}, temp: nil},
				{data: []extendedItemInfo{{index: 3, positions: []int{4}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
		},
		{
			name: "multiple sets with overlapping indices",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 3, positions: []int{4}}}, temp: nil},
				{data: []extendedItemInfo{{index: 1, positions: []int{3}}, {index: 2, positions: []int{5}}}, temp: nil},
				{data: []extendedItemInfo{{index: 2, positions: []int{6}}, {index: 4, positions: []int{7}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{5}}, {index: 3, positions: []int{4}}, {index: 4, positions: []int{7}}},
				temp: nil,
			},
		},
		{
			name: "complex case with multiple indices and positions",
			sets: []*ItemInfoSet{
				{data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 3, positions: []int{4}}, {index: 5, positions: []int{6}}}, temp: nil},
				{data: []extendedItemInfo{{index: 2, positions: []int{3}}, {index: 4, positions: []int{5}}, {index: 6, positions: []int{7}}}, temp: nil},
				{data: []extendedItemInfo{{index: 1, positions: []int{8}}, {index: 3, positions: []int{9}}, {index: 7, positions: []int{10}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}, {index: 4, positions: []int{5}}, {index: 5, positions: []int{6}}, {index: 6, positions: []int{7}}, {index: 7, positions: []int{10}}},
				temp: nil,
			},
		},
		{
			name: "sets with temp fields",
			sets: []*ItemInfoSet{
				{data: nil, temp: []itemInfo{{index: 1, position: 2}}},
				{data: []extendedItemInfo{{index: 2, positions: []int{3}}}, temp: nil},
			},
			want: &ItemInfoSet{
				data: []extendedItemInfo{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.sets...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}
