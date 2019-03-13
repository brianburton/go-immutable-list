package immutableList

type leafNode struct {
	contents []Object
}

func createLeafNode(objectBuffer []Object, count int) node {
	contents := make([]Object, count)
	for i := 0; i < count; i++ {
		contents[i] = objectBuffer[i]
	}
	return leafNode{contents}
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

func (this leafNode) height() int {
	return 1
}

func (this leafNode) maxCompleteHeight() int {
	if len(this.contents) >= minPerNode {
		return 1
	} else {
		return 0
	}
}

func (this leafNode) visitNodesOfHeight(targetHeight int, proc nodeProcessor) {
	if targetHeight == 1 {
		proc(this)
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
