package main

type Input struct {
	Space bool
	Ctrl  bool
}

func (input *Input) Clear() {
	input.Space = false
}
