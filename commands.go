// ====================================================================================================================
// Author: Robert Jiang
// Name: tracktivity-discord-bot
// File Name: commands.go
// Description: Stores all necessary command handlers for Tracktivity
// Version: 0.5.0.0
// ====================================================================================================================

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	HELP_COMMAND    string = "Tracktivity Commands\n\n1. \\*tracktivity <mention> - see Tracktivity data for <mention>.\n2. \\*help - see valid Tracktivity commands."
	UNKNOWN_COMMAND string = "Tracktivity does not recognize that command! Use *help to see the list of valid Tracktivity commands!"
)

func TracktivityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogFuncEntered("TracktivityCommand")

	splits := strings.Split(m.Content, " ")
	mention := splits[1]

	if mention[:2] == "<@" && mention[len(mention)-1] == '>' {
		user_id := mention[2 : len(mention)-1]
		activity, exists := users[user_id]

		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does not exist or has no Tracktivity profile yet, please check and try again later!", mention))
			return
		}

		UpdateViews(activity)
		UpdateSearches(users[m.Author.ID])
		UpdateLastSearch(users[m.Author.ID])

		username := activity.username
		if username == "#" {
			username = NA
		}
		nick := activity.nick
		if nick == "" {
			nick = NA
		}
		user_type := activity.user_type
		var joined string
		if activity.joined.IsZero() {
			joined = NA
		} else {
			joined = GetJoinedTime(activity) + OR + GetTimeSinceJoined(activity) + AGO
		}
		roles_list := ""
		for i, role_id := range activity.roles {
			if i != len(activity.roles)-1 {
				role := roles[role_id]
				if role == "" {
					roles_list += NA + " | "
				} else {
					roles_list += role + " | "
				}

			} else {
				role := roles[role_id]
				if role == "" {
					roles_list += NA
				} else {
					roles_list += role
				}
			}
		}
		if roles_list == "" {
			roles_list = NA
		}
		searches := activity.searches
		var searchTime string
		if activity.lastSearch.IsZero() {
			searchTime = NA
		} else {
			searchTime = GetLastSearchTime(activity) + OR + GetTimeSinceLastSearch(activity) + AGO
		}
		views := activity.views
		var viewTime string
		if activity.lastView.IsZero() {
			viewTime = NA
		} else {
			viewTime = GetLastViewTime(activity) + OR + GetTimeSinceLastView(activity) + AGO
		}
		totalMsgs := activity.msgs
		var msgTime string
		if activity.lastMsg.IsZero() {
			msgTime = NA
		} else {
			msgTime = GetLastMsgTime(activity) + OR + GetTimeSinceLastMsg(activity) + AGO
		}
		totalRctns := activity.rctns
		var rctnTime string
		if activity.lastRctn.IsZero() {
			rctnTime = NA
		} else {
			rctnTime = GetLastRctnTime(activity) + OR + GetTimeSinceLastRctn(activity) + AGO
		}
		status := activity.status
		var statusTime string
		if activity.lastStatus.IsZero() {
			statusTime = NA
		} else {
			statusTime = GetLastStatusTime(activity) + OR + GetTimeSinceLastStatus(activity) + AGO
		}
		app := activity.app
		appType := activity.appType
		var appTime string
		if activity.lastApp.IsZero() {
			appTime = NA
		} else {
			appTime = GetLastAppTime(activity) + OR + GetTimeSinceLastApp(activity) + AGO
		}
		var lastApp string
		if app == NA || appType == NA || appTime == NA {
			lastApp = NA
		} else {
			lastApp = appType + SINCE + appTime
		}
		muted := activity.muted
		deafened := activity.deafened
		stream := activity.stream
		video := activity.video
		boostTime := GetLastBoostTime(activity)
		timeout := GetTimeoutDuration(activity)

		report := ""
		report += mention + "'s Tracktivity @ " + time.Now().Local().Format(time.UnixDate) + ":\n\n"
		report += "Username - " + username + "\n"
		report += "Nickname - " + nick + "\n"
		report += "User Type - " + user_type + "\n"
		report += "Joined - " + joined + "\n"
		report += "Role(s) - " + roles_list + "\n\n"
		report += "Tracktivity Searches - " + fmt.Sprint(searches) + "\n"
		report += "Last Search - " + searchTime + "\n"
		report += "Tracktivity Views - " + fmt.Sprint(views) + "\n"
		report += "Last View - " + viewTime + "\n"
		report += "Total Messages - " + fmt.Sprint(totalMsgs) + "\n"
		report += "Last Message - " + msgTime + "\n"
		report += "Total Reactions - " + fmt.Sprint(totalRctns) + "\n"
		report += "Last Reaction - " + rctnTime + "\n\n"
		report += "Status - " + status + "\n"
		report += "Last Status - " + statusTime + "\n"
		report += "App - " + app + "\n"
		report += "Last App - " + lastApp + "\n"
		report += "Muted - " + muted + "\n"
		report += "Deafened - " + deafened + "\n"
		report += "Sharing Screen - " + stream + "\n"
		report += "Sharing Video - " + video + "\n"
		report += "Boosting - " + boostTime + "\n"
		report += "Timeout - " + timeout + "\n"

		s.ChannelMessageSend(m.ChannelID, report)

		UpdateLastView(activity)
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not a valid mention, please check and try again!", mention))
	}

}

func HelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogFuncEntered("HelpCommand")

	s.ChannelMessageSend(m.ChannelID, HELP_COMMAND)
}

func UnknownCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogFuncEntered("UnknownCommand")

	s.ChannelMessageSend(m.ChannelID, UNKNOWN_COMMAND)
}
