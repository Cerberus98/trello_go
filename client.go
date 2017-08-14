package main

import (
	"flag"
	"fmt"
	"os"
	trello "trello_api/trello"
)

type TrelloCommand struct {
	BaseURL *string
	ApiKey  *string
	Token   *string
	BoardId *string
}

func initFlags() (string, *TrelloCommand) {
	var boardConfig = TrelloCommand{}
	var cardsConfig = TrelloCommand{}
	var membersConfig = TrelloCommand{}
	var commandConfig *TrelloCommand
	var commandFlagSet *flag.FlagSet

	getBoardCmd := flag.NewFlagSet("board_get", flag.ExitOnError)
	getCardsCmd := flag.NewFlagSet("cards_get", flag.ExitOnError)
	getMembersCmd := flag.NewFlagSet("members_get", flag.ExitOnError)

	set_flags := func(cmd *TrelloCommand, flagset *flag.FlagSet) {
		cmd.BaseURL = flagset.String("baseurl", "https://api.trello.com", "Trello API Endpoint")
		cmd.ApiKey = flagset.String("apikey", os.Getenv("TRELLO_API_KEY"), "Your Trello API Key")
		cmd.Token = flagset.String("token", os.Getenv("TRELLO_API_TOKEN"), "Your Trello API Token")
	}
	set_flags(&boardConfig, getBoardCmd)
	set_flags(&cardsConfig, getCardsCmd)
	set_flags(&membersConfig, getMembersCmd)

	switch os.Args[1] {
	case "board_get":
		getBoardCmd.Parse(os.Args[2:])
		commandConfig = &boardConfig
		commandFlagSet = getBoardCmd
	case "cards_get":
		getCardsCmd.Parse(os.Args[2:])
		commandConfig = &cardsConfig
		commandFlagSet = getCardsCmd
	case "members_get":
		getMembersCmd.Parse(os.Args[2:])
		commandConfig = &membersConfig
		commandFlagSet = getMembersCmd
	default:
		fmt.Println("Unrecognized Command")
		os.Exit(1)
	}

	if *commandConfig.ApiKey == "" || commandFlagSet.NArg() == 0 {
		commandFlagSet.PrintDefaults()
		os.Exit(1)
	}

	if getBoardCmd.Parsed() {
	} else if getCardsCmd.Parsed() {
	} else if getMembersCmd.Parsed() {
	}

	if getBoardCmd.Parsed() || getCardsCmd.Parsed() || getMembersCmd.Parsed() {
		commandConfig.BoardId = &commandFlagSet.Args()[0]
	}
	return os.Args[1], commandConfig
}

func main() {
	cmd, cfg := initFlags()

	trelloApi := trello.TrelloApiV1{Key: *cfg.ApiKey, Token: *cfg.Token, BaseUrl: *cfg.BaseURL}
	board := trelloApi.GetBoard(*cfg.BoardId)
	switch cmd {
	case "cards_get":
		fmt.Println(board.GetCards())
	case "members_get":
		fmt.Println(board.GetMembers())
	case "board_get":
		fmt.Println(board)
	}
}
