package imageprocessing

import (
	"image"

	"gocv.io/x/gocv"
)

// 1920*1080
// TODO: variable window dimensions
var startX = 476
var startY = 408
var rewardSizeX = 242
var rewardSizeY = 54

func GetRewardBoxes(img gocv.Mat) (rewards []gocv.Mat, err error) {
	// subdivide image into reward choices
	// TODO: how to detect <4 players? line detection -> rectangle detection
	rewardsCount := 4
	rewards = make([]gocv.Mat, rewardsCount)

	for i := range rewards {
		rewardStartX := startX + i*rewardSizeX
		rewards[i] = img.Region(image.Rect(rewardStartX, startY, rewardStartX+rewardSizeX, startY+rewardSizeY))
	}

	return
}
