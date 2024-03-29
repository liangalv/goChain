package structures

import (
	"errors"
	"golang.org/x/crypto/sha3"
	"math"
)

/*
Why not MPT?
MPT are preferred for state management, but there are too many considerations to do to implement such a thing
MPT Benefits include:
-Bounded Length: Prevents DoS attacks by disallowing attacks to generate too much tree depth
-Much Faster recalculation in the event of a value change
-Root calculation should be contingent on the data, not on the update order (not the case in Merkle Trees)


Design considerations:
-Perserve the root so that you don't have to recalculate everytime, will be set to nil on unsafe Add/Delete
-Keep array for transactionHashes for ease of update and delete O(N) preserves order as well
-Keep Tree so we can construct MerkleProofs
-Block constructions will happen synchronously (per elected node) so no need for synchronization primitives
*/

type TransactionTree struct {
	merkleRoot        [32]byte
	transactionHashes [][32]byte
	tree              [][32]byte
}

// Construct MerkleTree
func (tt *TransactionTree) Construct() error {
	if len(tt.transactionHashes) == 0 {
		return errors.New("tree contains no transactions")
	}
	//Ensure that we have an an even amount of leaves
	if len(tt.transactionHashes)%2 != 0 {
		tt.transactionHashes = append(tt.transactionHashes, tt.transactionHashes[len(tt.transactionHashes)-1])
	}
	//Need to generate a copy to ensure that there is no data loss modifying currentLevel
	//We will use a make function to ensure that we are not overallocating capacity
	currentLevel := make([][32]byte, len(tt.transactionHashes))
	for len(currentLevel) > 1 {
		//Need to copy every level of tree, to reconstruct proof later
		tt.tree = append(tt.tree, currentLevel...)
		hashPairs(&currentLevel)
	}
	tt.merkleRoot = currentLevel[0]
	return nil
}

/**/
// Construct Merkle Proof for param transactionHash
func (tt *TransactionTree) ConstructProof(transactionHash [32]byte) (*[][32]byte, error) {
	//Verify the membership of the transactionHash and return the index to construct proof
	ind, ok := tt.verifyMembership(transactionHash)
	if !ok {
		return nil, errors.New("transaction was not found in tree")
	}
	//If the Merkle Root was computed, this implies that transactionHashes reflects an even number of transactions and is therefore up to date
	if len(tt.tree) == 0 {
		return nil, errors.New("merkle root was never computed, proof could not be constructed")
	}
	proof := emptyProofArray(len(tt.transactionHashes))
	//we can use the length of 'proof' here cause the make function populates the array with default values
	totalSize := len(tt.transactionHashes)
	for i := 0; i < len(proof); i++ {
		//set path
		proof[i] = tt.tree[ind^1]
		//Calculate Parent Index
		ind = tt.calculateParentInd(ind, totalSize)
		//Set new totalsize for parent index calculation
		totalSize += totalSize / 2
	}
	return &proof, nil
}

// Add a transaction to the list of transactionHashes (does not recalculate MerkleRoot)
func (tt *TransactionTree) Add(transactionHash [32]byte) error {
	//Ensure that duplicate transactions are not added
	//ensure membership before attempting add
	_, ok := tt.verifyMembership(transactionHash)
	if ok {
		return errors.New("transaction is already in the tree")
	}
	//Add element
	tt.transactionHashes = append(tt.transactionHashes, transactionHash)
	//Reset the Merkle root
	tt.merkleRoot = [32]byte{}
	return nil
}

// Delete a transaction from the transactionHashes (does not recalculate MerkleRoot)
func (tt *TransactionTree) Delete(transactionHash [32]byte) error {
	//ensure membership before attempting delete
	ind, ok := tt.verifyMembership(transactionHash)
	if !ok {
		return errors.New("transaction was not found in the tree")
	}
	//Remove the element
	tt.transactionHashes = append(tt.transactionHashes[:ind], tt.transactionHashes[ind+1:]...)
	//Reset the Merkle root
	tt.merkleRoot = [32]byte{}
	return nil
}

// Safe Methods for Add/Del (recalculates MerkleRoot)
// Calls Add and Recalculates MerkleRoot
func (tt *TransactionTree) SafeAdd(transactionHash [32]byte) error {
	ok := tt.Add(transactionHash)
	if ok != nil {
		return ok
	}
	//Recalculate the MerkleRoot
	err := tt.Construct()
	if err != nil {
		return err
	}
	return nil
}

func (tt *TransactionTree) SafeDelete(transactionHash [32]byte) error {
	//Attempt to delete the transaction
	err := tt.Delete(transactionHash)
	if err != nil {
		return err
	}
	//Recalculate the MerkleRoot
	err = tt.Construct()
	if err != nil {
		return err
	}
	return nil
}

func (tt *TransactionTree) VerifyProof(th [32]byte, proof [][32]byte) (bool, error) {
	if len(tt.merkleRoot) == 0 {
		return false, errors.New("merkle root has not been computed yet")
	}
	for _, hash := range proof {
		th = sha3.Sum256(append(th[:], hash[:]...))
	}
	return (th == tt.merkleRoot), nil
}

func (tt *TransactionTree) ResetTree() {
	tt.merkleRoot = [32]byte{}
	tt.transactionHashes = [][32]byte{}
	tt.tree = [][32]byte{}
}

// Db Methods
// Helper Methods
// Check if an element is in the transactionHashes array
func (tt *TransactionTree) verifyMembership(transactionHash [32]byte) (int, bool) {
	for i, hash := range tt.transactionHashes {
		if hash == transactionHash {
			return i, true
		}
	}
	return -1, false
}

// Constructs an empty proof array with (num of levels) - 1 capacity omitting rootLevel
func emptyProofArray(length int) [][32]byte {
	return make([][32]byte, int(math.Sqrt(float64(length))))
}

// calculates the parent index for a pair of hashes
func (tt *TransactionTree) calculateParentInd(ind int, totalSize int) int {
	siblingInd := ind ^ 1
	origInd := ind
	//Offset if we're not on level 0
	if totalSize != len(tt.transactionHashes) {
		siblingInd -= totalSize
		origInd -= totalSize
	}
	offset := int(math.Min(float64(origInd), float64(siblingInd)))/2 + 1
	return totalSize - 1 + offset
}

// Hash Pairs on the same level
func hashPairs(children *[][32]byte) {
	//Unwrap array
	currLevel := *children
	//Ensure that there are an even number of leaves
	if len(currLevel)%2 != 0 {
		currLevel = append(currLevel, currLevel[len(currLevel)-1])
	}
	//parentHashes array will always only be half the size of the currLevel
	parentHashes := make([][32]byte, len(currLevel)/2)
	for i := 0; i < len(currLevel); i += 2 {
		//concatenate the pair
		pair := append(currLevel[i][:], currLevel[i+1][:]...)
		//hash and add to parentHashes
		//append potentially resizes capacity so put it in the index
		parentHashes[i/2] = sha3.Sum256(pair)
	}
	//Cannot just reassign "children" to a new pointer value, modify underlying slice to pass value back
	*children = parentHashes
}
