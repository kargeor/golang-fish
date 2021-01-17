package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IntArray []int
type Move [2]int
type Piece byte
type Board [120]Piece
type PieceToIntArray []IntArray

const FEN_INITIAL = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const PIECE_P = 0
const PIECE_N = 1
const PIECE_B = 2
const PIECE_R = 3
const PIECE_Q = 4
const PIECE_K = 5
const PIECE_IS_LOWER = 1 << 3
const PIECE_IS_EMPTY = 1 << 4
const PIECE_IS_INVALID = 1 << 5

var pst = PieceToIntArray{
	// PIECE_P
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 100, 100, 100, 100, 100, 100, 100, 0, 0, 178, 183, 186, 173, 202, 182, 185, 190, 0, 0, 107, 129, 121, 144, 140, 131, 144, 107, 0, 0, 83, 116, 98, 115, 114, 100, 115, 87, 0, 0, 74, 103, 110, 109, 106, 101, 100, 77, 0, 0, 78, 109, 105, 89, 90, 98, 103, 81, 0, 0, 69, 108, 93, 63, 64, 86, 103, 69, 0, 0, 100, 100, 100, 100, 100, 100, 100, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_N
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 214, 227, 205, 205, 270, 225, 222, 210, 0, 0, 277, 274, 380, 244, 284, 342, 276, 266, 0, 0, 290, 347, 281, 354, 353, 307, 342, 278, 0, 0, 304, 304, 325, 317, 313, 321, 305, 297, 0, 0, 279, 285, 311, 301, 302, 315, 282, 280, 0, 0, 262, 290, 293, 302, 298, 295, 291, 266, 0, 0, 257, 265, 282, 280, 282, 280, 257, 260, 0, 0, 206, 257, 254, 256, 261, 245, 258, 211, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_B
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 261, 242, 238, 244, 297, 213, 283, 270, 0, 0, 309, 340, 355, 278, 281, 351, 322, 298, 0, 0, 311, 359, 288, 361, 372, 310, 348, 306, 0, 0, 345, 337, 340, 354, 346, 345, 335, 330, 0, 0, 333, 330, 337, 343, 337, 336, 320, 327, 0, 0, 334, 345, 344, 335, 328, 345, 340, 335, 0, 0, 339, 340, 331, 326, 327, 326, 340, 336, 0, 0, 313, 322, 305, 308, 306, 305, 310, 310, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_R
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 514, 508, 512, 483, 516, 512, 535, 529, 0, 0, 534, 508, 535, 546, 534, 541, 513, 539, 0, 0, 498, 514, 507, 512, 524, 506, 504, 494, 0, 0, 479, 484, 495, 492, 497, 475, 470, 473, 0, 0, 451, 444, 463, 458, 466, 450, 433, 449, 0, 0, 437, 451, 437, 454, 454, 444, 453, 433, 0, 0, 426, 441, 448, 453, 450, 436, 435, 426, 0, 0, 449, 455, 461, 484, 477, 461, 448, 447, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_Q
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 935, 930, 921, 825, 998, 953, 1017, 955, 0, 0, 943, 961, 989, 919, 949, 1005, 986, 953, 0, 0, 927, 972, 961, 989, 1001, 992, 972, 931, 0, 0, 930, 913, 951, 946, 954, 949, 916, 923, 0, 0, 915, 914, 927, 924, 928, 919, 909, 907, 0, 0, 899, 923, 916, 918, 913, 918, 913, 902, 0, 0, 893, 911, 929, 910, 914, 914, 908, 891, 0, 0, 890, 899, 898, 916, 898, 893, 895, 887, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_K
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 60004, 60054, 60047, 59901, 59901, 60060, 60083, 59938, 0, 0, 59968, 60010, 60055, 60056, 60056, 60055, 60010, 60003, 0, 0, 59938, 60012, 59943, 60044, 59933, 60028, 60037, 59969, 0, 0, 59945, 60050, 60011, 59996, 59981, 60013, 60000, 59951, 0, 0, 59945, 59957, 59948, 59972, 59949, 59953, 59992, 59950, 0, 0, 59953, 59958, 59957, 59921, 59936, 59968, 59971, 59968, 0, 0, 59996, 60003, 59986, 59950, 59943, 59982, 60013, 60004, 0, 0, 60017, 60030, 59997, 59986, 60006, 59999, 60040, 60018, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

const A1, H1, A8, H8 = 91, 98, 21, 28

const N, E, S, W = -10, 1, 10, -1

var directions = PieceToIntArray{
	// PIECE_P
	{N, N + N, N + W, N + E},
	// PIECE_N
	{N + N + E, E + N + E, E + S + E, S + S + E, S + S + W, W + S + W, W + N + W, N + N + W},
	// PIECE_B
	{N + E, S + E, S + W, N + W},
	// PIECE_R
	{N, E, S, W},
	// PIECE_Q
	{N, E, S, W, N + E, S + E, S + W, N + W},
	// PIECE_K
	{N, E, S, W, N + E, S + E, S + W, N + W},
}

const MATE_LOWER = 50710
const MATE_UPPER = 69290

const TABLE_SIZE = 1e7

var SETTING_QS_LIMIT = 219
var SETTING_EVAL_ROUGHNESS = 13
var SETTING_MAX_DEPTH = 50

type Position struct {
	board  Board
	score  int
	wc, bc [2]bool
	ep, kp int
}

func (board *Board) contains(p Piece) bool {
	for _, v := range board {
		if v == p {
			return true
		}
	}

	return false
}

func (p Piece) isupper() bool {
	return p != PIECE_IS_INVALID && p != PIECE_IS_EMPTY && (p&PIECE_IS_LOWER) == 0
}

func (p Piece) islower() bool {
	return p != PIECE_IS_INVALID && p != PIECE_IS_EMPTY && (p&PIECE_IS_LOWER) != 0
}

func (p Piece) isspace() bool {
	return p == PIECE_IS_INVALID
}

func (p Piece) swapcase() Piece {
	if p == PIECE_IS_INVALID || p == PIECE_IS_EMPTY {
		return p
	}

	return p ^ PIECE_IS_LOWER
}

func (p Piece) String() string {
	switch p {
	case PIECE_IS_EMPTY:
		return "."
	case PIECE_IS_INVALID:
		return " "
	case PIECE_P:
		return "P"
	case PIECE_N:
		return "N"
	case PIECE_B:
		return "B"
	case PIECE_R:
		return "R"
	case PIECE_Q:
		return "Q"
	case PIECE_K:
		return "K"
	case PIECE_P | PIECE_IS_LOWER:
		return "p"
	case PIECE_N | PIECE_IS_LOWER:
		return "n"
	case PIECE_B | PIECE_IS_LOWER:
		return "b"
	case PIECE_R | PIECE_IS_LOWER:
		return "r"
	case PIECE_Q | PIECE_IS_LOWER:
		return "q"
	case PIECE_K | PIECE_IS_LOWER:
		return "k"
	}

	return "?"
}

func MakePiece(b byte) Piece {
	switch b {
	case '.':
		return PIECE_IS_EMPTY
	case ' ':
		return PIECE_IS_INVALID
	case 'P':
		return PIECE_P
	case 'N':
		return PIECE_N
	case 'B':
		return PIECE_B
	case 'R':
		return PIECE_R
	case 'Q':
		return PIECE_Q
	case 'K':
		return PIECE_K
	case 'p':
		return PIECE_P | PIECE_IS_LOWER
	case 'n':
		return PIECE_N | PIECE_IS_LOWER
	case 'b':
		return PIECE_B | PIECE_IS_LOWER
	case 'r':
		return PIECE_R | PIECE_IS_LOWER
	case 'q':
		return PIECE_Q | PIECE_IS_LOWER
	case 'k':
		return PIECE_K | PIECE_IS_LOWER
	}

	fmt.Printf("Bad piece byte %d\n", b)
	return 0xFF
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (self *Position) print() {
	line := 8

	for i := A8; i <= H1; i += S {
		fmt.Printf("  %d ", line)
		for j := 0; j < 8; j++ {
			fmt.Printf("%s ", self.board[i+j])
		}
		fmt.Printf("\n")
		line--
	}
	fmt.Printf("    a b c d e f g h\n\n")
}

func (self *Position) gen_moves(yield func(m Move) bool) {
	for i, p := range self.board {
		if !(p.isupper()) {
			continue
		}
		for _, d := range directions[p] {
			for j := (i + d); ; j += d {
				q := self.board[j]
				if q.isspace() || q.isupper() {
					break
				}

				if p == PIECE_P && (d == N || d == N+N) && q != PIECE_IS_EMPTY {
					break
				}

				if p == PIECE_P && d == N+N && (i < A1+N || self.board[i+N] != PIECE_IS_EMPTY) {
					break
				}

				if p == PIECE_P && (d == N+W || d == N+E) && q == PIECE_IS_EMPTY && j != self.ep && j != self.kp && j != self.kp-1 && j != self.kp+1 {
					break
				}

				if yield(Move{i, j}) {
					return
				}

				if q.islower() || p == PIECE_P || p == PIECE_N || p == PIECE_K {
					break
				}

				if i == A1 && self.board[j+E] == PIECE_K && self.wc[0] {
					if yield(Move{j + E, j + W}) {
						return
					}
				}

				if i == H1 && self.board[j+W] == PIECE_K && self.wc[1] {
					if yield(Move{j + W, j + E}) {
						return
					}
				}
			}
		}
	}
}

func (self *Position) sorted_moves() []Move {
	var result []Move
	self.gen_moves(func(m Move) bool {
		result = append(result, m)
		return false
	})

	sort.Slice(result, func(i, j int) bool {
		return self.value(result[i]) > self.value(result[j])
	})

	return result
}

func (self *Position) is_dead() bool {
	result := false
	self.gen_moves(func(m Move) bool {
		if self.value(m) >= MATE_LOWER {
			result = true
			return true
		}
		return false
	})

	return result
}

func (self *Position) rotate() Position {
	board := self.board
	score := -self.score
	wc := self.bc
	bc := self.wc
	ep := 0
	kp := 0

	if self.ep != 0 {
		ep = 119 - self.ep
	}

	if self.kp != 0 {
		kp = 119 - self.kp
	}

	// might miss swapcase if len() is odd!
	for i, j := 0, len(board)-1; i < j; i, j = i+1, j-1 {
		board[i], board[j] = board[j].swapcase(), board[i].swapcase()
	}

	return Position{board, score, wc, bc, ep, kp}
}

func (self *Position) nullmove() Position {
	result := self.rotate()
	result.ep = 0
	result.kp = 0
	return result
}

func (self *Position) move(move Move) Position {
	i, j := move[0], move[1]
	p := self.board[i]
	board := self.board
	wc, bc, ep, kp := self.wc, self.bc, 0, 0
	score := self.score + self.value(move)

	board[j] = board[i]
	board[i] = PIECE_IS_EMPTY

	if i == A1 {
		wc = [2]bool{false, wc[1]}
	}
	if i == H1 {
		wc = [2]bool{wc[0], false}
	}
	if j == A8 {
		bc = [2]bool{bc[0], false}
	}
	if j == H8 {
		bc = [2]bool{false, bc[1]}
	}

	if p == PIECE_K {
		wc = [2]bool{false, false}
		if abs(j-i) == 2 {
			kp = (i + j) / 2
			if j < i {
				board[A1] = PIECE_IS_EMPTY
			} else {
				board[H1] = PIECE_IS_EMPTY
			}
			board[kp] = PIECE_R
		}
	}

	if p == PIECE_P {
		if A8 <= j && j <= H8 {
			board[j] = PIECE_Q
		}
		if j-i == 2*N {
			ep = i + N
		}
		if j == self.ep {
			board[j+S] = PIECE_IS_EMPTY
		}
	}

	position := Position{board, score, wc, bc, ep, kp}
	return position.rotate()
}

func (self *Position) value(move Move) int {
	i, j := move[0], move[1]
	p, q := self.board[i], self.board[j]

	score := pst[p][j] - pst[p][i]

	if q.islower() {
		score += pst[q.swapcase()][119-j]
	}

	if abs(j-self.kp) < 2 {
		score += pst[PIECE_K][119-j]
	}

	if p == PIECE_K && abs(i-j) == 2 {
		score += pst[PIECE_R][(i+j)/2]
		if j < i {
			score -= pst[PIECE_R][A1]
		} else {
			score -= pst[PIECE_R][H1]
		}
	}

	if p == PIECE_P {
		if A8 <= j && j <= H8 {
			score += pst[PIECE_Q][j] - pst[PIECE_P][j]
		}
		if j == self.ep {
			score += pst[PIECE_P][119-(j+S)]
		}
	}

	return score
}

var digit = regexp.MustCompile(`\d`)
var slash = regexp.MustCompile(`/`)

func parseFEN(fen string) Position {
	parts := strings.Split(fen, " ")
	board, color, castling, enpas := parts[0], parts[1], parts[2], parts[3]
	board = digit.ReplaceAllStringFunc(board, func(str string) string {
		count, _ := strconv.Atoi(str)
		return strings.Repeat(".", count)
	})
	board = "                     " + slash.ReplaceAllString(board, "  ") + "                     "

	if len(board) != 120 {
		fmt.Printf("FEN parse failed [%s]\n", fen)
		return Position{}
	}

	var parsed_board Board
	for i := 0; i < 120; i++ {
		parsed_board[i] = MakePiece(board[i])
	}

	wc := [2]bool{strings.Contains(castling, "Q"), strings.Contains(castling, "K")}
	bc := [2]bool{strings.Contains(castling, "k"), strings.Contains(castling, "q")}

	ep := 0
	if enpas != "-" {
		/////TODO.....
		// ep = sunfish.parse(enpas) if enpas != '-' else 0
	}

	/////TODO.....
	// score = sum(sunfish.pst[p][i] for i,p in enumerate(board) if p.isupper())
	// score -= sum(sunfish.pst[p.upper()][119-i] for i,p in enumerate(board) if p.islower())
	score := 0

	pos := Position{
		board: parsed_board,
		score: score,
		wc:    wc,
		bc:    bc,
		ep:    ep,
		kp:    0,
	}

	if color == "w" {
		return pos
	}
	return pos.rotate()
}

type Entry struct {
	lower, upper int
}

type PDR struct {
	pos   Position
	depth int
	root  bool
}

type Searcher struct {
	tp_score map[PDR]Entry
	tp_move  map[Position]Move
	nodes    int
}

type ScoreMove struct {
	valid bool
	score int
	move  Move
}

func NewSearcher() *Searcher {
	return &Searcher{
		tp_score: make(map[PDR]Entry),
		tp_move:  make(map[Position]Move),
		nodes:    0,
	}
}

func (self *Searcher) bound(pos Position, gamma int, depth int, root bool) int {
	self.nodes += 1
	depth = max(depth, 0)

	if pos.score <= -MATE_LOWER {
		return -MATE_UPPER
	}

	entry, entry_found := self.tp_score[PDR{pos, depth, root}]
	if !entry_found {
		entry = Entry{-MATE_UPPER, MATE_UPPER}
	}

	if entry.lower >= gamma {
		if !root {
			return entry.lower
		}
		if _, found := self.tp_move[pos]; found {
			return entry.lower
		}
	}

	if entry.upper < gamma {
		return entry.upper
	}

	moves := func(yield func(sm ScoreMove) bool) {
		if depth > 0 && !root {
			if pos.board.contains(PIECE_R) ||
				pos.board.contains(PIECE_B) ||
				pos.board.contains(PIECE_N) ||
				pos.board.contains(PIECE_Q) {
				if yield(ScoreMove{
					valid: false,
					score: -self.bound(pos.nullmove(), 1-gamma, depth-3, false),
				}) {
					return
				}
			}
		}

		if depth == 0 {
			if yield(ScoreMove{
				valid: false,
				score: pos.score,
			}) {
				return
			}
		}

		killer, killer_found := self.tp_move[pos]
		if killer_found && (depth > 0 || pos.value(killer) >= SETTING_QS_LIMIT) {
			if yield(ScoreMove{
				valid: true,
				move:  killer,
				score: -self.bound(pos.move(killer), 1-gamma, depth-1, false),
			}) {
				return
			}
		}

		for _, move := range pos.sorted_moves() {
			if depth > 0 || pos.value(move) >= SETTING_QS_LIMIT {
				if yield(ScoreMove{
					valid: true,
					move:  move,
					score: -self.bound(pos.move(move), 1-gamma, depth-1, false),
				}) {
					return
				}
			}
		}
	}

	best := -MATE_UPPER
	moves(func(sm ScoreMove) bool {
		best = max(best, sm.score)
		if best >= gamma {
			if len(self.tp_move) > TABLE_SIZE {
				fmt.Printf("info string tp_move table clear\n")
				self.tp_move = make(map[Position]Move)
			}

			if sm.valid {
				self.tp_move[pos] = sm.move
			} else {
				delete(self.tp_move, pos)
			}

			return true
		}

		return false
	})

	if best < gamma && best < 0 && depth > 0 {
		all_is_dead := true
		pos.gen_moves(func(m Move) bool {
			pos_move := pos.move(m)
			if !(pos_move.is_dead()) {
				all_is_dead = false
				return true
			}
			return false
		})
		if all_is_dead {
			pos_nullmove := pos.nullmove()
			in_check := pos_nullmove.is_dead()
			if in_check {
				best = -MATE_UPPER
			} else {
				best = 0
			}
		}
	}

	if len(self.tp_score) > TABLE_SIZE {
		fmt.Printf("info string tp_score table clear\n")
		self.tp_score = make(map[PDR]Entry)
	}

	if best >= gamma {
		self.tp_score[PDR{pos, depth, root}] = Entry{best, entry.upper}
	}

	if best < gamma {
		self.tp_score[PDR{pos, depth, root}] = Entry{entry.lower, best}
	}

	return best
}

type SearchResult struct {
	depth int
	move  Move
	score int
	nodes int
}

func (self *Searcher) search(pos Position, yield func(r SearchResult) bool) {
	self.nodes = 0

	for depth := 1; depth < 1000; depth++ {
		lower, upper := -MATE_UPPER, MATE_UPPER
		for lower < upper-SETTING_EVAL_ROUGHNESS {
			gamma := (lower + upper + 1) / 2
			score := self.bound(pos, gamma, depth, true)
			if score >= gamma {
				lower = score
			} else {
				upper = score
			}
		}

		self.bound(pos, lower, depth, true)

		if yield(SearchResult{
			depth: depth,
			move:  self.tp_move[pos],
			score: self.tp_score[PDR{pos, depth, true}].lower,
			nodes: self.nodes,
		}) {
			return
		}
	}
}

func (m Move) String() string {
	from := m[0] - A8
	from1 := from % 10
	from2 := from / 10

	to := m[1] - A8
	to1 := to % 10
	to2 := to / 10

	return fmt.Sprintf("%c%d%c%d", from1+'a', 8-from2, to1+'a', 8-to2)
}

func (m Move) rotate() Move {
	return Move{119 - m[0], 119 - m[1]}
}

func parseMove(str string) (Move, bool) {
	if len(str) < 4 {
		return Move{}, false
	}

	tbytes := []byte(str)
	m0 := int(tbytes[0] - 'a')
	m1 := int(tbytes[1] - '1')
	m2 := int(tbytes[2] - 'a')
	m3 := int(tbytes[3] - '1')

	if m0 < 0 || m0 > 7 || m1 < 0 || m1 > 7 || m2 < 0 || m2 > 7 || m3 < 0 || m3 > 7 {
		return Move{}, false
	}

	return Move{A1 + m0*E + m1*N, A1 + m2*E + m3*N}, true
}

func main() {
	interactiveFlagPtr := flag.Bool("i", false, "interactive mode (default is uci)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	searcher := NewSearcher()

	pos := parseFEN(FEN_INITIAL)

	if *interactiveFlagPtr {
		for true {
			pos.print()

			if pos.score <= -MATE_LOWER {
				fmt.Printf("You lost\n")
				return
			}

			fmt.Printf("Your move: ")
			text, _ := reader.ReadString('\n')

			move, move_parse_valid := parseMove(text)
			if !move_parse_valid {
				continue
			}
			valid := false

			pos.gen_moves(func(m Move) bool {
				if m == move {
					valid = true
					return true
				}
				return false
			})

			if !valid {
				continue
			}

			fmt.Printf("Your move = %s\n", move)

			pos = pos.move(move)
			rotated := pos.rotate()
			rotated.print()

			if pos.score <= -MATE_LOWER {
				fmt.Printf("You won!\n")
				return
			}

			start := time.Now()
			var bestResult SearchResult
			searcher.search(pos, func(r SearchResult) bool {
				elapsed := time.Since(start)
				fmt.Printf("(%s) depth=%d score=%d move=[%s]\n", elapsed, r.depth, r.score, r.move.rotate())
				bestResult = r
				return r.depth >= 8
			})

			if bestResult.score == MATE_UPPER {
				fmt.Printf("Checkmate!\n")
			}

			fmt.Printf("\nMy Move: depth=%d score=%d move=[%s]\n\n", bestResult.depth, bestResult.score, bestResult.move.rotate())

			pos = pos.move(bestResult.move)
		}
	} else {
		white_turn := true
		for true {
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			switch {
			case strings.HasPrefix(command, "quit"):
				return
			case strings.HasPrefix(command, "ucinewgame"):
				searcher = NewSearcher()
				pos = parseFEN(FEN_INITIAL)
			case strings.HasPrefix(command, "uci"):
				fmt.Printf("id name GoLangFish\n")
				fmt.Printf("id author kargeor & Sunfish Contributors\n")
				fmt.Printf("option name SETTING_MAX_DEPTH type spin default %d min 1 max 9999\n", SETTING_MAX_DEPTH)
				fmt.Printf("option name SETTING_QS_LIMIT type spin default %d min 1 max 9999\n", SETTING_QS_LIMIT)
				fmt.Printf("option name SETTING_EVAL_ROUGHNESS type spin default %d min 1 max 9999\n", SETTING_EVAL_ROUGHNESS)
				fmt.Printf("uciok\n")
			case strings.HasPrefix(command, "setoption"):
				// TODO......
			case strings.HasPrefix(command, "isready"):
				fmt.Printf("readyok\n")
			case strings.HasPrefix(command, "position"):
				parts := strings.Split(command, " ")
				for i := 1; i < len(parts); i++ {
					part := parts[i]
					switch {
					case strings.HasPrefix(part, "startpos"):
						pos = parseFEN(FEN_INITIAL)
						white_turn = true
					case strings.HasPrefix(part, "moves"):
						for i++; i < len(parts); i++ {
							move, move_ok := parseMove(parts[i])
							if move_ok {
								if white_turn {
									pos = pos.move(move)
									white_turn = false
								} else {
									pos = pos.move(move.rotate())
									white_turn = true
								}
							} else {
								fmt.Printf("info string Failed to parse move [%s]\n", parts[i])
							}
						}
					case strings.HasPrefix(part, "fen"):
						// TODO......
						// don't forget to set [white_turn]
						fmt.Printf("info string Failed to parse FEN\n")
					}
				}
			case strings.HasPrefix(command, "go"):
				wtime := 60000
				btime := 60000
				movestogo := 10
				parts := strings.Split(command, " ")
				for i := 1; i < len(parts); i++ {
					part := parts[i]
					switch {
					case strings.HasPrefix(part, "wtime"):
						i++
						wtime, _ = strconv.Atoi(parts[i])
					case strings.HasPrefix(part, "btime"):
						i++
						btime, _ = strconv.Atoi(parts[i])
					case strings.HasPrefix(part, "movestogo"):
						i++
						movestogo, _ = strconv.Atoi(parts[i])
					}
				}

				time_left_msec := btime / movestogo
				if white_turn {
					time_left_msec = wtime / movestogo
				}
				time_left_msec = max(0, time_left_msec-250) // safety margin

				start := time.Now()
				var bestResult SearchResult

				searcher.search(pos, func(r SearchResult) bool {
					elapsed_ms := time.Since(start).Milliseconds()

					pv := bestResult.move
					if !white_turn {
						pv = pv.rotate()
					}

					fmt.Printf("info depth %d score cp %d nodes %d time %d pv %s\n", r.depth, r.score, r.nodes, elapsed_ms, pv)
					bestResult = r
					return r.depth >= SETTING_MAX_DEPTH || elapsed_ms > int64(time_left_msec)
				})

				if white_turn {
					fmt.Printf("bestmove %s\n", bestResult.move)
				} else {
					fmt.Printf("bestmove %s\n", bestResult.move.rotate())
				}
			}
		}
	}
}
