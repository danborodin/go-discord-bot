package main

import "github.com/bwmarrin/discordgo"

type command struct {
	Use         string
	Description string
	Run         func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
}
