package payindb

import (
	"time"

	"github.com/Yeremi528/itudy-back/payin"
)

type dbBalance struct {
	ID        int
	Correo    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toCoreBalance(dbbalance dbBalance) payin.Pay {
	return payin.Pay{}
}
