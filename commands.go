package main

import (
	"container/list"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	maxMasgListSize = 100
)

var (
	appCommand  string = "!go"
	commandList []command
	msgIDList   = list.New()
)

var helpCommand = command{
	Use:         "help",
	Description: "List all commands",
	Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		if len(args) == 0 {
			var resMsg string
			for i, v := range commandList {
				resMsg += strconv.Itoa(i+1) + " -> " + v.Use + " - " + v.Description + "\n"
			}
			s.ChannelMessageSend(m.ChannelID, resMsg)
		}
		return nil
	},
}

var clearCommand = command{
	Use:         "clear",
	Description: "Clear n messages",
	Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		if len(args) > 0 {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				log.Println(err)
			}
			if n < 0 {
				n = 0
			}
			n++ // add curent msg
			if n > msgIDList.Len() {
				n = msgIDList.Len()
			}
			for i := 0; i < n; i++ {
				s.ChannelMessageDelete(m.ChannelID, fmt.Sprintf("%v", msgIDList.Front().Value))
				msgIDList.Remove(msgIDList.Front())
			}
		}
		return nil
	},
}

var quoteCommand = command{
	Use:         "quote",
	Description: "send a quote",
	Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		var (
			csvQuote        = 1
			csvQuoteAuthor  = 2
			csvQuoteKeyword = 3
		)

		if len(args) > 0 {
			csvFile, err := os.Open("data/quotes.csv")
			if err != nil {
				log.Println(err)
				return err
			}
			defer csvFile.Close()

			r := csv.NewReader(csvFile)
			r.Comma = '|'

			keyword := strings.ToLower(args[rand.Intn(len(args))])
			quotes := make([]string, 0)

			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Println(err)
					return err
				}
				csvKeywords := strings.Fields(record[csvQuoteKeyword])
				for _, v := range csvKeywords {
					if v == keyword {
						msg := "\"" + record[csvQuote] + "\"" + "\n" + record[csvQuoteAuthor]
						quotes = append(quotes, msg)
					}
				}
			}
			if len(quotes) > 0 {
				s.ChannelMessageSend(m.ChannelID, quotes[rand.Intn(len(quotes))])
			}
		}
		return nil
	},
}

var printCommand = command{
	Use:         "print",
	Description: "print text",
	Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		var msg string
		for i := range args {
			msg += args[i]
			msg += " "
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil
	},
}

func init() {
	commandList = append(commandList, helpCommand)
	commandList = append(commandList, clearCommand)
	commandList = append(commandList, quoteCommand)

	commandList = append(commandList, printCommand)
}

//help - done
//clear n - done
//citate - done
//reputation
//play song
