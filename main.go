package main

import (
	"checkers/assets"
	"checkers/game"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	side = 50
	text = 35
	edge = 10
)

var (
	turn    game.Tile
	board   game.Board
	moves   [][][2]int
	objects [54]fyne.CanvasObject
	tiles   [50]*Tile
)

type layout struct{}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(10*side+2*edge, 10*side+2*edge+text)
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	table := min(containerSize.Width, containerSize.Height-text)
	piece := (table - 2*edge) / 10

	objects[0].Move(fyne.NewPos((containerSize.Width-table)/2, 0))
	objects[0].Resize(fyne.NewSquareSize(text))

	objects[1].Move(objects[0].Position().AddXY(text, 0))
	objects[0].Resize(fyne.NewSquareSize(text))

	objects[2].Move(objects[0].Position().AddXY(0, text))
	objects[2].Resize(fyne.NewSquareSize(table))

	objects[3].Move(objects[2].Position().AddXY(edge, edge))
	objects[3].Resize(fyne.NewSquareSize(table - 2*edge))

	for _, obj := range objects[4:] {
		obj.Move(objects[3].Position().AddXY(float32(obj.(*Tile).jdx-1)*piece, float32(obj.(*Tile).idx-1)*piece))
		obj.Resize(fyne.NewSquareSize(piece))
	}
}

type Change int

const (
	Equal Change = iota
	Reset
	Focus
	Taken
)

type Tile struct {
	widget.Icon

	idx     int
	jdx     int
	change  Change
	inFocus bool
}

func (t *Tile) tapEmpty() {
	kdx := slices.IndexFunc(tiles[:], func(t *Tile) bool {
		return t.inFocus && board[t.idx][t.jdx] > 2
	})

	x0, x1, sx := tiles[kdx].idx, t.idx, 1
	if x1 < x0 {
		sx = -1
	}
	y0, y1, sy := tiles[kdx].jdx, t.jdx, 1
	if y1 < y0 {
		sy = -1
	}

	board[x1][y1] = board[x0][y0]
	board[x0][y0] = game.Empty

	for x0, y0 = x0+sx, y0+sy; x0 != x1; x0, y0 = x0+sx, y0+sy {
		if board[x0][y0] != game.Empty {
			board[x0][y0] = game.Taken
			tiles[5*(x0-1)+(y0-1)/2].change = Taken
			break
		}
	}

	switch len(moves[0]) {
	case 2:
		for _, t := range tiles {
			switch {
			case board[t.idx][t.jdx] == game.WPawn && t.idx == 1:
				board[t.idx][t.jdx] = game.WKing
				t.change = Reset
			case board[t.idx][t.jdx] == game.BPawn && t.idx == 10:
				board[t.idx][t.jdx] = game.BKing
				t.change = Reset
			case board[t.idx][t.jdx] == game.Taken:
				board[t.idx][t.jdx] = game.Empty
				t.change = Reset
			case t.inFocus:
				t.change = Reset
			}
		}
		turn = 8 - turn
		objects[1].(*widget.Label).SetText(turn.String() + " plays")
		moves = board.CheckMoves(turn)
		if len(moves) == 0 {
			turn = 8 - turn
			objects[1].(*widget.Label).SetText(turn.String() + " wins")
		}
	default:
		moves = slices.DeleteFunc(moves, func(move [][2]int) bool {
			return move[1][0] != t.idx || move[1][1] != t.jdx
		})
		for idx, move := range moves {
			moves[idx] = move[1:]
		}
		t.inFocus = false
		t.tapPiece()
	}
}

func (t *Tile) tapPiece() {
	for _, t := range tiles {
		if t.inFocus {
			t.change = Reset
		}
	}
	if t.inFocus {
		return
	}
	t.change = Focus
	for _, move := range moves {
		if move[0][0] != t.idx || move[0][1] != t.jdx {
			continue
		}
		tiles[5*(move[1][0]-1)+(move[1][1]-1)/2].change = Focus
	}
}

func (t *Tile) changeReset() {
	switch board[t.idx][t.jdx] {
	case game.Empty:
		t.SetResource(assets.Brown)
	case game.WPawn:
		t.SetResource(assets.WPawn)
	case game.WKing:
		t.SetResource(assets.WKing)
	case game.BPawn:
		t.SetResource(assets.BPawn)
	case game.BKing:
		t.SetResource(assets.BKing)
	}
	t.inFocus = false
}

func (t *Tile) changeFocus() {
	switch board[t.idx][t.jdx] {
	case game.Empty:
		t.SetResource(assets.Focus)
	case game.WPawn:
		t.SetResource(assets.WPawnFocus)
	case game.WKing:
		t.SetResource(assets.WKingFocus)
	case game.BPawn:
		t.SetResource(assets.BPawnFocus)
	case game.BKing:
		t.SetResource(assets.BKingFocus)
	}
	t.inFocus = true
}

func (t *Tile) changeTaken() {
	t.SetResource(assets.Taken)
	t.inFocus = false
}

func (t *Tile) Tapped(*fyne.PointEvent) {
	if len(moves) == 0 || board[t.idx][t.jdx] == game.Taken || board[t.idx][t.jdx].IsEnemy(turn) {
		return
	}
	switch board[t.idx][t.jdx] {
	case game.Empty:
		if t.inFocus {
			t.tapEmpty()
		}
	default:
		t.tapPiece()
	}
	for _, t := range tiles {
		switch t.change {
		case Reset:
			t.changeReset()
		case Focus:
			t.changeFocus()
		case Taken:
			t.changeTaken()
		}
		t.change = Equal
	}
}

func newTile(idx, jdx int) *Tile {
	tile := &Tile{
		idx:    idx,
		jdx:    jdx,
		change: Equal,
	}
	tile.ExtendBaseWidget(tile)
	switch board[idx][jdx] {
	case game.BPawn:
		tile.SetResource(assets.BPawn)
	case game.Empty:
		tile.SetResource(assets.Brown)
	case game.WPawn:
		tile.SetResource(assets.WPawn)
	}
	return tile
}

func setupGame() {
	turn = game.WPawn
	board = game.NewGame()
	moves = board.CheckMoves(turn)

	objects[0] = widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() { setupGame() })
	objects[1] = widget.NewLabel("White plays")
	objects[2] = canvas.NewImageFromResource(assets.Table)
	objects[3] = canvas.NewImageFromResource(assets.Beige)

	oSlice := objects[4:4:cap(objects)]
	tSlice := tiles[0:0:cap(tiles)]

	for idx := 1; idx <= 10; idx += 1 {
		for jdx := 1 + idx%2; jdx <= 10; jdx += 2 {
			tSlice = append(tSlice, newTile(idx, jdx))
			oSlice = append(oSlice, tSlice[len(tSlice)-1])
		}
	}

}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Damas")

	myWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyF5:
			setupGame()
			myWindow.Content().Refresh()
		case fyne.KeyF11:
			myWindow.SetFullScreen(!myWindow.FullScreen())
		}
	})

	setupGame()

	myWindow.SetContent(container.New(&layout{}, objects[:]...))
	myWindow.ShowAndRun()
}
