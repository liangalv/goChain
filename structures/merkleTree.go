package structures

// Interface defintion of MerkleTree
type merkleTree interface {
	//Add object to Tree
	Add([32]byte) error
	//Delete object from tree
	Delete([32]byte) error
	//Construct tree-array and store the root hash
	Construct() error
	//Verify the Merkle proof given the transaction
	VerifyProof([32]byte, [][32]byte) bool
	ResetTree()
}
