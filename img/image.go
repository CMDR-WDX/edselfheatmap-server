package img

import (
	"EDSelfHeatmap/data"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strconv"
)

var imageInstance *image.RGBA

// every pixel has 10 "actual" pixels to avoid blur
var initialState map[data.IntPoint]int
var offset data.IntPoint

func InitFromEnv(initial map[data.IntPoint]int) {
	//LOWER_X = 10000
	//LOWER_Y = 40000
	//UPPER_X = 13380
	//UPPER_Y = 42000
	lowerX, err := strconv.Atoi(os.Getenv("LOWER_X"))
	if err != nil {
		log.Fatalln(err)
	}
	lowerY, err := strconv.Atoi(os.Getenv("LOWER_Y"))
	if err != nil {
		log.Fatalln(err)
	}
	upperX, err := strconv.Atoi(os.Getenv("UPPER_X"))
	if err != nil {
		log.Fatalln(err)
	}
	upperY, err := strconv.Atoi(os.Getenv("UPPER_Y"))
	if err != nil {
		log.Fatalln(err)
	}
	offset = data.IntPoint{
		X: lowerX / 10,
		Y: lowerY / 10,
	}
	initialState = initial

	// Dividing by 10 and then multiplying by 10 IS NOT A BUG. Every px gets 10 real px to avoid blur
	Init(lowerX/10, upperX/10, lowerY/10, upperY/10)

}

func Init(minX int, maxX int, minY int, maxY int) {

	deltaPoint := image.Point{
		X: (maxX - minX) * 10,
		Y: (maxY - minY) * 10,
	}
	imageInstance = image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: deltaPoint,
	})
	draw.Draw(imageInstance, imageInstance.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)

	for val, _ := range initialState {
		setPixelVisitedCount(initialState[val], val.X, val.Y)
	}
}

func setPixelVisitedCount(count int, x int, y int) {
	// if out of bounds, ignore
	x -= offset.X
	y -= offset.Y
	if x < 0 || y < 0 {
		return
	}
	if x > imageInstance.Bounds().Dx() || y > imageInstance.Bounds().Dy() {
		return
	}

	colorValue := GetColorForCount(count)
	// Invert y
	//minY := imageInstance.Bounds().Min.Y
	//maxY := imageInstance.Bounds().Max.Y

	//deltaToMax := maxY - y
	y = imageInstance.Bounds().Dy()/10 - y

	log.Printf("Received count %d for (%d, %d)\n", count, x, y)
	rectangleForPixel := image.Rect(x*10, y*10, (x+1)*10, (y+1)*10)
	draw.Draw(imageInstance, rectangleForPixel, &image.Uniform{C: colorValue}, imageInstance.Bounds().Min, draw.Src)

}

func Increment(x int, y int) {
	point := data.IntPoint{X: x, Y: y}
	value, ok := initialState[point]
	if !ok {
		value = 0
	}
	value++
	initialState[point] = value
	setPixelVisitedCount(value, x, y)
}

func MakeImage() *bytes.Buffer {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, imageInstance)
	if err != nil {
		println(err)
		return nil
	}
	return buf
}
