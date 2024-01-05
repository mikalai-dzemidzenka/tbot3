package main

import "C"
import (
	"fmt"
	"github.com/sampgo/sampgo"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
	"tbot3/tbot"
)

const (
	DialogAttachIndex = iota
	DialogAttachIndexSelection
	DialogAttachEditReplace
	DialogAttachModelSelection
	DialogAttachBoneSelection
)

func clamp[T constraints.Float | constraints.Integer](f, low, high T) T {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

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

	sampgo.On("playerSpawn", func(p sampgo.Player) bool {
		if tbot.IsBot(p.ID) {
			botNum := tbot.Players[p.ID].BotNumber
			if tbot.IsNicksVisible {
				tbot.AttachBotNick(botNum)
			}
		}
		return true
	})

	sampgo.On("playerEditAttachedObject", func(p sampgo.Player, response int, index int, modelid int, boneid int,
		offx float32, offy float32, offz float32,
		rx float32, ry float32, rz float32,
		sx float32, sy float32, sz float32) bool {

		const sizeLimit = 20
		sx = clamp(sx, -sizeLimit, sizeLimit)
		sy = clamp(sy, -sizeLimit, sizeLimit)
		sz = clamp(sz, -sizeLimit, sizeLimit)

		sampgo.SendClientMessage(p.ID, -1, fmt.Sprintf("SetPlayerAttachedObject(playerid,%d,%d,%d,%g,%g,%g,%g,%g,%g,%g,%g,%g)",
			index, modelid, boneid, offx, offy, offz, rx, ry, rz, sx, sy, sz))

		sampgo.SetPlayerAttachedObject(p.ID, index, modelid, boneid, offx, offy, offz, rx, ry, rz, sx, sy, sz, 0, 0)
		tbot.Players[p.ID].Attachments[index] = tbot.Attachment{Index: index, Modelid: modelid, Boneid: boneid, Offx: offx, Offy: offy, Offz: offz, Rx: rx, Ry: ry, Rz: rz, Sx: sx, Sy: sy, Sz: sz}

		sampgo.SendClientMessage(p.ID, -1, "You finished editing an attached object")

		return true
	})

	sampgo.On("dialogResponse", func(p sampgo.Player, dialogid int, response int, listitem int, inputtext string) bool {
		switch dialogid {
		case DialogAttachIndexSelection:
			if response == 0 {
				return true
			}

			if sampgo.IsPlayerAttachedObjectSlotUsed(p.ID, listitem) {
				sampgo.ShowPlayerDialog(p.ID, DialogAttachEditReplace, sampgo.DialogStyleMsgbox,
					"{FF0000}Attachment Modification",
					"Do you wish to edit the attachment in that slot, or delete it?",
					"Edit", "Delete")
			} else {
				var sb strings.Builder
				for _, att := range tbot.DefaultAttachments {
					sb.WriteString(att.Name + "\n")
				}
				sampgo.ShowPlayerDialog(p.ID, DialogAttachModelSelection, sampgo.DialogStyleList,
					"{FF0000}Attachment Modification - Model Selection", sb.String(), "Select", "Cancel")
			}

			tbot.Players[p.ID].CurrentAttachmentIndex = listitem
		case DialogAttachEditReplace:
			idx := tbot.Players[p.ID].CurrentAttachmentIndex
			if response == 1 {
				sampgo.EditAttachedObject(p.ID, idx)
			} else {
				sampgo.RemovePlayerAttachedObject(p.ID, idx)
				tbot.Players[p.ID].Attachments[idx] = tbot.Attachment{}
			}
			tbot.Players[p.ID].CurrentAttachmentIndex = tbot.NoAttachments
		case DialogAttachModelSelection:
			if response == 0 {
				tbot.Players[p.ID].CurrentAttachmentIndex = tbot.NoAttachments
				return true
			}
			// TODO ???
			//if(GetPVarInt(playerid, "AttachmentUsed") == 1) EditAttachedObject(playerid, listitem);
			// else

			tbot.Players[p.ID].CurrentAttachmentModel = tbot.DefaultAttachments[listitem].ID

			var sb strings.Builder
			for _, bone := range tbot.Bones {
				sb.WriteString(bone + "\n")
			}
			sampgo.ShowPlayerDialog(p.ID, DialogAttachBoneSelection, sampgo.DialogStyleList,
				"{FF0000}Attachment Modification - Bone Selection", sb.String(), "Select", "Cancel")

		case DialogAttachBoneSelection:
			if response == 1 {
				sampgo.SetPlayerAttachedObject(p.ID,
					tbot.Players[p.ID].CurrentAttachmentIndex,
					tbot.Players[p.ID].CurrentAttachmentModel,
					listitem+1,
					0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, 0, 0)
				sampgo.EditAttachedObject(p.ID, tbot.Players[p.ID].CurrentAttachmentIndex)
				sampgo.SendClientMessage(p.ID, -1, fmt.Sprintf("Editing slot %d. Hint: Use {FFFF00}~k~~PED_SPRINT~{FFFFFF} to look around.", tbot.Players[p.ID].CurrentAttachmentIndex))
			}

			tbot.Players[p.ID].CurrentAttachmentIndex = tbot.NoAttachments
			tbot.Players[p.ID].CurrentAttachmentModel = tbot.NoAttachments

		default:
			return false
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
		case "ta": // TODO show dialog same
			var sb strings.Builder
			for i := 0; i < sampgo.MaxPlayerAttachedObjects; i++ {
				if sampgo.IsPlayerAttachedObjectSlotUsed(p.ID, i) {
					sb.WriteString(strconv.Itoa(i) + " (Used)\n")
				} else {
					sb.WriteString(strconv.Itoa(i) + "\n")
				}
			}

			sampgo.ShowPlayerDialog(p.ID, DialogAttachIndexSelection, sampgo.DialogStyleList, "{FF0000}Attachment Modification - Index Selection", sb.String(), "Select", "Cancel")
		case "tacopy":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /tacopy <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.Tacopy(p.ID, botNum)
		case "taclear":
			if len(tokens) < 2 {
				sampgo.SendClientMessage(p.ID, 0xFFAA00, "usage: /taclear <bot_id>")
				return true
			}
			botNum, _ := strconv.Atoi(tokens[1])

			tbot.Taclear(botNum)
		case "tchain":
		// ??? think about design
		case "thelp":
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tlist /tkick /tgkick /tdel /tgdel /tnicks")
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tg /tbot /tsingle /tgsingle")
			sampgo.SendClientMessage(p.ID, 0xFFAA00, "/tgbot /tsave /tload /ta /tacopy /taclear")
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
