package scrooge

import (
	"crypto/rsa"

	"scrooge/cryptoutil"
)

var testDebugOutput bool = false

type PersonWallet struct {
	name    string
	priKey  *rsa.PrivateKey
	utxos   []*UTXO
	toutput []*TOutput
}

func testInit() (*UTXOPool, []*PersonWallet) {

	wallets := make([]*PersonWallet, 0, 3)

	wallets = append(wallets, &PersonWallet{})
	walletIdx := 0

	wallets[walletIdx].name = "Alice"
	wallets[walletIdx].priKey = cryptoutil.GetPrivateKey()
	wallets[walletIdx].utxos = []*UTXO{
		&UTXO{TxHash: "txhash#1", Index: 0},
		&UTXO{TxHash: "txhash#1", Index: 1},
	}
	wallets[walletIdx].toutput = []*TOutput{
		&TOutput{Value: 10.5, Address: wallets[walletIdx].priKey.PublicKey},
		&TOutput{Value: 1, Address: wallets[walletIdx].priKey.PublicKey},
	}

	wallets = append(wallets, &PersonWallet{})
	walletIdx++
	wallets[walletIdx].name = "Bob"
	wallets[walletIdx].priKey = cryptoutil.GetPrivateKey()
	wallets[walletIdx].utxos = []*UTXO{
		&UTXO{TxHash: "txhash#1", Index: 2},
		&UTXO{TxHash: "txhash#1", Index: 3},
	}
	wallets[walletIdx].toutput = []*TOutput{
		&TOutput{Value: 2.5, Address: wallets[walletIdx].priKey.PublicKey},
		&TOutput{Value: 11.2, Address: wallets[walletIdx].priKey.PublicKey},
	}

	wallets = append(wallets, &PersonWallet{})
	walletIdx++
	wallets[walletIdx].name = "Charlie"
	wallets[walletIdx].priKey = cryptoutil.GetPrivateKey()

	wallets = append(wallets, &PersonWallet{})
	walletIdx++
	wallets[walletIdx].name = "David"
	wallets[walletIdx].priKey = cryptoutil.GetPrivateKey()

	utxopool := NewUTXOPool()
	for _, tmpWallet := range wallets {
		for idx, tmpUtxo := range tmpWallet.utxos {
			utxopool.AddUTXO(*tmpUtxo, tmpWallet.toutput[idx])
		}
	}

	return utxopool, wallets
}

func hGetWalletFor(wallets []*PersonWallet, name string) *PersonWallet {
	for _, wallet := range wallets {
		if wallet.name == name {
			return wallet
		}
	}
	return nil
}

func hToAddSignature(myTx *Transaction, prKey *rsa.PrivateKey, txIdx int) {
	rawData := myTx.GetRawDataToSign(txIdx)
	signature, err := cryptoutil.RSASign(prKey, rawData)
	if err == nil {
		myTx.AddSignature(signature, txIdx)
		//fmt.Printf("txIdx:%v\npubKey:%x\nsignature:%x\n", txIdx, prKey.PublicKey, signature)
	}
}
