package main

import (
	"brainwallet/db"
	"brainwallet/model"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"log"
	"net/http"
)

type BlockchainResult struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount int `json:"funded_txo_count"`
		FundedTxoSum   int `json:"funded_txo_sum"`
		SpentTxoCount  int `json:"spent_txo_count"`
		SpentTxoSum    int `json:"spent_txo_sum"`
		TxCount        int `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount int `json:"funded_txo_count"`
		FundedTxoSum   int `json:"funded_txo_sum"`
		SpentTxoCount  int `json:"spent_txo_count"`
		SpentTxoSum    int `json:"spent_txo_sum"`
		TxCount        int `json:"tx_count"`
	} `json:"mempool_stats"`
}

func checkWallet(wallet model.Wallet, walletDb *db.File) error {
	url := fmt.Sprintf("https://blockstream.info/api/address/%s", wallet.BtcAddress)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result BlockchainResult
	if err := sonic.Unmarshal(body, &result); err != nil {
		return err
	}

	wallet.TxCount = result.ChainStats.TxCount
	wallet.Balance = result.ChainStats.FundedTxoSum - result.ChainStats.SpentTxoSum

	if wallet.Balance > 0 || wallet.TxCount > 0 {
		log.Println(wallet.String())
		if err := walletDb.Insert(wallet); err != nil {
			return err
		}
	}
	return nil
}
