// +build wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

const initial = "                     rnbqkbnr  pppppppp  ........  ........  ........  ........  PPPPPPPP  RNBQKBNR                     "

var document, chessboardDiv, logDiv js.Value
var squareDivs []js.Value
var pos *Position
var searcher *Searcher
var moveFrom = 0

func getElementById(name string) js.Value {
	return document.Call("getElementById", name)
}

func createDiv(parent js.Value) js.Value {
	newDiv := document.Call("createElement", "div")
	parent.Call("appendChild", newDiv)

	return newDiv
}

func addClickHandler(div js.Value, handler func(int, int), arg1, arg2 int) {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler(arg1, arg2)
		return nil
	})
	div.Call("addEventListener", "click", cb)
}

func updateChessBoard(pos *Position) {
	divCounter := 0

	for i := A8; i <= H1; i += S {
		for j := 0; j < 8; j++ {
			v := pos.board[i+j]
			if v == PIECE_IS_EMPTY {
				squareDivs[divCounter].Set("innerHTML", "&nbsp;")
			} else {
				squareDivs[divCounter].Set("innerText", v.String())
			}
			divCounter++
		}
	}
}

func log(str string) {
	logDiv.Set("innerText", str)
}

func clearSelected() {
	for _, div := range squareDivs {
		div.Set("className", "")
	}
}

func clickHandler(i, j int) {
	if moveFrom == 0 {
		moveFrom = A8 + i*S + j*E
		available := 0

		pos.gen_moves(func(m Move) bool {
			if m[0] == moveFrom {
				available++

				i1 := (m[1] - A8) / 10
				j1 := (m[1] - A8) % 10

				squareDivs[i1*8+j1].Set("className", "selected")
			}
			return false
		})

		log(fmt.Sprintf("Moves available %d\n", available))
	} else {
		moveTo := A8 + i*S + j*E
		move_valid := false
		var move Move

		pos.gen_moves(func(m Move) bool {
			if m[0] == moveFrom && m[1] == moveTo {
				move_valid = true
				move = m
				return true
			}
			return false
		})

		moveFrom = 0
		clearSelected()

		if move_valid {
			pos = pos.move(move)
			rotated := pos.rotate()
			updateChessBoard(rotated)

			start := time.Now()
			var bestResult SearchResult
			searcher.search(pos, func(r SearchResult) bool {
				elapsed := time.Since(start)
				log(fmt.Sprintf("(%s) depth=%d score=%d move=[%s]\n", elapsed, r.depth, r.score, r.move.rotate()))
				bestResult = r
				return r.depth >= 8
			})
		}
	}
}

func main() {
	document = js.Global().Get("document")
	chessboardDiv = getElementById("chessboard")
	logDiv = getElementById("logbox")

	for i := 0; i < 8; i++ {
		p := createDiv(chessboardDiv)
		for j := 0; j < 8; j++ {
			div := createDiv(p)
			squareDivs = append(squareDivs, div)
			addClickHandler(div, clickHandler, i, j)
		}
	}

	var initial_board Board
	for i := 0; i < 120; i++ {
		initial_board[i] = MakePiece(initial[i])
	}

	pos = &Position{
		board: initial_board,
		score: 0,
		wc:    [2]bool{true, true},
		bc:    [2]bool{true, true},
		ep:    0,
		kp:    0,
	}

	searcher = NewSearcher()
	updateChessBoard(pos)

	// wait forever (for callbacks)
	<-make(chan bool)
}
