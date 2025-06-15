package search

import (
	"reflect"
	"testing"
)

func Test_itemInfoSet_Intersect(t *testing.T) {
	type fields struct {
		Data           []itemInfo
		shouldCopyNext bool
	}
	type args struct {
		other []itemInfo
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected []itemInfo
	}{
		{
			name: "empty set should remain empty",
			fields: fields{
				Data:           []itemInfo{},
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 2}},
			},
			expected: []itemInfo{},
		},
		{
			name: "nil other should result in empty set",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}},
				shouldCopyNext: false,
			},
			args: args{
				other: nil,
			},
			expected: []itemInfo{},
		},
		{
			name: "nil data should be replaced with other",
			fields: fields{
				Data:           nil,
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 2}},
			},
			expected: []itemInfo{{index: 1, position: 2}},
		},
		{
			name: "no matching indices should result in empty set",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}},
				shouldCopyNext: true,
			},
			args: args{
				other: []itemInfo{{index: 2, position: 3}},
			},
			expected: []itemInfo{},
		},
		{
			name: "matching indices but no valid positions should result in empty set",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 5}},
				shouldCopyNext: true,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			expected: []itemInfo{},
		},
		{
			name: "matching indices with valid positions should be included",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}},
				shouldCopyNext: true,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			expected: []itemInfo{{index: 1, position: 3}},
		},
		{
			name: "multiple matching indices with valid positions",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}},
				shouldCopyNext: true,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			expected: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}, {index: 3, position: 4}},
				shouldCopyNext: true,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 1}, {index: 1, position: 2}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			expected: []itemInfo{{index: 1, position: 3}, {index: 3, position: 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &itemInfoSet{
				Data:           tt.fields.Data,
				shouldCopyNext: tt.fields.shouldCopyNext,
			}
			s.Intersect(tt.args.other)
			if !reflect.DeepEqual(s.Data, tt.expected) {
				t.Errorf("Intersect() = %v, want %v", s.Data, tt.expected)
			}
		})
	}
}

func Test_itemInfoSet_Intersection(t *testing.T) {
	type fields struct {
		Data           []itemInfo
		shouldCopyNext bool
	}
	type args struct {
		other []itemInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *itemInfoSet
	}{
		{
			name: "basic intersection",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}},
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
			},
			want: &itemInfoSet{
				Data:           []itemInfo{{index: 1, position: 3}, {index: 2, position: 4}},
				shouldCopyNext: false,
			},
		},
		{
			name: "no matching indices",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}},
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 2, position: 3}},
			},
			want: &itemInfoSet{
				Data:           []itemInfo{},
				shouldCopyNext: false,
			},
		},
		{
			name: "matching indices but no valid positions",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 5}},
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 3}},
			},
			want: &itemInfoSet{
				Data:           []itemInfo{},
				shouldCopyNext: false,
			},
		},
		{
			name: "complex case with multiple indices and positions",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}, {index: 2, position: 3}, {index: 3, position: 4}},
				shouldCopyNext: false,
			},
			args: args{
				other: []itemInfo{{index: 1, position: 1}, {index: 1, position: 3}, {index: 2, position: 2}, {index: 3, position: 5}},
			},
			want: &itemInfoSet{
				Data:           []itemInfo{{index: 1, position: 3}, {index: 3, position: 5}},
				shouldCopyNext: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &itemInfoSet{
				Data:           tt.fields.Data,
				shouldCopyNext: tt.fields.shouldCopyNext,
			}
			if got := s.Intersection(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_itemInfoSet_IsEmpty(t *testing.T) {
	type fields struct {
		Data           []itemInfo
		shouldCopyNext bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil data should not be empty",
			fields: fields{
				Data:           nil,
				shouldCopyNext: false,
			},
			want: false,
		},
		{
			name: "empty slice should be empty",
			fields: fields{
				Data:           []itemInfo{},
				shouldCopyNext: false,
			},
			want: true,
		},
		{
			name: "non-empty slice should not be empty",
			fields: fields{
				Data:           []itemInfo{{index: 1, position: 2}},
				shouldCopyNext: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &itemInfoSet{
				Data:           tt.fields.Data,
				shouldCopyNext: tt.fields.shouldCopyNext,
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
		want *itemInfoSet
	}{
		{
			name: "should create new empty set",
			want: &itemInfoSet{
				Data:           nil,
				shouldCopyNext: false,
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
