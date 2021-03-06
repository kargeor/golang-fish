// +build wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

const initial = "                     rnbqkbnr  pppppppp  ........  ........  ........  ........  PPPPPPPP  RNBQKBNR                     "

const (
	CLICK_SQUARE = iota
	CLICK_NEW_GAME_WHITE
	CLICK_NEW_GAME_BLACK
	CLICK_UNDO
)

type Event struct {
	event_type, x, y int
}

var document, chessboardDiv, logDiv js.Value
var squareDivs []js.Value
var pos *Position
var searcher *Searcher
var moveFrom = 0
var events = make(chan Event)
var spinner js.Value

func getElementById(name string) js.Value {
	return document.Call("getElementById", name)
}

func createDiv(parent js.Value) js.Value {
	newDiv := document.Call("createElement", "div")
	parent.Call("appendChild", newDiv)

	return newDiv
}

func addClickHandler(div js.Value, event Event) {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		select {
		case events <- event:
		default:
			// ignore event if already full
		}
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

func setSpinnerVisible(v bool) {
	if v {
		spinner.Get("style").Set("display", "")
	} else {
		spinner.Get("style").Set("display", "none")
	}
}

func clearSelected() {
	for _, div := range squareDivs {
		div.Set("className", "")
	}
}

func squareClickHandler(i, j int) {
	if pos.score <= -MATE_LOWER {
		log("You lost")
		return
	}

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

		waitForJs()

		if move_valid {
			pos = pos.move(move)
			rotated := pos.rotate()
			updateChessBoard(rotated)

			setSpinnerVisible(true)
			waitForJs()

			if pos.score <= -MATE_LOWER {
				log("You won!")
				setSpinnerVisible(false)
				return
			}

			start := time.Now()
			var bestResult SearchResult
			searcher.search(pos, func(r SearchResult) bool {
				elapsed := time.Since(start)
				log(fmt.Sprintf("(%s) depth=%d score=%d move=[%s]\n", elapsed, r.depth, r.score, r.move.rotate()))
				waitForJs()
				bestResult = r
				return r.depth >= 7
			})

			if bestResult.score == MATE_UPPER {
				log("Checkmate!")
			}

			pos = pos.move(bestResult.move)
			updateChessBoard(pos)
			setSpinnerVisible(false)
		}
	}
}

func newGame(playFirst bool) {
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

	if playFirst {
		var bestResult SearchResult
		searcher.search(pos, func(r SearchResult) bool {
			bestResult = r
			return true
		})

		pos = pos.move(bestResult.move)
	}

	updateChessBoard(pos)
}

func main() {
	document = js.Global().Get("document")
	chessboardDiv = getElementById("chessboard")
	logDiv = getElementById("logbox")
	spinner = getElementById("spinner")

	for i := 0; i < 8; i++ {
		p := createDiv(chessboardDiv)
		for j := 0; j < 8; j++ {
			div := createDiv(p)
			squareDivs = append(squareDivs, div)
			addClickHandler(div, Event{CLICK_SQUARE, i, j})
		}
	}

	addClickHandler(getElementById("playW"), Event{event_type: CLICK_NEW_GAME_WHITE})
	addClickHandler(getElementById("playB"), Event{event_type: CLICK_NEW_GAME_BLACK})
	addClickHandler(getElementById("undoMove"), Event{event_type: CLICK_UNDO})

	newGame(false)

	// wait for events
	for event := range events {
		switch event.event_type {
		case CLICK_SQUARE:
			squareClickHandler(event.x, event.y)
		case CLICK_NEW_GAME_WHITE:
			newGame(false)
		case CLICK_NEW_GAME_BLACK:
			newGame(true)
		case CLICK_UNDO:
		}
	}
}

var waitForJsTimeOutChan = make(chan bool)
var jsTimeoutCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	waitForJsTimeOutChan <- true
	return nil
})

func waitForJs() {
	js.Global().Call("setTimeout", jsTimeoutCallback, 10)
	<-waitForJsTimeOutChan
}
