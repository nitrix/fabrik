package main

import "github.com/veandco/go-sdl2/sdl"

type Bone struct {
	head sdl.FPoint
	tail sdl.FPoint
	width float64
	length float64
}
