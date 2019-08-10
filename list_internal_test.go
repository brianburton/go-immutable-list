package immutableList

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAppend(t *testing.T) {
	list := Create()
	list = list.Append(val(1)).Append(val(2)).Append(val(3))
	validateList(t, list, 3)

	for i := 4; i <= 1024; i++ {
		list = list.Append(val(i))
	}
	validateList(t, list, 1024)
}

func TestInsert(t *testing.T) {
	list := Create()
	list = list.Insert(0, val(1)).Insert(1, val(512))

	for i := 2; i <= 256; i++ {
		list = list.Insert(i-1, val(i))
		list = list.Insert(i, val(513-i))
	}

	validateList(t, list, 512)
}

func TestInsert2(t *testing.T) {
	list := Create()
	expected := make([]Object, 0)
	nextValue := 1
	for list.Size() < 4096 {
		size := list.Size()
		increment := size / 3
		if increment < 3 {
			increment = 3
		}
		for index := 0; index <= size/2; index += increment {
			value := val(nextValue)
			nextValue += 1

			list = list.Insert(size-index, value)
			expected = insertToSlice(expected, size-index, value)
			list = list.Insert(index, value)
			expected = insertToSlice(expected, index, value)
			validateList3(t, list, expected)
			size = list.Size()
		}
	}
}

func TestSelect(t *testing.T) {
	list := Create()
	for i := 1; i <= 1024; i++ {
		list = list.Append(val(i))
	}
	list = list.Select(func(obj Object) bool {
		i, _ := strconv.ParseInt(obj.(string), 10, 64)
		return i <= 512
	})
	validateList(t, list, 512)
}

func TestSlice(t *testing.T) {
	list := Create()
	for i := 1; i <= 400; i++ {
		list = list.Append(val(i))
	}
	sliced := list.Slice(0, 123)
	for i, v := range sliced {
		if val(i+1) != v.(string) {
			t.Error(fmt.Sprintf("slice expected %v/%s but got %v/%s", i, val(i+1), i, v))
		}
	}
	sliced = list.Slice(0, 0)
	if len(sliced) != 0 {
		t.Error(fmt.Sprintf("slice expected empty but got size %d", len(sliced)))
	}
}

func TestDelete(t *testing.T) {
	list := createListForTest(1, 1024)
	for i := 1023; i >= 0; i-- {
		list = list.Delete(i)
		validateList(t, list, i)
	}
	list = createListForTest(1, 512)
	list = list.AppendList(createListForTest(1, 512))
	for i := 1; i <= 512; i++ {
		list = list.Delete(0)
	}
	validateList(t, list, 512)
	list = createListForTest(1, 800)
	list = list.AppendList(createListForTest(10000, 12000))
	list = list.AppendList(createListForTest(801, 2000))
	for i := 1000; i >= 1; i-- {
		list = list.Delete(800 + i)
	}
	for i := 1000; i >= 0; i-- {
		list = list.Delete(800 + i)
	}
	validateList(t, list, 2000)
}

func TestHead(t *testing.T) {
	list := createListForTest(1, 1234)
	for i := list.Size(); i >= 0; i-- {
		truncated := list.Head(i)
		validateList(t, truncated, i)
	}
	list = createListForTest(1, maxValuesPerLeaf)
	for i := list.Size(); i >= 0; i-- {
		list = list.Head(i)
		validateList(t, list, i)
	}
}

func TestTail(t *testing.T) {
	list := createListForTest(1, 1234)
	for i := list.Size(); i >= 0; i-- {
		truncated := list.Tail(i)
		validateList2(t, truncated, i+1, list.Size())
	}
	list = createListForTest(1, maxValuesPerLeaf)
	for i := 1; i <= maxValuesPerLeaf; i++ {
		list = list.Tail(1)
		validateList2(t, list, i+1, maxValuesPerLeaf)
	}
}

func TestSubList(t *testing.T) {
	list := createListForTest(1, 1122)
	offset := 0
	limit := list.Size()
	for offset < limit {
		sub := list.SubList(offset, limit)
		validateList2(t, sub, offset+1, limit)
		limit--
		sub = list.SubList(offset, limit)
		validateList2(t, sub, offset+1, limit)
		offset++
	}
	list = createListForTest(1, 257)
	for i := 0; i <= list.Size(); i++ {
		sub := list.SubList(i, list.Size())
		validateList2(t, sub, i+1, list.Size())
	}
	for i := 0; i <= list.Size(); i++ {
		sub := list.SubList(0, i)
		validateList2(t, sub, 1, i)
	}
}

