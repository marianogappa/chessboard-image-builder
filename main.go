package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

var pieceToPath = map[byte]string{
	'q': "imgs/bQ.png",
	'k': "imgs/bK.png",
	'b': "imgs/bB.png",
	'n': "imgs/bN.png",
	'r': "imgs/bR.png",
	'p': "imgs/bP.png",
	'Q': "imgs/wQ.png",
	'K': "imgs/wK.png",
	'B': "imgs/wB.png",
	'N': "imgs/wN.png",
	'R': "imgs/wR.png",
	'P': "imgs/wP.png",
	'*': "imgs/background.png",
}

var flagDebug = flag.Bool("debug", false, "Debug.")
var flagFen = flag.String("fen", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", "Board in FEN Notation.")

func main() {
	flag.Parse()
	var (
		fen   = *flagFen
		lines = strings.Split(fen, "/")
		grid  [8][8]byte

		// Set these parameters manually while running with -debug flag to match the grid
		boardMinX  = 40
		boardMinY  = 15
		boardMaxX  = 783
		boardMaxY  = 760
		cellWidth  = 93
		cellHeight = 93

		finalImageWidth = 800
		pathToImage     = make(map[string]image.Image, 13)
	)

	// Read images into memory
	for _, path := range pieceToPath {
		var resizeWidth int
		if path != "imgs/background.png" {
			resizeWidth = int(math.Floor(float64(cellWidth) * 0.8))
		}

		img, err := readImage(path, resizeWidth)
		if err != nil {
			log.Fatalf("Bug! This image cannot be found or read: %v\n", path)
		}
		pathToImage[path] = img
	}

	// Construct grid of pieces
	for y, line := range lines {
		x := 0
		for i := range line {
			switch line[i] {
			case '1', '2', '3', '4', '5', '6', '7', '8':
				n, err := strconv.Atoi(string(line[i]))
				if err != nil {
					log.Fatalf("Your FEN sucks: %v; read invalid number: %v\n", fen, line[i])
				}
				x += n
			case 'q', 'k', 'b', 'n', 'r', 'p', 'Q', 'K', 'B', 'N', 'R', 'P':
				if x >= 8 {
					log.Fatalf("Your FEN sucks: %v; x is over or equal to 7: %v\n", fen, x)
				}
				grid[y][x] = line[i]
				x++
			default:
				log.Fatalf("Your FEN sucks: %v; read invalid character: %v\n", fen, line[i])
			}
		}
	}

	baseImage := resize.Resize(uint(finalImageWidth), 0, pathToImage[pieceToPath['*']], resize.Lanczos3)
	finalImage := image.NewRGBA(baseImage.Bounds())
	draw.Draw(finalImage, baseImage.Bounds(), baseImage, image.Point{}, draw.Src)

	// Use the debug flag to adjust the parameters to match the grid
	if *flagDebug {
		drawRect(boardMinX, boardMinY, boardMaxX, boardMaxY, finalImage, color.RGBA{255, 0, 0, 255})
		for y := range grid {
			for x := range grid[y] {
				drawRect(boardMinX+x*cellWidth, boardMinY+y*cellHeight, boardMinX+(x+1)*cellWidth, boardMinY+(y+1)*cellHeight, finalImage, color.RGBA{255, 0, 0, 255})
			}
		}
	}

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] != 0 {
				overImg := pathToImage[pieceToPath[grid[y][x]]]
				pieceBounds := overImg.Bounds()
				pieceWidth := pieceBounds.Max.X - pieceBounds.Min.X
				pieceHeight := pieceBounds.Max.Y - pieceBounds.Min.Y
				pos := image.Pt(boardMinX+x*cellWidth+cellWidth/2-pieceWidth/2, boardMinY+y*cellHeight+cellHeight/2-pieceHeight/2)
				draw.Draw(finalImage, overImg.Bounds().Add(pos), overImg, image.Point{}, draw.Over)
			}
		}
	}

	if err := png.Encode(os.Stdout, finalImage); err != nil {
		log.Fatal(err)
	}
}

func readImage(path string, resizeWidth int) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	if resizeWidth > 0 {
		img = resize.Resize(uint(resizeWidth), 0, img, resize.Bilinear)
	}

	return img, nil
}

// hLine draws a horizontal line
func hLine(x1, y, x2 int, img *image.RGBA, c color.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, c)
	}
}

// vLine draws a veritcal line
func vLine(x, y1, y2 int, img *image.RGBA, c color.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, c)
	}
}

// Rect draws a rectangle utilizing hLine() and vLine()
func drawRect(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	hLine(x1, y1, x2, img, c)
	hLine(x1, y2, x2, img, c)
	vLine(x1, y1, y2, img, c)
	vLine(x2, y1, y2, img, c)
}
