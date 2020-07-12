package main

import "github.com/veandco/go-sdl2/sdl"

type Arm struct {
	anchor sdl.FPoint
	bones []*Bone
}
