package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

var digit = regexp.MustCompile(`\d`)
var slash = regexp.MustCompile(`/`)

func parseFEN(fen string) *Position {
	parts := strings.Split(fen, " ")
	board, color, castling, enpas := parts[0], parts[1], parts[2], parts[3]
	board = digit.ReplaceAllStringFunc(board, func(str string) string {
		count, _ := strconv.Atoi(str)
		return strings.Repeat(".", count)
	})
	board = "                     " + slash.ReplaceAllString(board, "  ") + "                     "

	if len(board) != 120 {
		fmt.Printf("FEN parse failed [%s]\n", fen)
		return nil
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
		return &pos
	}
	return pos.rotate()
}

func main() {
	golang_fish_init()

	interactiveFlagPtr := flag.Bool("i", false, "interactive mode (default is uci)")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		fmt.Printf("**CPU Profile Active**\n")
	}

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
			text := "a2a4" // used for profiling
			if *cpuprofile == "" {
				text, _ = reader.ReadString('\n')
			}

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
				return r.depth >= 9
			})

			if bestResult.score == MATE_UPPER {
				fmt.Printf("Checkmate!\n")
			}

			fmt.Printf("\nMy Move: depth=%d score=%d move=[%s]\n\n", bestResult.depth, bestResult.score, bestResult.move.rotate())

			pos = pos.move(bestResult.move)

			if *cpuprofile != "" {
				// just do on move for profiling
				return
			}
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
				// fmt.Printf("option name SETTING_MAX_DEPTH type spin default %d min 1 max 9999\n", SETTING_MAX_DEPTH)
				// fmt.Printf("option name SETTING_QS_LIMIT type spin default %d min 1 max 9999\n", SETTING_QS_LIMIT)
				// fmt.Printf("option name SETTING_EVAL_ROUGHNESS type spin default %d min 1 max 9999\n", SETTING_EVAL_ROUGHNESS)
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
					return /*r.depth >= SETTING_MAX_DEPTH ||*/ elapsed_ms > int64(time_left_msec)
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
