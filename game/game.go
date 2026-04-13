package game

import (
	"slices"
)

type Tile int

const (
	Empty Tile = iota
	Table
	Taken
	WPawn
	WKing
	BPawn
	BKing
)

func (t Tile) String() string {
	switch t {
	case WPawn:
		return "White"
	case BPawn:
		return "Black"
	}
	return ""
}

func (t Tile) IsAlly(other Tile) bool {
	return t > 2 && other > 2 && (9-2*t)*(9-2*other) > 0
}

func (t Tile) IsEnemy(other Tile) bool {
	return t > 2 && other > 2 && (9-2*t)*(9-2*other) < 0
}

func (t Tile) maxSteps() int {
	return 9 - 8*int(t%2)
}

type Board [12][12]Tile

func NewGame() Board {
	return Board{
		{Table, Table, Table, Table, Table, Table, Table, Table, Table, Table, Table, Table},
		{Table, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Table},
		{Table, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, Table},
		{Table, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Table},
		{Table, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, BPawn, Empty, Table},
		{Table, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Table},
		{Table, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty, Table},
		{Table, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Table},
		{Table, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, Table},
		{Table, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Table},
		{Table, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, WPawn, Empty, Table},
		{Table, Table, Table, Table, Table, Table, Table, Table, Table, Table, Table, Table},
	}
}

func (b Board) checkJump(idx, jdx int) [][][2]int {
	moves := [][][2]int{}

	for vStep := -1; vStep < 3; vStep += 2 {
		for hStep := -1; hStep < 3; hStep += 2 {
			k := 1
			for k < b[idx][jdx].maxSteps() && b[idx+k*vStep][jdx+k*hStep] == Empty {
				k += 1
			}
			if !b[idx][jdx].IsEnemy(b[idx+k*vStep][jdx+k*hStep]) {
				continue
			}
			takenIdx, takenJdx := idx+k*vStep, jdx+k*hStep
			for k += 1; k <= b[idx][jdx].maxSteps()+1 && b[idx+k*vStep][jdx+k*hStep] == Empty; k += 1 {
				b[idx+k*vStep][jdx+k*hStep] = b[idx][jdx]
				b[idx][jdx] = Empty
				takenTile := b[takenIdx][takenJdx]
				b[takenIdx][takenJdx] = Taken

				switch future := b.checkJump(idx+k*vStep, jdx+k*hStep); len(future) == 0 {
				case true:
					moves = append(moves, [][2]int{{idx + k*vStep, jdx + k*hStep}})
				case false:
					for _, path := range future {
						moves = append(moves, append([][2]int{{idx + k*vStep, jdx + k*hStep}}, path...))
					}
				}

				b[idx][jdx] = b[idx+k*vStep][jdx+k*hStep]
				b[idx+k*vStep][jdx+k*hStep] = Empty
				b[takenIdx][takenJdx] = takenTile
			}
		}
	}

	return moves
}

func (b Board) checkWalk(idx, jdx int) [][2]int {
	moves := [][2]int{}

	for vStep := -1; vStep < 3; vStep += 2 {
		if b[idx][jdx] == WPawn && vStep == 1 || b[idx][jdx] == BPawn && vStep == -1 {
			continue
		}
		for hStep := -1; hStep < 3; hStep += 2 {
			for k := 1; k <= b[idx][jdx].maxSteps() && b[idx+k*vStep][jdx+k*hStep] == Empty; k += 1 {
				moves = append(moves, [2]int{idx + k*vStep, jdx + k*hStep})
			}
		}
	}

	return moves
}

func (b Board) CheckMoves(turn Tile) [][][2]int {
	moves := [][][2]int{}

	for idx := 0; idx <= 10; idx += 1 {
		for jdx := 1 + idx%2; jdx <= 10; jdx += 2 {
			if turn.IsAlly(b[idx][jdx]) {
				for _, path := range b.checkJump(idx, jdx) {
					moves = append(moves, append([][2]int{{idx, jdx}}, path...))
				}
			}
		}
	}

	if len(moves) > 0 {
		maxMove := len(slices.MaxFunc(moves, func(new, old [][2]int) int {
			return len(new) - len(old)
		}))
		return slices.DeleteFunc(moves, func(move [][2]int) bool {
			return len(move) < maxMove
		})
	}

	for idx := 0; idx <= 10; idx += 1 {
		for jdx := 1 + idx%2; jdx <= 10; jdx += 2 {
			if turn.IsAlly(b[idx][jdx]) {
				for _, path := range b.checkWalk(idx, jdx) {
					moves = append(moves, [][2]int{{idx, jdx}, path})
				}
			}
		}
	}

	return moves
}
