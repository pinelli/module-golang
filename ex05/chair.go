package main

type Chair struct {
	id   int
	free chan bool
}

func NewChair(id int) Chair {
	free := make(chan bool, 1)
	free <- true
	return Chair{id, free}
}

func (this Chair) take() bool {
	select {
	case <-this.free:
		return true
	default:
		return false
	}
}

func (this Chair) release() {
	this.free <- true
}
