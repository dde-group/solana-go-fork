package ws

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ParsedBlockAccountResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value struct {
		Slot  uint64                           `json:"slot"`
		Err   interface{}                      `json:"err,omitempty"`
		Block *rpc.GetParsedBlockAccountResult `json:"block,omitempty"`
	} `json:"value"`
}

// NOTE: Unstable, disabled by default
//
// Subscribe to receive notification anytime a new block is Confirmed or Finalized.
//
// **This subscription is unstable and only available if the validator was started
// with the `--rpc-pubsub-enable-block-subscription` flag. The format of this
// subscription may change in the future**
func (cl *Client) ParsedBlockAccountSubscribe(
	filter BlockSubscribeFilter,
	opts *BlockSubscribeOpts,
) (*ParsedBlockAccountsSubscription, error) {
	var params []interface{}
	if filter != nil {
		switch v := filter.(type) {
		case BlockSubscribeFilterAll:
			params = append(params, "all")
		case *BlockSubscribeFilterMentionsAccountOrProgram:
			params = append(params, rpc.M{"mentionsAccountOrProgram": v.Pubkey})
		}
	}
	if opts != nil {
		obj := make(rpc.M)
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				// Valid encodings:
				// solana.EncodingJSON, // TODO
				solana.EncodingJSONParsed, // TODO
				//solana.EncodingBase58,
				//solana.EncodingBase64,
				//solana.EncodingBase64Zstd,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = rpc.TransactionDetailsAccounts
		}
		if opts.Rewards != nil {
			obj["rewards"] = opts.Rewards
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	genSub, err := cl.subscribe(
		params,
		nil,
		"blockSubscribe",
		"blockUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res ParsedBlockAccountResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &ParsedBlockAccountsSubscription{
		sub: genSub,
	}, nil
}

type ParsedBlockAccountsSubscription struct {
	sub *Subscription
}

func (sw *ParsedBlockAccountsSubscription) Key() string {
	return sw.sub.Key()
}

func (sw *ParsedBlockAccountsSubscription) Recv() (*ParsedBlockAccountResult, error) {
	select {
	case d := <-sw.sub.stream:
		return d.(*ParsedBlockAccountResult), nil
	case err := <-sw.sub.err:
		return nil, err
	}
}

func (sw *ParsedBlockAccountsSubscription) Err() <-chan error {
	return sw.sub.err
}

func (sw *ParsedBlockAccountsSubscription) Response() <-chan *ParsedBlockAccountResult {
	typedChan := make(chan *ParsedBlockAccountResult, 1)
	go func(ch chan *ParsedBlockAccountResult) {
		// TODO: will this subscription yield more than one result?
		d, ok := <-sw.sub.stream
		if !ok {
			return
		}
		ch <- d.(*ParsedBlockAccountResult)
	}(typedChan)
	return typedChan
}

func (sw *ParsedBlockAccountsSubscription) Unsubscribe() {
	sw.sub.Unsubscribe()
}
