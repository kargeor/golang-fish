package main

import (
	"fmt"
	"sort"
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
const PIECE_NOT_PIECE = PIECE_IS_EMPTY | PIECE_IS_INVALID

var piece = IntArray{
	// P
	100,
	// N
	280,
	// B
	320,
	// R
	479,
	// Q
	929,
	// K
	60_000,
}

var pst = PieceToIntArray{
	// PIECE_P
	{0, 0, 0, 0, 0, 0, 0, 0,
		78, 83, 86, 73, 102, 82, 85, 90,
		7, 29, 21, 44, 40, 31, 44, 7,
		-17, 16, -2, 15, 14, 0, 15, -13,
		-26, 3, 10, 9, 6, 1, 0, -23,
		-22, 9, 5, -11, -10, -2, 3, -19,
		-31, 8, -7, -37, -36, -14, 3, -31,
		0, 0, 0, 0, 0, 0, 0, 0},
	// PIECE_N
	{-66, -53, -75, -75, -10, -55, -58, -70,
		-3, -6, 100, -36, 4, 62, -4, -14,
		10, 67, 1, 74, 73, 27, 62, -2,
		24, 24, 45, 37, 33, 41, 25, 17,
		-1, 5, 31, 21, 22, 35, 2, 0,
		-18, 10, 13, 22, 18, 15, 11, -14,
		-23, -15, 2, 0, 2, 0, -23, -20,
		-74, -23, -26, -24, -19, -35, -22, -69},
	// PIECE_B
	{-59, -78, -82, -76, -23, -107, -37, -50,
		-11, 20, 35, -42, -39, 31, 2, -22,
		-9, 39, -32, 41, 52, -10, 28, -14,
		25, 17, 20, 34, 26, 25, 15, 10,
		13, 10, 17, 23, 17, 16, 0, 7,
		14, 25, 24, 15, 8, 25, 20, 15,
		19, 20, 11, 6, 7, 6, 20, 16,
		-7, 2, -15, -12, -14, -15, -10, -10},
	// PIECE_R
	{35, 29, 33, 4, 37, 33, 56, 50,
		55, 29, 56, 67, 55, 62, 34, 60,
		19, 35, 28, 33, 45, 27, 25, 15,
		0, 5, 16, 13, 18, -4, -9, -6,
		-28, -35, -16, -21, -13, -29, -46, -30,
		-42, -28, -42, -25, -25, -35, -26, -46,
		-53, -38, -31, -26, -29, -43, -44, -53,
		-30, -24, -18, 5, -2, -18, -31, -32},
	// PIECE_Q
	{6, 1, -8, -104, 69, 24, 88, 26,
		14, 32, 60, -10, 20, 76, 57, 24,
		-2, 43, 32, 60, 72, 63, 43, 2,
		1, -16, 22, 17, 25, 20, -13, -6,
		-14, -15, -2, -5, -1, -10, -20, -22,
		-30, -6, -13, -11, -16, -11, -16, -27,
		-36, -18, 0, -19, -15, -15, -21, -38,
		-39, -30, -31, -13, -31, -36, -34, -42},
	// PIECE_K
	{4, 54, 47, -99, -99, 60, 83, -62,
		-32, 10, 55, 56, 56, 55, 10, 3,
		-62, 12, -57, 44, -67, 28, 37, -31,
		-55, 50, 11, -4, -19, 13, 0, -49,
		-55, -43, -52, -28, -51, -47, -8, -50,
		-47, -42, -43, -79, -64, -32, -29, -32,
		-4, 3, -14, -50, -57, -18, 13, 4,
		17, 30, -3, -14, 6, -1, 40, 18},
}

var golang_fish_init_done = false

func golang_fish_init() {
	if golang_fish_init_done {
		return
	}

	golang_fish_init_done = true

	// Pad tables and join piece and pst dictionaries.
	for i := range pst {
		result := IntArray{}

		// Add 20 zeros
		for j := 0; j < 20; j += 1 {
			result = append(result, 0)
		}

		for j := 0; j < 64; j += 8 {

			base := piece[i]

			result = append(result,
				0,
				base+pst[i][j+0],
				base+pst[i][j+1],
				base+pst[i][j+2],
				base+pst[i][j+3],
				base+pst[i][j+4],
				base+pst[i][j+5],
				base+pst[i][j+6],
				base+pst[i][j+7],
				0,
			)
		}

		// Add 20 zeros
		for j := 0; j < 20; j += 1 {
			result = append(result, 0)
		}

		pst[i] = result
	}
}

const A1, H1, A8, H8 = 91, 98, 21, 28

const N, E, S, W = -10, 1, 10, -1

// Lists of possible moves for each piece type.
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

/*
Mate value must be greater than 8*queen + 2*(rook+knight+bishop).
King value is set to twice this value such that if the opponent is
8 queens up, but we got the king, we still exceed MATE_VALUE.
When a MATE is detected, we'll set the score to MATE_UPPER - plies to get there.
E.g. Mate in 3 will be MATE_UPPER - 6.
*/
var MATE_LOWER = piece[PIECE_K] - 10*piece[PIECE_Q]
var MATE_UPPER = piece[PIECE_K] + 10*piece[PIECE_Q]

// Constants for tuning search
const QS = 40             // 0..300
const QS_A = 140          // 0..300
const EVAL_ROUGHNESS = 15 //0..50

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
	return (p&PIECE_IS_LOWER) == 0 && (p&PIECE_NOT_PIECE) == 0
}

