package immutableList

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBinaryAppend(t *testing.T) {
	expected := make([]Object, 0)
	var binary binaryNode = &emptyLeafNode{}
	for length := 0; length <= 4096; length += 1 {
		expected = insertToSlice(expected, length, val(length))
		binary = binary.insert(length, val(length))
		validateBinaryNode(t, binary, expected)
	}
}

func TestBinaryPrepend(t *testing.T) {
	expected := make([]Object, 0)
	var binary binaryNode = &emptyLeafNode{}
	for length := 0; length <= 4096; length += 1 {
		expected = insertToSlice(expected, 0, val(length))
		binary = binary.insert(0, val(length))
		validateBinaryNode(t, binary, expected)
	}
}

func TestBinaryInsert(t *testing.T) {
	expected := make([]Object, 0)
	var binary binaryNode = &emptyLeafNode{}
	expected = insertToSlice(expected, 0, val(0))
	binary = binary.insert(0, val(0))
	for length := 1; length <= 4096; length += 1 {
		index := rand.Intn(length)
		expected = insertToSlice(expected, index, val(length))
		binary = binary.insert(index, val(length))
		validateBinaryNode(t, binary, expected)
	}
}

func TestBinaryPop(t *testing.T) {
	b, e := binaryAppendLists(1024)
	for len(e) > 0 {
		var value Object
		value, b = b.pop()
		if value != e[0] {
			t.Error(fmt.Sprintf("incorrect value from pop(): expected=%v actual=%v", e[0], value))
		}
		e = deleteFromSlice(e, 0)
		validateBinaryNode(t, b, e)
	}
}

func TestBinarySet(t *testing.T) {
	b, e := binaryAppendLists(87)
	next := 1000
	for i := 0; i < len(e); i++ {
		b = b.set(i, val(next))
		e[i] = val(next)
	}
	validateBinaryNode(t, b, e)
}

func TestBinaryGetFirstLast(t *testing.T) {
	expected := make([]Object, 0)
	var binary binaryNode = &emptyLeafNode{}
	for length := 0; length <= 30; length += 1 {
		expected = insertToSlice(expected, length, val(length))
		binary = binary.insert(length, val(length))
		if expected[0] != binary.getFirst() {
			t.Error(fmt.Sprintf("incorrect value from getFirst(): expected=%v actual=%v", expected[0], binary.getFirst()))
		}
		if expected[len(expected)-1] != binary.getLast() {
			t.Error(fmt.Sprintf("incorrect value from getLast(): expected=%v actual=%v", expected[len(expected)-1], binary.getLast()))
		}
	}
}

func TestBinaryAppendList(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		alen := rand.Intn(loop)
		blen := rand.Intn(loop)
		ab, ae := binaryAppendLists(alen)
		bb, be := binaryAppendLists(blen)
		expected := append(ae, be...)
		binary := appendBinaryNodes(ab, bb)
		validateBinaryNode(t, binary, expected)
	}
}

func TestBinaryDelete(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		binary, expected := binaryAppendLists(loop)
		for binary.size() > 0 {
			index := rand.Intn(binary.size())
			binary = binary.delete(index)
			expected = deleteFromSlice(expected, index)
			validateBinaryNode(t, binary, expected)
		}
	}
}

func TestBinaryHead(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		binary, expected := binaryAppendLists(loop)
		for binary.size() > 0 {
			index := rand.Intn(binary.size() + 1)
			binary = binary.head(index)
			for len(expected) > index {
				expected = deleteFromSlice(expected, index)
			}
			validateBinaryNode(t, binary, expected)
		}
	}
}

func TestBinaryTail(t *testing.T) {
	for loop := 1; loop <= 500; loop += 1 {
		binary, expected := binaryAppendLists(loop)
		for binary.size() > 0 {
			index := rand.Intn(binary.size() + 1)
			binary = binary.tail(index)
			for i := 0; i < index; i++ {
				expected = deleteFromSlice(expected, 0)
			}
			validateBinaryNode(t, binary, expected)
		}
	}
}

func validateBinaryNode(t *testing.T, b binaryNode, e []Object) {
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

func binaryAppendLists(length int) (binaryNode, []Object) {
	expected := make([]Object, 0)
	var binary binaryNode = &emptyLeafNode{}
	for i := 0; i < length; i += 1 {
		value := val(i)
		expected = insertToSlice(expected, i, value)
		binary = binary.insert(i, value)
	}
	return binary, expected
}
