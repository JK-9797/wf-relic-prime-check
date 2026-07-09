package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

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

		// // TODO: REMOVE DEBUG
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

		// // TODO: REMOVE DEBUG
		// window.IMShow(hsvImg)
		// // wait for key press
		// window.WaitKey(0)

		colorMask := gocv.NewMat()
		defer colorMask.Close()

		// color bounds are BGR but in terms of the hsv-colored image
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

		// subdivide image into reward choices
		// TODO: how to detect <4 players? line detection -> rectangle detection
		rewardsCount := 4
		rewards := make([]gocv.Mat, rewardsCount)
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
		rewardByteSlices := make([][]byte, len(rewards))
		for i := range rewardByteSlices {
			rewardBytes, _ := gocv.IMEncode(".png", rewards[i])
			rewardByteSlices[i] = rewardBytes.GetBytes()
		}

		// get string for each reward choice
		tesseractClient := gosseract.NewClient()
		defer tesseractClient.Close()

		rewardStrings := make([]string, len(rewards))
		for i := range rewardByteSlices {
			tesseractClient.SetImageFromBytes(rewardByteSlices[i])
			rewardStrings[i], err = tesseractClient.Text()
			if err != nil {
				fmt.Printf("couldn't read text: %v", err)
				return
			}

			// format string for market request
			rewardStrings[i] = strings.ToLower(rewardStrings[i])
			rewardStrings[i] = strings.ReplaceAll(rewardStrings[i], "\n", " ")
			rewardStrings[i] = strings.ReplaceAll(rewardStrings[i], " ", "_")

			// // TODO: REMOVE DEBUG
			// fmt.Printf("----- REWARD %v -----\n", i)
			// fmt.Println(rewardStrings[i])
		}

		fmt.Printf("\n\n")

		// market request
		wfmClient := &http.Client{}

		// TODO: can't find a way to batch multiply queries into 1 request
		for i := range rewardStrings {
			// skip formas
			if rewardStrings[i] == "forma_blueprint" {
				fmt.Printf("%v: 0\n", rewardStrings[i])
				continue
			}

			reqURL := "https://api.warframe.market/v2/orders/item/" + rewardStrings[i] + "/top"

			var req *http.Request
			req, err = http.NewRequest("GET", reqURL, nil)
			if err != nil {
				fmt.Printf("couldn't create request: %v", err)
				return
			}

			req.Header.Add("lang", "en")
			req.Header.Add("platform", "pc")
			req.Header.Add("crossplay", "true")

			var resp *http.Response
			resp, err = wfmClient.Do(req)
			if err != nil {
				fmt.Printf("couldn't make request: %v", err)
				return
			}
			defer resp.Body.Close()

			// parse response
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("couldn't read response body: %v", err)
				return
			}

			var topOrdersResponse TopOrdersResponse
			err = json.Unmarshal(respBody, &topOrdersResponse)
			if err != nil {
				fmt.Printf("couldn't unmarshal response body: %v", err)
				return
			}

			priceSum := 0
			for _, sellOrder := range topOrdersResponse.TopOrders.Sells {
				priceSum += sellOrder.Platinum
			}

			priceAvg := float64(priceSum) / float64(len(topOrdersResponse.TopOrders.Sells))

			fmt.Printf("%v: %v\n", rewardStrings[i], priceAvg)

			// // TODO: REMOVE DEBUG
			// fmt.Printf("\n%v: %v", rewardStrings[i], topOrdersResponse.TopOrders.Sells)
		}

		// wait for key press to refresh
		window.WaitKey(0)
	}
}

type OrdersResponse struct {
	Version string  `json:"apiVersion"`
	Orders  []Order `json:"data"`
	Error   string  `json:"error"`
}

// "online in game" & visible only?
type TopOrdersResponse struct {
	Version   string    `json:"apiVersion"`
	TopOrders TopOrders `json:"data"`
	Error     string    `json:"error"`
}

type Order struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Platinum int    `json:"platinum"`
	Visible  bool   `json:"visible"`
}

type TopOrders struct {
	Sells []Order `json:"sell"`
	Buys  []Order `json:"buy"`
}
