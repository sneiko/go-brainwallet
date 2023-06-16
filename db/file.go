package db

import (
	"github.com/sneiko/go-brainwalle/model"
	"log"
	"os"
	"sync"
)

type File struct {
	path  string
	cache []model.Wallet
	mu    sync.Mutex
}

func InitFileDb(path string) *File {
	return &File{path: path}
}

func (f *File) GetWallets() []model.Wallet { return f.cache }

func (f *File) GetWalletsWithTx() []model.Wallet {
	wallets := make([]model.Wallet, 0)
	for _, w := range f.cache {
		if w.TxCount > 0 {
			wallets = append(wallets, w)
		}
	}
	return wallets
}

func (e *File) Insert(wallet model.Wallet) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cache = append(e.cache, wallet)

	f, err := os.OpenFile(e.path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(wallet.String() + "\n"); err != nil {
		return err
	}
	return nil
}
