package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type Board struct {
	Cells         [9][9]Cell
	minCandidates int
}

type Cell struct {
	value      int
	candidates map[int]bool
}

func NewBoard(cells [9][9]int) Board {
	board := Board{}
	for i, row := range cells {
		for j, cell := range row {
			board.Cells[i][j] = Cell{
				value: cell,
			}
		}
	}
	board.setCandidates()

	return board
}

func (b *Board) setCandidates() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b.Cells[i][j].candidates = make(map[int]bool)
			if b.Cells[i][j].value != 0 {
				continue
			}

			for n := 1; n <= 9; n++ {
				b.Cells[i][j].value = n
				if b.isValid(i, j) {
					b.Cells[i][j].candidates[n] = true
				}
			}
			b.Cells[i][j].value = 0
		}
	}
	// set minCandidates
	b.setMinCandidates()
}

func (b *Board) setMinCandidates() {
	b.minCandidates = 9
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if b.Cells[i][j].value != 0 {
				continue
			}
			if len(b.Cells[i][j].candidates) < b.minCandidates {
				b.minCandidates = len(b.Cells[i][j].candidates)
			}
		}
	}
}

func (b *Board) set(row, col, value int) {
	prevValue := b.Cells[row][col].value
	b.Cells[row][col].value = value

	b.Cells[row][col].candidates = make(map[int]bool)
	if value == 0 {
		for n := 1; n <= 9; n++ {
			b.Cells[row][col].value = n
			if b.isValid(row, col) {
				b.Cells[row][col].candidates[n] = true
			}
		}
		b.Cells[row][col].value = 0
	}

	// update row's candidates
	for _, cell := range b.Cells[row] {
		if cell.value != 0 {
			continue
		}

		if value == 0 {
			cell.candidates[prevValue] = true
		} else {
			delete(cell.candidates, value)
		}
	}

	// update col's candidates
	for i := 0; i < 9; i++ {
		cell := b.Cells[i][col]
		if cell.value != 0 {
			continue
		}

		if value == 0 {
			cell.candidates[prevValue] = true
		} else {
			delete(cell.candidates, value)
		}
	}

	// update block's candidates
	blockRow := row / 3
	blockCol := col / 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cell := b.Cells[blockRow*3+i][blockCol*3+j]
			if cell.value != 0 {
				continue
			}

			if value == 0 {
				cell.candidates[prevValue] = true
			} else {
				delete(cell.candidates, value)
			}
		}
	}

	b.setMinCandidates()
}

func (b Board) pretty() string {
	var buf bytes.Buffer

	for i, row := range b.Cells {
		for j, cell := range row {
			buf.WriteString(fmt.Sprintf(" %d", cell.value))
			if j%3 == 2 && j != 8 {
				buf.WriteString("|")
			}
		}
		buf.WriteString("\n")
		if i%3 == 2 && i != 8 {
			buf.WriteString("------+------+------\n")
		}
	}
	return buf.String()
}

func (b *Board) isValid(row, col int) bool {
	rowCellNums := [9]int{}
	for i, cell := range b.Cells[row] {
		rowCellNums[i] = cell.value
	}
	if isDuplicate(rowCellNums) {
		return false
	}

	colCellNums := [9]int{}
	for i := 0; i < 9; i++ {
		colCellNums[i] = b.Cells[i][col].value
	}
	if isDuplicate(colCellNums) {
		return false
	}

	blockCellNums := [9]int{}
	blockRow := row / 3
	blockCol := col / 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			blockCellNums[i*3+j] = b.Cells[blockRow*3+i][blockCol*3+j].value
		}
	}
	if isDuplicate(blockCellNums) {
		return false
	}

	return true
}

func isDuplicate(arr [9]int) bool {
	seen := make(map[int]bool)

	for _, num := range arr {
		if num == 0 {
			continue
		}
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}

func solve(board *Board, debug bool) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board.Cells[i][j].value != 0 {
				continue
			}

			if len(board.Cells[i][j].candidates) > board.minCandidates {
				continue
			}

			for n := range board.Cells[i][j].candidates {
				if debug {
					time.Sleep(300 * time.Millisecond)
				}

				board.set(i, j, n)
				if debug {
					fmt.Println(board.pretty())
				}
				if board.isValid(i, j) && solve(board, debug) {
					return true
				}
				board.set(i, j, 0)
				if debug {
					fmt.Println("👈backtracked")
					fmt.Println(board.pretty())
				}
			}
			return false
		}
	}
	return true
}

func main() {
	// handle --debug flag
	debug := false
	for _, arg := range os.Args[1:] {
		if arg == "--debug" {
			debug = true
		}
	}

	board := NewBoard(
		[9][9]int{
			{5, 3, 1, 0, 0, 0, 0, 0, 0},
			{6, 0, 4, 0, 0, 0, 0, 0, 0},
			{2, 9, 8, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 3},
			{0, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 0, 0, 0, 0, 0, 0, 0, 6},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 5},
			{0, 0, 0, 0, 0, 0, 0, 7, 9},
		},
	)

	fmt.Println("initial")
	println(board.pretty())

	solved := solve(&board, debug)

	if solved {
		fmt.Println("solved")
		fmt.Println(board.pretty())
	} else {
		fmt.Println("unsolvable")
	}
}
