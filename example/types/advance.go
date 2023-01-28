package example_types

import (
	"github.com/shopspring/decimal"
	"time"
)

type Advance struct {
	Name      TString
	Num       *TInt
	Count     Decimal
	Total     decimal.Decimal
	R         *Advance
	R2        []*Advance
	R3        [][]*Advance
	Cache1    map[string]any
	Cache2    map[string]map[int]any
	CreatedAt *time.Time
	UpdatedAt *Time
}

type Decimal decimal.Decimal
type Time = time.Time
type TInt int64
type TString string
