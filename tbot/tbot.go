package tbot

import (
	"encoding/json"
	"fmt"
	"github.com/sampgo/sampgo"
	"os"
)

type (
	botNumber = int
	playerID  = int
)

var (
	Players        = make(map[playerID]*PlayerInfo)
	Bots           = make(map[botNumber]*BotInfo)
	IsNicksVisible = false
)

func Tbready(id int) {
	Players[id].ready = true
	groupID := Players[id].BotGroupID

	var bots []*BotInfo
	for _, bot := range Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		if !bot.ready {
			return
		}
		bots = append(bots, bot)
	}

	for _, bot := range bots {
		if bot.Car != NoCar && bot.SeatID != 0 {
			sampgo.PutPlayerInVehicle(bot.id, bot.Car, bot.SeatID)
		}
		sampgo.SetPlayerSkin(bot.id, bot.Skin)

		sampgo.SendClientMessage(bot.id, 0x000002, " ")
	}
}

func TBot(id int, botNum int, isSingle bool) {
	if !IsRecording(id) {
		StartRecording(id, botNum, isSingle)
	} else {
		StopRecording(id)
	}
}

func Tgrs(groupID int) {
	for botNum, bot := range Bots {
		if bot.BotGroupID == groupID {
			Trs(botNum)
		}
	}
}

func InitRun() {
	for botNum, bot := range Bots {
		if !bot.IsSingle {
			Trs(botNum)
		}
	}
}

func SaveData() {
	data, err := json.Marshal(Bots)
	if err != nil {
		sampgo.Print(err.Error())
		return
	}
	f, _ := os.Create("scriptfiles/tbot_dump.json")
	f.Write(data)
	f.Close()
}

func ExportData() {
	data, err := os.ReadFile("scriptfiles/tbot_dump.json")
	if err != nil {
		sampgo.Print(err.Error())
		return
	}

	err = json.Unmarshal(data, &Bots)
	if err != nil {
		sampgo.Print(err.Error())
	}
}

func Trs(botNum int) {
	bot, ok := Bots[botNum]
	if !ok {
		sampgo.SendClientMessage(0, 0xFF0000, "bot doesn't exist")
		return
	}

	if bot.id != BotNotConnected {
		if bot.Car != NoCar && bot.SeatID != 0 {
			sampgo.PutPlayerInVehicle(bot.id, bot.Car, bot.SeatID)
		}
		sampgo.SetPlayerSkin(bot.id, bot.Skin)

		sampgo.SendClientMessage(bot.id, 0x000002, " ")
	} else {
		botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
		sampgo.ConnectNPC(botName, "tbot")
	}

}

func Tdelall() {
	for botNum := range Bots {
		Tkick(botNum)

		delete(Bots, botNum)
	}
}

func Tgdel(groupID int) {
	for botNum, bot := range Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		Tkick(botNum)

		delete(Bots, botNum)
	}
}

func Tdel(botNum int) {
	Tkick(botNum)
	delete(Bots, botNum)
}

func Tkick(botNum int) {
	bot, ok := Bots[botNum]
	if !ok {
		return
	}
	if bot.id != BotNotConnected {
		sampgo.Kick(bot.id)
	}
}

func Tgkick(groupID int) {
	for botNum, bot := range Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		Tkick(botNum)
	}
}

func Tlist() []string {
	list := make([]string, 0, len(Bots))
	for botNum, bot := range Bots {
		info := fmt.Sprintf("TBot%d: %s", botNum, bot.String())
		list = append(list, info)
	}
	return list
}

func Tbinit(id int) (recording string, recType int, isSingle int, groupID int) {
	bot, ok := Players[id]
	if !ok {
		sampgo.Print("tbinit: bot player not found!")
		return
	}

	if bot.Car != NoCar {
		recording = fmt.Sprintf("tbotcar%d", bot.Number)
		if bot.SeatID == sampgo.SeatDriver {
			recType = sampgo.PlayerRecordingTypeDriver
		} else {
			recType = 3 // Passenger
		}
	} else {
		recording = fmt.Sprintf("tbotfoot%d", bot.Number)
		recType = sampgo.PlayerRecordingTypeOnfoot
	}

	if bot.IsSingle {
		isSingle = 1
	} else {
		isSingle = 0
	}

	groupID = bot.BotGroupID

	return
}
