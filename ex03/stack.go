package stack

type Stack []int

func New() *Stack {
	return new(Stack)
}

func (this *Stack) Push(e int) {
	*this = append(*this, e)
}

func (this *Stack) Pop() (res int, sucess bool) {
	l := len(*this)

	if l < 1 {
		res, sucess = 0, false
	} else {
		res = (*this)[l-1]
		sucess = true
		*this = (*this)[:l-1]
	}

	return
}
