package engine

func isRoundWon(board [3][3]string) bool {
	//check rows
	for i := 0; i < 3; i++ {
		column := board[i]
		if column[0] == column[1] && column[1] == column[2] {
			if column[0] == "" { //empty strings would pass this 'if' check cause they'd all be equal to one another
				continue
			}
			return true
		}
	}

	//check cols
	for i := 0; i < 3; i++ {
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			if board[0][i] == "" { //empty strings would pass this 'if' check cause they'd all be equal to one another
				continue
			}
			return true
		}
	}

	//check diagdown
	j := 0
	if board[j][j] == board[j+1][j+1] && board[j+1][j+1] == board[j+2][j+2] {
		if board[j][j] != "" {
			return true
		}
	}

	//check diagup
	if board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		if board[0][2] != "" {
			return true
		}
	}

	return false
}

//effect a move
func effectMove(board [][3]string, avt string, m Move) {
	board[m.Row][m.Col] = avt
}

//check if a move is valid
func isValidMove(board [3][3]string, m Move) bool {
	//a move is invalid if there is a player on the position
	//check bounds first
	if m.Row < 0 || m.Row > 2 || m.Col < 0 || m.Col > 2 {
		return false
	}
	return board[m.Row][m.Col] == ""
}
