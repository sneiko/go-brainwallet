package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/sneiko/go-brainwalle/db"
	"github.com/sneiko/go-brainwalle/mnemonic"
	"github.com/sneiko/go-brainwalle/model"
	"github.com/sourcegraph/conc/pool"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func main() {
	walletDb := db.InitFileDb("./handled_wallets.txt")
	generator := mnemonic.InitGenerator("./dic.txt")

	p := pool.New().WithMaxGoroutines(2)
	for {
		mnemonic, err := generator.Generate()
		if err != nil {
			log.Fatalln(err)
		}

		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
		if err != nil {
			log.Fatalln(err)
		}
		privk, pubk := btcec.PrivKeyFromBytes(btcec.S256(), seed)
		address, err := btcutil.NewAddressPubKey(pubk.SerializeUncompressed(), &chaincfg.MainNetParams)
		if err != nil {
			log.Fatalln(err)
		}

		p.Go(func() {
			err := checkWallet(model.Wallet{
				Mnemonic:   mnemonic,
				PrivateKey: hex.EncodeToString(privk.Serialize()),
				PublicKey:  hex.EncodeToString(pubk.SerializeUncompressed()),
				BtcAddress: address.EncodeAddress(),
			}, walletDb)
			if err != nil {
				log.Fatalln(err)
			}
		})
	}
	p.Wait()
}
