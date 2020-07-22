package engine

import "testing"

func TestIsValidMove(t *testing.T) {
	testBoard := [3][3]string{{"o", "", ""}, {"", "x", ""}, {"o", "", ""}}
	//check a valid move
	validMove := Move{row: 1, col: 2}
	if isValidMove(testBoard, validMove) != true {
		t.Error("isValidMove, expected true for a valid move, got false")
	}

	//check an invalid move
	inValidMove := Move{row: 0, col: 0}
	if isValidMove(testBoard, inValidMove) != false {
		t.Error("isValidMove, expected false for an invalid move, got true")
	}
}
