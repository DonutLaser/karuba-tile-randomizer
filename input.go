package main

type Input struct {
	Space bool
	R     bool
}

func (input *Input) Clear() {
	input.Space = false
	input.R = false
}
