package dtos

type Operate struct {
	DayOperational []string `json:"day_operational,omitempty" enums:"senin,selasa,rabu,kamis,jum'at,sabtu,minggu"`
	Open           string   `json:"open,omitempty"`
	Close          string   `json:"close,omitempty"`
	Active         bool     `json:"active"`
}
