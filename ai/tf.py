#!/usr/bin/env python3


import tensorflow as tf
import numpy as np
import keras
import chess






def square_to_array(square):
    array = np.zeros((8, 8))
    # Handle string input (e.g., 'e4')
    if isinstance(square, str): square = chess.parse_square(square)
    # Convert square number (0-63) to rank and file
    rank = chess.square_rank(square)  # 0-7
    file = chess.square_file(square)  # 0-7
    # Set the corresponding position to 1
    array[rank, file] = 1
    return array

def piece_to_array(piece):
    array = np.zeros(12)
    if piece == None: return array
    piece_type = piece.piece_type - 1 # Get piece type (1-6 for PAWN through KING)
    if not piece.color: piece_type += 6 # If it's a black piece, offset by 6
    array[piece_type] = 1
    return array

def board_to_array(board):
    array = np.zeros((8, 8, 12))
    for rank in range(8):
        for file in range(8):
            square = chess.square(file, rank)
            piece = board.piece_at(square)
            if piece is not None: array[rank, file, :] = piece_to_array(piece)
    return array




model_from = keras.models.load_model("model_from.keras")
model_to = keras.models.load_model("model_to.keras")

board = chess.Board()


def print_board(board):
    ranks = range(8, 0, -1)
    files = 'abcdefgh'
    
    board_str = str(board)
    board_rows = board_str.split('\n')
    
    print("  " + " ".join(files))
    for rank, row in zip(ranks, board_rows):
        print(f"{rank} {row} {rank}")
    print("  " + " ".join(files))


def main():
    while True:
        print()
        print_board(board)
        print()

        if board.is_game_over():
            print("Game Over", board.outcome())
            break

        board_array = board_to_array(board).reshape(1, 8, 8, 12)
        pred_from = model_from.predict(board_array).reshape(64)
        pred_to = model_to.predict(board_array).reshape(64)

        possible_moves = []
        for m in board.legal_moves:
            i = m.from_square
            j = m.to_square
            score = pred_from[i] * pred_to[j]
            possible_moves.append( (i, j, score, m) )

        possible_moves = sorted(possible_moves, key=lambda x: x[2])[::-1]

        count = 4
        print("TOP MOVES:")
        for m in possible_moves:
            f = chess.square_name(m[0])
            t = chess.square_name(m[1])
            print('>> ', f, t, int(m[2]*1000))
            count -= 1
            if count == 0: break

        board.push(possible_moves[0][3])

        print()
        print_board(board)
        print()

        while True:
            move_input = input("Your Move: ")

            try:
                if len(move_input) == 4:  # Algebraic notation like 'e2e4'
                    move = chess.Move.from_uci(move_input)
                else:  # Standard algebraic notation
                    move = board.parse_san(move_input)

                if move in board.legal_moves:
                    board.push(move)
                    break
                else:
                    print("Illegal move. Try again.")
                    continue

            except chess.InvalidMoveError as e:
                print(f"Invalid move input: {e}")
                continue
            except chess.IllegalMoveError as e:
                print(f"IllegalMoveError: {e}")
                continue


def uci_loop():
    global board
    while True:
        line = input()
        if line.startswith('uci'):
            print("id name GeorgeOne")
            print("id author kargeor")
            print('uciok')
        elif line.startswith('isready'):
            print('readyok')
        elif line.startswith('position'):
            # Parse the position and set up the board
            parts = line.split()
            if parts[1] == 'startpos':
                board = chess.Board()
            elif parts[1] == 'fen':
                fen_string = ' '.join(parts[2:])
                board = chess.Board(fen=fen_string)
            else:
                # Parse move sequence and apply to the board
                for move_str in parts[2:]:
                    move = chess.Move.from_uci(move_str)
                    board.push(move)
        elif line.startswith('go'):
            # Play
            if board.turn == chess.BLACK: print("Error: Engine only plays as White.")

            board_array = board_to_array(board).reshape(1, 8, 8, 12)
            pred_from = model_from.predict(board_array).reshape(64)
            pred_to = model_to.predict(board_array).reshape(64)

            possible_moves = []
            for m in board.legal_moves:
                i = m.from_square
                j = m.to_square
                score = pred_from[i] * pred_to[j]
                possible_moves.append( (i, j, score, m) )

            possible_moves = sorted(possible_moves, key=lambda x: x[2])[::-1]
            
            print(f'bestmove {possible_moves[0][3]}')
        elif line.startswith('quit'):
            break

if __name__ == '__main__':
    uci_loop()


