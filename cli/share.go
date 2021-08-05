package cli

type Company struct {
	Id     uint   `pg:",pk"`
	Symbol string `pg:",unique"`
	Name   string
}

type Share struct {
	Id         uint     `pg:",pk"`
	Company    *Company `pg:"rel:has-one"`
	PCP        float32
	LTP        float32
	ImportDate string `pg:"import_date"`
}

type ShareAmount struct {
	Id         uint     `pg:",pk"`
	Company    *Company `pg:"rel:has-one"`
	ImportDate string   `pg:"import_date"`
	Amount     uint
}
