package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

// 1920*1080
// TODO: variable window dimensions
var startX = 476
var startY = 408
var rewardSizeX = 242
var rewardSizeY = 54

func DetectRewardBoxes(img gocv.Mat) (count int, err error) {
	// detect number of reward boxes using by finding rectangles
	matCanny := gocv.NewMat()
	matLines := gocv.NewMat()

	window := gocv.NewWindow("detected lines")
	defer window.Close()

	matGray := gocv.NewMat()
	matGray = img
	err = gocv.CvtColor(img, &matGray, gocv.ColorBGRToGray)

	window.IMShow(matGray)
	// wait for key press
	window.WaitKey(0)

	matBlur := gocv.NewMat()
	// gets highlighted reward
	// err = gocv.GaussianBlur(img, &matGray, image.Point{5, 5}, 1.4, 0, gocv.BorderDefault)
	// TODO: REMOVE testing
	// try drawing giant lines with hough and then re-run contours?
	err = gocv.GaussianBlur(matGray, &matBlur, image.Point{5, 5}, 0.0002, 0, gocv.BorderDefault)

	// matBlur = img

	window.IMShow(matBlur)
	// wait for key press
	window.WaitKey(0)

	// threshold filtering
	matThresh := gocv.NewMat()
	err = gocv.AdaptiveThreshold(matGray, &matThresh, 255, gocv.AdaptiveThresholdMean, gocv.ThresholdTrunc, 3, 2)

	window.IMShow(matThresh)
	// wait for key press
	window.WaitKey(0)

	// find edges
	// gocv.Canny(matBlur, &matCanny, 50, 200)
	gocv.Canny(matThresh, &matCanny, 5, 50)

	window.IMShow(matCanny)
	// wait for key press
	window.WaitKey(0)

	// contour algorithm test
	contours := gocv.FindContours(matCanny, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	contourImg := gocv.NewMat()
	contourImg = matBlur
	err = gocv.CvtColor(matBlur, &contourImg, gocv.ColorGrayToBGR)

	gocv.DrawContours(&contourImg, contours, -1, color.RGBA{0, 255, 200, 50}, int(gocv.Filled))

	window.IMShow(contourImg)
	// wait for key press
	window.WaitKey(0)

	// hough algorithm test
	gocv.HoughLinesP(matCanny, &matLines, 1, math.Pi/180, 80)

	fmt.Println(matLines.Cols())
	fmt.Println(matLines.Rows())
	for i := 0; i < matLines.Rows(); i++ {
		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&matBlur, pt1, pt2, color.RGBA{0, 255, 0, 50}, 5)
	}

	for {
		window.IMShow(matBlur)
		if window.WaitKey(10) >= 0 {
			break
		}
	}

	return
}

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
