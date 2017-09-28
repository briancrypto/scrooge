package scrooge

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var seed = time.Now().UTC().UnixNano()
var rng = rand.New(rand.NewSource(seed))

// Test 1: test handleTransactions() with simple and valid transactions
func TestHandleTxsSimpleValidTransactions(t *testing.T) {
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")
	bobWallet := hGetWalletFor(wallets, "Bob")
	charlieWallet := hGetWalletFor(wallets, "Charlie")

	// Scenario 1: start with validating one transaction
	// =================================================
	myTx := createTestTransaction(aliceWallet, []int{0}, []*PersonWallet{bobWallet, charlieWallet})
	possibleTxs := []*Transaction{myTx}

	txHandler := NewTxHandler(pool)
	txHandler.HandleTxs(possibleTxs)

	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)

	// Scenario 2: sub-sequently with 2 transactions
	// =============================================

	myTx = createTestTransaction(aliceWallet, []int{1}, []*PersonWallet{bobWallet})
	myTx2 := createTestTransaction(bobWallet, []int{0}, []*PersonWallet{aliceWallet, charlieWallet})

	possibleTxs = []*Transaction{myTx, myTx2}

	txHandler = NewTxHandler(pool)
	txHandler.HandleTxs(possibleTxs)

	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)
}

// Test 2: test handleTransactions() with simple but some invalid transactions because of invalid signatures
func TestHandleTxsMixtureValidandInvalidTransactions(t *testing.T) {
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")
	bobWallet := hGetWalletFor(wallets, "Bob")
	charlieWallet := hGetWalletFor(wallets, "Charlie")
	davidWallet := hGetWalletFor(wallets, "David")

	myTx := createTestTransaction(aliceWallet, []int{0}, []*PersonWallet{bobWallet, charlieWallet})
	myTx2 := createTestTransaction(aliceWallet, []int{1}, []*PersonWallet{bobWallet})
	myTx3 := createTestTransaction(bobWallet, []int{0}, []*PersonWallet{aliceWallet, charlieWallet})
	// invalidate signature for myTx3 by using someone else private key to sign input
	hToAddSignature(myTx3, davidWallet.priKey, 0)
	myTx3.Finalize()

	possibleTxs := []*Transaction{myTx, myTx2, myTx3}

	txHandler := NewTxHandler(pool)
	possibleTxs = txHandler.HandleTxs(possibleTxs)

	if len(possibleTxs) >= 3 {
		t.Fatalf("Result has %v tx, but it should be one less as one tx has to be removed due to invalid signature!", len(possibleTxs))
	}
	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)
}

// Test 3: test handleTransactions() with simple but some invalid transactions because of inputSum < outputSum
func TestHandleTxsWithTxOfInputSumLessThanOutputSum(t *testing.T) {
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")
	bobWallet := hGetWalletFor(wallets, "Bob")
	charlieWallet := hGetWalletFor(wallets, "Charlie")

	myTx := createTestTransaction(aliceWallet, []int{0}, []*PersonWallet{bobWallet, charlieWallet})
	myTx2 := createTestTransaction(aliceWallet, []int{1}, []*PersonWallet{bobWallet})
	myTx3 := createTestTransactionWithOutputExceedInput(bobWallet, []int{0}, []*PersonWallet{aliceWallet, charlieWallet}, 1)

	possibleTxs := []*Transaction{myTx, myTx2, myTx3}

	txHandler := NewTxHandler(pool)
	possibleTxs = txHandler.HandleTxs(possibleTxs)

	if len(possibleTxs) >= 3 {
		t.Fatalf("Result has %v tx, but it should be one less as one tx has sum of output greater than sum of input!", len(possibleTxs))
	}
	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)

}

// Test 4: test handleTransactions() with simple and valid transactions with some double spends
func TestHandleTxsWithDoubleSpend(t *testing.T) {
	// double spend in one transaction

	// double spend in multiple transaction
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")
	bobWallet := hGetWalletFor(wallets, "Bob")
	charlieWallet := hGetWalletFor(wallets, "Charlie")

	// alice using the same UTXO (TxHash: "txhash#1", Index: 0) in two separate transaction
	myTx := createTestTransaction(aliceWallet, []int{0, 0}, []*PersonWallet{charlieWallet, bobWallet})
	myTx2 := createTestTransaction(bobWallet, []int{0}, []*PersonWallet{aliceWallet, charlieWallet})

	possibleTxs := []*Transaction{myTx, myTx2}

	txHandler := NewTxHandler(pool)
	possibleTxs = txHandler.HandleTxs(possibleTxs)

	if len(possibleTxs) >= 2 {
		t.Fatalf("Result has %v tx, but it should be one less as one tx is double spend!", len(possibleTxs))
	}
	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)
}

// Test 5: test handleTransactions() with valid but some transactions are simple, some depend on other transactions
func TestHandleTxsDependsOnOtherTransactions(t *testing.T) {
}

