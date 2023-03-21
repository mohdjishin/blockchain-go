package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const MINING_DEFFICULTY = 3

// main block struct this is the base model of a block	 (single block)
type Block struct {
	nounce       int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

// creating a new block with passsing nounce and previous hash
func NewBlock(nounce int, previousHash [32]byte, trasaction []*Transaction) *Block {
	b := new(Block)
	b.nounce = nounce
	b.timestamp = int64(time.Now().UnixNano())
	b.previousHash = previousHash
	b.transactions = trasaction

	return b
}

// struct for blockchain
type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

// creating a new blockchian and returning a block
func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)

	bc.CreateBlock(0, b.Hash())
	return bc

}

// just a methord to marshal struct to json - this is a ad-on not imp. values inside block struct is private so getting  values using receiver method and updating
func (b *Block) jsonMarshal() ([]byte, error) {

	return json.Marshal(

		struct {
			Timestamp    int64          `json:"timestamp"`
			Nonce        int            `json:"nounce"`
			PreviousHash [32]byte       `json:"previous_hash"`
			Transactions []*Transaction `json:"transactions"`
		}{
			Timestamp:    b.timestamp,
			Nonce:        b.nounce,
			PreviousHash: b.previousHash,
			Transactions: b.transactions,
		})
}

// returning last block.
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// just a simple function go thrw Whole block chain ((chain)slice)
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain : %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))

		// block is type of Block type struct so we can use function with that struct - print is a function defined bellow with Block struct

		block.Print()

	}
	fmt.Printf("%s", strings.Repeat("*", 50))
}

// to adding  a transaction and appending to blockchain trasnation pool . transation pool is slice of tranasaction
func (bc *Blockchain) AddTransation(sender, recipient string, value float32) {

	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransationPool() []*Transaction {
	transaction := make([]*Transaction, 0)
	for _, v := range bc.transactionPool {
		transaction = append(transaction, NewTransaction(v.senderBlockChainAddress, v.recipientBlockChainAddress, v.value))

	}

	return transaction

}

// 000

// getting hash starting as 000 return true matches
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transaction []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{nounce: nonce,
		previousHash: previousHash,
		timestamp:    0,
		transactions: transaction}

	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	fmt.Println(guessHashStr)

	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transaction := bc.CopyTransationPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DEFFICULTY) {
		nonce = nonce + 1
	}
	return nonce
}

//  transation struct

type Transaction struct {
	senderBlockChainAddress    string
	recipientBlockChainAddress string
	value                      float32
}

// creating a transaction and return a transation object
func NewTransaction(sender, recipeient string, value float32) *Transaction {
	t := new(Transaction)
	t.senderBlockChainAddress = sender
	t.recipientBlockChainAddress = recipeient
	t.value = value
	return t

}

// simple methord to print tranastions
func (t *Transaction) Print() {
	fmt.Printf("%s \n", strings.Repeat("-", 50))
	fmt.Printf("sender_blockchain_address : %s\n", t.senderBlockChainAddress)
	fmt.Printf("recipient_blockchain_address : %s\n", t.recipientBlockChainAddress)
	fmt.Printf("value                        : %1f\n", t.value)
}

// struct to json to access private data from transation methord
func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockChainAddress    string  `json:"sender_blockchain_addr"`
		RecipientBlockChainAddress string  `json:"recipient_blockchain_addr"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockChainAddress:    t.senderBlockChainAddress,
		RecipientBlockChainAddress: t.recipientBlockChainAddress,
		Value:                      t.value,
	})

}

// function To create Hash
func (b *Block) Hash() [32]byte {
	m, _ := b.jsonMarshal()
	// fmt.Printf("%x \n", sha256.Sum256([]byte(m)))
	return sha256.Sum256([]byte(m))
}

// creating block using NewBlock function and appending to BlockChain (chain) slice
func (bc *Blockchain) CreateBlock(nounce int, previousHash [32]byte) *Block {
	b := NewBlock(nounce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}

	return b
}

// just a function to print blocks
func (b *Block) Print() {
	fmt.Printf("Timestamp		  : %d\n", b.timestamp)
	fmt.Printf("previous_hash    : %x\n", b.previousHash)
	fmt.Printf("Nounce			  : %d\n", b.nounce)

	for _, t := range b.transactions {
		t.Print()
	}
	// fmt.Printf("Transactions		  : %s\n", b.transactions)
}

// init function
func init() {
	log.SetPrefix("Blockchain :")

}

func main() {

	blockChain := NewBlockchain()
	blockChain.Print()

	fmt.Println()
	fmt.Println()
	fmt.Println()
	blockChain.AddTransation("A", "B", 10)
	prevHash := blockChain.LastBlock().Hash()

	nonce := blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, prevHash)

	fmt.Println()

	fmt.Println()

	blockChain.AddTransation("C", "D", 878778)
	blockChain.AddTransation("X", "Y", 43)
	prevHash = blockChain.LastBlock().Hash()
	nonce = blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, prevHash)

	blockChain.Print()
}
