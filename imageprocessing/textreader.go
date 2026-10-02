package imageprocessing

import (
	"fmt"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

func ReadTextBoxes(boxes []gocv.Mat, tesseractClient *gosseract.Client) (foundStrings []string, err error) {
	// convert mat to []byte
	boxByteSlices := make([][]byte, len(boxes))
	for i := range boxByteSlices {
		rewardBytes, _ := gocv.IMEncode(".png", boxes[i])
		boxByteSlices[i] = rewardBytes.GetBytes()
	}

	// run each box through text reader
	foundStrings = make([]string, len(boxes))
	for i := range boxByteSlices {
		tesseractClient.SetImageFromBytes(boxByteSlices[i])
		foundStrings[i], err = tesseractClient.Text()
		if err != nil {
			fmt.Printf("couldn't read text: %v", err)
			return
		}

		// format string for market request
		foundStrings[i] = strings.ToLower(foundStrings[i])
		foundStrings[i] = strings.ReplaceAll(foundStrings[i], "\n", " ")
		foundStrings[i] = strings.ReplaceAll(foundStrings[i], " ", "_")

		// // TODO: REMOVE DEBUG
		// fmt.Printf("----- REWARD %v -----\n", i)
		// fmt.Println(foundStrings[i])
	}

	return
}
