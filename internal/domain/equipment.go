package domain

type Device struct {
	Class string
	Cost  int
}

func NewDevice(class string) *Device {
	return &Device{
		Class: class,
		Cost:  GetDeviceInfo()[class],
	}
}

const (
	Pickaxe     = "Кирка"
	Ventilation = "Вентиляция"
	Trolleys    = "Вагонетка"
)

func GetDeviceInfo() map[string]int {
	return map[string]int{
		Pickaxe:     3000,
		Ventilation: 15000,
		Trolleys:    50000,
	}
}
