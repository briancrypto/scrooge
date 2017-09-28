//package main

package scrooge

import (
	"crypto/rsa"
	"fmt"

	"scrooge/cryptoutil"
)

type TxHandler struct {
	Pool *UTXOPool
}

func NewTxHandler(pool *UTXOPool) *TxHandler {
	return &TxHandler{Pool: pool}
}

/**
 * @return true if:
 * (1) all outputs claimed by {@code tx} are in the current UTXO pool,
 * (2) the signatures on each input of {@code tx} are valid,
 * (3) no UTXO is claimed multiple times by {@code tx},
 * (4) all of {@code tx}s output values are non-negative, and
 * (5) the sum of {@code tx}s input values is greater than or equal to the sum of its output
 *     values; and false otherwise.
 */
func (handler *TxHandler) IsValidTx(tx *Transaction) bool {
	txUTXOs := make(map[UTXO]bool)
	var inValueSum, outValueSum float64
	for inputIdx, txIn := range tx.Inputs {
		// (3) no UTXO is claimed multiple times by {@code tx},
		tmpUtxo := UTXO{TxHash: string(txIn.PrevTxHash), Index: txIn.OutputIdx}
		if _, ok := txUTXOs[tmpUtxo]; ok {
			fmt.Printf("%v is used as multiple input entry in transaction!\n", tmpUtxo)
			return false
		}
		txUTXOs[tmpUtxo] = true
		// (1) all outputs claimed by {@code tx} are in the current UTXO pool
		// what it actually means is whether the TxInput claimed existed in UTXO pool
		exist := handler.Pool.Contains(tmpUtxo)
		if !exist {
			fmt.Printf("UTXO %v does not exist!\n", tmpUtxo)
			return false
		}
		// (2) the signatures on each input of {@code tx} are valid,
		utxoTxOutput := handler.Pool.GetTxOutput(tmpUtxo)
		//rawData := tx.GetRawDataToSign(txIn.OutputIdx)
		rawData := tx.GetRawDataToSign(inputIdx)
		isValid := cryptoutil.RSAVerify(&utxoTxOutput.Address, rawData, txIn.Signature)
		if !isValid {
			fmt.Printf("rsa verification failed for Tx Input#%v!\ndata:%x\npub key:%x\nsign:%x\n", txIn.OutputIdx, rawData, utxoTxOutput.Address, txIn.Signature)
			return false
		}
		inValueSum += utxoTxOutput.Value
	}
	for _, txOut := range tx.Outputs {
		// (4) all of {@code tx}s output values are non-negative, and
		if txOut.Value < 0 {
			fmt.Printf("%v has negative value!\n", txOut)
			return false
		}
		outValueSum += txOut.Value
	}

	// (5) the sum of {@code tx}s input values is greater than or equal to the sum of its output
	// values; and false otherwise.
	if inValueSum < outValueSum {
		fmt.Printf("Sum of output value is greater than sum of input value for transaction with hash %x\n", tx.Hash)
		return false
	}

	return true
}

/**
 * Handles each epoch by receiving an unordered array of proposed transactions, checking each
 * transaction for correctness, returning a mutually valid array of accepted transactions, and
 * updating the current UTXO pool as appropriate.
 */
func (handler *TxHandler) HandleTxs(possibleTxs []*Transaction) []*Transaction {
	removedUTXOs := make([]*UTXO, 0, len(possibleTxs))
	for idx := 0; idx < len(possibleTxs); idx++ {
		tx := possibleTxs[idx]
		isValid := handler.IsValidTx(tx)
		if isValid {
			removedUTXOs = handler.removeInputFromUTXOPool(tx, removedUTXOs)
			handler.addOutputIntoUTXOPool(tx)
		} else {
			// remove tx from possibleTxs
			possibleTxs[idx] = possibleTxs[len(possibleTxs)-1]
			possibleTxs[len(possibleTxs)-1] = nil
			possibleTxs = possibleTxs[:len(possibleTxs)-1]
			idx--
		}
	}
	return possibleTxs
}

func (handler *TxHandler) removeInputFromUTXOPool(tx *Transaction, removedUTXOs []*UTXO) []*UTXO {
	for _, txInput := range tx.Inputs {
		removeUtxo := &UTXO{TxHash: string(txInput.PrevTxHash), Index: txInput.OutputIdx}
		handler.Pool.RemoveUTXO(*removeUtxo)
		removedUTXOs = append(removedUTXOs, removeUtxo)
	}
	return removedUTXOs
}

func (handler *TxHandler) addOutputIntoUTXOPool(tx *Transaction) {
	for outIdx := 0; outIdx < len(tx.Outputs); outIdx++ {
		tmpOutput := tx.Outputs[outIdx]
		tmpUtxo := &UTXO{TxHash: string(tx.Hash), Index: outIdx}
		handler.Pool.AddUTXO(*tmpUtxo, &tmpOutput)
	}
}

//(tx.Outputs)

func TestTxHandler() {

	pk1 := cryptoutil.GetPrivateKey()
	pk2 := cryptoutil.GetPrivateKey()
	pk3 := cryptoutil.GetPrivateKey()

	// create initial output in UTXOPool
	utxos := []*UTXO{&UTXO{TxHash: "txhash#1", Index: 0},
		&UTXO{TxHash: "txhash#1", Index: 1},
		&UTXO{TxHash: "txhash#1", Index: 2},
		&UTXO{TxHash: "txhash#2", Index: 0},
		&UTXO{TxHash: "txhash#2", Index: 1},
	}

	utxosOutput := []*TOutput{
		&TOutput{Value: 10.5, Address: pk1.PublicKey},
		&TOutput{Value: 15.5, Address: pk2.PublicKey},
		&TOutput{Value: 5.5, Address: pk3.PublicKey},
		&TOutput{Value: 1, Address: pk1.PublicKey},
		&TOutput{Value: 12.3, Address: pk2.PublicKey},
	}

	utxopool := NewUTXOPool()
	for idx, item := range utxos {
		utxopool.AddUTXO(*item, utxosOutput[idx])
	}

	newUtxos := utxopool.GetAllUTXO()
	for _, item := range newUtxos {
		fmt.Printf("item: %v\n", item)
	}

	txHandler := NewTxHandler(utxopool)

	myTx := NewTransaction()
	myTx.AddInput([]byte("txhash#1"), 0)

	myTx.AddInput([]byte("prev tx hash #1"), 1)
	myTx.AddOutput(5, pk1.PublicKey)
	myTx.AddOutput(2, pk1.PublicKey)
	myTx.AddOutput(3, pk2.PublicKey)
	myTx.AddOutput(1, pk3.PublicKey)
	toAddSignature(myTx, pk1, 0)
	toAddSignature(myTx, pk1, 1)
	myTx.Finalize()

	txHandler.IsValidTx(myTx)

}

func toAddSignature(myTx *Transaction, prKey *rsa.PrivateKey, txIdx int) {
	rawData := myTx.GetRawDataToSign(txIdx)
	signature, err := cryptoutil.RSASign(prKey, rawData)
	if err == nil {
		myTx.AddSignature(signature, txIdx)
		fmt.Printf("txIdx:%v\npubKey:%x\nsignature:%x\n", txIdx, prKey.PublicKey, signature)
	}
}
