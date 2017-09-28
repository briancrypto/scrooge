package scrooge

import "fmt"

type UTXO struct {
	TxHash string
	Index  int
}

func NewUTXO(txHash string, index int) *UTXO {
	return &UTXO{TxHash: txHash, Index: index}
}

// implemented non reference comparison
func (utxo *UTXO) Equal(other *UTXO) bool {
	if (utxo.Index == other.Index) && (utxo.TxHash == other.TxHash) {
		return true
	}
	return false
}

func TestUTXOComparsion() {
	utxo1 := &UTXO{TxHash: "txhash#1", Index: 1}
	utxo1mirror := &UTXO{TxHash: "txhash#1", Index: 1}
	utxo2 := &UTXO{TxHash: "txhash#2", Index: 1}
	utxo3 := &UTXO{TxHash: "txhash#1", Index: 3}

	// note slice cannot be compared with ==, can only compare with nil
	fmt.Printf("utxo1 equal utxo1mirror: %v\n", utxo1 == utxo1mirror)
	fmt.Printf("utxo1 NOT equal utxo2: %v\n", utxo1 == utxo2)
	fmt.Printf("utxo1 NOT equal utxo3: %v\n", utxo1 == utxo3)
	fmt.Printf("utxo2 NOT equal utxo3: %v\n", utxo2 == utxo3)

	// Equal function
	fmt.Printf("utxo1 equal utxo1mirror (using Equal fn): %v\n", utxo1.Equal(utxo1mirror))
	fmt.Printf("utxo1 NOT equal utxo2(using Equal fn): %v\n", utxo1.Equal(utxo2))
	fmt.Printf("utxo1 NOT equal utxo3(using Equal fn): %v\n", utxo1.Equal(utxo3))
	fmt.Printf("utxo2 NOT equal utxo3(using Equal fn): %v\n", utxo2.Equal(utxo3))
}
