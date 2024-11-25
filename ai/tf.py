import tensorflow as tf
import numpy as np
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


