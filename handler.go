package main

import (
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	addMessageToList(m.ID) // add new message id to list

	if m.Author.ID == s.State.User.ID {
		return
	}

	if isValidAppCommand(m.Content) {
		err := executeCommand(s, m, m.Content)
		if err != nil {
			log.Println(err)
		}
		return
	}

	//quotes
	randomQuotes(s, m, m.Content)
}

func isValidAppCommand(line string) bool {
	fields := strings.Fields(line)
	if len(fields) > 0 {
		if fields[0] == "!go" && len(fields) >= 2 {
			return true
		}
	}
	return false
}

func executeCommand(s *discordgo.Session, m *discordgo.MessageCreate, line string) error {
	cmd, err := getCommand(line)
	if err != nil {
		return err
	}
	args := getArgs(line)
	err = cmd.Run(s, m, args)
	if err != nil {
		return err
	}
	return err
}

func getCommand(line string) (*command, error) {
	fields := strings.Fields(line)

	for _, v := range commandList {
		if fields[1] == v.Use {
			return &v, nil
		}
	}
	return nil, errors.New("Command not found")
}

func getArgs(line string) []string {
	fields := strings.Fields(line)
	var args []string
	for i := 2; i < len(fields); i++ {
		args = append(args, fields[i])
	}
	return args
}

func addMessageToList(id string) {

	msgIDList.PushFront(id)

	if msgIDList.Len() >= maxMasgListSize {
		msgIDList.Remove(msgIDList.Back())
	}
}

func randomQuotes(s *discordgo.Session, m *discordgo.MessageCreate, line string) error {
	cmd, err := getCommand(appCommand + " " + quoteCommand.Use)

	args := strings.Fields(line)
	err = cmd.Run(s, m, args)
	if err != nil {
		return err
	}
	return err
}
