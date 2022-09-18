// ====================================================================================================================
// Author: Robert Jiang
// Name: tracktivity-discord-bot
// File Name: discord.go
// Description: Implements the Activity struct, which stores all Tracktivity user information
// Version: 0.5.0.0
// ====================================================================================================================

package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	YES     string = "YES"
	NO      string = "NO"
	UNKNOWN string = "UNKNOWN"
	NA      string = "N/A"

	ONLINE    string = "ONLINE"
	OFFLINE   string = "OFFLINE"
	INVISIBLE string = "INVISIBLE"
	IDLE      string = "IDLE"
	DND       string = "DO NOT DISTURB"

	USER  string = "USER"
	BOT   string = "BOT"
	OR    string = " or "
	AGO   string = " ago"
	SINCE string = " since "
)

type Activity struct {
	user_id    string
	username   string
	nick       string
	user_type  string
	joined     time.Time
	roles      []string
	searches   uint64
	lastSearch time.Time
	views      uint64
	lastView   time.Time
	msgs       uint64
	lastMsg    time.Time
	rctns      uint64
	lastRctn   time.Time
	status     string
	lastStatus time.Time
	app        string
	appType    string
	lastApp    time.Time
	muted      string
	deafened   string
	stream     string
	video      string
	lastBoost  *time.Time
	timeout    *time.Time
}

func NewActivity(user_id string, username string, nick string, joined time.Time, roles []string, lastBoost *time.Time, timeout *time.Time) *Activity {
	return &Activity{
		user_id:    user_id,
		username:   username,
		nick:       nick,
		user_type:  NA,
		joined:     joined,
		roles:      roles,
		searches:   0,
		lastSearch: time.Time{},
		views:      0,
		lastView:   time.Time{},
		msgs:       0,
		lastMsg:    time.Time{},
		rctns:      0,
		lastRctn:   time.Time{},
		status:     UNKNOWN,
		lastStatus: time.Time{},
		app:        NA,
		appType:    NA,
		lastApp:    time.Time{},
		muted:      NA,
		deafened:   NA,
		stream:     NA,
		video:      NA,
		lastBoost:  lastBoost,
		timeout:    timeout,
	}
}

func UpdateUsername(activity *Activity, username string) {
	activity.username = username
}

func UpdateNick(activity *Activity, nick string) {
	activity.nick = nick
}

func UpdateUserType(activity *Activity, bot bool) {
	if bot {
		activity.user_type = BOT
	} else {
		activity.user_type = USER
	}
}

func UpdateJoined(activity *Activity, joined time.Time) {
	activity.joined = joined
}

func UpdateRoles(activity *Activity, roles []string) {
	activity.roles = roles
}

func UpdateSearches(activity *Activity) {
	activity.searches += 1
}

func UpdateLastSearch(activity *Activity) {
	activity.lastSearch = time.Now()
}

func UpdateViews(activity *Activity) {
	activity.views += 1
}

func UpdateLastView(activity *Activity) {
	activity.lastView = time.Now()
}

func UpdateMsgs(activity *Activity, msgs uint64) {
	if msgs <= 0 {
		activity.msgs = 0
		activity.lastMsg = time.Time{}
	} else {
		activity.msgs = msgs
	}
}

func UpdateLastMsg(activity *Activity) {
	activity.lastMsg = time.Now()
}

func UpdateRctns(activity *Activity, rctns uint64) {
	if rctns <= 0 {
		activity.rctns = 0
		activity.lastRctn = time.Time{}
	} else {
		activity.rctns = rctns
	}
}

func UpdateLastRctn(activity *Activity) {
	activity.lastRctn = time.Now()
}

func UpdateStatus(activity *Activity, status discordgo.Status) {
	switch status {
	case discordgo.StatusOnline:
		activity.status = ONLINE
	case discordgo.StatusOffline:
		activity.status = OFFLINE
	case discordgo.StatusInvisible:
		activity.status = INVISIBLE
	case discordgo.StatusIdle:
		activity.status = IDLE
	case discordgo.StatusDoNotDisturb:
		activity.status = DND
	default:
		activity.status = UNKNOWN
	}
}

func UpdateLastStatus(activity *Activity) {
	activity.lastStatus = time.Now()
}

func UpdateApp(activity *Activity, app string) {
	activity.app = app
}

func UpdateAppType(activity *Activity, appType string) {
	activity.appType = appType
}

func UpdateLastApp(activity *Activity, lastApp time.Time) {
	activity.lastApp = lastApp
}

func UpdateMuted(activity *Activity, muted bool) {
	if muted {
		activity.muted = YES
	} else {
		activity.muted = NO
	}
}

func UpdateDeafened(activity *Activity, deafened bool) {
	if deafened {
		activity.deafened = YES
	} else {
		activity.deafened = NO
	}
}

func UpdateStream(activity *Activity, stream bool) {
	if stream {
		activity.stream = YES
	} else {
		activity.stream = NO
	}
}

func UpdateVideo(activity *Activity, video bool) {
	if video {
		activity.video = YES
	} else {
		activity.video = NO
	}
}

func UpdateLastBoost(activity *Activity, lastBoost *time.Time) {
	activity.lastBoost = lastBoost
}

func UpdateTimeout(activity *Activity, timeout *time.Time) {
	activity.timeout = timeout
}

func GetJoinedTime(activity *Activity) string {
	return activity.joined.Local().Format(time.UnixDate)
}

func GetTimeSinceJoined(activity *Activity) string {
	return time.Since(activity.joined).String()
}

func GetLastSearchTime(activity *Activity) string {
	return activity.lastSearch.Local().Format(time.UnixDate)
}

func GetTimeSinceLastSearch(activity *Activity) string {
	return time.Since(activity.lastSearch).String()
}

func GetLastViewTime(activity *Activity) string {
	return activity.lastView.Local().Format(time.UnixDate)
}

func GetTimeSinceLastView(activity *Activity) string {
	return time.Since(activity.lastView).String()
}

func GetLastMsgTime(activity *Activity) string {
	return activity.lastMsg.Local().Format(time.UnixDate)
}

func GetTimeSinceLastMsg(activity *Activity) string {
	return time.Since(activity.lastMsg).String()
}

func GetLastRctnTime(activity *Activity) string {
	return activity.lastRctn.Local().Format(time.UnixDate)
}

func GetTimeSinceLastRctn(activity *Activity) string {
	return time.Since(activity.lastRctn).String()
}

func GetLastStatusTime(activity *Activity) string {
	return activity.lastStatus.Local().Format(time.UnixDate)
}

func GetTimeSinceLastStatus(activity *Activity) string {
	return time.Since(activity.lastStatus).String()
}

func GetLastAppTime(activity *Activity) string {
	return activity.lastApp.Local().Format(time.UnixDate)
}

func GetTimeSinceLastApp(activity *Activity) string {
	return time.Since(activity.lastApp).String()
}

func GetLastBoostTime(activity *Activity) string {
	if activity.lastBoost == nil {
		return NA
	}
	return activity.lastBoost.Format(time.UnixDate) + OR + time.Since(*activity.lastBoost).String()
}

func GetTimeoutDuration(activity *Activity) string {
	if activity.timeout == nil {
		return NA
	}

	duration := time.Since(*activity.timeout)
	if duration.Minutes() <= 0 {
		return duration.Abs().String()
	} else {
		return NA
	}
}
