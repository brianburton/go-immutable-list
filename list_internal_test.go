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
}

func TestBuilder(t *testing.T) {
	builder := CreateBuilder()
	validateList(t, builder.Build(), 0)
	for i := 1; i <= maxPerNode; i++ {
		builder.Add(val(i))
		validateSize(t, builder.Size(), i)
		validateList(t, builder.Build(), i)
	}
	for i := maxPerNode + 1; i <= 200; i++ {
		builder.Add(val(i))
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

func TestSet(t *testing.T) {
	list := createListForTest(90, 1090)
	for i := 0; i < list.Size(); i++ {
		list = list.Set(i, val(i+1))
	}
	list = list.Set(1001, val(1002))
	validateList(t, list, 1002)
}

func createListForTest(firstValue int, lastValue int) List {
	builder := CreateBuilder()
	for i := firstValue; i <= lastValue; i++ {
		builder.Add(val(i))
	}
	return builder.Build()
}

func validateList(t *testing.T, list List, size int) {
	if list.Size() != size {
		t.Error(fmt.Sprintf("expected size %d but got %v", size, list.Size()))
	}
	ei := 0
	list.Visit(0, list.Size(), func(index int, obj Object) {
		if index != ei || obj.(string) != val(index+1) {
			t.Error(fmt.Sprintf("visitor expected %v/%s but got %v/%s", ei, val(ei+1), index, obj))
		}
		ei += 1
	})
	if ei != list.Size() {
		t.Error(fmt.Sprintf("expected count %d but got %v", list.Size(), ei))
	}
	for i := 0; i < size; i++ {
		expected := val(i + 1)
		actual := list.Get(i)
		if actual != expected {
			t.Error(fmt.Sprintf("get expected %v/%s but got %v/%s", i, expected, i, actual))
		}
	}
}

func validateSize(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Error(fmt.Sprintf("expected size %d but got %v", expected, actual))
	}
}

func val(index int) string {
	return fmt.Sprintf("%v", index)
}
