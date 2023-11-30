package main

import "C"
import (
	"fmt"
	"github.com/sampgo/sampgo"
	"strconv"
	"strings"
	"tbot2/tbot"
)

func init() {

	sampgo.On("goModeInit", func() bool {
		tbot.ExportData()
		tbot.InitRun()
		return true
	})

	sampgo.On("goModeExit", func() bool {
		tbot.SaveData()
		return true
	})

	sampgo.On("playerConnect", func(p sampgo.Player) bool {
		if tbot.IsBot(p.ID) {
			tbot.ConnectBot(p.ID)
		} else {
			tbot.ConnectPlayer(p.ID)
		}
		return true
	})

	sampgo.On("playerDisconnect", func(p sampgo.Player, reason int) bool {
		name := p.GetName()

		if strings.HasPrefix(name, tbot.BotPrefix) && sampgo.IsPlayerNPC(p.ID) {
			botID, err := strconv.Atoi(strings.TrimPrefix(name, tbot.BotPrefix))
			if err != nil {
				_ = sampgo.Print(fmt.Sprintf("get bot id: %v", err))
				return true
			}

			tbot.DisconnectBot(botID)
			return true
		}

		if tbot.IsRecording(p.ID) {
			tbot.StopRecording(p.ID)
		}

		return true
	})

	sampgo.On("playerCommandText", func(p sampgo.Player, cmd string) bool {
		tokens := strings.Split(cmd, " ")
		command := tokens[0][1:]
		switch command {
		default: // unknown command
			return false
		case "tsave":
			tbot.SaveData()
		case "tchain":
		// ??? think about design
		case "thelp":
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tlist /tkick /tgkick /tdel /tgdel /tnicks")
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tg /tbot /tsingle /tgsingle /tgbot")
		case "tlist":
			botList := tbot.Tlist()
			for _, bot := range botList {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, bot)
			}
		case "tkick":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tkick <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.Tkick(botNum)
		case "tgkick":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tgkick <group_id>")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			tbot.Tgkick(groupID)
		case "tdel":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tdel <bot_id>")
				return true
			}
			if tokens[1] == "all" {
				tbot.Tdelall()
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.Tdel(botNum)
		case "tgdel":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tgdel <group_id>")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			tbot.Tgdel(groupID)
		case "tnicks":
			//TODO
		case "trs":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /trs <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.Trs(botNum)
		case "tgrs":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tgrs <group_id>")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			tbot.Tgrs(groupID)
		case "tg":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, fmt.Sprintf("player group: %d", tbot.GetPlayerGroup(p.ID)))
				return true
			}
			group, _ := strconv.Atoi(tokens[1])

			tbot.SetPlayerGroup(p.ID, group)
		case "tbot":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tbot <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.TBot(p.ID, botNum, false)
		case "tsingle":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tsingle <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.TBot(p.ID, botNum, true)
		case "tgbot":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tgbot <group_id>")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			players := tbot.GetPlayersInGroup(groupID)

			for _, player := range players {
				botNum, ok := tbot.GetFreeBotNum()
				if !ok {
					sampgo.SendClientMessage(p.ID, 0xFFAA00, "no more free bots")
					return true
				}
				tbot.TBot(player, botNum, false)
			}
		case "tgsingle":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tgsingle <group_id>")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			players := tbot.GetPlayersInGroup(groupID)

			for _, player := range players {
				botNum, ok := tbot.GetFreeBotNum()
				if !ok {
					sampgo.SendClientMessage(p.ID, 0xFFAA00, "no more free bots")
					return true
				}
				tbot.TBot(player, botNum, true)
			}
		// BOT ONLY
		case "tbinit":
			if !sampgo.IsPlayerNPC(p.ID) {
				return false
			}

			playback, recType, isSingle, groupID := tbot.Tbinit(p.ID)
			sampgo.SendClientMessage(p.ID, 0x000001, fmt.Sprintf("%s %d %d %d", playback, recType, isSingle, groupID))
		case "tbready":
			if !sampgo.IsPlayerNPC(p.ID) {
				return false
			}

			tbot.Tbready(p.ID)
		}
		return true
	})
}

func main() {}
