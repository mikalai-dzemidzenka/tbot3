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
	t := tbot.New()

	sampgo.On("goModeInit", func() bool {
		t.ExportData()
		t.InitRun()
		return true
	})

	sampgo.On("goModeExit", func() bool {
		t.SaveData()
		return true
	})

	sampgo.On("playerConnect", func(p sampgo.Player) bool {
		if t.IsBot(p.ID) {
			t.ConnectBot(p.ID)
		} else {
			t.ConnectPlayer(p.ID)
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

			t.DisconnectBot(botID)
			return true
		}

		if t.IsRecording(p.ID) {
			t.StopRecording(p.ID)
		}

		return true
	})

	sampgo.On("playerCommandText", func(p sampgo.Player, cmd string) bool {
		sampgo.SendClientMessage(p.ID, 0xff0000, cmd)
		tokens := strings.Split(cmd, " ")
		switch tokens[0][1:] {
		case "tsave":
			t.SaveData()
		case "tchain":
		// ??? think about design
		case "thelp":
			sampgo.SendClientMessage(p.ID, 0xFF0000, "/tlist /tkick /tgkick /tdel /tgdel /tnicks")
			sampgo.SendClientMessage(p.ID, 0xFF0000, "/tg /tbot /tsingle /tgsingle /tgbot /tgsingle")
		case "tlist":
			botList := t.Tlist()
			for _, bot := range botList {
				sampgo.SendClientMessage(p.ID, 0xFF0000, bot)
			}
		case "tkick":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			t.Tkick(botNum)
		case "tgkick":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			t.Tgkick(groupID)
		case "tdel":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			if tokens[1] == "all" {
				t.Tdelall()
			}
			botNum, _ := strconv.Atoi(tokens[1])

			t.Tdel(botNum)
		case "tgdel":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			t.Tgdel(groupID)
		case "tnicks":
			//TODO
		case "trs":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			t.Trs(botNum)
		case "tgrs":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			t.Tgrs(groupID)
		case "tg":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, fmt.Sprintf("player group: %d", t.GetPlayerGroup(p.ID)))
				return true
			}
			group, _ := strconv.Atoi(tokens[1])

			t.SetPlayerGroup(p.ID, group)
		case "tbot":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			t.TBot(p.ID, botNum, false)
		case "tsingle":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			t.TBot(p.ID, botNum, true)
		case "tgbot":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			players := t.GetPlayersInGroup(groupID)

			for _, player := range players {
				botNum, ok := t.GetFreeBotNum()
				if !ok {
					sampgo.SendClientMessage(p.ID, 0xFF0000, "no more free bots")
					return true
				}
				t.TBot(player, botNum, false)
			}
		case "tgsingle":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}
			groupID, _ := strconv.Atoi(tokens[1])

			players := t.GetPlayersInGroup(groupID)

			for _, player := range players {
				botNum, ok := t.GetFreeBotNum()
				if !ok {
					sampgo.SendClientMessage(p.ID, 0xFF0000, "no more free bots")
					return true
				}
				t.TBot(player, botNum, true)
			}
		// BOT ONLY
		case "tbinit":
			if !sampgo.IsPlayerNPC(p.ID) {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}

			playback, recType, isSingle, groupID := t.Tbinit(p.ID)
			sampgo.SendClientMessage(p.ID, 0x000001, fmt.Sprintf("%s %d %d %d", playback, recType, isSingle, groupID))
		case "tbready":
			if !sampgo.IsPlayerNPC(p.ID) {
				sampgo.SendClientMessage(p.ID, 0xFF0000, "roflan")
				return true
			}

			t.Tbready(p.ID)
		default:
			return false
		}
		return true
	})
}

func main() {}