func TestDeleteRange(t *testing.T) {
	a := createListForTest(1, 256)
	b := createListForTest(1000, 1002)
	c := createListForTest(257, 500)

	list := a.AppendList(b).AppendList(c)

	del := list.DeleteRange(a.Size(), a.Size()+b.Size())
	validateList(t, del, 500)

	del = list.DeleteRange(0, list.Size()-118)
	validateList2(t, del, 383, 500)

	del = list.DeleteRange(200, list.Size())
	validateList(t, del, 200)

	a = createListForTest(1, 10)
	b = createListForTest(1999, 2222)
	c = createListForTest(11, 888)
	list = a.AppendList(b).AppendList(c)
	del = list.DeleteRange(a.Size(), a.Size()+b.Size())
	validateList(t, del, 888)

	del = a.DeleteRange(0, a.Size())
	validateList(t, del, 0)

	del = a.DeleteRange(2, 2)
	validateList(t, del, a.Size())
}

func TestBuilder(t *testing.T) {
	builder := CreateBuilder()
	validateList(t, builder.Build(), 0)
	for i := 1; i <= maxValuesPerLeaf; i++ {
		builder.Add(val(i))
		validateSize(t, builder.Size(), i)
		validateList(t, builder.Build(), i)
	}
	for i := maxValuesPerLeaf + 1; i <= 200; i++ {
		builder.Add(val(i))
		validateSize(t, builder.Size(), i)
		validateList(t, builder.Build(), i)
	}
	validateSize(t, builder.Size(), 200)
	validateList(t, builder.Build(), 200)
	for i := 201; i <= 512; i++ {
		builder.Add(val(i))
	}
	validateSize(t, builder.Size(), 512)
	validateList(t, builder.Build(), 512)
	for i := 513; i <= 700; i++ {
		builder.Add(val(i))
	}
	validateSize(t, builder.Size(), 700)
	validateList(t, builder.Build(), 700)
}

func TestAppendList(t *testing.T) {
	for totalSize := 1; totalSize <= 100; totalSize++ {
		for firstSize := 0; firstSize <= totalSize; firstSize++ {
			first := createListForTest(1, firstSize)
			second := createListForTest(firstSize+1, totalSize)
			merged := first.AppendList(second)
			validateList(t, merged, totalSize)
		}
	}

	var merged List
	merged = createListForTest(1, 1000).AppendList(createListForTest(1001, 2000))
	validateList(t, merged, 2000)
	merged = createListForTest(1, 1000).AppendList(createListForTest(1001, 20000))
	validateList(t, merged, 20000)
	merged = createListForTest(1, 6000).AppendList(createListForTest(6001, 17758))
	validateList(t, merged, 17758)
	merged = createListForTest(1, 65).AppendList(createListForTest(66, 60000))
	validateList(t, merged, 60000)
}

func TestAppendList2(t *testing.T) {
	firstSize := 872
	first := createListForTestDirectly(1, firstSize)
	for secondSize := 0; secondSize <= firstSize; secondSize++ {
		totalSize := firstSize + secondSize
		second := createListForTest(firstSize+1, totalSize)
		merged := first.AppendList(second)
		validateList(t, merged, totalSize)
	}
}

func TestPrependList(t *testing.T) {
	secondSize := 900
	for firstSize := 0; firstSize <= secondSize; firstSize++ {
		totalSize := firstSize + secondSize
		first := createListForTestDirectly(1, firstSize)
		second := createListForTestReverseDirectly(firstSize+1, totalSize)
		merged := first.AppendList(second)
		validateList(t, merged, totalSize)
	}
}

func TestIterator(t *testing.T) {
	for length := 0; length <= 1024; length++ {
		actual := copyList(createListForTest(1, length))
		validateList(t, actual, length)
	}
}

func TestStackOps(t *testing.T) {
	stack := Create().Push(val(4)).Push(val(3)).Push(val(2)).Push(val(1))
	popped := Create()
	for !stack.IsEmpty() {
		var value Object
		value, stack = stack.Pop()
		popped = popped.Append(value)
	}
	validateList(t, popped, 4)

	stack = Create()
	for i := 500; i >= 1; i-- {
		stack = stack.Push(val(i))
	}
	popped = Create()
	for !stack.IsEmpty() {
		var value Object
		value, stack = stack.Pop()
		popped = popped.Append(value)
	}
	validateList(t, popped, 500)
}

