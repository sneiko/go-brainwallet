package mnemonic

import (
	"github.com/tyler-smith/go-bip39"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Generator struct {
	dic []string
}

func InitGenerator(pathToDic string) *Generator {
	dic, err := os.ReadFile(pathToDic)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to read dic"))
	}
	words := strings.Split(string(dic), "\n")

	return &Generator{dic: words}
}

func (g *Generator) Generate() (string, error) {
	//bip39.SetWordList(g.dic)
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}
