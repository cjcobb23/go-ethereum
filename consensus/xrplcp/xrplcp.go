package main

import (
	"github.com/ethereum/go-ethereum/consensus"
	"fmt"
	"crypto/ecdsa"
	"crypto/sha256"
)

func main() {
	fmt.Println("vim-go")
}

var thresholds int[4] {50,65,70,95}
var timeout int = 1
var quorum = 80


type XRPLConsensus {

	myPubKey ecdsa.PublicKey
	myPrivKey ecdsa.PrivateKey
	unl map[ecdsa.PublicKey]bool
	
	props map[ecdsa.PublicKey]Proposal
	vals map[common.Hash]map[ecdsa.PublicKey]bool
	lastVals map[ecdsa.PublicKey]common.Hash
	transactions []common.Hash
	pendingTransactions []common.Hash
	prevLedger common.Hash
	validatedLedger common.Hash
	maxSeq uint32
	curRound int
}

type Proposal {
	pubKey ecdsa.PublicKey
	prevLedger common.hash
	transactions map[common.Hash]bool
	int round
	sig []byte
}

type struct Validation {

	pubKey ecdsa.PublicKey
	blockHash common.Hash
	uint32 seq
	sig []byte

}

func (p Proposal) String() string
{
	
	return fmt.Sprintf("%v round: %v sig: %v transactions %v",
		pubKey,round,sig,prevLedger,transactions)
}

func (v Validation) Sign(priv ecdsa.PrivateKey) error {

	s := fmt.Sprintf("%v %v %v",v.pubKey,v.blockHash,v.seq)
	h := sha256.Sum256([]byte(s))
	sig, err := ecdsa.SignASN1(rand.Reader,&priv,h[:])
	if err == nil
	{
		v.sig = sig
	}
    return err
}

func (p Proposal) Sign(priv ecdsa.PrivateKey) error
{
	s := fmt.Sprintf("%v %v %v %v",p.pubKey,p.prevLedger,p.transactions,p.round)

	h := sha256.Sum256([]byte(s))
	sig, err := ecdsa.SignASN1(rand.Reader,&priv,h[:])
	if err == nil
	{
		p.sig = sig
	}
    return err
}

func (c XRPLConsensus*) receive(v Validation) {

	if _,ok := c.unl[v.pubKey]; ok {

		c.vals[v.blockHash][v.pubKey] = true
		c.lastVals[v.pubKey] = v.blockHash
		if v.seq > maxSeq {

			count := 0
			for p := range vals[v.block] {
				count++
			}
	
			if count * 100 > len(c.unl) * quorum {
				c.validatedLedger = v.blockHash
			}
		}
	}
}


func (c XRPLConsensus*) receive(p Proposal)
{
	if _,ok := c.unl[p.pubKey]; !ok
	{
		fmt.Println("Ignoring untrusted proposal : ", p)
		return
	}
	if p.prevLedger != c.prevLedger
	{
		fmt.Println("Ignoring proposal for different ledger : ", p, c.prevLedger)
		return
	}
	oldProp, ok := c.props[p.pubKey]
	if !ok
	{
		c.props[p.pubKey] = p
		fmt.Println("Received initial proposal : ",p)
		
	}
	else if oldProp.round < p.round
	{
		c.props[p.pubKey] = p
		fmt.Println("Updating proposal : ", p)
	}
	else
	{
		fmt.Println("Ignoring old proposal : ", p)
	}
}

func (c XRPLConsensus*) start(l common.Hash)
{
	c.prevLedger = l
	c.curRound = 0
	c.transactions = c.pendingTransactions
	c.props = make(map[ecdsa.PublicKey]bool)
}

func apply(l common.Hash, txns []common.Hash) (common.Hash, uint32) {

	return nil, 0
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

func (c XRPLConsensus*) preferredLedger() {

	blocks := make([]common.Hash)
	for p,v := range c.lastVals {

		blocks = append(blocks,v)
	}
	eca := earliestCommonAncestor(blocks)
	done := false

	children := children(eca)
	for done := false; !done && len(children) > 0; {
	
		sort.Slice(children,func(i, j common.Hash) bool {
			
			return support(i) > support(j)
		})


		
	}

}

func (c XRPLConsensus*) update()
{
	if c.prevLedger != preferredLedger() {

		start(preferredLedger())
	}
	else {
		updatePosition()
		if checkConsensus() {

			newLedger, seq := apply(c.prevLedger,c.transactions)

			if seq > c.maxSeq {
				v := Validation{c.myPubKey,newLedger,seq}
				v.Sign(c.myPrivKey)
				broadcast(v)
				c.maxSeq = seq
			}
			c.start(newLedger)
		}
	}

}

func (c XRPLConsensus*) perferredLedger() common.Hash {

	return nil
}

func (c XRPLConsensus*) checkConsensus() bool
{
	matching := 1
	for p := range c.props
	{
	
		bool mismatch = false
		for txn := range p.transactions
		{
	
			if _, ok := c.transactions[txn]; !ok {
				mismatch = true
				break
			}
		}
		if !mismatch && len(c.transactions) == len(p.transactions) {
			matching++
		}
	}
	return matching * 100 >= quorum * len(c.unl) 
}



func (c XRPLConsensus*) updatePosition()
{
	txns := make(map[common.Hash]int)
	for v,p := range c.props
	{

		for txn := p.transactions
		{
			txns[txn]++
		}
	}

	thresh := thresholds[c.curRound]

	for txn,sup := range txns
	{
		if sup * 100 < thresh
		{
			delete(txns, txn)
		}
	}

	i := 0
	c.transactions = make([]common.Hash,len(txns))
	for txn := range txns
	{
		c.transactions[i] = txn
		i++
	}
	c.curRound++
	myProp := Proposal{myKey,c.prevLedger,c.transactions,c.curRound}
	myProp.Sign(c.myPrivKey)

	broadcast(myProp)
}

func broadcast(p Proposal) error {
	return nil
}

func broadcast(v Validation) error {
	return nil
}
