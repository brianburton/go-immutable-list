package immutableList

type leafNode struct {
	contents []Object
}

func (this leafNode) size() int {
	return len(this.contents)
}

func (this leafNode) get(index int) Object {
	if index >= 0 && index < len(this.contents) {
		return this.contents[index]
	} else {
		return nil
	}
}

func (this leafNode) append(value Object) (node, node) {
	return this.insert(len(this.contents), value)
}

func (this leafNode) insert(indexBefore int, value Object) (node, node) {
	currentSize := len(this.contents)
	if currentSize < maxPerNode {
		return leafNode{contents: insertObject(indexBefore, value, this.contents)}, nil
	} else {
		first, second := splitInsertObject(indexBefore, value, this.contents)
		return leafNode{contents: first}, leafNode{contents: second}
	}
}

func (this leafNode) forEach(proc Processor) {
	for _, value := range this.contents {
		proc(value)
	}
}

func (this leafNode) visit(start int, limit int, v Visitor) {
	for i := start; i < limit; i++ {
		v(i, this.contents[i])
	}
}

func insertObject(insertIndex int, extra Object, from []Object) []Object {
	newObjects := make([]Object, len(from)+1)

	index := 0
	insert := func(obj Object) {
		newObjects[index] = obj
		index++
	}
	for _, v := range from {
		if index == insertIndex {
			insert(extra)
		}
		insert(v)
	}
	if index == insertIndex {
		insert(extra)
	}
	return newObjects
}

func splitInsertObject(insertIndex int, extra Object, from []Object) ([]Object, []Object) {
	newLen := len(from) + 1
	secondLen := newLen / 2
	firstLen := newLen - secondLen
	first := make([]Object, firstLen)
	second := make([]Object, secondLen)

	index := 0
	insert := func(obj Object) {
		if index < firstLen {
			first[index] = obj
		} else {
			second[index-firstLen] = obj
		}
		index++
	}

	for _, value := range from {
		if index == insertIndex {
			insert(extra)
		}
		insert(value)
	}
	if index == insertIndex {
		insert(extra)
	}
	return first, second
}
