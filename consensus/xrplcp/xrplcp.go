package consensus

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var thresholds = [4]int{50, 65, 70, 95}
var timeout int = 1
var quorum = 80

type XRPLConsensus struct {
	myPubKey  ecdsa.PublicKey
	myPrivKey ecdsa.PrivateKey
	unl       map[ecdsa.PublicKey]bool

	props               map[ecdsa.PublicKey]Proposal
	vals                map[common.Hash]map[ecdsa.PublicKey]bool
	lastVals            map[ecdsa.PublicKey]common.Hash
	transactions        map[common.Hash]bool
	pendingTransactions map[common.Hash]bool
	prevLedger          common.Hash
	validatedLedger     common.Hash
	maxSeq              uint32
	curRound            int
}

type Proposal struct {
	pubKey       ecdsa.PublicKey
	prevLedger   common.Hash
	transactions map[common.Hash]bool
	round        int
	sig          []byte
}

type Validation struct {
	pubKey    ecdsa.PublicKey
	blockHash common.Hash
	seq       uint32
	sig       []byte
}

func (p Proposal) String() string {

	return fmt.Sprintf("%v round: %v sig: %v transactions %v",
		p.pubKey, p.round, p.sig, p.prevLedger, p.transactions)
}

func (v Validation) Sign(priv ecdsa.PrivateKey) error {

	s := fmt.Sprintf("%v %v %v", v.pubKey, v.blockHash, v.seq)
	h := sha256.Sum256([]byte(s))
	sig, err := ecdsa.SignASN1(rand.Reader, &priv, h[:])
	if err == nil {
		v.sig = sig
	}
	return err
}

func (p Proposal) Sign(priv ecdsa.PrivateKey) error {
	s := fmt.Sprintf("%v %v %v %v", p.pubKey, p.prevLedger, p.transactions, p.round)

	h := sha256.Sum256([]byte(s))
	sig, err := ecdsa.SignASN1(rand.Reader, &priv, h[:])
	if err == nil {
		p.sig = sig
	}
	return err
}

func (c *XRPLConsensus) receiveValidation(v Validation) {

	if _, ok := c.unl[v.pubKey]; ok {

		c.vals[v.blockHash][v.pubKey] = true
		c.lastVals[v.pubKey] = v.blockHash
		if v.seq > c.maxSeq {

			count := len(c.vals[v.blockHash])

			if count*100 >= len(c.unl)*quorum {
				c.validatedLedger = v.blockHash
			}
		}
	}
}

func (c *XRPLConsensus) IsValidated(blockHash common.Hash) bool {

	count := len(c.vals[blockHash])
	return count*100 >= len(c.unl)*quorum
}

func (c *XRPLConsensus) receiveProposal(p Proposal) {
	if _, ok := c.unl[p.pubKey]; !ok {
		fmt.Println("Ignoring untrusted proposal : ", p)
		return
	}
	if p.prevLedger != c.prevLedger {
		fmt.Println("Ignoring proposal for different ledger : ", p, c.prevLedger)
		return
	}
	oldProp, ok := c.props[p.pubKey]
	if !ok {
		c.props[p.pubKey] = p
		fmt.Println("Received initial proposal : ", p)

	} else if oldProp.round < p.round {
		c.props[p.pubKey] = p
		fmt.Println("Updating proposal : ", p)
	} else {
		fmt.Println("Ignoring old proposal : ", p)
	}
}

func (c *XRPLConsensus) Start(l common.Hash) {
	c.prevLedger = l
	c.curRound = 0
	c.transactions = c.pendingTransactions
	c.props = make(map[ecdsa.PublicKey]Proposal)
}

func apply(l common.Hash, txns map[common.Hash]bool) (common.Hash, uint32) {

	ret := common.Hash{}
	return ret, 0
}

func earliestCommonAncestor(blocks []common.Hash) common.Hash {
	return blocks[0]
}

func children(block common.Hash) []common.Hash {

	return nil
}

func support(block common.Hash) int {

	return 1
}

/*
func (c *XRPLConsensus) preferredLedger() {

	blocks := make([]common.Hash, 0)
	for _, v := range c.lastVals {

		blocks = append(blocks, v)
	}
	eca := earliestCommonAncestor(blocks)
	done := false

	children := children(eca)
	for done := false; !done && len(children) > 0; {

		sort.Slice(children, func(i, j int) bool {

			return support(children[i]) > support(children[j])
		})

	}

}
*/

func (c *XRPLConsensus) update() {
	if c.prevLedger != c.preferredLedger() {

		c.start(c.preferredLedger())
	} else {
		c.updatePosition()
		if c.checkConsensus() {

			newLedger, seq := apply(c.prevLedger, c.transactions)

			if seq > c.maxSeq {
				v := Validation{c.myPubKey, newLedger, seq, nil}
				v.Sign(c.myPrivKey)
				broadcastValidation(v)
				c.maxSeq = seq
			}
			c.start(newLedger)
		}
	}

}

func (c *XRPLConsensus) preferredLedger() common.Hash {

	return c.prevLedger
}

func (c *XRPLConsensus) checkConsensus() bool {
	matching := 1
	for _, prop := range c.props {

		mismatch := false
		for txn := range prop.transactions {

			if _, ok := c.transactions[txn]; !ok {
				mismatch = true
				break
			}
		}
		if !mismatch && len(c.transactions) == len(prop.transactions) {
			matching++
		}
	}
	return matching*100 >= quorum*len(c.unl)
}

func (c *XRPLConsensus) updatePosition() {
	txns := make(map[common.Hash]int)
	for _, p := range c.props {

		for txn := range p.transactions {
			txns[txn]++
		}
	}

	thresh := thresholds[c.curRound]

	for txn, sup := range txns {
		if sup*100 < thresh {
			delete(txns, txn)
		}
	}

	c.transactions = make(map[common.Hash]bool, len(txns))
	for txn := range txns {
		c.transactions[txn] = true
	}
	c.curRound++
	myProp := Proposal{c.myPubKey, c.prevLedger, c.transactions, c.curRound, nil}
	myProp.Sign(c.myPrivKey)

	broadcastProposal(myProp)
}

func broadcastProposal(p Proposal) error {
	return nil
}

func broadcastValidation(v Validation) error {
	return nil
}
