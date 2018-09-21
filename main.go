package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func splitPixels(img *image.RGBA) {
	arrSize := (len(img.Pix) / 4)
	a := make([][]uint8, arrSize)

	// Helps us keep track of grouping RGBA values together
	pixelIndex := 0
	// Helps us maintain RGBA sequence R = 0, G = 1, B = 2, A = 3
	position := 0

	fmt.Println("Image Pixel length")
	fmt.Println(len(img.Pix))
	fmt.Println(arrSize)
	for i := 0; i < len(img.Pix); i++ {
		// If we are at the start of a new pixel then we need to create a slice
		// to store our RGBA values
		if position == 0 {
			// Create a slice for 4 elements
			a[pixelIndex] = make([]uint8, 4)
		}
		// fmt.Printf("Iteration: %d, Val: %d \n", i, img.Pix[i])
		// Lets add the current value to the appropriate slice index and for the correct RGBA position
		a[pixelIndex][position] = img.Pix[i]

		// Keep track of our RGBA position
		position++

		// We have one large array of uint8 values that are in the order of
		// R, G, B, A ... since we are starting on a zero based index we have
		// R = 0, G = 1, B = 2, A = 3 This means that we can modulus on (i + 1) % 4 and
		// be able to easily chunk up our RGBA values for each pixel.
		// This will allow us to easily export for distribution analysis
		if i > 0 && (i+1)%4 == 0 {
			// If we evaluate this to true then we are currently on the 'A' of our R, G, B, A values
			// This means that we can wrap up this current array and move along

			// Increment to our next slice position
			pixelIndex++
			// Reset our position back to 0 ... or the 'R'
			position = 0
		}
	}
	writePixels(a)
	// fmt.Println(a)
}

func writePixels(values [][]uint8) {
	file, _ := os.Create("result.csv")
	//checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range values {
		data := make([]string, 4)
		data[0] = fmt.Sprint(value[0])
		data[1] = fmt.Sprint(value[1])
		data[2] = fmt.Sprint(value[2])
		data[3] = fmt.Sprint(value[3])
		writer.Write(data)
		//checkError("Cannot write to file", err)
	}
}

func main() {
	// Capture each displays.
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
	}

	var all image.Rectangle = image.Rect(0, 0, 0, 0)

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		all = bounds.Union(all)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		// fmt.Println(img)
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		splitPixels(img)
		// save(img, fileName)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}

	// Capture all desktop region into an image.
	// fmt.Printf("%v\n", all)
	// img, err := screenshot.Capture(all.Min.X, all.Min.Y, all.Dx(), all.Dy())
	// if err != nil {
	// 	panic(err)
	// }
	// save(img, "all.png")
}
