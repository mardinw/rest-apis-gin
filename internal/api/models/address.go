package models

type Provinces struct {
	ID   uint8
	Name string
}

type Regencies struct {
	ID          uint16
	ProvincesID uint8
	Provinces   Provinces `gorm:"foreignKey:ProvincesID;references:ID"`
	Name        string
}

type Districts struct {
	ID        uint16
	RegencyID uint16
	Regency   Regencies `gorm:"foreignKey:RegencyID;references:ID"`
	Name      string
}

type Villages struct {
	ID         uint16
	DistrictID uint16
	District   Districts `gorm:"foreignKey:DistrictID;references:ID"`
	Name       string
}

type Address struct {
	ID         int64
	Line1      string
	Line2      string
	PostalCode string
	VillageID  uint16
	Village    Villages `gorm:"foreignKey:VillageID;references:ID"`
	Latitude   float64
	Longitude  float64
	UserID     string
}
