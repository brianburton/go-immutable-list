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
