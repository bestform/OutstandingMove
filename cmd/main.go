package main

import (
	"chessBot"
	"fmt"
)

func main() {
	b := chessBot.NewBoard()
	fmt.Println(b[0])
	fmt.Println(b[15])
	fmt.Println(b[63])
}