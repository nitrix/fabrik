package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatalln("Unable to initialize SDL:", err)
	}
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600, 0)
	if err != nil {
		log.Fatalln("Unable to create SDL window:", err)
	}

	defer func() {
		_ = window.Destroy()
	}()

	defer func() {
		_ = renderer.Destroy()
	}()

	window.SetTitle("Learn IK")
	window.Show()

	target := sdl.FPoint{}
	chain := make([]*Bone, 0)

	for i := 0; i < 30; i++ {
		chain = append(chain, &Bone{
			head: sdl.FPoint{X: 400, Y: 300},
			tail: sdl.FPoint{X: 400, Y: 300},
			width: 10,
			length: 30,
		})
	}

	mainLoop:
	for {
		event := sdl.WaitEvent()

		switch e := event.(type) {
		case *sdl.MouseMotionEvent:
			target.X = float32(e.X)
			target.Y = float32(e.Y)
		case *sdl.QuitEvent:
			break mainLoop
		}

		_ = renderer.SetDrawColor(0, 0, 0, 0)
		_ = renderer.Clear()
		_ = renderer.SetDrawColor(255, 255, 255, 0)
		_ = drawFrame(renderer, chain, target)
		renderer.Present()
	}
}

func drawFrame(renderer *sdl.Renderer, chain []*Bone, target sdl.FPoint) error {
	_ = renderer.DrawPoint(400, 300)

	for k, bone := range chain { // FIXME: Irrelevant for now, single bone, not iterative.
		// First bone reaches the target.
		if k == 0 {
			boneReach(chain[0], target)
		} else {
			boneReach(bone, chain[k-1].tail)
		}

		_ = drawBone(renderer, *bone)
	}

	return nil
}

func drawBone(renderer *sdl.Renderer, bone Bone) error {
	return renderer.DrawLine(int32(bone.head.X), int32(bone.head.Y), int32(bone.tail.X), int32(bone.tail.Y))
}

func boneReach(bone *Bone, target sdl.FPoint) {
	// Given a bone and a target, set the bone's head to the target.
	bone.head = target

	// Notice that this stretches the bone. We fix that by sliding the tail along the new bone.

	// Calculate the stretched distance.
	sdx := bone.tail.X - target.X
	sdy := bone.tail.Y - target.Y
	distance := math.Sqrt(float64(sdx)*float64(sdx) + float64(sdy)*float64(sdy))
	scale := float32(bone.length / distance)

	// Scale the new tail based on the distance from the target.
	bone.tail = sdl.FPoint{
		X: target.X + sdx * scale,
		Y: target.Y + sdy * scale,
	}
}