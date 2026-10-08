package imageprocessing

import (
	"cmp"
	"image"
	"image/color"
	"slices"

	"gocv.io/x/gocv"
)

const rewardsMax = 4

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

func FindRewards(img gocv.Mat) (rewardImages []gocv.Mat, err error) {
	// only need to scan a small region of the screen capture
	subImg := img.Region(image.Rect(startX, startY, startX+rewardsMax*rewardSizeX, startY+rewardSizeY))

	//TODO: remove debug
	// window := gocv.NewWindow("count test")
	// defer window.Close()
	// window.IMShow(subImg)
	// // wait for key press
	// window.WaitKey(0)

	// filter images to ease text processing
	colorMask, err := MakeColorMask(subImg)
	if err != nil {
		return
	}
	defer colorMask.Close()

	// dilate all letter edges to turn clusters of words into single contours
	dilateImg := gocv.NewMat()
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{11, 11})
	err = gocv.Dilate(colorMask, &dilateImg, kernel)
	if err != nil {
		return
	}

	contourImg := colorMask.Clone()
	contours := gocv.FindContours(dilateImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	err = gocv.DrawContours(&contourImg, contours, -1, color.RGBA{0, 255, 200, 50}, int(gocv.Filled))
	if err != nil {
		return
	}

	//TODO: remove debug
	// window.IMShow(contourImg)
	// // wait for key press
	// window.WaitKey(0)

	// define an image region for each contour
	boundingBoxes := make([]image.Rectangle, 0, rewardsMax)
	for i := range contours.ToPoints() {
		boundingBoxes = append(boundingBoxes, gocv.BoundingRect(contours.At(i)))
	}

	// sort rewards left to right
	slices.SortFunc(boundingBoxes, func(a image.Rectangle, b image.Rectangle) int { return cmp.Compare(a.Min.X, b.Min.X) })

	rewardImages = make([]gocv.Mat, len(boundingBoxes))
	for i, boundingBox := range boundingBoxes {
		rewardImages[i] = colorMask.Region(boundingBox)
	}

	return
}