// Test 6: test handleTransactions() with valid and simple but some transactions take inputs from non-exisiting utxo's
func TestHandleTxsOnNonExistenceUTXO(t *testing.T) {
	// using a UTXO that has been previously used.
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")
	bobWallet := hGetWalletFor(wallets, "Bob")
	charlieWallet := hGetWalletFor(wallets, "Charlie")

	// Scenario 1: start with validating one transaction
	// =================================================
	myTx := createTestTransaction(aliceWallet, []int{0}, []*PersonWallet{charlieWallet})
	possibleTxs := []*Transaction{myTx}

	txHandler := NewTxHandler(pool)
	txHandler.HandleTxs(possibleTxs)

	// Scenario 2: sub-sequently with 2 transactions
	// =============================================

	// alice reusing a previously used UTXO
	myTx = createTestTransaction(aliceWallet, []int{0}, []*PersonWallet{bobWallet})
	myTx2 := createTestTransaction(bobWallet, []int{0}, []*PersonWallet{aliceWallet, charlieWallet})

	possibleTxs = []*Transaction{myTx, myTx2}

	txHandler = NewTxHandler(pool)
	possibleTxs = txHandler.HandleTxs(possibleTxs)

	if len(possibleTxs) >= 2 {
		t.Fatalf("Result has %v tx, but it should be one less as one tx is using a UTXO that does not exist!", len(possibleTxs))
	}

	assertInputRemovedFromUTXOPool(txHandler.Pool, possibleTxs, t)
	assertOutputAddedToUTXOPool(txHandler.Pool, possibleTxs, t)
}

func assertInputRemovedFromUTXOPool(pool *UTXOPool, txs []*Transaction, t *testing.T) {
	for _, tx := range txs {
		for _, tInput := range tx.Inputs {
			inputStillInUTXO := pool.Contains(UTXO{TxHash: string(tInput.PrevTxHash), Index: tInput.OutputIdx})
			if inputStillInUTXO == true {
				t.Errorf("Input did not remove successfully from UTXOPool! input PrevTxHash:%v, PrevTxOutputIdx:%v", tInput.PrevTxHash, tInput.OutputIdx)
			}
		}
	}
}

func assertOutputAddedToUTXOPool(pool *UTXOPool, txs []*Transaction, t *testing.T) {
	for _, tx := range txs {
		for outIdx, _ := range tx.Outputs {
			outputAdded := pool.Contains(UTXO{TxHash: string(tx.Hash), Index: outIdx})
			if outputAdded == false {
				t.Errorf("Output is not added successfully to UTXOPool! input PrevTxHash:%v, PrevTxOutputIdx:%v", tx.Hash, outIdx)
			}
		}
	}
}

func createTestTransactionWithOutputExceedInput(wallet *PersonWallet, inputsIdx []int, receiverWallets []*PersonWallet, numberOutputOutMoreThanIn int) *Transaction {
	myTx := NewTransaction()
	totalUtxoValue := 0.0
	for iIdx := 0; iIdx < len(inputsIdx); iIdx++ {
		myTx.AddInput([]byte(wallet.utxos[inputsIdx[iIdx]].TxHash), wallet.utxos[inputsIdx[iIdx]].Index)
		totalUtxoValue += wallet.toutput[inputsIdx[iIdx]].Value
	}
	if testDebugOutput {
		fmt.Printf("New Transaction: total utxo value to be spend:%v\n", totalUtxoValue)
	}

	//numberOutputOutMoreThanIn = int(math.Min(float64(len(receiverWallets)), float64(numberOutputOutMoreThanIn)))
	//for oIdx := 0; oIdx < len(receiverWallets); oIdx, numberOutputOutMoreThanIn := oIdx+1, numberOutputOutMoreThanIn-1 {
	for oIdx := 0; oIdx < len(receiverWallets); oIdx, numberOutputOutMoreThanIn = oIdx+1, numberOutputOutMoreThanIn-1 {
		outputValue := 0.0
		// to create invalid transaction, in which sum of output is greater than sum of input
		if numberOutputOutMoreThanIn > 0 {
			totalUtxoValue *= (1 + rng.Float64())
		}
		// for the last output, use remaining totalUtxoValue
		if oIdx >= len(receiverWallets)-1 {
			outputValue = totalUtxoValue
		} else {
			percentage := rng.Float64()
			outputValue = percentage * totalUtxoValue
			totalUtxoValue -= outputValue
		}
		myTx.AddOutput(outputValue, receiverWallets[oIdx].priKey.PublicKey)
		if testDebugOutput {
			fmt.Printf("New Transaction: Output %v value:%v\n", oIdx, outputValue)
		}
	}
	for iIdx := 0; iIdx < len(inputsIdx); iIdx++ {
		hToAddSignature(myTx, wallet.priKey, iIdx)
	}
	myTx.Finalize()
	return myTx
}

func createTestTransaction(wallet *PersonWallet, inputsIdx []int, receiverWallets []*PersonWallet) *Transaction {
	return createTestTransactionWithOutputExceedInput(wallet, inputsIdx, receiverWallets, 0)
}
