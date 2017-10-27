package flatten

import (
	"reflect"
	"testing"
)

// This implements a flatten function for arbitrarily nested integer lists.
// I've chosen Go as the implementation language as it's the one I'm most
// comfortable in right now, but in a way this problem is a bad fit for the
// language. As a statically typed language something like an arbitrarily
// nested list of integers just doesn't happen very often in the wild, and
// if it did the correct solution would be to refactor that part instead. :)
// As it is, this requires some gymnastics using interface{} types (Go's
// most generic type) which is ugly, but in the end it implements a simple
// recursive "flattener" and contains tests for the most common cases. This
// is runnable by saving this as flatten_test.go in a directory and running
// "go test -v" in that directory.

// flatten returns l as a flattened list of integers. The parameter l may be
// of concrete type []int, int or []interface{}. Parameters of type
// []interface{} are recursively flattened. If the argument or any of it's
// elements are not of the mentioned types, ok=false is returned.
func flatten(l interface{}) (ints []int, ok bool) {
	switch t := l.(type) {
	case nil:
		// nil is an instance of the empty slice, which is already flat.
		return []int{}, true

	case []int:
		// A slice of ints is already flat. Return it as is.
		return t, true

	case int:
		// An int just needs to be wrapped into a []int to be a flat slice.
		return []int{t}, true

	case []interface{}:
		// A slice of interface{} needs flattening of each element. If it's
		// zero length, we can just return the empty slice and be done with it.
		if len(t) == 0 {
			return []int{}, true
		}

		// For the more general case we use a recursive approach - the
		// flattened list is the concatenation of the flattened first
		// element and the flattened rest of the list. In either case we
		// bail out early if the element is not flattenable.
		first, ok := flatten(t[0])
		if !ok {
			return nil, false
		}

		rest, ok := flatten(t[1:])
		if !ok {
			return nil, false
		}

		return append(first, rest...), true

	default:
		// The parameter is of some type we don't support.
		return nil, false
	}
}

func TestFlatten(t *testing.T) {
	cases := []struct {
		in  interface{} // parameter to flatten
		out []int       // expected flatten "ints" output
		ok  bool        // expected flatten "ok" output
	}{
		// [] => []
		{[]int{}, []int{}, true},

		// nil => []
		{nil, []int{}, true},

		// [[], nil] => []
		{[]interface{}{[]int{}, nil}, []int{}, true},

		// 42 => [42]
		{42, []int{42}, true},

		// [1, 2, 3] => [1, 2, 3]
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},

		// [1, [2, 3]] => [1, 2, 3]
		{
			[]interface{}{1, []int{2, 3}},
			[]int{1, 2, 3},
			true,
		},

		// [[1, 2,], 3] => [1, 2, 3]
		{
			[]interface{}{[]int{1, 2}, 3},
			[]int{1, 2, 3},
			true,
		},

		// [1, [2, 3], 4] => [1, 2, 3, 4]
		{
			[]interface{}{1, []int{2, 3}, 4},
			[]int{1, 2, 3, 4},
			true,
		},

		// [1, [2, [3], 4], 5] => [1, 2, 3, 4, 5]
		{
			[]interface{}{1, []interface{}{2, []interface{}{3}, 4}, 5},
			[]int{1, 2, 3, 4, 5},
			true,
		},

		// [1, nil, [nil], [2, [], [], [3], nil], 5] => [1, 2, 3, 5]
		// (nils and empty slices are trimmed as part of the flattening)
		{
			[]interface{}{1, nil, []interface{}{nil}, []interface{}{2, []interface{}{}, []int{}, []interface{}{3}, nil}, 5},
			[]int{1, 2, 3, 5},
			true,
		},

		// "foo" => not flattenable (not integer)
		{"foo", nil, false},

		// 42.5  => not flattenable (not integer)
		{42.5, nil, false},

		// [[1, 2, [3]], 4] => [1, 2, 3, 4]
		{
			[]interface{}{[]interface{}{1, 2, []int{3}, 4}},
			[]int{1, 2, 3, 4},
			true,
		},
	}

	for i, c := range cases {
		out, ok := flatten(c.in)
		if ok != c.ok {
			t.Errorf("Case %d: got ok=%v, expected %v", i, ok, c.ok)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("Case %d: got out=%#v, expected %#v", i, out, c.out)
		}
	}
}