func (p Piece) islower() bool {
	return (p&PIECE_IS_LOWER) != 0 && (p&PIECE_NOT_PIECE) == 0
}

func (p Piece) is_invalid_space() bool {
	return p == PIECE_IS_INVALID
}

func (p Piece) swapcase() Piece {
	if (p & PIECE_NOT_PIECE) != 0 {
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

// iif returns trueValue if condition is true, otherwise falseValue
func iif[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func (self *Position) print() {
	line := 8

	fmt.Printf("     a b c d e f g h\n")
	for i := A8; i <= H1; i += S {
		fmt.Printf("  %d  ", line)
		for j := 0; j < 8; j++ {
			fmt.Printf("%s ", self.board[i+j])
		}
		fmt.Printf(" %d\n", line)
		line--
	}
	fmt.Printf("     a b c d e f g h\n\n")
}

func (self *Position) gen_moves(yield func(m Move) bool) {
	for i, p := range self.board {
		if !(p.isupper()) {
			continue
		}
		for _, d := range directions[p] {
			for j := (i + d); ; j += d {
				q := self.board[j]
				if q.is_invalid_space() || q.isupper() {
					break
				}

				if p == PIECE_P {
					if (d == N || d == N+N) && q != PIECE_IS_EMPTY {
						break
					}

					if d == N+N && (i < A1+N || self.board[i+N] != PIECE_IS_EMPTY) {
						break
					}

					if (d == N+W || d == N+E) && q == PIECE_IS_EMPTY && j != self.ep && j != self.kp && j != self.kp-1 && j != self.kp+1 {
						break
					}
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

func (self *Position) rotate() *Position {
	pos := &Position{}

	pos.board = self.board
	pos.score = -self.score
	pos.wc = self.bc
	pos.bc = self.wc
	pos.ep = 0
	pos.kp = 0

	if self.ep != 0 {
		pos.ep = 119 - self.ep
	}

	if self.kp != 0 {
		pos.kp = 119 - self.kp
	}

	// rotate & swap case
	pos.board[21], pos.board[98] = pos.board[98].swapcase(), pos.board[21].swapcase()
	pos.board[22], pos.board[97] = pos.board[97].swapcase(), pos.board[22].swapcase()
	pos.board[23], pos.board[96] = pos.board[96].swapcase(), pos.board[23].swapcase()
	pos.board[24], pos.board[95] = pos.board[95].swapcase(), pos.board[24].swapcase()
	pos.board[25], pos.board[94] = pos.board[94].swapcase(), pos.board[25].swapcase()
	pos.board[26], pos.board[93] = pos.board[93].swapcase(), pos.board[26].swapcase()
	pos.board[27], pos.board[92] = pos.board[92].swapcase(), pos.board[27].swapcase()
	pos.board[28], pos.board[91] = pos.board[91].swapcase(), pos.board[28].swapcase()
	pos.board[31], pos.board[88] = pos.board[88].swapcase(), pos.board[31].swapcase()
	pos.board[32], pos.board[87] = pos.board[87].swapcase(), pos.board[32].swapcase()
	pos.board[33], pos.board[86] = pos.board[86].swapcase(), pos.board[33].swapcase()
	pos.board[34], pos.board[85] = pos.board[85].swapcase(), pos.board[34].swapcase()
	pos.board[35], pos.board[84] = pos.board[84].swapcase(), pos.board[35].swapcase()
	pos.board[36], pos.board[83] = pos.board[83].swapcase(), pos.board[36].swapcase()
	pos.board[37], pos.board[82] = pos.board[82].swapcase(), pos.board[37].swapcase()
	pos.board[38], pos.board[81] = pos.board[81].swapcase(), pos.board[38].swapcase()
	pos.board[41], pos.board[78] = pos.board[78].swapcase(), pos.board[41].swapcase()
	pos.board[42], pos.board[77] = pos.board[77].swapcase(), pos.board[42].swapcase()
	pos.board[43], pos.board[76] = pos.board[76].swapcase(), pos.board[43].swapcase()
	pos.board[44], pos.board[75] = pos.board[75].swapcase(), pos.board[44].swapcase()
	pos.board[45], pos.board[74] = pos.board[74].swapcase(), pos.board[45].swapcase()
	pos.board[46], pos.board[73] = pos.board[73].swapcase(), pos.board[46].swapcase()
	pos.board[47], pos.board[72] = pos.board[72].swapcase(), pos.board[47].swapcase()
	pos.board[48], pos.board[71] = pos.board[71].swapcase(), pos.board[48].swapcase()
	pos.board[51], pos.board[68] = pos.board[68].swapcase(), pos.board[51].swapcase()
	pos.board[52], pos.board[67] = pos.board[67].swapcase(), pos.board[52].swapcase()
	pos.board[53], pos.board[66] = pos.board[66].swapcase(), pos.board[53].swapcase()
	pos.board[54], pos.board[65] = pos.board[65].swapcase(), pos.board[54].swapcase()
	pos.board[55], pos.board[64] = pos.board[64].swapcase(), pos.board[55].swapcase()
	pos.board[56], pos.board[63] = pos.board[63].swapcase(), pos.board[56].swapcase()
	pos.board[57], pos.board[62] = pos.board[62].swapcase(), pos.board[57].swapcase()
	pos.board[58], pos.board[61] = pos.board[61].swapcase(), pos.board[58].swapcase()

	return pos
}

func (self *Position) nullmove() *Position {
	pos := self.rotate()
	pos.ep = 0
	pos.kp = 0
	return pos
}

func (self *Position) move(move Move) *Position {
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

/*
Let s* be the "true" score of the sub-tree we are searching.
The method returns r, where:
- if gamma >  s* then s* <= r < gamma  (A better upper bound).
- if gamma <= s* then gamma <= r <= s* (A better lower bound).
can_null = !root = Default TRUE.
*/
func (self *Searcher) bound(pos *Position, gamma int, depth int, can_null bool) int {
	self.nodes += 1

	/*
		Depth <= 0 is QSearch. Here any position is searched as deeply as is needed for
		calmness, and from this point on there is no difference in behaviour depending on
		depth, so so there is no reason to keep different depths in the transposition table.
	*/
	depth = max(depth, 0)

	/*
		Sunfish is a king-capture engine, so we should always check if we
		still have a king. Notice since this is the only termination check,
		the remaining code has to be comfortable with being mated, stalemated
		or able to capture the opponent king.
	*/
	if pos.score <= -MATE_LOWER {
		return -MATE_UPPER
	}

	/*
		Look in the table if we have already searched this position before.
		We also need to be sure, that the stored search was over the same
		nodes as the current search.
	*/
	entry, entry_found := self.tp_score[PDR{*pos, depth, can_null}]
	if !entry_found {
		entry = Entry{-MATE_UPPER, MATE_UPPER}
	}

	if entry.lower >= gamma {
		return entry.lower
	}

	if entry.upper < gamma {
		return entry.upper
	}

	// TODO: Let's not repeat positions
	// if can_null and depth > 0 and pos in self.history: return 0

	// Generator of moves to search in order.
	moves := func(yield func(sm ScoreMove) bool) {
		// First try not moving at all. We only do this if there is at least one major piece left on the board.
		if depth > 2 && can_null && abs(pos.score) < 500 {
			if yield(ScoreMove{
				valid: false,
				score: -self.bound(pos.nullmove(), 1-gamma, depth-3, false),
			}) {
				return
			}
		}

		// For QSearch we have a different kind of null-move, namely we can just stop
		// and not capture anything else.
		if depth == 0 {
			if yield(ScoreMove{
				valid: false,
				score: pos.score,
			}) {
				return
			}
		}

		// Look for the strongest move from last time, the hash-move.
		killer, killer_found := self.tp_move[*pos]

		/*
			If there isn't one, try to find one with a more shallow search.
			This is known as Internal Iterative Deepening (IID). We set
			can_null=True, since we want to make sure we actually find a move.
		*/
		if !killer_found && depth > 2 {
			self.bound(pos, gamma, depth-3, false)
			killer, killer_found = self.tp_move[*pos]
		}

		// If depth == 0 we only try moves with high intrinsic score (captures and
		// promotions). Otherwise we do all moves. This is called quiescent search.
		val_lower := QS - depth*QS_A

		/*
			Only play the move if it would be included at the current val-limit,
			since otherwise we'd get search instability.
			We will search it again in the main loop below, but the tp will fix
			things for us.
		*/
		if killer_found && pos.value(killer) >= val_lower {
			if yield(ScoreMove{
				valid: true,
				move:  killer,
				score: -self.bound(pos.move(killer), 1-gamma, depth-1, true),
			}) {
				return
			}
		}

		// Then all the other moves

		sorted_moves := make([]Move, 0, 64)
		pos.gen_moves(func(m Move) bool {
			sorted_moves = append(sorted_moves, m)
			return false
		})

		sort.Slice(sorted_moves, func(i, j int) bool {
			return pos.value(sorted_moves[i]) > pos.value(sorted_moves[j])
		})

		for i := 0; i < len(sorted_moves); i++ {
			move := sorted_moves[i]
			val := pos.value(move)

			// Quiescent search
			if val < val_lower {
				break
			}

			// If the new score is less than gamma, the opponent will for sure just
			// stand pat, since ""pos.score + val < gamma === -(pos.score + val) >= 1-gamma""
			// This is known as futility pruning.
			if depth <= 1 && pos.score+val < gamma {
				// Need special case for MATE, since it would normally be caught before standing pat.
				if yield(ScoreMove{
					valid: true,
					move:  move,
					score: pos.score + iif(val < MATE_LOWER, val, MATE_UPPER),
				}) {
					return
				}
				// We can also break, since we have ordered the moves by value, so it can't get any better than this.
				break
			}

			if yield(ScoreMove{
				valid: true,
				move:  move,
				score: -self.bound(pos.move(move), 1-gamma, depth-1, true),
			}) {
				return
			}
		}
	}

	// Run through the moves, shortcutting when possible.
	best := -MATE_UPPER
	moves(func(sm ScoreMove) bool {
		best = max(best, sm.score)
		if best >= gamma {
			// Save the move for pv construction and killer heuristic.
			if sm.valid {
				self.tp_move[*pos] = sm.move
			}

			// break
			return true
		}

		// keep returning results
		return false
	})

	// Stalemate checking is a bit tricky...

	// This is too expensive to test at depth == 0
	if depth > 2 && best == -MATE_UPPER {
		flipped := pos.nullmove()

		// Hopefully this is already in the TT because of null-move.
		in_check := self.bound(flipped, MATE_UPPER, 0, true) == MATE_UPPER
		best = iif(in_check, -MATE_LOWER, 0)
	}

	// Table part 2

	if best >= gamma {
		self.tp_score[PDR{*pos, depth, can_null}] = Entry{best, entry.upper}
	}

	if best < gamma {
		self.tp_score[PDR{*pos, depth, can_null}] = Entry{entry.lower, best}
	}

	return best
}

type SearchResult struct {
	depth int
	gamma int
	score int
	move  Move
}

// Iterative deepening MTD-bi search
func (self *Searcher) search(pos *Position, yield func(r SearchResult) bool) {
	self.nodes = 0
	// self.history = set(history)
	self.tp_score = make(map[PDR]Entry)

	gamma := 0

	// In finished games, we could potentially go far enough to cause a recursion
	// limit exception. Hence we bound the ply. We also can't start at 0, since
	// that's quiscent search, and we don't always play legal moves there.
	for depth := 1; depth < 1000; depth++ {
		// The inner loop is a binary search on the score of the position.
		// Inv: lower <= score <= upper
		// 'while lower != upper' would work, but it's too much effort to spend
		// on what's probably not going to change the move played.
		lower, upper := -MATE_LOWER, MATE_LOWER

		for lower < upper-EVAL_ROUGHNESS {
			score := self.bound(pos, gamma, depth, false)
			if score >= gamma {
				lower = score
			} else {
				upper = score
			}

			if yield(SearchResult{
				depth: depth,
				gamma: gamma,
				score: score,
				move:  self.tp_move[*pos],
			}) {
				return
			}

			gamma = (lower + upper + 1) / 2
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
		debug("parseMove: len must be 4+")
		return Move{}, false
	}

	tbytes := []byte(str)
	m0 := int(tbytes[0] - 'a')
	m1 := int(tbytes[1] - '1')
	m2 := int(tbytes[2] - 'a')
	m3 := int(tbytes[3] - '1')

	if m0 < 0 || m0 > 7 || m1 < 0 || m1 > 7 || m2 < 0 || m2 > 7 || m3 < 0 || m3 > 7 {
		debug("parseMove: bad number")
		return Move{}, false
	}

	return Move{A1 + m0*E + m1*N, A1 + m2*E + m3*N}, true
}
