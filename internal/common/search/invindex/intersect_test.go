package invindex

import (
	"reflect"
	"testing"
)

func Test_ItemInfoSet_Intersect(t *testing.T) {
	type fields struct {
		data []indexPositions
		temp []indexPosition
	}
	type args struct {
		other []indexPosition
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected []indexPositions
	}{
		{
			name: "empty set should remain empty",
			fields: fields{
				data: []indexPositions{},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 2}},
			},
			expected: []indexPositions{},
		},
		{
			name: "nil other should result in empty set",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: nil,
			},
			expected: []indexPositions{},
		},
		{
			name: "nil data with nil temp should store other in temp",
			fields: fields{
				data: nil,
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 2}},
			},
			expected: nil, // data remains nil, temp is set to other
		},
		{
			name: "nil data with existing temp should intersect temp with other",
			fields: fields{
				data: nil,
				temp: []indexPosition{{index: 1, position: 2}},
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}},
			},
			expected: []indexPositions{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "no matching indices should result in empty set",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 2, position: 3}},
			},
			expected: []indexPositions{},
		},
		{
			name: "matching indices but no valid positions should result in empty set",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{5}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}},
			},
			expected: []indexPositions{},
		},
		{
			name: "matching indices with valid positions should be included",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}},
			},
			expected: []indexPositions{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			expected: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			expected: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{4, 5}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &indexPositionSet{
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
		data []indexPositions
		temp []indexPosition
	}
	type args struct {
		other []indexPosition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *indexPositionSet
	}{
		{
			name: "basic intersection",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
				temp: nil,
			},
		},
		{
			name: "no matching indices",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 2, position: 3}},
			},
			want: &indexPositionSet{
				data: []indexPositions{},
				temp: nil,
			},
		},
		{
			name: "matching indices but no valid positions",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{5}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}},
			},
			want: &indexPositionSet{
				data: []indexPositions{},
				temp: nil,
			},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
			args: args{
				other: []indexPosition{{index: 1, position: 1}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{4, 5}}},
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
				other: []indexPosition{{index: 1, position: 2}},
			},
			want: &indexPositionSet{
				data: nil,
				temp: []indexPosition{{index: 1, position: 2}},
			},
		},
		{
			name: "nil data with existing temp should intersect temp with other",
			fields: fields{
				data: nil,
				temp: []indexPosition{{index: 1, position: 2}},
			},
			args: args{
				other: []indexPosition{{index: 1, position: 3}},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2, 3}}},
				temp: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &indexPositionSet{
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
		data []indexPositions
		temp []indexPosition
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
				data: []indexPositions{},
				temp: nil,
			},
			want: true,
		},
		{
			name: "non-empty data slice should not be empty",
			fields: fields{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
			want: false,
		},
		{
			name: "nil data with empty temp slice should be empty",
			fields: fields{
				data: nil,
				temp: []indexPosition{},
			},
			want: true,
		},
		{
			name: "nil data with non-empty temp slice should not be empty",
			fields: fields{
				data: nil,
				temp: []indexPosition{{index: 1, position: 2}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &indexPositionSet{
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
		want *indexPositionSet
	}{
		{
			name: "should create new empty set",
			want: &indexPositionSet{
				data: nil,
				temp: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newIndexPositionSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newIndexPositionSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intersectNew(t *testing.T) {
	tests := []struct {
		name string
		cur  []indexPosition
		next []indexPosition
		want []indexPositions
	}{
		{
			name: "empty current slice should return empty slice",
			cur:  []indexPosition{},
			next: []indexPosition{{index: 1, position: 2}},
			want: []indexPositions{},
		},
		{
			name: "empty next slice should return empty slice",
			cur:  []indexPosition{{index: 1, position: 2}},
			next: []indexPosition{},
			want: []indexPositions{},
		},
		{
			name: "no matching indices should return empty slice",
			cur:  []indexPosition{{index: 1, position: 2}},
			next: []indexPosition{{index: 2, position: 3}},
			want: []indexPositions{},
		},
		{
			name: "matching indices but no valid positions should return empty slice",
			cur:  []indexPosition{{index: 1, position: 5}},
			next: []indexPosition{{index: 1, position: 3}},
			want: []indexPositions{},
		},
		{
			name: "matching indices with valid positions should be included",
			cur:  []indexPosition{{index: 1, position: 2}},
			next: []indexPosition{{index: 1, position: 3}},
			want: []indexPositions{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			cur:  []indexPosition{{index: 1, position: 2}, {index: 2, position: 3}},
			next: []indexPosition{{index: 1, position: 3}, {index: 2, position: 4}},
			want: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			cur:  []indexPosition{{index: 1, position: 2}, {index: 2, position: 3}, {index: 3, position: 4}},
			next: []indexPosition{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			want: []indexPositions{
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
		cur  []indexPositions
		next []indexPosition
		want []indexPositions
	}{
		{
			name: "empty current slice should return empty slice",
			cur:  []indexPositions{},
			next: []indexPosition{{index: 1, position: 2}},
			want: []indexPositions{},
		},
		{
			name: "empty next slice should return empty slice",
			cur:  []indexPositions{{index: 1, positions: []int{2}}},
			next: []indexPosition{},
			want: []indexPositions{},
		},
		{
			name: "no matching indices should return empty slice",
			cur:  []indexPositions{{index: 1, positions: []int{2}}},
			next: []indexPosition{{index: 2, position: 3}},
			want: []indexPositions{},
		},
		{
			name: "matching indices but no valid positions should return empty slice",
			cur:  []indexPositions{{index: 1, positions: []int{5}}},
			next: []indexPosition{{index: 1, position: 3}},
			want: []indexPositions{},
		},
		{
			name: "matching indices with valid positions should be included",
			cur:  []indexPositions{{index: 1, positions: []int{2}}},
			next: []indexPosition{{index: 1, position: 3}},
			want: []indexPositions{{index: 1, positions: []int{2, 3}}},
		},
		{
			name: "multiple matching indices with valid positions",
			cur:  []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
			next: []indexPosition{{index: 1, position: 3}, {index: 2, position: 4}},
			want: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 2, positions: []int{3, 4}}},
		},
		{
			name: "complex case with multiple indices and positions",
			cur:  []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{3, 4}}},
			next: []indexPosition{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			want: []indexPositions{{index: 1, positions: []int{2, 3}}, {index: 3, positions: []int{3, 4, 5}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]indexPositions, min(len(tt.cur), len(tt.next)))
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
		sets []*indexPositionSet
		want *indexPositionSet
	}{
		{
			name: "no sets should return empty set",
			sets: []*indexPositionSet{},
			want: &indexPositionSet{
				data: []indexPositions{},
				temp: nil,
			},
		},
		{
			name: "single empty set should return empty set",
			sets: []*indexPositionSet{
				{data: []indexPositions{}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{},
				temp: nil,
			},
		},
		{
			name: "multiple empty sets should return empty set",
			sets: []*indexPositionSet{
				{data: []indexPositions{}, temp: nil},
				{data: []indexPositions{}, temp: nil},
				{data: []indexPositions{}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{},
				temp: nil,
			},
		},
		{
			name: "single set with one item should return that item",
			sets: []*indexPositionSet{
				{data: []indexPositions{{index: 1, positions: []int{2}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}},
				temp: nil,
			},
		},
		{
			name: "single set with multiple items should return all items",
			sets: []*indexPositionSet{
				{data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
				temp: nil,
			},
		},
		{
			name: "multiple sets with no overlapping indices",
			sets: []*indexPositionSet{
				{data: []indexPositions{{index: 1, positions: []int{2}}}, temp: nil},
				{data: []indexPositions{{index: 2, positions: []int{3}}}, temp: nil},
				{data: []indexPositions{{index: 3, positions: []int{4}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}},
				temp: nil,
			},
		},
		{
			name: "multiple sets with overlapping indices",
			sets: []*indexPositionSet{
				{data: []indexPositions{{index: 1, positions: []int{2}}, {index: 3, positions: []int{4}}}, temp: nil},
				{data: []indexPositions{{index: 1, positions: []int{3}}, {index: 2, positions: []int{5}}}, temp: nil},
				{data: []indexPositions{{index: 2, positions: []int{6}}, {index: 4, positions: []int{7}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{5}}, {index: 3, positions: []int{4}}, {index: 4, positions: []int{7}}},
				temp: nil,
			},
		},
		{
			name: "complex case with multiple indices and positions",
			sets: []*indexPositionSet{
				{data: []indexPositions{{index: 1, positions: []int{2}}, {index: 3, positions: []int{4}}, {index: 5, positions: []int{6}}}, temp: nil},
				{data: []indexPositions{{index: 2, positions: []int{3}}, {index: 4, positions: []int{5}}, {index: 6, positions: []int{7}}}, temp: nil},
				{data: []indexPositions{{index: 1, positions: []int{8}}, {index: 3, positions: []int{9}}, {index: 7, positions: []int{10}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}, {index: 3, positions: []int{4}}, {index: 4, positions: []int{5}}, {index: 5, positions: []int{6}}, {index: 6, positions: []int{7}}, {index: 7, positions: []int{10}}},
				temp: nil,
			},
		},
		{
			name: "sets with temp fields",
			sets: []*indexPositionSet{
				{data: nil, temp: []indexPosition{{index: 1, position: 2}}},
				{data: []indexPositions{{index: 2, positions: []int{3}}}, temp: nil},
			},
			want: &indexPositionSet{
				data: []indexPositions{{index: 1, positions: []int{2}}, {index: 2, positions: []int{3}}},
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
