package engine

//check if a move is valid
func isValidMove(board [3][3]string, m Move) bool {
	//a move is invalid if there is a player on the position
	return board[m.row][m.col] == ""
}
