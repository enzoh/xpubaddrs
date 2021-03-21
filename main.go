package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	argXPubKey := flag.String("xpub", "", "Extended public key. (required)")
	argN := flag.Int("n", 8, "Number of addresses.")
	flag.Parse()
	if *argXPubKey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	err := run(*argXPubKey, *argN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(argXPubKey string, argN int) error {
	masterXPubKey, err := hdkeychain.NewKeyFromString(argXPubKey)
	if err != nil {
		return err
	}
	addresses, err := generate(masterXPubKey, argN)
	if err != nil {
		return err
	}
	for _, address := range addresses {
		fmt.Println(pretty(address))
	}
	return nil
}

func generate(
	masterXPubKey *hdkeychain.ExtendedKey,
	n int,
) ([]common.Address, error) {
	masterXPubKey0, err := masterXPubKey.Child(0)
	if err != nil {
		return nil, err
	}
	result := make([]common.Address, 0)
	for i := 0; i < n; i++ {
		childXPubKey, err := masterXPubKey0.Child(uint32(i))
		if err != nil {
			return nil, err
		}
		pubKey, err := childXPubKey.ECPubKey()
		if err != nil {
			return nil, err
		}
		address := crypto.PubkeyToAddress(*pubKey.ToECDSA())
		result = append(result, address)
	}
	return result, nil
}

func pretty(address common.Address) string {
	return strings.ToLower(strings.TrimPrefix(address.String(), "0x"))
}
