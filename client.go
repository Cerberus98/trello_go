package main

import (
	"fmt"
	"os"
	trello "trello_api/trello"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("client.go <endpoint url> <api key> <token>")
		return
	}

	boardId := os.Args[1]
	baseUrl := os.Args[2]
	apiKey := os.Args[3]
	token := os.Args[4]

	trelloApi := trello.TrelloApiV1{Key: apiKey, Token: token, BaseUrl: baseUrl}
	board := trelloApi.GetBoard(boardId)
	fmt.Println(board)
	members := trelloApi.GetMembers(boardId)
	fmt.Println(members)
}
