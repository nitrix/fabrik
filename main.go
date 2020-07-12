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

	// Target.
	target := sdl.FPoint{}

	// Arm and bones.
	arm := Arm{}
	arm.anchor = sdl.FPoint{X: 400, Y: 600}
	arm.bones = make([]*Bone, 0)
	for i := 0; i < 5; i++ {
		arm.bones = append(arm.bones, &Bone{
			head: sdl.FPoint{X: 400, Y: 600},
			tail: sdl.FPoint{X: 400, Y: 600},
			width: 10,
			length: 100,
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
		_ = drawFrame(renderer, arm, target)
		renderer.Present()
	}
}

func drawFrame(renderer *sdl.Renderer, arm Arm, target sdl.FPoint) error {
	_ = updateArm(arm, target)
	drawArm(renderer, arm)

	return nil
}

func drawArm(renderer *sdl.Renderer, arm Arm) {
	for _, bone := range arm.bones {
		_ = drawBone(renderer, *bone)
	}
}

func updateArm(arm Arm, target sdl.FPoint) error {
	// Update forward
	for k := range arm.bones {
		// First bone reaches the target, the rest follows.
		if k == 0 {
			boneReach(&arm.bones[k].head, &arm.bones[k].tail, arm.bones[k].length, target)
		} else {
			boneReach(&arm.bones[k].head, &arm.bones[k].tail, arm.bones[k].length, arm.bones[k-1].tail)
		}
	}

	// Update backward
	for k := len(arm.bones) - 1; k >= 0; k-- {
		if k == len(arm.bones) - 1 {
			boneReach(&arm.bones[k].tail, &arm.bones[k].head, arm.bones[k].length, arm.anchor)
		} else {
			boneReach(&arm.bones[k].tail, &arm.bones[k].head, arm.bones[k].length, arm.bones[k+1].head)
		}
	}

	return nil
}

func drawBone(renderer *sdl.Renderer, bone Bone) error {
	_ = renderer.SetDrawColor(255, 0, 0, 255)
	_ = renderer.DrawLine(int32(bone.head.X), int32(bone.head.Y), int32(bone.tail.X), int32(bone.tail.Y))

	N := bone.width
	x1 := float64(bone.head.X)
	y1 := float64(bone.head.Y)
	x2 := float64(bone.tail.X)
	y2 := float64(bone.tail.Y)
	dx := x1-x2
	dy := y1-y2
	dist := math.Sqrt(dx*dx + dy*dy)
	dx /= dist
	dy /= dist
	x3 := x1 + (N/2)*dy
	y3 := y1 - (N/2)*dx
	x4 := x2 + (N/2)*dy
	y4 := y2 - (N/2)*dx
	x5 := x2 - (N/2)*dy
	y5 := y2 + (N/2)*dx
	x6 := x1 - (N/2)*dy
	y6 := y1 + (N/2)*dx

	_ = renderer.SetDrawColor(255, 255, 255, 255)
	_ = renderer.DrawLineF(float32(x3), float32(y3), float32(x4), float32(y4))
	_ = renderer.DrawLineF(float32(x5), float32(y5), float32(x6), float32(y6))
	_ = renderer.DrawLineF(float32(x3), float32(y3), float32(x6), float32(y6))
	_ = renderer.DrawLineF(float32(x4), float32(y4), float32(x5), float32(y5))

	return nil
}

func boneReach(boneHead *sdl.FPoint, boneTail *sdl.FPoint, boneLength float64, target sdl.FPoint) {
	// Given a bone and a target, set the bone's head to the target.
	*boneHead = target

	// Notice that this stretches the bone. We fix that by sliding the tail along the new bone.

	// Calculate the stretched distance.
	sdx := boneTail.X - target.X
	sdy := boneTail.Y - target.Y
	distance := math.Sqrt(float64(sdx)*float64(sdx) + float64(sdy)*float64(sdy))
	scale := float32(boneLength / distance)

	// Scale the new tail based on the distance from the target.
	*boneTail = sdl.FPoint{
		X: target.X + sdx * scale,
		Y: target.Y + sdy * scale,
	}
}
