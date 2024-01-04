package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
	"os"
)

var Players = make(map[playerID]*PlayerInfo)

const (
	DefaultGroupID       = 0
	NoRecordingBotNumber = -1
	BotNotConnected      = -1
	NoCar                = -1
)

type PlayerInfo struct {
	ID                  int
	PlayerGroupID       int
	RecordingBotNumber  int
	RecordingPlayerType int
	*BotInfo
}

func GetPlayerGroup(id int) int {
	player, ok := Players[id]
	if !ok {
		return 0
	}
	return player.PlayerGroupID
}

func SetPlayerGroup(id int, groupID int) {
	Players[id].PlayerGroupID = groupID
}

func ConnectPlayer(id int) {
	Players[id] = &PlayerInfo{
		ID:                  id,
		PlayerGroupID:       DefaultGroupID,
		RecordingBotNumber:  NoRecordingBotNumber,
		RecordingPlayerType: sampgo.PlayerRecordingTypeNone,
	}
}

func GetPlayersInGroup(groupID int) []int {
	var res []int
	for _, p := range Players {
		if p.BotInfo != nil {
			continue
		}

		if p.PlayerGroupID == groupID {
			res = append(res, p.ID)
		}
	}
	return res
}

func IsRecording(id int) bool {
	p, ok := Players[id]
	if !ok {
		return false
	}
	return p.RecordingPlayerType != sampgo.PlayerRecordingTypeNone
}

func IsRecordingConflict(botNum int) bool {
	bot, ok := Bots[botNum]
	if !ok {
		return false
	}
	return bot.recordingPlayerID != NoRecordingBotNumber
}

func StartRecording(id int, botNum int, isSingle bool) {
	if IsRecordingConflict(botNum) {
		sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("someone else is recording bot %d", botNum))
		return
	}

	if IsBotConnected(botNum) {
		sampgo.Kick(Bots[botNum].id)
	}

	Players[id].RecordingBotNumber = botNum

	sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("start recording bot %d", botNum))

	state := sampgo.GetPlayerState(id)
	switch state {
	case sampgo.PlayerStateDriver, sampgo.PlayerStatePassenger:
		vehID := sampgo.GetPlayerVehicleID(id)
		Bots[botNum] = &BotInfo{
			BotNumber:      botNum,
			BotGroupID:     Players[id].PlayerGroupID,
			Skin:           sampgo.GetPlayerSkin(id),
			CarInfo:        Vehs[vehID],
			SeatID:         sampgo.GetPlayerVehicleSeat(id),
			IsSingle:       isSingle,
			State:          sampgo.PlayerRecordingTypeDriver,
			BotRuntimeInfo: NewBotRuntimeInfo(vehID),
		}

		Players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeDriver
		sampgo.StartRecordingPlayerData(id, sampgo.PlayerRecordingTypeDriver, fmt.Sprintf("tbotcar%d", botNum))
	default:
		Bots[botNum] = &BotInfo{
			BotNumber:      botNum,
			BotGroupID:     Players[id].PlayerGroupID,
			Skin:           sampgo.GetPlayerSkin(id),
			CarInfo:        CarInfo{},
			SeatID:         0,
			IsSingle:       isSingle,
			State:          sampgo.PlayerRecordingTypeOnfoot,
			BotRuntimeInfo: NewBotRuntimeInfo(NoCar),
		}

		Players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeOnfoot
		sampgo.StartRecordingPlayerData(id, sampgo.PlayerRecordingTypeOnfoot, fmt.Sprintf("tbotfoot%d", botNum))
	}
}

func StopRecording(id int) {
	botNum := Players[id].RecordingBotNumber

	sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("stop recording %d bot", botNum))
	sampgo.StopRecordingPlayerData(id)
	os.Remove(fmt.Sprintf("npcmodes/recordings/tbotfoot%d.rec", botNum))
	os.Remove(fmt.Sprintf("npcmodes/recordings/tbotcar%d.rec", botNum))

	switch Players[id].RecordingPlayerType {
	case sampgo.PlayerRecordingTypeOnfoot:
		os.Rename(
			fmt.Sprintf("scriptfiles/tbotfoot%d.rec", botNum),
			fmt.Sprintf("npcmodes/recordings/tbotfoot%d.rec", botNum),
		)
		if !Bots[botNum].IsSingle {
			botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
			sampgo.ConnectNPC(botName, "tbot")
		}
	case sampgo.PlayerRecordingTypeDriver:
		os.Rename(
			fmt.Sprintf("scriptfiles/tbotcar%d.rec", botNum),
			fmt.Sprintf("npcmodes/recordings/tbotcar%d.rec", botNum),
		)
		if !Bots[botNum].IsSingle {
			botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
			sampgo.ConnectNPC(botName, "tbot")
		}
	}

	Bots[botNum].recordingPlayerID = NoRecordingBotNumber
	Players[id].RecordingBotNumber = NoRecordingBotNumber
	Players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeNone
}
