package engine

//check if a move is valid
func isValidMove(board [3][3]string, m Move) bool {
	//a move is invalid if there is a player on the position
	//check bounds first
	if m.row < 0 || m.row > 2 || m.col < 0 || m.col > 2 {
		return false
	}
	return board[m.row][m.col] == ""
}