func TestQueueOps(t *testing.T) {
	s := Create().Append(val(1)).Append(val(2)).Append(val(3)).Append(val(4))
	popped := Create()
	for !s.IsEmpty() {
		var value Object
		value, s = s.Pop()
		popped = popped.Append(value)
	}
	validateList(t, popped, 4)
}

func TestFirstLast(t *testing.T) {
	list := createListForTest(1, 1)
	validateValue(t, 1, list.GetFirst())
	validateValue(t, 1, list.GetLast())

	list = list.Append(val(2))
	validateValue(t, 1, list.GetFirst())
	validateValue(t, 2, list.GetLast())

	list = createListForTest(1, 500)
	validateValue(t, 1, list.GetFirst())
	validateValue(t, 500, list.GetLast())
}

func TestDeleteAll(t *testing.T) {
	for length := 16; length <= 512; length += maxValuesPerLeaf {
		increment := length / 11
		for index := 0; index < length; index += increment {
			if index < length {
				deleteAllImpl(t, length, index)
				deleteAllImpl(t, length, length-index)
			}
		}
	}
}

func deleteAllImpl(t *testing.T, length int, index int) {
	b := CreateBuilder()
	expected := make([]Object, length)
	for i := 0; i < length; i++ {
		value := val(i)
		b.Add(value)
		expected[i] = value
	}
	actual := b.Build()
	validateList3(t, actual, expected)
	for len(expected) > 0 {
		if index >= len(expected) {
			index = len(expected) - 1
		}
		actual = actual.Delete(index)
		expected = deleteFromSlice(expected, index)
		validateList3(t, actual, expected)
	}
}

func TestPopAll(t *testing.T) {
	length := 1
	for length <= 4096 {
		popAllImpl(t, length)
		length = length * maxValuesPerLeaf
	}
}

func popAllImpl(t *testing.T, length int) {
	b := CreateBuilder()
	expected := make([]Object, length)
	for i := 0; i < length; i++ {
		value := val(i)
		b.Add(value)
		expected[i] = value
	}
	actual := b.Build()
	validateList3(t, actual, expected)
	expectedValue := 0
	for len(expected) > 0 {
		var deletedValue Object
		deletedValue, actual = actual.Pop()
		expected = deleteFromSlice(expected, 0)
		validateValue(t, expectedValue, deletedValue)
		expectedValue += 1
		validateList3(t, actual, expected)
	}
}

func copyList(list List) List {
	answer := Create()
	for i := list.FwdIterate(); i.Next(); {
		answer = answer.Append(i.Get())
	}
	return answer
}

func TestSet(t *testing.T) {
	list := createListForTest(90, 1090)
	for i := 0; i < list.Size(); i++ {
		list = list.Set(i, val(i+1))
	}
	list = list.Set(1001, val(1002))
	validateList(t, list, 1002)
}

func TestInsertList(t *testing.T) {
	prefix := createListForTestInsertList(val(1), 0)
	middle := createListForTestInsertList(val(2), 0)
	//suffix := createListForTestInsertList(val(3), 0)

	inserted := prefix.InsertList(0, middle)
	validateInsertList(t, inserted, 0, 0, 0)

	middle = createListForTestInsertList(val(2), 100)
	inserted = prefix.InsertList(0, middle)
	validateInsertList(t, inserted, 0, 100, 0)

	prefix = createListForTestInsertList(val(1), 300)
	middle = createListForTestInsertList(val(2), 500)
	inserted = prefix.InsertList(300, middle)
	validateInsertList(t, inserted, 300, 500, 0)

	prefix = createListForTestInsertList(val(1), 3)
	suffix := createListForTestInsertList(val(3), maxValuesPerLeaf-1)
	middle = createListForTestInsertList(val(2), maxValuesPerLeaf-1)
	inserted = prefix.AppendList(suffix).InsertList(3, middle)
	validateInsertList(t, inserted, 3, maxValuesPerLeaf-1, maxValuesPerLeaf-1)

	prefix = createListForTestInsertList(val(1), 300)
	suffix = createListForTestInsertList(val(3), 300)
	middle = createListForTestInsertList(val(2), 500)
	inserted = prefix.AppendList(suffix).InsertList(300, middle)
	validateInsertList(t, inserted, 300, 500, 300)
}

func createListForTest(firstValue int, lastValue int) List {
	builder := CreateBuilder()
	for i := firstValue; i <= lastValue; i++ {
		builder.Add(val(i))
	}
	return builder.Build()
}

