#!/usr/bin/env python3

import chess
import chess.engine

# import logging
# logging.basicConfig(level=logging.DEBUG)


engine1 = chess.engine.SimpleEngine.popen_uci('./stockfish.macos.arm64.bin')

# Skill Level type spin default 20 min 0 max 20
# UCI_Elo type spin default 1320 min 1320 max 3190
engine1.configure({"Skill Level": 1, "UCI_Elo": 1321})

engine2 = chess.engine.SimpleEngine.popen_uci('./george_one_uci')

board = chess.Board()
while not board.is_game_over():
    # engine 2 plays white
    # engine 1 plays black
    e = engine1 if board.turn == chess.BLACK else engine2
    result = e.play(board, chess.engine.Limit(time=0.5))
    print(result.move)
    board.push(result.move)
    print(board)

engine1.quit()
engine2.quit()

print("Game Over", board.outcome())
if board.outcome().winner == chess.WHITE: print("WHITE WON")
if board.outcome().winner == chess.BLACK: print("BLACK WON")

