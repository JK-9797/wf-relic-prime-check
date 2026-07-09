package main

import (
	"fmt"
	"image"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

func main() {
	// TODO:
	var err error

	// TODO: input file path
	// filePath := "/mnt/c/Users/jkwok/Documents/misc/warframe/sc_test.png"
	filePath := "/mnt/c/Users/jkwok/Documents/misc/warframe/reward_sample.png"
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("couldnt read image")
		return
	}
	defer img.Close()

	window := gocv.NewWindow("Hello")
	defer window.Close()

	for {
		// open image
		img = gocv.IMRead(filePath, gocv.IMReadColor)
		if img.Empty() {
			fmt.Printf("couldnt read image")
			return
		}

		// subdivide image into reward choices
		// TODO: how to detect <4 players? line detection -> rectangle detection
		rewards := make([]gocv.Mat, 4)
		// 1920*1080
		startX := 476
		startY := 220
		rewardSizeX := 240
		rewardSizeY := 242
		for i := range rewards {
			rewardStartX := startX + i*rewardSizeX
			rewards[i] = img.Region(image.Rect(rewardStartX, startY, rewardStartX+rewardSizeX, startY+rewardSizeY))
		}

		// filter images to ease text processing
		// TODO: should it just filter the entire original image?

		// convert mat to []byte
		rewardByteSlices := make([][]byte, 4)
		for i := range rewardByteSlices {
			rewardBytes, _ := gocv.IMEncode(".png", rewards[i])
			rewardByteSlices[i] = rewardBytes.GetBytes()
		}

		// get string for each reward choice
		tesseractClient := gosseract.NewClient()
		defer tesseractClient.Close()

		rewardStrings := make([]string, 4)
		for i := range rewardByteSlices {
			tesseractClient.SetImageFromBytes(rewardByteSlices[i])
			rewardStrings[i], err = tesseractClient.Text()
			if err != nil {
				fmt.Printf("couldn't read text")
				return
			}

			// TODO: REMOVE
			fmt.Printf("----- REWARD %v -----\n", i)
			fmt.Println(rewardStrings[i])
		}

		// market request

		// parse response

		// TODO: testing edited images, REMOVE
		window.IMShow(rewards[3])
		// constant refresh
		window.WaitKey(1)
		// wait for key press
		// window.WaitKey(0)
	}
}
