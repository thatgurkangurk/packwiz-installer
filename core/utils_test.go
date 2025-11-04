package core

import (
	"cmp"
	"slices"
	"testing"
)

func Test_diffSlice(t *testing.T) {
	type args[E cmp.Ordered] struct {
		old []E
		new []E
	}
	tests := []struct {
		name        string
		args        args[int]
		wantAdded   []int
		wantRemoved []int
	}{

		{
			name: "all-removed",
			args: args[int]{
				old: []int{1, 2, 3},
				new: []int{},
			},
			wantRemoved: []int{1, 2, 3},
			wantAdded:   []int{},
		},
		{
			name: "all-added",
			args: args[int]{
				old: []int{},
				new: []int{1, 2, 3},
			},
			wantRemoved: []int{},
			wantAdded:   []int{1, 2, 3},
		},
		{

			name: "old-pivot",
			args: args[int]{
				old: []int{2},
				new: []int{0, 1, 3, 4},
			},
			wantRemoved: []int{2},
			wantAdded:   []int{0, 1, 3, 4},
		},
		{
			name: "old-pivot2",
			args: args[int]{
				old: []int{2},
				new: []int{0, 1, 2, 3, 4},
			},
			wantRemoved: []int{},
			wantAdded:   []int{0, 1, 3, 4},
		},
		{
			name: "new-pivot",
			args: args[int]{
				old: []int{0, 1, 3, 4},
				new: []int{2},
			},
			wantRemoved: []int{0, 1, 3, 4},
			wantAdded:   []int{2},
		},
		{
			name: "new-pivot2",
			args: args[int]{
				old: []int{0, 1, 2, 3, 4},
				new: []int{2},
			},
			wantRemoved: []int{0, 1, 3, 4},
			wantAdded:   []int{},
		},
		{
			name: "W",
			args: args[int]{
				old: []int{1, 2, 3, 4, 5, 6},
				new: []int{2, 5},
			},
			wantRemoved: []int{1, 3, 4, 6},
			wantAdded:   []int{},
		},
		{
			name: "M",
			args: args[int]{
				old: []int{2, 5},
				new: []int{1, 2, 3, 4, 5, 6},
			},
			wantRemoved: []int{},
			wantAdded:   []int{1, 3, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdded, gotRemoved, _ := diffSlice(tt.args.old, tt.args.new)
			if slices.Compare(gotAdded, tt.wantAdded) != 0 {
				t.Errorf("diffSlice() gotAdded = %v, want %v", gotAdded, tt.wantAdded)
			}
			if slices.Compare(gotRemoved, tt.wantRemoved) != 0 {
				t.Errorf("diffSlice() gotRemoved = %v, want %v", gotRemoved, tt.wantRemoved)
			}
		})
	}
}
