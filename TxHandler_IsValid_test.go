package scrooge

import "testing"

// (1) all outputs claimed by {@code tx} are in the current UTXO pool,
func TestIsValidOutputClaimedAreInUTXOPool(t *testing.T) {
}

// (2) the signatures on each input of {@code tx} are valid,
func TestIsValidSignatureOnEachInputAreValid(t *testing.T) {
}

// (3) no UTXO is claimed multiple times by {@code tx},
func TestIsValidNoUTXOClaimedMultipleTimes(t *testing.T) {
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")

	myTx := NewTransaction()
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddOutput(10, aliceWallet.priKey.PublicKey)
	myTx.AddOutput(10, aliceWallet.priKey.PublicKey)
	hToAddSignature(myTx, aliceWallet.priKey, 0)
	hToAddSignature(myTx, aliceWallet.priKey, 1)
	myTx.Finalize()

	txHandler := NewTxHandler(pool)
	validTx := txHandler.IsValidTx(myTx)
	if validTx == true {
		t.Errorf("IsValidTx=%v", validTx)
	}
}

// (4) all of {@code tx}s output values are non-negative
func TestIsValidAllOutputAreNonNegative(t *testing.T) {
	pool, wallets := testInit()
	// fmt.Println(pool.H)
	// fmt.Printf("%v,%x,%v,%v\n", wallets[0].name, wallets[0].priKey, wallets[0].toutput, wallets[0].utxos)

	aliceWallet := hGetWalletFor(wallets, "Alice")

	myTx := NewTransaction()
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddOutput(-5, aliceWallet.priKey.PublicKey)
	myTx.AddOutput(5.5, aliceWallet.priKey.PublicKey)
	hToAddSignature(myTx, aliceWallet.priKey, 0)
	myTx.Finalize()

	txHandler := NewTxHandler(pool)
	validTx := txHandler.IsValidTx(myTx)
	if validTx == true {
		t.Errorf("IsValidTx=%v, negative output is NOT expected.", validTx)
	}
}

// (5) the sum of {@code tx}s input values is greater than or equal to the sum of its output
//     values; and false otherwise.
func TestIsValidSumOfInputValuesGreaterThanSumOfOutputValues(t *testing.T) {
	pool, wallets := testInit()

	aliceWallet := hGetWalletFor(wallets, "Alice")

	// Case 1: Sum Input = Sum Output
	myTx := NewTransaction()
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddOutput(5, aliceWallet.priKey.PublicKey)
	myTx.AddOutput(5.5, aliceWallet.priKey.PublicKey)
	hToAddSignature(myTx, aliceWallet.priKey, 0)
	myTx.Finalize()

	txHandler := NewTxHandler(pool)
	validTx := txHandler.IsValidTx(myTx)
	if validTx == false {
		t.Errorf("IsValidTx=%v", validTx)
	}

	// Case 2: Sum Input > Sum Output
	myTx = NewTransaction()
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddOutput(5, aliceWallet.priKey.PublicKey)
	myTx.AddOutput(4.5, aliceWallet.priKey.PublicKey)
	hToAddSignature(myTx, aliceWallet.priKey, 0)
	myTx.Finalize()

	validTx = txHandler.IsValidTx(myTx)
	if validTx == false {
		t.Errorf("IsValidTx=%v", validTx)
	}

	// Case 3: Sum Input < Sum Output
	myTx = NewTransaction()
	myTx.AddInput([]byte(aliceWallet.utxos[0].TxHash), aliceWallet.utxos[0].Index)
	myTx.AddOutput(10, aliceWallet.priKey.PublicKey)
	myTx.AddOutput(4.5, aliceWallet.priKey.PublicKey)
	hToAddSignature(myTx, aliceWallet.priKey, 0)
	myTx.Finalize()

	validTx = txHandler.IsValidTx(myTx)
	if validTx == true {
		t.Errorf("IsValidTx=%v", validTx)
	}
}
