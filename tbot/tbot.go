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
	vehid     = int
)

func init() {
	InitBots()
}

func Tnicks() {
	if IsNicksVisible {
		IsNicksVisible = false
		for botNum := range Bots {
			DetachBotNick(botNum)
		}
	} else {
		IsNicksVisible = true
		for botNum := range Bots {
			AttachBotNick(botNum)
		}
	}
}

func InitBots() {
	for _, b := range Bots {
		if b.CarModel != 0 {
			veh := sampgo.CreateVehicle(b.CarModel, b.X, b.Y, b.Z, b.Angle, b.Color[0], b.Color[1], 60, false)
			AddCar(veh, b.CarInfo)
			b.BotRuntimeInfo = NewBotRuntimeInfo(veh)
		} else {
			b.BotRuntimeInfo = NewBotRuntimeInfo(NoCar)
		}
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

func SaveData(name string) {
	data, err := json.Marshal(Bots)
	if err != nil {
		sampgo.Print(err.Error())
		return
	}
	f, _ := os.Create(fmt.Sprintf("scriptfiles/%s.json", name))
	f.Write(data)
	f.Close()
}

func LoadData(name string) error {
	data, err := os.ReadFile(fmt.Sprintf("scriptfiles/%s.json", name))
	if err != nil {
		return err
	}

	Tdelall()

	for veh := range Vehs {
		sampgo.DestroyVehicle(veh)
	}

	Bots = make(map[botNumber]*BotInfo)
	err = json.Unmarshal(data, &Bots)
	if err != nil {
		return err
	}

	InitBots()
	return nil
}

func Trs(botNum int) {
	_, ok := Bots[botNum]
	if !ok {
		sampgo.SendClientMessage(0, 0xFF0000, "bot doesn't exist")
		return
	}

	Tkick(botNum)

	botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
	sampgo.ConnectNPC(botName, "tbot")
}

func Tdelall() {
	for botNum := range Bots {
		Tdel(botNum)
	}
}

func Tgdel(groupID int) {
	for botNum, bot := range Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		Tdel(botNum)
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

func Tbinit(id int) {
	bot, ok := Players[id]
	if !ok {
		sampgo.Print("tbinit: bot player not found!")
		return
	}

	var recording string
	var recType int
	if bot.car != NoCar {
		recording = fmt.Sprintf("tbotcar%d", bot.BotNumber)
		recType = sampgo.PlayerRecordingTypeDriver
	} else {
		recording = fmt.Sprintf("tbotfoot%d", bot.BotNumber)
		recType = sampgo.PlayerRecordingTypeOnfoot
	}

	var isSingle int
	if bot.IsSingle {
		isSingle = 1
	} else {
		isSingle = 0
	}

	sampgo.SendClientMessage(id, 0x000001, fmt.Sprintf("%s %d %d %d", recording, recType, isSingle, bot.BotGroupID))
}

func Tbready(id int) {
	bot := Players[id].BotInfo
	if bot.car != NoCar {
		sampgo.PutPlayerInVehicle(bot.id, bot.car, 0)
	}
	sampgo.SetPlayerSkin(bot.id, bot.Skin)

	sampgo.SendClientMessage(bot.id, 0x000002, " ")
}
