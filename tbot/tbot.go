package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
)

const (
	BotPrefix = "TBot"
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
		sampgo.SendClientMessage(0, 0xFF0000, "ROFL 0")
		return
	}

	if bot.ID != BotNotConnected {
		sampgo.Kick(bot.ID)
	}

	botName := fmt.Sprintf("%s%d", BotPrefix, botNum)

	var script string
	switch bot.Car {
	case NoCar:
		if bot.IsSingle {
			script = fmt.Sprintf("tbotfootsingle%d", botNum)
		} else {
			script = fmt.Sprintf("tbotfoot%d", botNum)
		}
	default:
		if bot.IsSingle {
			script = fmt.Sprintf("tbotcarsingle%d", botNum)
		} else {
			script = fmt.Sprintf("tbotcar%d", botNum)
		}
	}
	sampgo.ConnectNPC(botName, script)
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
