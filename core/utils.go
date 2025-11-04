// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package core

import (
	"cmp"
	"slices"
)

// diffSlice calculates elements added, removed and unchanged between two slices.
// The elements must be cmp.Ordered. Neither input slice is mutated (we clone them).
//
// Returns:
//
//	added   - items present in new but not in old
//	removed - items present in old but not in new
//	unchanged - items present in both
func diffSlice[S ~[]E, E cmp.Ordered](old, new S) (added, removed, unchanged S) {
	// quick empty cases; return clones to avoid aliasing caller data
	if len(old) == 0 && len(new) == 0 {
		return nil, nil, nil
	}
	if len(old) == 0 {
		return slices.Clone(new), nil, nil
	}
	if len(new) == 0 {
		return nil, slices.Clone(old), nil
	}

	// clone inputs so we don't mutate callers' slices when sorting
	o := slices.Clone(old)
	n := slices.Clone(new)

	// ensure sorted
	if !slices.IsSorted(o) {
		slices.Sort(o)
	}
	if !slices.IsSorted(n) {
		slices.Sort(n)
	}

	added = make(S, 0)
	removed = make(S, 0)
	unchanged = make(S, 0)

	i, j := 0, 0
	for i < len(o) && j < len(n) {
		ov := o[i]
		nv := n[j]
		switch {
		case cmp.Compare(ov, nv) < 0: // ov < nv => ov removed
			removed = append(removed, ov)
			i++
		case cmp.Compare(ov, nv) == 0: // equal => unchanged
			unchanged = append(unchanged, nv)
			i++
			j++
		default: // ov > nv => nv added
			added = append(added, nv)
			j++
		}
	}

	// remaining tails
	for ; i < len(o); i++ {
		removed = append(removed, o[i])
	}
	for ; j < len(n); j++ {
		added = append(added, n[j])
	}

	return added, removed, unchanged
}

// diffSliceFunc is the same as diffSlice but accepts a custom comparator function.
// cmpFn must return a negative number if a < b, zero if a == b, positive if a > b.
func diffSliceFunc[S ~[]E, E any](old, new S, cmpFn func(a, b E) int) (added, removed, unchanged S) {
	// quick empty cases; return clones to avoid aliasing caller data
	if len(old) == 0 && len(new) == 0 {
		return nil, nil, nil
	}
	if len(old) == 0 {
		return slices.Clone(new), nil, nil
	}
	if len(new) == 0 {
		return nil, slices.Clone(old), nil
	}

	// clone inputs so we don't mutate callers' slices when sorting
	o := slices.Clone(old)
	n := slices.Clone(new)

	// ensure sorted with the provided comparator
	if !slices.IsSortedFunc(o, cmpFn) {
		slices.SortFunc(o, cmpFn)
	}
	if !slices.IsSortedFunc(n, cmpFn) {
		slices.SortFunc(n, cmpFn)
	}

	added = make(S, 0)
	removed = make(S, 0)
	unchanged = make(S, 0)

	i, j := 0, 0
	for i < len(o) && j < len(n) {
		ov := o[i]
		nv := n[j]
		c := cmpFn(ov, nv)
		switch {
		case c < 0: // ov < nv => ov removed
			removed = append(removed, ov)
			i++
		case c == 0: // equal => unchanged
			unchanged = append(unchanged, nv)
			i++
			j++
		default: // ov > nv => nv added
			added = append(added, nv)
			j++
		}
	}

	// remaining tails
	for ; i < len(o); i++ {
		removed = append(removed, o[i])
	}
	for ; j < len(n); j++ {
		added = append(added, n[j])
	}

	return added, removed, unchanged
}