func createListForTestDirectly(firstValue int, lastValue int) List {
	answer := Create()
	for i := firstValue; i <= lastValue; i++ {
		answer = answer.Append(val(i))
	}
	return answer
}

func createListForTestInsertList(value Object, length int) List {
	answer := CreateBuilder()
	for i := 1; i <= length; i++ {
		answer.Add(value)
	}
	return answer.Build()
}

func createListForTestReverseDirectly(firstValue int, lastValue int) List {
	answer := Create()
	for i := lastValue; i >= firstValue; i-- {
		answer = answer.Insert(0, val(i))
	}
	return answer
}

func validateValue(t *testing.T, expected int, actual Object) {
	e := val(expected)
	if actual != e {
		t.Error(fmt.Sprintf("expected %v but got %v", expected, actual))
	}
}

func validateList(t *testing.T, list List, size int) {
	validateList2(t, list, 1, size)
}

func validateList2(t *testing.T, list List, first int, last int) {
	size := last - first + 1
	if list.Size() != size {
		t.Error(fmt.Sprintf("expected size %d but got %v", size, list.Size()))
	}
	ei := 0
	list.Visit(0, list.Size(), func(index int, obj Object) {
		if index != ei || obj.(string) != val(index+first) {
			t.Error(fmt.Sprintf("visitor expected %v/%s but got %v/%s", ei, val(ei+1), index, obj))
		}
		ei += 1
	})
	if ei != list.Size() {
		t.Error(fmt.Sprintf("expected count %d but got %v", list.Size(), ei))
	}
	for i := 0; i < size; i++ {
		expected := val(i + first)
		actual := list.Get(i)
		if actual != expected {
			t.Error(fmt.Sprintf("get expected %v/%s but got %v/%s", i, expected, i, actual))
		}
	}
	list.checkInvariants(func(message string) {
		t.Error(message)
	})
}

func validateList3(t *testing.T, list List, expected []Object) {
	size := len(expected)
	if list.Size() != size {
		t.Error(fmt.Sprintf("expected size %d but got %v", size, list.Size()))
	}
	ei := 0
	list.Visit(0, list.Size(), func(index int, obj Object) {
		if index != ei || obj.(string) != expected[index] {
			t.Error(fmt.Sprintf("visitor expected %v/%s but got %v/%s", ei, val(ei+1), index, obj))
		}
		ei += 1
	})
	if ei != list.Size() {
		t.Error(fmt.Sprintf("expected count %d but got %v", list.Size(), ei))
	}
	for i := 0; i < size; i++ {
		actual := list.Get(i)
		if actual != expected[i] {
			t.Error(fmt.Sprintf("get expected %v/%s but got %v/%s", i, expected[i], i, actual))
		}
	}
	list.checkInvariants(func(message string) {
		t.Error(message)
	})
}

func validateInsertList(t *testing.T, list List, prefixLength int, insertLength int, suffixLength int) {
	size := prefixLength + insertLength + suffixLength
	if list.Size() != size {
		t.Error(fmt.Sprintf("expected size %d but got %v", size, list.Size()))
		return
	}

	offset := 0
	validateRegion(t, "1", list, offset, prefixLength)

	offset += prefixLength
	validateRegion(t, "2", list, offset, offset+insertLength)

	offset += insertLength
	validateRegion(t, "3", list, offset, offset+suffixLength)
}

func validateSize(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Error(fmt.Sprintf("expected size %d but got %v", expected, actual))
	}
}

func val(index int) string {
	return fmt.Sprintf("%v", index)
}

func validateRegion(t *testing.T, expected Object, list List, offset int, limit int) {
	failed := false
	for i := offset; i < limit; i++ {
		actual := list.Get(i)
		if actual != expected {
			if !failed {
				failed = true
				t.Error(fmt.Sprintf("expected value %v at index %d but got %v", expected, i, actual))
			}
		}
	}
}

func insertToSlice(slice []Object, index int, value string) []Object {
	answer := make([]Object, len(slice)+1)
	copy(answer[0:], slice[0:index])
	answer[index] = value
	copy(answer[index+1:], slice[index:])
	return answer
}

func deleteFromSlice(slice []Object, index int) []Object {
	answer := make([]Object, len(slice)-1)
	copy(answer[0:], slice[0:index])
	copy(answer[index:], slice[index+1:])
	return answer
}
