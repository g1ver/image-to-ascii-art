package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type AsciiBitmap []string

// These values represent increasing levels of gray
const gValueCharacters string = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please select a file")
		return
	}

	imgFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("File error: ", err)
		return
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Image decoding error: ", err)
		return
	}

	ascii := imgToAscii(img)

	txtFile, err := os.Create(os.Args[1] + ".txt")
	if err != nil {
		fmt.Println("File error: ", err)
		return
	}
	defer txtFile.Close()

	err = writeAsciiToFile(ascii, txtFile)
	if err != nil {
		fmt.Println("Write error: ", err)
		return
	}
}

// Convert image to grayscale then map each value to an ascii characer
func imgToAscii(img image.Image) AsciiBitmap {
	var ascii AsciiBitmap
	for y := 0; y < img.Bounds().Max.Y; y++ {
		var vertical string
		for x := 0; x < img.Bounds().Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y))
			char := grayscaleToCharacter(c)
			vertical += char
		}
		ascii = append(ascii, vertical)
	}
	return ascii
}

func grayscaleToCharacter(c color.Color) string {
	y, _, _, _ := c.RGBA()
	charIndex := intMap(int(y), 0, 65535, 0, len(gValueCharacters)-1)
	return string(gValueCharacters[charIndex])
}

func intMap(y, minIn, maxIn, minOut, maxOut int) int {
	return (y-minIn)*(maxOut-minOut)/(maxIn-minOut) + minOut
}

func writeAsciiToFile(ab AsciiBitmap, f *os.File) error {
	for _, s := range ab {
		_, err := f.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
