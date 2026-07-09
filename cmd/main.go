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

		// // TODO: REMOVE
		// window.IMShow(img)
		// // wait for key press
		// window.WaitKey(0)

		// filter images to ease text processing

		// build text color mask
		// hsv images are better for color differentiation
		hsvImg := gocv.NewMat()
		defer hsvImg.Close()
		err = gocv.CvtColor(img, &hsvImg, gocv.ColorRGBToHSV)
		if err != nil {
			fmt.Printf("couldn't convert image to hsv: %v", err)
			return
		}

		// TODO: REMOVE
		window.IMShow(hsvImg)
		// wait for key press
		window.WaitKey(0)

		colorMask := gocv.NewMat()
		defer colorMask.Close()

		// color bounds are BGR but in terms of the hsv-colored image
		err = gocv.InRangeWithScalar(hsvImg, gocv.NewScalar(90.0, 110.0, 180.0, 0.0), gocv.NewScalar(97.0, 118.0, 190.0, 0.0), &colorMask)
		if err != nil {
			fmt.Printf("couldn't create color mask: %v", err)
			return
		}

		// TODO: REMOVE
		window.IMShow(colorMask)
		// wait for key press
		window.WaitKey(0)

		// don't actually care about the color, so we can just continue processing on the mask itself

		// filteredImg := gocv.NewMat()
		// defer filteredImg.Close()
		// err = gocv.BitwiseAnd(img, colorMask, &filteredImg)
		// if err != nil {
		// 	fmt.Printf("couldn't apply mask: %v", err)
		// 	return
		// }

		// // TODO: REMOVE
		// window.IMShow(filteredImg)
		// // wait for key press
		// window.WaitKey(0)

		// subdivide image into reward choices
		// TODO: how to detect <4 players? line detection -> rectangle detection
		rewards := make([]gocv.Mat, 4)
		// 1920*1080
		startX := 476
		startY := 408
		rewardSizeX := 242
		rewardSizeY := 54
		for i := range rewards {
			rewardStartX := startX + i*rewardSizeX
			rewards[i] = colorMask.Region(image.Rect(rewardStartX, startY, rewardStartX+rewardSizeX, startY+rewardSizeY))
		}

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
				fmt.Printf("couldn't read text: %v", err)
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
		// wait for key press
		window.WaitKey(0)
	}
}
