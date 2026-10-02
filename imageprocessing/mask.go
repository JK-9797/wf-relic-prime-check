package imageprocessing

import (
	"fmt"

	"gocv.io/x/gocv"
)

func MakeColorMask(img gocv.Mat) (colorMask gocv.Mat, err error) {
	// build text color mask
	// HSV images are better for color differentiation
	hsvImg := gocv.NewMat()
	defer hsvImg.Close()
	err = gocv.CvtColor(img, &hsvImg, gocv.ColorRGBToHSV)
	if err != nil {
		fmt.Printf("couldn't convert image to hsv: %v", err)
		return
	}

	// // TODO: REMOVE DEBUG
	// window.IMShow(hsvImg)
	// // wait for key press
	// window.WaitKey(0)

	colorMask = gocv.NewMat()

	// color bounds are BGR but in terms of the HSV-colored image
	// why?
	err = gocv.InRangeWithScalar(hsvImg, gocv.NewScalar(90.0, 110.0, 180.0, 0.0), gocv.NewScalar(97.0, 118.0, 190.0, 0.0), &colorMask)
	if err != nil {
		fmt.Printf("couldn't create color mask: %v", err)
		return
	}

	// // TODO: REMOVE DEBUG
	// window.IMShow(colorMask)
	// // wait for key press
	// window.WaitKey(0)

	// don't actually care about the color, so we can just continue processing on the mask itself
	// filteredImg := gocv.NewMat()
	// defer filteredImg.Close()
	// err = gocv.BitwiseAnd(img, colorMask, &filteredImg)
	// if err != nil {
	// 	fmt.Printf("couldn't apply mask: %v", err)
	// 	return
	// }

	// // TODO: REMOVE DEBUG
	// window.IMShow(filteredImg)
	// // wait for key press
	// window.WaitKey(0)

	return
}
