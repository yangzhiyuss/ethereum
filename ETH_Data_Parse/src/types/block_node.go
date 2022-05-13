package types

type BlockNode struct {
	blockHeight    uint64
	blockTxns      BlockTxns
	logErcs        LogErcs
	adressStatuses AddressStatuses
}

type BlockTxns []*BlockTxn
type LogErcs []*LogErc
type AddressStatuses []*AddressStatus

func (bn BlockNode) setBlockHeight(blockHeight uint64) {
	bn.blockHeight = blockHeight
}

func (bn BlockNode) BlockHeight() uint64 {
	return bn.blockHeight
}

func (bn BlockNode) SetBlockTxns(blockTxns BlockTxns) {
	bn.blockTxns = blockTxns
}

func (bn BlockNode) BlockTxns() BlockTxns {
	return bn.blockTxns
}

func (bn BlockNode) SetLogErcs(logErcs LogErcs) {
	bn.logErcs = logErcs
}

func (bn BlockNode) LogErcs() LogErcs {
	return bn.logErcs
}

func (bn BlockNode) SetAddressStatuses(addressStatuses AddressStatuses) {
	bn.adressStatuses = addressStatuses
}

func (bn BlockNode) AddressStatuses() AddressStatuses {
	return bn.adressStatuses
}
