package model

import (
	"fmt"
	"time"
)

type Wallet struct {
	Mnemonic   string
	PublicKey  string
	PrivateKey string
	BtcAddress string
	Balance    int
	TxCount    int
}

func (w *Wallet) String() string {
	return fmt.Sprintf("[%s] %s %s %s %s %d %d", time.Now().UTC(), w.Mnemonic, w.BtcAddress, w.PrivateKey, w.PrivateKey, w.Balance, w.TxCount)
}
