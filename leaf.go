package immutableList

type leafNode struct {
	contents []Object
}

func createLeafNode(objectBuffer []Object, count int) node {
	contents := make([]Object, count)
	copy(contents, objectBuffer)
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

func (this leafNode) set(index int, value Object) node {
	newContents := make([]Object, len(this.contents))
	copy(newContents, this.contents)
	newContents[index] = value
	return leafNode{newContents}
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
	if this.isComplete() {
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
	copy(newObjects[0:], from[0:insertIndex])
	newObjects[insertIndex] = extra
	copy(newObjects[insertIndex+1:], from[insertIndex:])
	return newObjects
}

func splitInsertObject(insertIndex int, extra Object, from []Object) ([]Object, []Object) {
	newContents := insertObject(insertIndex, extra, from)
	newLen := len(newContents)
	secondLen := newLen / 2
	firstLen := newLen - secondLen
	first := make([]Object, firstLen)
	copy(first, newContents[0:firstLen])
	second := make([]Object, secondLen)
	copy(second, newContents[firstLen:])
	return first, second
}

func (this leafNode) isComplete() bool {
	return len(this.contents) >= minPerNode
}

func (this leafNode) mergeWith(other node) node {
	otherLeaf := other.(leafNode)
	myLen := len(this.contents)
	otherLen := len(otherLeaf.contents)
	newLen := myLen + otherLen
	newContents := make([]Object, newLen)
	copy(newContents[0:], this.contents)
	copy(newContents[myLen:], otherLeaf.contents)
	return leafNode{newContents}
}

func (this leafNode) delete(index int) node {
	oldLen := len(this.contents)
	if oldLen == 1 {
		return emptyNode{}
	} else {
		newContents := make([]Object, oldLen-1)
		copy(newContents[0:], this.contents[0:index])
		copy(newContents[index:], this.contents[index+1:])
		return leafNode{newContents}
	}
}
