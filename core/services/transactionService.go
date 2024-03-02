package services

import (
	"context"
	"github.com/liangalv/goChain/core"
	. "github.com/liangalv/goChain/core/types"
)

type TransactionService struct{}

func (ts *TransactionService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*TransactionResponse, error) {
	//validate transaction and throw it into the channel that sucks it up in the mempool
	//the channel that sends to the mempool should not be buffered, it's already inherently a buffer

	return &TransactionResponse{Status: Status_SUCCESS}, nil
}

// TODO: Basic Validation for incoming transactions
func basicValidate(req *CreateAccountRequest) {

}
