package core

import (
	"container/heap"
	"log"
	"sync"
)

/*
TODO:
mempools should have a finite size on each node
-we can have a hard cap on the number of transactions on any one given mempool, there's too much overhead when keeping
https://bitcoin.stackexchange.com/questions/96068/what-if-the-mempool-exceeds-300-mb
consult this article on dynamic mempool sizing limits for futher implementation
track of the amount of memory consumed by an individual object
we have to find a way to see if a transaction has already been seen, but not in an in memory fashion
-is there a way to do this with minmal overhead is LevelDb the only solution?
*/

// A memPool implements a priority queue using gas as the priority to determine which transactions are added to the creation of a block
type MemPool struct {
	mux          sync.Mutex
	transactions []*Transaction //each pointer is 8 bytes
	validator    *Account
	idToTransMap map[[32]byte]*Transaction //each entry has 32 byte key and 8 byte pointer
}

func NewMemPool(acc *Account) *MemPool {
	return &MemPool{
		validator:    acc,
		transactions: []*Transaction{},
		idToTransMap: map[[32]byte]*Transaction{},
	}
}

func (mp *MemPool) AddTransactionToPool(t *Transaction) {
	heap.Push(mp, t)
}

func (mp *MemPool) RemoveHighestGasTransaction() (*Transaction, error) {
	return heap.Pop(mp).(*Transaction), nil
}

// Implement Heap Interface
func (mp *MemPool) Pop() any {
	mp.mux.Lock()
	defer mp.mux.Unlock()
	if len(mp.transactions) == 0 {
		log.Printf("No transactions found in mempool")
		return nil
	}
	lastIndex := len(mp.transactions) - 1
	//Get transaction and dereference
	trans := mp.transactions[lastIndex]
	mp.transactions[lastIndex] = nil //For memory leak prevention
	trans.index = -1                 //so that it can't be used for references in mempool

	//Configure new Slice
	mp.transactions = mp.transactions[0:lastIndex]

	//Delete map entry
	delete(mp.idToTransMap, trans.ID)

	return trans
}

func (mp *MemPool) Push(x any) {
	mp.mux.Lock()
	defer mp.mux.Unlock()
	trans, ok := x.(*Transaction)
	if !ok {
		log.Println("Attempted to add non-transaction object to mempool")
		return
	}

	n := len(mp.transactions)
	trans.index = n
	mp.transactions = append(mp.transactions, trans)
	mp.idToTransMap[trans.ID] = trans
}

//Implement Sort Interface for heap.Interface

func (mp *MemPool) Len() int {
	return len(mp.transactions)
}

func (mp *MemPool) Less(i, j int) bool {
	//we want the highest gas cost to be popped off first so we use greater than >
	return mp.transactions[i].gas > mp.transactions[j].gas
}

func (mp *MemPool) Swap(i, j int) {
	mp.transactions[i], mp.transactions[j] = mp.transactions[j], mp.transactions[i]
	mp.transactions[i].index, mp.transactions[j].index = i, j
}
