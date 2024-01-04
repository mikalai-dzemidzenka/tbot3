package tbot

var Vehs = make(map[vehid]CarInfo)

type CarInfo struct {
	CarModel       int
	X, Y, Z, Angle float32
	Color          [2]int
}

func RemoveCar(vehid int) {
	delete(Vehs, vehid)
}

func AddCar(vehid int, carInfo CarInfo) {
	Vehs[vehid] = carInfo
}
