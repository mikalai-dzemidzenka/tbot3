package main

import "C"
import (
	"fmt"
	"github.com/sampgo/sampgo"
	"strconv"
	"strings"
	"tbot2/tbot"
)

var tick uint = 1
var tickRate uint = 100

func init() {

	sampgo.On("goModeInit", func() bool {
		err := tbot.LoadData("tbot_dump")
		if err != nil {
			sampgo.Print(fmt.Sprintf("failed to load tbot_dump.json: %s", err.Error()))
		}
		for botNum, bot := range tbot.Bots {
			if !bot.IsSingle {
				tbot.Trs(botNum)
			}
		}
		return true
	})

	//sampgo.On("tick", func() {
	//	tick++
	//})

	sampgo.On("playerSpawn", func(p sampgo.Player) bool {
		if tbot.IsBot(p.ID) {
			botNum := tbot.Players[p.ID].BotNumber
			if tbot.IsNicksVisible {
				tbot.AttachBotNick(botNum)
			}
		}
		return true
	})

	sampgo.On("goModeExit", func() bool {
		tbot.SaveData("tbot_dump")
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
		case "veh":
			if len(tokens) < 4 {
				sampgo.SendClientMessage(p.ID, -1, "USAGE: /veh [VehicleID][Color1][Color2]")
				return true
			}
			carModel, err := strconv.Atoi(tokens[1])
			if err != nil {
				sampgo.SendClientMessage(p.ID, -1, "VehicleID must be a number")
				return true
			}

			if carModel < 400 || carModel > 611 {
				sampgo.SendClientMessage(p.ID, 0x0000FF, "Wrong vehicle ID input! (400-611)")
				return true
			}

			color1, err := strconv.Atoi(tokens[2])
			if err != nil {
				sampgo.SendClientMessage(p.ID, -1, "Color1 must be a number")
				return true
			}
			if color1 > 255 || color1 < 0 {
				sampgo.SendClientMessage(p.ID, 0x0000FF, "Wrong color input! (0-255)")
				return true
			}

			color2, err := strconv.Atoi(tokens[3])
			if err != nil {
				sampgo.SendClientMessage(p.ID, -1, "Color2 must be a number")
				return true
			}
			if color2 > 255 || color2 < 0 {
				sampgo.SendClientMessage(p.ID, 0x0000FF, "Wrong color input! (0-255)")
				return true
			}

			var x, y, z float32
			var angle float32
			sampgo.GetPlayerPos(p.ID, &x, &y, &z)
			sampgo.GetPlayerFacingAngle(p.ID, &angle)
			vehid := sampgo.CreateVehicle(carModel, x, y, z, angle, color1, color2, 60, false)
			sampgo.PutPlayerInVehicle(p.ID, vehid, 0)
			tbot.AddCar(vehid, tbot.CarInfo{CarModel: carModel, X: x, Y: y, Z: z, Angle: angle, Color: [2]int{color1, color2}})

		case "deleteveh":
			vehid := sampgo.GetPlayerVehicleID(p.ID)
			if vehid != 0 {
				tbot.RemoveCar(vehid)
				sampgo.DestroyVehicle(vehid)
			}
		case "tchain":
		// ??? think about design
		case "thelp":
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tlist /tkick /tgkick /tdel /tgdel /tnicks")
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tg /tbot /tsingle /tgsingle")
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tgbot /tsave /tload")
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
		case "tload":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFFFFF, "usage: /tload <filename>")
				return true
			}

			err := tbot.LoadData(tokens[1])
			if err != nil {
				sampgo.SendClientMessage(p.ID, 0xFFFFFF, "usage: /tload <filename>")
				return true
			}

			for botNum := range tbot.Bots {
				tbot.Trs(botNum)
			}

			sampgo.SendClientMessage(p.ID, 0xFFFFFF, "success!")
		case "tsave":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tsave <filename>")
				return true
			}

			tbot.SaveData(tokens[1])

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
			tbot.Tnicks()
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

			tbot.Tbinit(p.ID)
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
