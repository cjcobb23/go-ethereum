package consensus

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"
)

type XRPLEngine struct {
	consensus XRPLConsensus
}

func (xrp *XRPLEngine) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

func (xrp *XRPLEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {

	if xrp.consensus.IsValidated(header.Hash()) {
		return nil
	} else {
		return errors.New("not validated")
	}
}

func (xrp *XRPLEngine) VerifyHeaders(
	chain consensus.ChainHeaderReader,
	header []*types.Header,
	seals []bool,
) (chan<- struct{}, <-chan error) {

	return make(chan struct{}), make(chan error)
}

//TODO is there something else this function should do? I don't know if the below things are even necessary
func (xrp *XRPLEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Coinbase = common.Address{}
	header.Nonce = types.BlockNonce{}
	header.MixDigest = common.Hash{}
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64())
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	return nil
}

func (xrp *XRPLEngine) Finalize(
	chain consensus.ChainHeaderReader,
	header *types.Header,
	state *state.StateDB,
	txs []*types.Transaction,
	uncles []*types.Header) {

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (xrp *XRPLEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Finalize block
	xrp.Finalize(chain, header, state, txs, uncles)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil)), nil
}

func (xrp *XRPLEngine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

	xrp.Start(block.Hash())

	return nil
}
func (xrp *XRPLEngine) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	enc := []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	}
	if header.BaseFee != nil {
		enc = append(enc, header.BaseFee)
	}
	rlp.Encode(hasher, enc)
	hasher.Sum(hash[:0])
	return hash
}

func (xrp *XRPLEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

// TODO flesh out the API
type API struct {
	xrp *XRPLEngine
}

// APIs implements consensus.Engine, returning the user facing RPC APIs.
func (xrp *XRPLEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "xrp",
			Version:   "1.0",
			Service:   &API{xrp},
			Public:    true,
		},
	}
}

// Close closes the exit channel to notify all backend threads exiting.
func (xrp *XRPLEngine) Close() error {
	// TODO do whatever needs to be done
	return nil
}
