package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
	"os"
)

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

func (t *T) GetPlayerGroup(id int) int {
	player, ok := t.players[id]
	if !ok {
		return 0
	}
	return player.PlayerGroupID
}

func (t *T) SetPlayerGroup(id int, groupID int) {
	t.players[id].PlayerGroupID = groupID
}

func (t *T) ConnectPlayer(id int) {
	t.players[id] = &PlayerInfo{
		ID:                  id,
		PlayerGroupID:       DefaultGroupID,
		RecordingBotNumber:  NoRecordingBotNumber,
		RecordingPlayerType: sampgo.PlayerRecordingTypeNone,
	}
}

func (t *T) IsRecording(id int) bool {
	p, ok := t.players[id]
	if !ok {
		return false
	}
	return p.RecordingPlayerType != sampgo.PlayerRecordingTypeNone
}

func (t *T) IsRecordingConflict(botNum int) bool {
	bot, ok := t.bots[botNum]
	if !ok {
		return false
	}
	return bot.RecordingPlayerID != NoRecordingBotNumber
}

func (t *T) StartRecording(id int, botNum int, isSingle bool) {
	if t.IsRecordingConflict(botNum) {
		sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("someone else is recording bot %d", botNum))
		return
	}

	if t.IsBotConnected(botNum) {
		sampgo.Kick(t.bots[botNum].ID)
	}

	t.players[id].RecordingBotNumber = botNum

	sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("start recording bot %d", botNum))

	state := sampgo.GetPlayerState(id)
	switch state {
	case sampgo.PlayerStateDriver, sampgo.PlayerStatePassenger:
		t.bots[botNum] = &BotInfo{
			ID:                BotNotConnected,
			BotGroupID:        t.players[id].PlayerGroupID,
			Skin:              sampgo.GetPlayerSkin(id),
			Car:               sampgo.GetPlayerVehicleID(id),
			SeatID:            sampgo.GetPlayerVehicleSeat(id),
			IsSingle:          isSingle,
			RecordingPlayerID: id,
		}

		t.players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeDriver
		sampgo.StartRecordingPlayerData(id, sampgo.PlayerRecordingTypeDriver, fmt.Sprintf("tbotcar%d", botNum))
	default:
		t.bots[botNum] = &BotInfo{
			ID:                BotNotConnected,
			BotGroupID:        t.players[id].PlayerGroupID,
			Skin:              sampgo.GetPlayerSkin(id),
			Car:               NoCar,
			SeatID:            0,
			IsSingle:          isSingle,
			RecordingPlayerID: id,
		}

		t.players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeOnfoot
		sampgo.StartRecordingPlayerData(id, sampgo.PlayerRecordingTypeOnfoot, fmt.Sprintf("tbotfoot%d", botNum))
	}
}

func (t *T) StopRecording(id int) {
	botNum := t.players[id].RecordingBotNumber

	sampgo.SendClientMessage(id, 0xFF0000, fmt.Sprintf("stop recording %d bot", botNum))
	sampgo.StopRecordingPlayerData(id)
	os.Remove(fmt.Sprintf("npcmodes/recordings/tbotfoot%d.rec", botNum))
	os.Remove(fmt.Sprintf("npcmodes/recordings/tbotcar%d.rec", botNum))

	switch t.players[id].RecordingPlayerType {
	case sampgo.PlayerRecordingTypeOnfoot:
		os.Rename(
			fmt.Sprintf("scriptfiles/tbotfoot%d.rec", botNum),
			fmt.Sprintf("npcmodes/recordings/tbotfoot%d.rec", botNum),
		)
		if !t.bots[botNum].IsSingle {
			botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
			script := fmt.Sprintf("tbotfoot%d", botNum)
			sampgo.ConnectNPC(botName, script)
		}
	case sampgo.PlayerRecordingTypeDriver:
		os.Rename(
			fmt.Sprintf("scriptfiles/tbotcar%d.rec", botNum),
			fmt.Sprintf("npcmodes/recordings/tbotcar%d.rec", botNum),
		)
		if !t.bots[botNum].IsSingle {
			botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
			script := fmt.Sprintf("tbotcar%d", botNum)
			sampgo.ConnectNPC(botName, script)
		}
	}

	t.bots[botNum].RecordingPlayerID = NoRecordingBotNumber
	t.players[id].RecordingBotNumber = NoRecordingBotNumber
	t.players[id].RecordingPlayerType = sampgo.PlayerRecordingTypeNone
}
