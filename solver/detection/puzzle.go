package detection

import (
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"
)

const (
	coeffRGBToGrayR = 0.299
	coeffRGBToGrayG = 0.587
	coeffRGBToGrayB = 0.114

	bitShiftRGBA = 8
)

var (
	sobelX = [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	sobelY = [][]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}
)

type PuzzleSolver struct {
	puzzle string
	piece  string
}

func NewPuzzleSolver(base64Puzzle, base64Piece string) *PuzzleSolver {
	return &PuzzleSolver{
		puzzle: base64Puzzle,
		piece:  base64Piece,
	}
}

func (p *PuzzleSolver) GetPosition(yPos int) (int, error) {
	puzzle, err := p.backgroundPreprocessing()

	if err != nil {
		return 0, err
	}

	piece, err := p.piecePreprocessing()

	if err != nil {
		return 0, err
	}

	matchX := matchTemplate(puzzle, piece, yPos)

	return matchX, nil
}

func (p *PuzzleSolver) backgroundPreprocessing() ([][]float64, error) {
	img, err := decodeBase64ToGray(p.puzzle)

	if err != nil {
		return nil, err
	}

	return sobelOperator(img), nil
}

func (p *PuzzleSolver) piecePreprocessing() ([][]float64, error) {
	img, err := decodeBase64ToGray(p.piece)

	if err != nil {
		return nil, err
	}

	return sobelOperator(img), nil
}

func decodeBase64ToGray(b64 string) ([][]float64, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))

	img, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	grayMatrix := make([][]float64, height)

	for y := 0; y < height; y++ {
		grayMatrix[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := coeffRGBToGrayR*float64(r>>bitShiftRGBA) + coeffRGBToGrayG*float64(g>>bitShiftRGBA) + coeffRGBToGrayB*float64(b>>bitShiftRGBA)
			grayMatrix[y][x] = gray
		}
	}

	return grayMatrix, nil
}

func sobelOperator(img [][]float64) [][]float64 {
	height := len(img)
	width := len(img[0])

	grad := make([][]float64, height)

	for i := range grad {
		grad[i] = make([]float64, width)
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var gx, gy float64

			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					gx += img[y+ky][x+kx] * float64(sobelX[ky+1][kx+1])
					gy += img[y+ky][x+kx] * float64(sobelY[ky+1][kx+1])
				}
			}

			grad[y][x] = gx*gx + gy*gy
		}
	}

	return grad
}

func matchTemplate(image [][]float64, template [][]float64, yPos int) int {
	imgHeight := len(image)
	imgWidth := len(image[0])
	tempHeight := len(template)
	tempWidth := len(template[0])

	if yPos < 0 || yPos >= imgHeight {
		return -1
	}

	maxX := 0
	maxScore := -math.MaxFloat64

	startX := 0
	endX := imgWidth - tempWidth

	for x := startX; x <= endX; x++ {
		var score float64

		for ty := 0; ty < tempHeight; ty++ {
			for tx := 0; tx < tempWidth; tx++ {
				score += image[yPos+ty][x+tx] * template[ty][tx]
			}
		}

		if score > maxScore {
			maxScore = score
			maxX = x
		}
	}

	return maxX
}
