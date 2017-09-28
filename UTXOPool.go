package scrooge

import "fmt"

type UTXOPool struct {
	H map[UTXO]*TOutput
}

func NewUTXOPool() *UTXOPool {
	return &UTXOPool{H: make(map[UTXO]*TOutput)}
}

func (pool *UTXOPool) AddUTXO(utxo UTXO, txOutput *TOutput) {
	pool.H[utxo] = txOutput
}

func (pool *UTXOPool) RemoveUTXO(utxo UTXO) {
	delete(pool.H, utxo)
}

func (pool *UTXOPool) GetTxOutput(utxo UTXO) *TOutput {
	return pool.H[utxo]
}

func (pool *UTXOPool) Contains(utxo UTXO) bool {
	_, ok := pool.H[utxo]
	return ok
}

func (pool *UTXOPool) GetAllUTXO() []UTXO {
	utxos := make([]UTXO, 0, len(pool.H))
	for key, _ := range pool.H {
		utxos = append(utxos, key)
	}
	return utxos
}

func TestUTXOPool() {
	utxo1 := &UTXO{TxHash: "txhash#1", Index: 1}
	utxo1mirror := &UTXO{TxHash: "txhash#1", Index: 1}
	utxo2 := &UTXO{TxHash: "txhash#2", Index: 1}
	utxo3 := &UTXO{TxHash: "txhash#1", Index: 3}

	utxopool := NewUTXOPool()
	utxopool.AddUTXO(*utxo1, nil)
	utxopool.AddUTXO(*utxo1mirror, nil)
	utxopool.AddUTXO(*utxo2, nil)
	utxopool.AddUTXO(*utxo3, nil)

	utxos := utxopool.GetAllUTXO()
	for i, item := range utxos {
		fmt.Printf("%v)%v\n", i, item)
	}

	fmt.Printf("find previously added utxo1 in utxopool: %v\n", utxopool.Contains(*utxo1))
	fmt.Printf("find utxo1 through new object in utxopool: %v\n", utxopool.Contains(UTXO{TxHash: "txhash#1", Index: 1}))

	utxopool.RemoveUTXO(*utxo1)
	fmt.Printf("find after deletion of *utxo1: %v\n", utxopool.Contains(UTXO{TxHash: "txhash#1", Index: 1}))
	utxopool.RemoveUTXO(*&UTXO{TxHash: "txhash#2", Index: 1})
	fmt.Printf("find after deletion of UTXO{TxHash: \"txhash#2\", Index: 1}: %v\n", utxopool.Contains(UTXO{TxHash: "txhash#2", Index: 1}))

}
