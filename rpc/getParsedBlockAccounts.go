package rpc

import (
	"context"
	"github.com/gagliardetto/solana-go"
)

type GetParsedBlockAccountResult struct {
	// The blockhash of this block.
	Blockhash solana.Hash `json:"blockhash"`

	// The blockhash of this block's parent;
	// if the parent block is not available due to ledger cleanup,
	// this field will return "11111111111111111111111111111111".
	PreviousBlockhash solana.Hash `json:"previousBlockhash"`

	// The slot index of this block's parent.
	ParentSlot uint64 `json:"parentSlot"`

	// Present if "full" transaction details are requested.
	Transactions []TransactionAccountParsedWithMeta `json:"transactions"`

	// Present if "signatures" are requested for transaction details;
	// an array of signatures, corresponding to the transaction order in the block.
	Signatures []solana.Signature `json:"signatures"`

	// Present if rewards are requested.
	Rewards []BlockReward `json:"rewards"`

	// Estimated production time, as Unix timestamp (seconds since the Unix epoch).
	// Nil if not available.
	BlockTime *solana.UnixTimeSeconds `json:"blockTime"`

	// The number of blocks beneath this block.
	BlockHeight *uint64 `json:"blockHeight"`
}

type TransactionAccountParsedWithMeta struct {
	Transaction *ParsedTransactionAccounts `json:"transaction"`
	Meta        *ParsedTransactionMeta     `json:"meta"`
	Version     TransactionVersion         `json:"version"`
}

func (cl *Client) GetParsedBlockAccountsWithOpts(
	ctx context.Context,
	slot uint64,
	opts *GetParsedBlockOpts,
) (out *GetParsedBlockAccountResult, err error) {

	obj := M{
		"encoding": solana.EncodingJSONParsed,
	}

	if opts != nil {
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = opts.TransactionDetails
		}
		if opts.Rewards != nil {
			obj["rewards"] = opts.Rewards
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
	}

	params := []interface{}{slot, obj}

	err = cl.rpcClient.CallForInto(ctx, &out, "getBlock", params)

	if err != nil {
		return nil, err
	}
	if out == nil {
		// Block is not confirmed.
		return nil, ErrNotConfirmed
	}
	return
}
