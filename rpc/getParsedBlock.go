package rpc

import (
	"context"
	"github.com/gagliardetto/solana-go"
)

type GetParsedBlockOpts struct {
	// Level of transaction detail to return.
	// If parameter not provided, the default detail level is "full".
	//
	// This parameter is optional.
	TransactionDetails TransactionDetailsType

	// Whether to populate the rewards array.
	// If parameter not provided, the default includes rewards.
	//
	// This parameter is optional.
	Rewards *bool

	// "processed" is not supported.
	// If parameter not provided, the default is "finalized".
	//
	// This parameter is optional.
	Commitment CommitmentType

	// Max transaction version to return in responses.
	// If the requested block contains a transaction with a higher version, an error will be returned.
	MaxSupportedTransactionVersion *uint64
}

type GetParsedBlockResult struct {
	// The blockhash of this block.
	Blockhash solana.Hash `json:"blockhash"`

	// The blockhash of this block's parent;
	// if the parent block is not available due to ledger cleanup,
	// this field will return "11111111111111111111111111111111".
	PreviousBlockhash solana.Hash `json:"previousBlockhash"`

	// The slot index of this block's parent.
	ParentSlot uint64 `json:"parentSlot"`

	// Present if "full" transaction details are requested.
	Transactions []TransactionParsedWithMeta `json:"transactions"`

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

type TransactionParsedWithMeta struct {
	Transaction *ParsedTransaction
	Meta        *ParsedTransactionMeta
}

func (cl *Client) GetParsedBlockWithOpts(
	ctx context.Context,
	slot uint64,
	opts *GetParsedBlockOpts,
) (out *GetParsedBlockResult, err error) {

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
