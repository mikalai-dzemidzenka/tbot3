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
type T struct {
	players        map[playerID]*PlayerInfo
	Bots           map[botNumber]*BotInfo
	IsNicksVisible bool
}

func New() *T {

	// TODO on initialization fill cache
	//file, err := os.OpenFile(fmt.Sprintf("scriptfiles/tbotskin%d.cfg", botNum), os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	sampgo.Print(fmt.Sprintf("failed to read bot file: %v", err))
	//	return
	//}
	//defer file.Close()
	//
	//skinid, err := io.ReadAll(file)
	//if err != nil {
	//	sampgo.Print(fmt.Sprintf("failed to read bot skin: %v", err))
	//	return
	//}
	//skinID, err := strconv.Atoi(string(skinid))
	//if err != nil {
	//	sampgo.Print(fmt.Sprintf("failed to parse bot skin: %v", err))
	//	return
	//}
	//
	//if t.Bots[botNum].IsSingle {
	//		file, err := os.Create(fmt.Sprintf("scriptfiles/tbotsingle%d.cfg", botNum))
	//		if err != nil {
	//			sampgo.Print(fmt.Sprintf("create tbot single file: %v", err))
	//			return
	//		}
	//		file.WriteString(strconv.Itoa(t.Bots[botNum].BotGroupID))
	//		file.Close()
	//	}

	t := &T{
		players:        make(map[playerID]*PlayerInfo),
		Bots:           make(map[botNumber]*BotInfo),
		IsNicksVisible: false,
	}
	// TODO init from file/sqlite
	return t
}

func (t *T) Tbready(id int) {
	t.players[id].ready = true
	groupID := t.players[id].BotGroupID

	var bots []*BotInfo
	for _, bot := range t.Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		if !bot.ready {
			return
		}
		bots = append(bots, bot)
	}

	for _, bot := range bots {
		sampgo.SendClientMessage(bot.id, 0x000002, " ")
		if bot.Car != NoCar && bot.SeatID != 0 {
			sampgo.PutPlayerInVehicle(bot.id, bot.Car, bot.SeatID)
		}
		sampgo.SetPlayerSkin(bot.id, bot.Skin)
	}
}

func (t *T) TBot(id int, botNum int, isSingle bool) {
	if !t.IsRecording(id) {
		t.StartRecording(id, botNum, isSingle)
	} else {
		t.StopRecording(id)
	}
}

func (t *T) Tgrs(groupID int) {
	for botNum, bot := range t.Bots {
		if bot.BotGroupID == groupID {
			t.Trs(botNum)
		}
	}
}

func (t *T) ExportData() {
	data, err := os.ReadFile("scriptfiles/tbot_dump.txt")
	if err != nil {
		sampgo.Print(err.Error())
		return
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		sampgo.Print(err.Error())
	}
}

func (t *T) InitRun() {
	for botNum, bot := range t.Bots {
		if !bot.IsSingle {
			t.Trs(botNum)
		}
	}
}

func (t *T) SaveData() {
	data, err := json.Marshal(t)
	if err != nil {
		sampgo.Print(err.Error())
		return
	}
	f, _ := os.Create("scriptfiles/tbot_dump.txt")
	f.Write(data)
	f.Close()
}

func (t *T) Trs(botNum int) {
	_, ok := t.Bots[botNum]
	if !ok {
		sampgo.SendClientMessage(0, 0xFF0000, "bot doesn't exist")
		return
	}

	t.Tkick(botNum)

	botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
	sampgo.ConnectNPC(botName, "tbot")
}

func (t *T) Tdelall() {
	for botNum := range t.Bots {
		t.Tkick(botNum)

		delete(t.Bots, botNum)
	}
}

func (t *T) Tgdel(groupID int) {
	for botNum, bot := range t.Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		t.Tkick(botNum)

		delete(t.Bots, botNum)
	}
}

func (t *T) Tdel(botNum int) {
	t.Tkick(botNum)
	delete(t.Bots, botNum)
}

func (t *T) Tkick(botNum int) {
	bot, ok := t.Bots[botNum]
	if !ok {
		return
	}
	if bot.id != BotNotConnected {
		sampgo.Kick(bot.id)
	}
}

func (t *T) Tgkick(groupID int) {
	for botNum, bot := range t.Bots {
		if bot.BotGroupID != groupID {
			continue
		}
		t.Tkick(botNum)
	}
}

func (t *T) Tlist() []string {
	list := make([]string, 0, len(t.Bots))
	for botNum, bot := range t.Bots {
		info := fmt.Sprintf("TBot%d: %s", botNum, bot.String())
		list = append(list, info)
	}
	return list
}

func (t *T) Tbinit(id int) (recording string, recType int, isSingle int, groupID int) {
	bot, ok := t.players[id]
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
