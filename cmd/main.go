package main

import (
	"bufio"
	"chessBot/engine"
	"chessBot/uci"
	"fmt"
	"math/rand"
	"os"
)

func main() {


	engine.Log("Welcome to Outstanding Move! Waiting for commands...\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		engine.Log(text)
		stmnts, err := uci.Parse(text)
		if err != nil {
			engine.Log("Error when parsing input:\n")
			engine.Log(err.Error())
			continue
		}

		for _, stmnt := range stmnts {
			engine.Log("<- " + string(stmnt.Kind))
			switch stmnt.Kind {
			case uci.UciStatementKind:
				engine.Send("uciok")
			case uci.IsReadyStatementKind:
				engine.Send("readyok")
			case uci.UciNewGameStatementKind:
			case uci.PositionStatementKind:
				engine.Log(fmt.Sprintf("%+v", stmnt.Position))
				err = engine.InitBoard(stmnt.Position)
				if err != nil {
					engine.Log("error initializing board: " + err.Error())
				}
				engine.Log("Current Board:")
				engine.Log(engine.CurrentBoard.String())
				engine.CurrentBoard.SwitchSide() // todo: do that on engine level.
			case uci.GoStatementKind:
				moves := engine.CalculatePossibleMoves()
				move := moves[rand.Int() % len(moves)]
				engine.Send("bestmove " + move.String())

			}
		}
	}

}

