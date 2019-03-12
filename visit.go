package immutableList

// receive a value and return true to terminate loop or false to continue loop
type VisitProc func(index int, value Object) bool

type Visitable interface {
	Visit(processor VisitProc)
}

type ReduceProc func(acc Object, val Object) Object

func Reduce(source Visitable, initialValue Object, proc ReduceProc) Object {
	sum := initialValue
	source.Visit(func(_ int, value Object) bool {
		sum = proc(sum, value)
		return false
	})
	return sum
}
