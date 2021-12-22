package main

type KeyState uint8

const (
	KeyStateNone KeyState = iota
	KeyStateWentDown
	KeyStateDown
	KeyStateWentUp
)

type Input struct {
	Space bool
	R     KeyState
}

func (input *Input) Clear() {
	input.Space = false
}
