// ====================================================================================================================
// Author: Robert Jiang
// Name: tracktivity-discord-bot
// File Name: discord.go
// Description: Launches the Tracktivity Discord Bot with the necessary intents, handlers, and configuration.
// Version: 0.5.0.0
// ====================================================================================================================

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token   = "" // insert token here
	version = "v0.5.0.0"
)

func ConnectToDiscord() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageDelete)
	dg.AddHandler(messageDeleteBulk)
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageReactionRemove)
	dg.AddHandler(presenceUpdate)
	dg.AddHandler(userUpdate)
	dg.AddHandler(voiceStateUpdate)
	dg.AddHandler(guildMemberAdd)
	dg.AddHandler(guildMemberRemove)
	dg.AddHandler(guildMemberUpdate)
	dg.AddHandler(guildRoleCreate)
	dg.AddHandler(guildRoleUpdate)
	dg.AddHandler(guildRoleDelete)

	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildMessageTyping |
		discordgo.IntentsGuildPresences |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildVoiceStates |
		discordgo.IntentGuilds

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	log.Println("Tracktivity Bot " + fmt.Sprint(version) + " Launched!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func LogFuncEntered(fn string) {
	log.Println(fn + " handler entered!")
}
