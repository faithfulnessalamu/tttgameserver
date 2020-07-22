package engine

import "testing"

func TestIsRoundWon(t *testing.T) {
	testBoardCol := [3][3]string{{"o", "", ""}, {"o", "x", ""}, {"o", "", "x"}}
	if !isRoundWon(testBoardCol) {
		t.Error("isRoundWon, expected true got false")
	}

	testBoardRow := [3][3]string{{"o", "o", "o"}, {"x", "x", ""}, {"", "", ""}}
	if !isRoundWon(testBoardRow) {
		t.Error("isRoundWon, expected true got false")
	}

	testBoardDiagDown := [3][3]string{{"o", "x", ""}, {"x", "o", ""}, {"", "", "o"}}
	if !isRoundWon(testBoardDiagDown) {
		t.Error("isRoundWon, expected true got false")
	}

	testBoardDiagUp := [3][3]string{{"", "x", "o"}, {"x", "o", ""}, {"o", "", ""}}
	if !isRoundWon(testBoardDiagUp) {
		t.Error("isRoundWon, expected true got false")
	}

	testBoardFalse := [3][3]string{{"", "x", "o"}, {"x", "o", "o"}, {"", "", ""}}
	if isRoundWon(testBoardFalse) {
		t.Error("isRoundWon, expected false, got true")
	}
}

func TestEffectMove(t *testing.T) {
	testBoard := [3][3]string{{"o", "", ""}, {"", "x", ""}, {"o", "", ""}}
	validMove := Move{Row: 1, Col: 2}
	effectMove(testBoard[:], "x", validMove)
	if testBoard[validMove.Row][validMove.Col] != "x" {
		t.Error("effectMove, expected insertion into board")
	}
}

func TestIsValidMove(t *testing.T) {
	testBoard := [3][3]string{{"o", "", ""}, {"", "x", ""}, {"o", "", ""}}
	//check a valid move
	validMove := Move{Row: 1, Col: 2}
	if isValidMove(testBoard, validMove) != true {
		t.Error("isValidMove, expected true for a valid move, got false")
	}

	//check an invalid move
	inValidMove := Move{Row: 0, Col: 0}
	if isValidMove(testBoard, inValidMove) != false {
		t.Error("isValidMove, expected false for an invalid move, got true")
	}

	//check out of bounds
	inValidMove = Move{Row: 4, Col: -3}
	if isValidMove(testBoard, inValidMove) != false {
		t.Error("isValidMove, expected false for an invalid move, got true")
	}
}
