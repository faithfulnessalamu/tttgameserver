package engine

//check if a move is valid
func isValidMove(board [3][3]string, m Move) bool {
	//a move is invalid if there is a player on the position
	//check bounds first
	if m.Row < 0 || m.Row > 2 || m.Col < 0 || m.Col > 2 {
		return false
	}
	return board[m.Row][m.Col] == ""
}
