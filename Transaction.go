package scrooge

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"scrooge/cryptoutil"
)

var debugOutput bool = false

type TOutput struct {
	Value   float64
	Address rsa.PublicKey
}

type TInput struct {
	PrevTxHash []byte
	OutputIdx  int
	Signature  []byte
}

type Transaction struct {
	Hash    []byte
	Inputs  []TInput
	Outputs []TOutput
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) AddInput(prevTxHash []byte, outputIdx int) {
	tx.Inputs = append(tx.Inputs, TInput{PrevTxHash: prevTxHash, OutputIdx: outputIdx})
	if debugOutput {
		fmt.Println(len(tx.Inputs))
	}
}

func (tx *Transaction) AddOutput(value float64, address rsa.PublicKey) {
	tx.Outputs = append(tx.Outputs, TOutput{Value: value, Address: address})
}

func (tx *Transaction) AddSignature(signature []byte, idx int) {
	//        inputs.get(index).addSignature(signature);
	if idx < len(tx.Inputs) {
		tx.Inputs[idx].Signature = signature
	}

}

func (tx *Transaction) GetRawDataToSign(idx int) []byte {
	if idx > tx.NumInputs() {
		return nil
	}
	// get the ith input - PrevTxHash
	input := tx.Inputs[idx]
	sigData := bytes.NewBuffer(input.PrevTxHash)
	if debugOutput {
		fmt.Printf("input.PrevTxHash: %v \n", sigData.Bytes())
	}
	// get the ith input - OutputIdx
	binary.Write(sigData, binary.BigEndian, int32(input.OutputIdx))
	if debugOutput {
		fmt.Printf("input.PrevTxHash + input.OutputIdx: %v \n", sigData.Bytes())
	}

	// get all the output
	for _, out := range tx.Outputs {
		// add output[i].Value
		sigData.Write(FloatToByte(out.Value))
		// add output[i].Address
		sigData.Write(cryptoutil.GetPEMPublicKey(out.Address))
	}
	if debugOutput {
		fmt.Printf("len: %v  cap:%v\n", sigData.Len(), sigData.Cap())
	}

	return sigData.Bytes()

}

func (tx *Transaction) GetRawTx() []byte {

	var rawData bytes.Buffer
	for _, in := range tx.Inputs {
		// get prevTxHash
		rawData.Write(in.PrevTxHash)
		// get outputIdx
		binary.Write(&rawData, binary.BigEndian, int32(in.OutputIdx))
		// get signature
		rawData.Write(in.Signature)
	}

	for _, out := range tx.Outputs {
		// get value
		rawData.Write(FloatToByte(out.Value))
		// get address
		rawData.Write(cryptoutil.GetPEMPublicKey(out.Address))
	}
	if debugOutput {
		fmt.Printf("[GetRawTx()]len: %v  cap:%v\n", rawData.Len(), rawData.Cap())
	}
	return rawData.Bytes()
	//	return []byte("some raw transaction....")
}

func (tx *Transaction) NumInputs() int {
	return len(tx.Inputs)
}

func (tx *Transaction) NumOutputs() int {
	return len(tx.Outputs)
}

func (tx *Transaction) Finalize() {
	tx.Hash = cryptoutil.HashSha256(tx.GetRawTx())
}

func FloatToByte(f float64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, f)
	return buf.Bytes()
}

func TestTransaction() {

	myHex := []byte("0615487ebeff81fc55effa0305a8c87663bb99cf7dc0e55b78212341f5d35026")
	myHex0InByte := make([]byte, hex.DecodedLen(len(myHex)))
	hex.Decode(myHex0InByte, myHex)
	fmt.Printf("hash in byte: %v\n", myHex0InByte)

	pk1 := cryptoutil.GetPrivateKey()
	pk2 := cryptoutil.GetPrivateKey()
	pk3 := cryptoutil.GetPrivateKey()

	myTx := NewTransaction()
	myTx.AddInput(myHex0InByte, 0)
	myTx.AddInput([]byte("prev tx hash #1"), 1)
	myTx.AddInput([]byte("prev tx hash #2"), 0)
	// myTx.AddOutput(5, cryptoutil.GetPEMPublicKey(pk1.PublicKey))
	// myTx.AddOutput(2, cryptoutil.GetPEMPublicKey(pk1.PublicKey))
	// myTx.AddOutput(3, cryptoutil.GetPEMPublicKey(pk2.PublicKey))
	// myTx.AddOutput(1, cryptoutil.GetPEMPublicKey(pk3.PublicKey))
	myTx.AddOutput(5, pk1.PublicKey)
	myTx.AddOutput(2, pk1.PublicKey)
	myTx.AddOutput(3, pk2.PublicKey)
	myTx.AddOutput(1, pk3.PublicKey)
	fmt.Printf("inputs: %v\n", myTx.NumInputs())
	fmt.Printf("outputs: %v\n", myTx.NumOutputs())

	rawDataToSign := myTx.GetRawDataToSign(0)
	fmt.Printf("len:%v cap:%v\n", len(rawDataToSign), cap(rawDataToSign))
	fmt.Printf("GetRawDataToSign: %v \n", rawDataToSign)
	fmt.Printf("GetRawDataToSign: %x \n", rawDataToSign)

	myTx.Finalize()
	fmt.Printf("hash in byte: %v\n", myTx.Hash)
	fmt.Printf("hash in hex: %x\n", myTx.Hash)
}
