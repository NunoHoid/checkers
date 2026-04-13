package assets

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"fyne.io/fyne/v2"
)

const (
	side = 255
	core = 100
	ring = 5
)

var (
	tableEdgeColour = color.RGBA{150, 75, 0, 255}
	beigeTileColour = color.RGBA{250, 225, 200, 255}
	brownTileColour = color.RGBA{200, 150, 100, 255}

	whiteCoreColour = color.RGBA{230, 230, 230, 255}
	blackCoreColour = color.RGBA{75, 75, 75, 255}
	pieceRingColour = color.RGBA{25, 25, 25, 255}

	focusTileColour = color.RGBA{100, 200, 100, 255}
	takenTileColour = color.RGBA{250, 125, 125, 255}
)

func bresenham(asset *image.RGBA, colour color.RGBA, x0, y0, x1, y1 int) {
	dx, sx := x1-x0, 1
	if dx < 0 {
		dx = -dx
		sx = -1
	}
	dy, sy := y1-y0, -1
	if dy > 0 {
		dy = -dy
		sy = 1
	}
	for e0, e1 := dx+dy, 2*(dx+dy); ; e1 = 2 * e0 {
		asset.Set(x0, y0, colour)
		if x0 == x1 && y0 == y1 {
			return
		}
		if e1 >= dy {
			x0 += sx
			e0 += dy
		}
		if e1 <= dx {
			y0 += sy
			e0 += dx
		}
	}
}

func drawSquare(colour color.RGBA) func(*image.RGBA) {
	return func(asset *image.RGBA) {
		for x := range side {
			for y := range side {
				asset.SetRGBA(x, y, colour)
			}
		}
	}
}

func drawCircle(colour color.RGBA) func(*image.RGBA) {
	return func(asset *image.RGBA) {
		for dx := -core; dx <= core; dx += 1 {
			for dy := -core; dy <= core; dy += 1 {
				switch {
				case dx*dx+dy*dy <= (core-ring)*(core-ring):
					asset.SetRGBA(dx+side/2, dy+side/2, colour)
				case dx*dx+dy*dy <= core*core:
					asset.SetRGBA(dx+side/2, dy+side/2, pieceRingColour)
				}
			}
		}
	}
}

func drawCrown(colour color.RGBA) func(*image.RGBA) {
	return func(asset *image.RGBA) {
		baseX, baseY := 40, 45
		centerX, centerY := 25, 15
		topX, topY := 60, 35

		points := [...][2]int{
			{0, -topY},
			{-centerX, centerY},
			{-topX, -topY},
			{-baseX, baseY},
			{baseX, baseY},
			{topX, -topY},
			{centerX, centerY},
			{0, -topY},
		}

		for idx := range len(points) - 1 {
			bresenham(asset, colour,
				side/2+points[idx][0],
				side/2+points[idx][1],
				side/2+points[idx+1][0],
				side/2+points[idx+1][1],
			)
		}

		for queue := [][2]int{{side / 2, side/2 + baseY - 1}}; len(queue) > 0; {
			last := queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			asset.Set(last[0], last[1], colour)
			for _, pair := range [...][2]int{{-1, 0}, {1, 0}, {0, -1}} {
				if asset.At(last[0]+pair[0], last[1]+pair[1]) != colour {
					queue = append(queue, [2]int{last[0] + pair[0], last[1] + pair[1]})
				}
			}
		}
	}
}

func drawAsset(drawFunc ...func(*image.RGBA)) []byte {
	asset := image.NewRGBA(image.Rect(0, 0, side, side))
	for _, df := range drawFunc {
		df(asset)
	}
	buffer := bytes.NewBuffer(nil)
	png.Encode(buffer, asset)
	return buffer.Bytes()
}

var Table = fyne.NewStaticResource("Table", drawAsset(
	drawSquare(tableEdgeColour),
))

var Beige = fyne.NewStaticResource("Beige", drawAsset(
	drawSquare(beigeTileColour),
))

var Brown = fyne.NewStaticResource("Brown", drawAsset(
	drawSquare(brownTileColour),
))

var Focus = fyne.NewStaticResource("Focus", drawAsset(
	drawSquare(focusTileColour),
))

var Taken = fyne.NewStaticResource("Focus", drawAsset(
	drawSquare(takenTileColour),
))

var WPawn = fyne.NewStaticResource("WPawn", drawAsset(
	drawSquare(brownTileColour), drawCircle(whiteCoreColour),
))

var WPawnFocus = fyne.NewStaticResource("WPawnFocus", drawAsset(
	drawSquare(focusTileColour), drawCircle(whiteCoreColour),
))

var WKing = fyne.NewStaticResource("WKing", drawAsset(
	drawSquare(brownTileColour), drawCircle(whiteCoreColour), drawCrown(blackCoreColour),
))

var WKingFocus = fyne.NewStaticResource("WKingFocus", drawAsset(
	drawSquare(focusTileColour), drawCircle(whiteCoreColour), drawCrown(blackCoreColour),
))

var BPawn = fyne.NewStaticResource("BPawn", drawAsset(
	drawSquare(brownTileColour), drawCircle(blackCoreColour),
))

var BPawnFocus = fyne.NewStaticResource("BPawnFocus", drawAsset(
	drawSquare(focusTileColour), drawCircle(blackCoreColour),
))

var BKing = fyne.NewStaticResource("BKing", drawAsset(
	drawSquare(brownTileColour), drawCircle(blackCoreColour), drawCrown(whiteCoreColour),
))

var BKingFocus = fyne.NewStaticResource("BKingFocus", drawAsset(
	drawSquare(focusTileColour), drawCircle(blackCoreColour), drawCrown(whiteCoreColour),
))
