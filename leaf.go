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
	if len(this.contents) < maxPerNode {
		return leafNode{contents: appendObject(this.contents, value)}, nil
	} else {
		first, second := splitAppendObject(this.contents, value)
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

func appendObject(from []Object, extra Object) []Object {
	oldLen := len(from)
	newObjects := make([]Object, oldLen+1)
	for i, v := range from {
		newObjects[i] = v
	}
	newObjects[oldLen] = extra
	return newObjects
}

func splitAppendObject(from []Object, extra Object) ([]Object, []Object) {
	oldLen := len(from)
	newLen := oldLen + 1
	firstLen := newLen - (oldLen / 2)
	secondLen := newLen - firstLen

	first := make([]Object, firstLen)
	second := make([]Object, secondLen)
	for i, v := range from {
		if i < firstLen {
			first[i] = v
		} else {
			second[i-firstLen] = v
		}
	}
	second[secondLen-1] = extra
	return first, second
}
