package engine

import (
	"chessBot/board"
	"chessBot/fen"
	"chessBot/uci"
	"fmt"
	"log"
	"os"
)

var CurrentBoard *board.Board
var LogOutput *os.File

func init() {
	var err error
	LogOutput, err = os.Create("outstandingMove.log")
	if err != nil {
		log.Fatal("error opening log file:", err)
	}
}

func Send(msg string) {
	fmt.Println(msg)
}

func Log(msg string) {
	LogOutput.WriteString(msg + "\n")
	LogOutput.Sync()
}

func InitBoard(stmnt *uci.PositionStatement) error {
	initString := fen.STARTPOSSTRING

	if stmnt.IsFen {
		initString = stmnt.FenString
	}
	var err error
	CurrentBoard, err = fen.FenToBoard(initString)
	if err != nil {
		return err
	}

	for _, moveString := range stmnt.Moves {
		Log("moving " + moveString)
		move := board.MoveFromString(moveString)
		CurrentBoard.Move(move)
	}

	return nil
}
