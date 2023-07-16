package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
)

type (
	botNumber = int
	playerID  = int
)
type T struct {
	players        map[playerID]*PlayerInfo
	bots           map[botNumber]*BotInfo
	isNicksVisible bool
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
	//if t.bots[botNum].IsSingle {
	//		file, err := os.Create(fmt.Sprintf("scriptfiles/tbotsingle%d.cfg", botNum))
	//		if err != nil {
	//			sampgo.Print(fmt.Sprintf("create tbot single file: %v", err))
	//			return
	//		}
	//		file.WriteString(strconv.Itoa(t.bots[botNum].BotGroupID))
	//		file.Close()
	//	}

	t := &T{
		players:        make(map[playerID]*PlayerInfo),
		bots:           make(map[botNumber]*BotInfo),
		isNicksVisible: false,
	}
	// TODO init from file/sqlite
	return t
}

func (t *T) Tbready(id int) {
	t.players[id].Ready = true
	groupID := t.players[id].BotGroupID

	var bots []*BotInfo
	for _, bot := range t.bots {
		if bot.BotGroupID != groupID {
			continue
		}
		if !bot.Ready {
			return
		}
		bots = append(bots, bot)
	}

	for _, bot := range bots {
		sampgo.SendClientMessage(bot.ID, 0x000002, " ")
		if bot.Car != NoCar && bot.SeatID != 0 {
			sampgo.PutPlayerInVehicle(bot.ID, bot.Car, bot.SeatID)
		}
		sampgo.SetPlayerSkin(bot.ID, bot.Skin)
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
	for botNum, bot := range t.bots {
		if bot.BotGroupID == groupID {
			t.Trs(botNum)
		}
	}
}

func (t *T) Trs(botNum int) {
	bot, ok := t.bots[botNum]
	if !ok {
		sampgo.SendClientMessage(0, 0xFF0000, "bot doesn't exist")
		return
	}

	if bot.ID != BotNotConnected {
		sampgo.Kick(bot.ID)
	}

	botName := fmt.Sprintf("%s%d", BotPrefix, botNum)
	sampgo.ConnectNPC(botName, "tbot")
}

func (t *T) Tdelall() {
	for botNum, bot := range t.bots {
		if bot.ID != BotNotConnected {
			t.DisconnectBot(botNum)
		}

		delete(t.bots, botNum)
	}
}

func (t *T) Tgdel(groupID int) {
	for botNum, bot := range t.bots {
		if bot.BotGroupID != groupID {
			continue
		}
		t.Tkick(botNum)

		delete(t.bots, botNum)
	}
}

func (t *T) Tdel(botNum int) {
	t.Tkick(botNum)
	delete(t.bots, botNum)
}

func (t *T) Tkick(botNum int) {
	bot, ok := t.bots[botNum]
	if !ok {
		return
	}
	if bot.ID != BotNotConnected {
		t.DisconnectBot(botNum)
	}
}

func (t *T) Tgkick(groupID int) {
	for botNum, bot := range t.bots {
		if bot.BotGroupID != groupID {
			continue
		}
		t.Tkick(botNum)
	}
}

func (t *T) Tlist() []string {
	list := make([]string, 0, len(t.bots))
	for botNum, bot := range t.bots {
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
