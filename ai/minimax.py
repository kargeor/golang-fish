#!/usr/bin/env python3


import chess


def is_terminal(board):
    return board.is_game_over()


def evaluate(board):
    piece_values = {
        chess.PAWN: 100,
        chess.KNIGHT: 300,
        chess.BISHOP: 350,
        chess.ROOK: 500,
        chess.QUEEN: 900,
        chess.KING: 0  # King is invaluable, handled separately
    }

    # Positional bonuses for each piece
    pawn_table = [
        0, 0, 0, 0, 0, 0, 0, 0,
        5, 10, 10, -20, -20, 10, 10, 5,
        5, -5, -10, 0, 0, -10, -5, 5,
        0, 0, 0, 20, 20, 0, 0, 0,
        5, 5, 10, 25, 25, 10, 5, 5,
        10, 10, 20, 30, 30, 20, 10, 10,
        50, 50, 50, 50, 50, 50, 50, 50,
        0, 0, 0, 0, 0, 0, 0, 0
    ]

    knight_table = [
        -50, -40, -30, -30, -30, -30, -40, -50,
        -40, -20, 0, 0, 0, 0, -20, -40,
        -30, 0, 10, 15, 15, 10, 0, -30,
        -30, 5, 15, 20, 20, 15, 5, -30,
        -30, 0, 15, 20, 20, 15, 0, -30,
        -30, 5, 10, 15, 15, 10, 5, -30,
        -40, -20, 0, 5, 5, 0, -20, -40,
        -50, -40, -30, -30, -30, -30, -40, -50
    ]

    bishop_table = [
        -20, -10, -10, -10, -10, -10, -10, -20,
        -10, 0, 0, 0, 0, 0, 0, -10,
        -10, 0, 5, 10, 10, 5, 0, -10,
        -10, 5, 5, 10, 10, 5, 5, -10,
        -10, 0, 10, 10, 10, 10, 0, -10,
        -10, 10, 10, 10, 10, 10, 10, -10,
        -10, 5, 0, 0, 0, 0, 5, -10,
        -20, -10, -10, -10, -10, -10, -10, -20
    ]

    rook_table = [
        0, 0, 0, 5, 5, 0, 0, 0,
        -5, 0, 0, 0, 0, 0, 0, -5,
        -5, 0, 0, 0, 0, 0, 0, -5,
        -5, 0, 0, 0, 0, 0, 0, -5,
        -5, 0, 0, 0, 0, 0, 0, -5,
        -5, 0, 0, 0, 0, 0, 0, -5,
        5, 10, 10, 10, 10, 10, 10, 5,
        0, 0, 0, 0, 0, 0, 0, 0
    ]

    queen_table = [
        -20, -10, -10, -5, -5, -10, -10, -20,
        -10, 0, 0, 0, 0, 0, 0, -10,
        -10, 0, 5, 5, 5, 5, 0, -10,
        -5, 0, 5, 5, 5, 5, 0, -5,
        0, 0, 5, 5, 5, 5, 0, -5,
        -10, 5, 5, 5, 5, 5, 0, -10,
        -10, 0, 5, 0, 0, 0, 0, -10,
        -20, -10, -10, -5, -5, -10, -10, -20
    ]

    king_middle_game_table = [
        -30, -40, -40, -50, -50, -40, -40, -30,
        -30, -40, -40, -50, -50, -40, -40, -30,
        -30, -40, -40, -50, -50, -40, -40, -30,
        -30, -40, -40, -50, -50, -40, -40, -30,
        -20, -30, -30, -40, -40, -30, -30, -20,
        -10, -20, -20, -20, -20, -20, -20, -10,
        20, 20, 0, 0, 0, 0, 20, 20,
        20, 30, 10, 0, 0, 10, 30, 20
    ]

    # Combine positional values
    positional_tables = {
        chess.PAWN: pawn_table,
        chess.KNIGHT: knight_table,
        chess.BISHOP: bishop_table,
        chess.ROOK: rook_table,
        chess.QUEEN: queen_table,
        chess.KING: king_middle_game_table  # Use the king's middle-game table for now
    }

    # Material and positional evaluation
    score = 0
    for piece_type in piece_values:
        for square in board.pieces(piece_type, chess.WHITE):
            score += piece_values[piece_type]
            score += positional_tables[piece_type][square]
        for square in board.pieces(piece_type, chess.BLACK):
            score -= piece_values[piece_type]
            score -= positional_tables[piece_type][square]

    # Checkmate, stalemate, and check evaluation
    if board.is_checkmate():
        score = float('inf') if board.turn == chess.WHITE else float('-inf')
    elif board.is_stalemate():
        score = 0
    elif board.is_check():
        score += 25 if board.turn == chess.BLACK else -25

    return score



def get_children(board):
    children = []
    for move in board.legal_moves:
        # Create a copy of the board to apply the move
        new_board = board.copy()
        new_board.push(move)
        children.append(new_board)
    return children



def minimax_with_ab_pruning(node, depth, alpha, beta, maximizing_player):
    if depth == 0 or is_terminal(node):
        return evaluate(node), node
    if maximizing_player:
        max_eval = float('-inf')
        best_node = None
        for child in get_children(node):
            eval, _ = minimax_with_ab_pruning(child, depth - 1, alpha, beta, False)
            if eval > max_eval:
                max_eval = eval
                best_node = child
            alpha = max(alpha, eval)
            if beta <= alpha:  # Prune
                break
        return max_eval, best_node
    else:
        min_eval = float('inf')
        best_node = None
        for child in get_children(node):
            eval, _ = minimax_with_ab_pruning(child, depth - 1, alpha, beta, True)
            if eval < min_eval:
                min_eval = eval
                best_node = child
            beta = min(beta, eval)
            if beta <= alpha:  # Prune
                break
        return min_eval, best_node




board = chess.Board()

while True:
    line = input()
    if line.startswith('ucinewgame'):
        board = chess.Board()
    elif line.startswith('uci'):
        print("id name GeorgeMiniMax")
        print("id author kargeor")
        print('uciok')
    elif line.startswith('isready'):
        print('readyok')
    elif line.startswith('position startpos'):
        board = chess.Board()
        for move_str in line.split():
            if move_str == "position" or move_str == "startpos" or move_str == "moves": continue
            board.push(chess.Move.from_uci(move_str))
    elif line.startswith('position'):
        # Parse the position and set up the board
        parts = line.split()
        if parts[1] == 'startpos':
            board = chess.Board()
        elif parts[1] == 'fen':
            fen_string = ' '.join(parts[2:])
            fen_string = fen_string.split("moves")[0].strip()
            board = chess.Board(fen=fen_string)
            fen_string = ' '.join(parts[2:])
            if "moves" in fen_string:
                moves = fen_string.split("moves")[1].strip()
                for move_str in moves.split():
                    move = chess.Move.from_uci(move_str)
                    board.push(move)
        else:
            # Parse move sequence and apply to the board
            for move_str in parts[2:]:
                move = chess.Move.from_uci(move_str)
                board.push(move)
    elif line.startswith('go'):
        # Play
        best_value, best_board = minimax_with_ab_pruning(board, depth=4, alpha=float('-inf'), beta=float('inf'), maximizing_player=(board.turn == chess.WHITE))
        move = best_board.pop()
        
        print(f'info pv {move}')
        print(f'bestmove {move}')
    elif line.startswith('quit'):
        break
