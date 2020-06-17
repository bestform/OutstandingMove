package main

import (
	"bufio"
	"chessBot/uci"
	"fmt"
	"log"
	"os"
)

func main() {

	logfile, err := os.Create("outstandingMove.log")
	if err != nil {
		log.Fatal("error opening log file:", err)
	}
	logfile.WriteString("Welcome to Outstanding Move! Waiting for commands...\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		logger(logfile, text)
		stmnts, err := uci.Parse(text)
		if err != nil {
			logger(logfile, "Error when parsing input:\n")
			logger(logfile, err.Error())
			continue
		}

		for _, stmnt := range stmnts {
			logger(logfile, "Statement: " + fmt.Sprintf("%+v", &stmnt))
			switch stmnt.Kind {
			case uci.UciStatementKind:
				logger(logfile, "received uci statement")
				fmt.Print("uciok")
			case uci.IsReadyStatementKind:
				logger(logfile, "received isready statement")
				fmt.Print("readyok")
			case uci.UciNewGameStatementKind:
				logger(logfile, "received ucinewgame statement")
			case uci.PositionStatementKind:
				logger(logfile, "received position statement")
				logger(logfile, fmt.Sprintf("%+v", stmnt.Position))
			}
		}
	}

}

func logger(f *os.File, msg string) {
	f.WriteString(msg + "\n")
	f.Sync()
}