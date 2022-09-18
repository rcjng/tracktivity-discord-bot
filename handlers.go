// ====================================================================================================================
// Author: Robert Jiang
// Name: tracktivity-discord-bot
// File Name: handlers.go
// Description: Stores all necessary event handlers for Tracktivity
// Version: 0.5.0.0
// ====================================================================================================================

package main

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	users       = make(map[string]*Activity)
	msgs        = make(map[string]*Activity)
	rctns       = make(map[string]*Activity)
	roles       = make(map[string]string)
	recentMsgs  []string
	recentRctns []string
	RECENT_CAP  = 10000
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogFuncEntered("MessageCreate")

	user_id := m.Author.ID
	msg_id := m.ID

	// Ignore if message is from Tracktivity Bot
	if user_id == s.State.User.ID {
		return
	}

	activity, exists := users[user_id]
	if !exists {
		activity = NewActivity(user_id, m.Author.String(), m.Member.Nick, m.Member.JoinedAt, m.Member.Roles, m.Member.PremiumSince, m.Member.CommunicationDisabledUntil)
		users[user_id] = activity
	}

	UpdateUserType(activity, m.Author.Bot)
	UpdateMsgs(activity, activity.msgs+1)
	UpdateLastMsg(activity)

	if len(recentMsgs) < RECENT_CAP {
		recentMsgs = append(recentMsgs, msg_id)
	} else {
		delete(msgs, recentMsgs[0])
		recentMsgs = append(recentMsgs[1:], msg_id)
	}

	msgs[msg_id] = activity

	if strings.HasPrefix(m.Content, "*") {
		if strings.HasPrefix(m.Content, "*tracktivity") && len(strings.Split(m.Content, " ")) == 2 {
			TracktivityCommand(s, m)
		} else if strings.HasPrefix(m.Content, "*help") && len(strings.Split(m.Content, " ")) == 1 {
			HelpCommand(s, m)
		} else {
			UnknownCommand(s, m)
		}
	}
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	LogFuncEntered("MessageDelete")

	msg_id := m.Message.ID

	for i, id := range recentMsgs {
		if msg_id == id {
			activity := msgs[msg_id]
			UpdateMsgs(activity, activity.msgs-1)

			delete(msgs, msg_id)
			recentMsgs = append(recentMsgs[:i], recentMsgs[i+1:]...)

			break
		}
	}
}

func messageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk) {
	LogFuncEntered("MessageDeleteBulk")

	for _, msg_id := range m.Messages {
		for i, id := range recentMsgs {
			if msg_id == id {
				activity := msgs[msg_id]
				UpdateMsgs(activity, activity.msgs-1)

				delete(msgs, msg_id)
				recentMsgs = append(recentMsgs[:i], recentMsgs[i+1:]...)

				break
			}
		}
	}
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	LogFuncEntered("MessageReactionAdd")

	user_id := r.UserID
	rctn_id := r.MessageID

	// if reaction is from bot
	if user_id == s.State.User.ID {
		return
	}

	activity, exists := users[user_id]
	if !exists {
		activity = NewActivity(user_id, r.Member.User.String(), r.Member.Nick, r.Member.JoinedAt, r.Member.Roles, r.Member.PremiumSince, r.Member.CommunicationDisabledUntil)
		users[user_id] = activity
	}

	UpdateUserType(activity, r.Member.User.Bot)
	UpdateRctns(activity, activity.rctns+1)
	UpdateLastRctn(activity)

	if len(recentRctns) < RECENT_CAP {
		recentRctns = append(recentRctns, rctn_id)
	} else {
		delete(rctns, recentMsgs[0])
		recentRctns = append(recentRctns[1:], rctn_id)
	}

	rctns[rctn_id] = activity
}

func messageReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	LogFuncEntered("MessageReactionRemove")

	rctn_id := r.MessageID

	for i, id := range recentRctns {
		if rctn_id == id {
			activity := rctns[rctn_id]
			UpdateRctns(activity, activity.rctns-1)

			delete(rctns, rctn_id)
			recentRctns = append(recentRctns[:i], recentRctns[i+1:]...)

			break
		}
	}
}

func presenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	LogFuncEntered("PresenceUpdate")

	user_id := p.User.ID

	activity, exists := users[user_id]
	if !exists {
		activity = NewActivity(user_id, p.Presence.User.String(), "", time.Time{}, make([]string, 0), nil, nil)
		users[user_id] = activity
	}

	UpdateStatus(activity, p.Status)
	UpdateLastStatus(activity)

	activities := p.Activities

	if len(activities) > 0 {
		UpdateApp(activity, activities[0].Name)
		UpdateLastApp(activity, activities[0].CreatedAt)
		switch activities[0].Type {
		case 0:
			UpdateAppType(activity, "Gaming")
		case 1:
			UpdateAppType(activity, "Streaming")
		case 2:
			UpdateAppType(activity, "Listening")
		case 3:
			UpdateAppType(activity, "Watching")
		case 4:
			UpdateAppType(activity, "Custom Activity-ing")
		case 5:
			UpdateAppType(activity, "Competing")
		default:
			UpdateAppType(activity, UNKNOWN)
		}
	} else {
		UpdateApp(activity, NA)
		UpdateAppType(activity, NA)
		UpdateLastApp(activity, time.Time{})
	}
}

func userUpdate(s *discordgo.Session, u *discordgo.UserUpdate) {
	LogFuncEntered("UserUpdate")
	user_id := u.User.ID

	activity, exists := users[user_id]
	if !exists {
		activity = NewActivity(user_id, u.String(), "", time.Time{}, make([]string, 0), nil, nil)
		users[user_id] = activity
	}

	UpdateUsername(activity, u.User.String())
}

func voiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	LogFuncEntered("VoiceStateUpdate")

	user_id := v.VoiceState.Member.User.ID

	activity, exists := users[user_id]
	if !exists {
		activity = NewActivity(user_id, v.Member.User.String(), v.Member.Nick, v.Member.JoinedAt, v.Member.Roles, v.Member.PremiumSince, v.Member.CommunicationDisabledUntil)
		users[user_id] = activity
	}

	UpdateUserType(activity, v.Member.User.Bot)
	UpdateMuted(activity, v.VoiceState.SelfMute || v.VoiceState.Mute)
	UpdateDeafened(activity, v.VoiceState.SelfDeaf || v.VoiceState.Deaf)
	UpdateStream(activity, v.VoiceState.SelfStream)
	UpdateVideo(activity, v.VoiceState.SelfVideo)
}

func guildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	LogFuncEntered("GuildMemberAdd")

	user_id := g.User.ID
	username := g.User.String()
	nick := g.Member.Nick
	bot := g.Member.User.Bot
	joined := g.Member.JoinedAt
	roles := g.Member.Roles
	lastBoost := g.Member.PremiumSince
	timeout := g.Member.CommunicationDisabledUntil

	activity, exists := users[user_id]
	if !exists {
		users[user_id] = NewActivity(user_id, username, nick, joined, roles, lastBoost, timeout)
	} else {
		UpdateUsername(activity, username)
		UpdateNick(activity, nick)
		UpdateUserType(activity, bot)
		UpdateRoles(activity, roles)
		UpdateLastBoost(activity, lastBoost)
		UpdateTimeout(activity, timeout)
	}
}

func guildMemberRemove(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	LogFuncEntered("GuildMemberRemove")

	delete(users, g.User.ID)
}

func guildMemberUpdate(s *discordgo.Session, g *discordgo.GuildMemberUpdate) {
	LogFuncEntered("GuildMemberUpdate")

	user_id := g.User.ID
	username := g.User.String()
	nick := g.Member.Nick
	bot := g.Member.User.Bot
	joined := g.Member.JoinedAt
	roles := g.Member.Roles
	lastBoost := g.Member.PremiumSince
	timeout := g.Member.CommunicationDisabledUntil

	activity, exists := users[user_id]
	if !exists {
		users[user_id] = NewActivity(user_id, username, nick, joined, roles, lastBoost, timeout)
	} else {
		UpdateUsername(activity, username)
		UpdateNick(activity, nick)
		UpdateUserType(activity, bot)
		UpdateRoles(activity, roles)
		UpdateLastBoost(activity, lastBoost)
		UpdateTimeout(activity, timeout)
	}
}

func guildRoleCreate(s *discordgo.Session, g *discordgo.GuildRoleCreate) {
	LogFuncEntered("GuildRoleCreate")

	roles[g.GuildRole.Role.ID] = g.GuildRole.Role.Name
}

func guildRoleUpdate(s *discordgo.Session, g *discordgo.GuildRoleUpdate) {
	LogFuncEntered("GuildRoleUpdate")

	delete(roles, g.GuildRole.Role.ID)
	roles[g.GuildRole.Role.ID] = g.GuildRole.Role.Name
}

func guildRoleDelete(s *discordgo.Session, g *discordgo.GuildRoleDelete) {
	LogFuncEntered("GuildRoleDelete")

	delete(roles, g.RoleID)
}
