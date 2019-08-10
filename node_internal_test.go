package immutableList

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNodeAppend(t *testing.T) {
	expected := make([]Object, 0)
	var list node = &emptyNode{}
	for length := 0; length <= 4096; length += 1 {
		expected = insertToSlice(expected, length, val(length))
		list = list.insert(length, val(length))
		validateNode(t, list, expected)
	}
}

func TestNodePrepend(t *testing.T) {
	expected := make([]Object, 0)
	var list node = &emptyNode{}
	for length := 0; length <= 4096; length += 1 {
		expected = insertToSlice(expected, 0, val(length))
		list = list.insert(0, val(length))
		validateNode(t, list, expected)
	}
}

func TestNodeInsert(t *testing.T) {
	expected := make([]Object, 0)
	var list node = &emptyNode{}
	expected = insertToSlice(expected, 0, val(0))
	list = list.insert(0, val(0))
	for length := 1; length <= 4096; length += 1 {
		index := rand.Intn(length)
		expected = insertToSlice(expected, index, val(length))
		list = list.insert(index, val(length))
		validateNode(t, list, expected)
	}
}

func TestNodePop(t *testing.T) {
	b, e := listAppendLists(1024)
	for len(e) > 0 {
		var value Object
		value, b = b.pop()
		if value != e[0] {
			t.Error(fmt.Sprintf("incorrect value from pop(): expected=%v actual=%v", e[0], value))
		}
		e = deleteFromSlice(e, 0)
		validateNode(t, b, e)
	}
}

func TestNodeSet(t *testing.T) {
	b, e := listAppendLists(87)
	next := 1000
	for i := 0; i < len(e); i++ {
		b = b.set(i, val(next))
		e[i] = val(next)
	}
	validateNode(t, b, e)
}

func TestNodeGetFirstLast(t *testing.T) {
	expected := make([]Object, 0)
	var list node = &emptyNode{}
	for length := 0; length <= 30; length += 1 {
		expected = insertToSlice(expected, length, val(length))
		list = list.insert(length, val(length))
		if expected[0] != list.getFirst() {
			t.Error(fmt.Sprintf("incorrect value from getFirst(): expected=%v actual=%v", expected[0], list.getFirst()))
		}
		if expected[len(expected)-1] != list.getLast() {
			t.Error(fmt.Sprintf("incorrect value from getLast(): expected=%v actual=%v", expected[len(expected)-1], list.getLast()))
		}
	}
}

func TestNodeAppendList(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		alen := rand.Intn(loop)
		blen := rand.Intn(loop)
		ab, ae := listAppendLists(alen)
		bb, be := listAppendLists(blen)
		expected := append(ae, be...)
		list := appendNodes(ab, bb)
		validateNode(t, list, expected)
	}
}

func TestNodeDelete(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		list, expected := listAppendLists(loop)
		for list.size() > 0 {
			index := rand.Intn(list.size())
			list = list.delete(index)
			expected = deleteFromSlice(expected, index)
			validateNode(t, list, expected)
		}
	}
}

func TestNodeHead(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		list, expected := listAppendLists(loop)
		for list.size() > 0 {
			index := rand.Intn(list.size() + 1)
			list = list.head(index)
			for len(expected) > index {
				expected = deleteFromSlice(expected, index)
			}
			validateNode(t, list, expected)
		}
	}
}

func TestNodeTail(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		list, expected := listAppendLists(loop)
		for list.size() > 0 {
			index := rand.Intn(list.size() + 1)
			list = list.tail(index)
			for i := 0; i < index; i++ {
				expected = deleteFromSlice(expected, 0)
			}
			validateNode(t, list, expected)
		}
	}
}

func TestNodeIterator(t *testing.T) {
	for length := 0; length <= 1024; length++ {
		expected := make([]Object, 0)
		var list node = &emptyNode{}
		for i := 0; i <= length; i += 1 {
			expected = insertToSlice(expected, i, val(i))
			list = list.append(val(i))
		}
		actual := createEmptyLeafNode()
		for i := createIterator(list); i.Next(); {
			actual = actual.insert(actual.size(), i.Get())
		}
		validateNode(t, actual, expected)
	}
}

func BenchmarkNodeGet1000(b *testing.B) {
	benchmarkNodeGet(1000, b)
}

func BenchmarkNodeGet10000(b *testing.B) {
	benchmarkNodeGet(10000, b)
}

func BenchmarkNodeGet100000(b *testing.B) {
	benchmarkNodeGet(100000, b)
}

func BenchmarkNodeDelete1000(b *testing.B) {
	benchmarkNodeDelete(1000, b)
}

func BenchmarkNodeDelete10000(b *testing.B) {
	benchmarkNodeDelete(10000, b)
}

func BenchmarkNodeDelete100000(b *testing.B) {
	benchmarkNodeDelete(100000, b)
}

func benchmarkNodeGet(size int, b *testing.B) {
	list := createNodeListForBenchmark(size)
	for i := 1; i <= b.N; i++ {
		list.get(i % size)
	}
}

func benchmarkNodeDelete(size int, b *testing.B) {
	orig := createNodeListForBenchmark(size)
	list := orig
	for i := 1; i <= b.N; i++ {
		list = list.delete(i % list.size())
		if list.size() == 0 {
			list = orig
		}
	}
}

func createNodeListForBenchmark(size int) node {
	list := createEmptyLeafNode()
	for i := 1; i <= size; i++ {
		list = list.append(val(i))
	}
	return list
}

func validateNode(t *testing.T, b node, e []Object) {
	if b.size() != len(e) {
		t.Error(fmt.Sprintf("incorrect size: b=%d e=%d", b.size(), len(e)))
	}
	for i := 0; i < b.size(); i += 1 {
		if b.get(i) != e[i] {
			t.Error(fmt.Sprintf("incorrect value: i=%d b=%v e=%v", i, b.get(i), e[i]))
		}
	}
	b.checkInvariants(func(message string) {
		t.Error(message)
	}, true)
}

func listAppendLists(length int) (node, []Object) {
	expected := make([]Object, 0)
	var list node = &emptyNode{}
	for i := 0; i < length; i += 1 {
		value := val(i)
		expected = insertToSlice(expected, i, value)
		list = list.insert(i, value)
	}
	return list, expected
}
