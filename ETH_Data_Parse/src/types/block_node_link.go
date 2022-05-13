package types

import (
	"errors"
	"sync"
)

/*
读取的数据存在BlockNodeQueue队列的结构体中，读是按照先进先出的原则,进行读取，队列不保证数据的有序性
*/
type BlockNodeQeueu struct {
	qMutex    sync.Mutex
	length    int
	maxLength int
	Head      *BlockNodes
	Tail      *BlockNodes
}

//一个节点存一百个块
type BlockNodes struct {
	nodes [NUM_PER_FETCH]*BlockNode
	Next  *BlockNodes
}

func (bns BlockNodes) SetNodes(nodes [NUM_PER_FETCH]*BlockNode) {
	bns.nodes = nodes
}

func (bns BlockNodes) Nodes() [NUM_PER_FETCH]*BlockNode {
	return bns.nodes
}

func (bnQeueu *BlockNodeQeueu) InitQueue(maxLength int) {
	bnQeueu = new(BlockNodeQeueu)

	bnQeueu.qMutex = sync.Mutex{}
	bnQeueu.maxLength = maxLength
	bnQeueu.length = 0

	//初始化队列头节点
	bnQeueu.Head = new(BlockNodes)
	bnQeueu.Head.Next = nil

	bnQeueu.Tail = bnQeueu.Head
}

//插入数据进入队列
func (bnQeueu *BlockNodeQeueu) InsertQueue(node *BlockNodes) error {
	if bnQeueu == nil || bnQeueu.Tail == nil {
		err := errors.New("The bnQueue is not init")
		return err
	}

	if bnQeueu.length >= bnQeueu.maxLength {
		err := errors.New("The maximum length of the bnQueue is exceeded")
		return err
	}
	//插入数据到队列
	bnQeueu.Tail.Next = node
	bnQeueu.Tail = node
	bnQeueu.Tail.Next = nil
	//队列长度加1
	bnQeueu.length++
	return nil
}

//从队列中删除数据，并取出删除的数据
func (bnQeueu *BlockNodeQeueu) RemoveQueue() (*BlockNodes, error) {
	var node *BlockNodes
	if bnQeueu == nil || bnQeueu.Head == nil {
		err := errors.New("The bnQueue is not init")
		return nil, err
	}

	if bnQeueu.length <= 0 {
		err := errors.New("Can't remove the bnQueue because the length is zero")
		return nil, err
	}
	//赋值给node
	node = bnQeueu.Head.Next
	//删除node
	bnQeueu.Head.Next = node.Next
	//斩断node与队列的关系,防止内存溢出
	node.Next = nil
	//队列长度减1
	bnQeueu.length--

	//假如队列长度为0,tail=head
	if bnQeueu.length == 0 {
		bnQeueu.Tail = bnQeueu.Head
	}

	return node, nil
}

//同步插入队列(队列级别锁，同一时刻只能对队列进行一次插入删除操作)
func (bnQeueu *BlockNodeQeueu) InsertQueueSync(node *BlockNodes) error {
	bnQeueu.qMutex.Lock()
	defer bnQeueu.qMutex.Unlock()
	if bnQeueu == nil || bnQeueu.Tail == nil {
		err := errors.New("The bnQueue is not init")
		return err
	}

	if bnQeueu.length >= bnQeueu.maxLength {
		err := errors.New("The maximum length of the bnQueue is exceeded")
		return err
	}
	//插入数据到队列
	bnQeueu.Tail.Next = node
	bnQeueu.Tail = node
	bnQeueu.Tail.Next = nil
	//队列长度加1
	bnQeueu.length++
	return nil
}

//同步删除队列的元素，并返回值(队列级别锁，同一时刻只能对队列进行一次插入删除操作)
func (bnQeueu *BlockNodeQeueu) RemoveQueueSync() (*BlockNodes, error) {
	bnQeueu.qMutex.Lock()
	defer bnQeueu.qMutex.Unlock()
	var node *BlockNodes
	if bnQeueu == nil || bnQeueu.Head == nil {
		err := errors.New("the bnQueue is not init")
		return nil, err
	}

	if bnQeueu.length <= 0 {
		err := errors.New("can't remove the bnQueue because the length is zero")
		return nil, err
	}
	//赋值给node
	node = bnQeueu.Head.Next
	//删除node
	bnQeueu.Head.Next = node.Next
	//斩断node与队列的关系,防止内存溢出
	node.Next = nil
	//队列长度减1
	bnQeueu.length--

	//假如队列长度为0,tail=head
	if bnQeueu.length == 0 {
		bnQeueu.Tail = bnQeueu.Head
	}

	return node, nil
}
