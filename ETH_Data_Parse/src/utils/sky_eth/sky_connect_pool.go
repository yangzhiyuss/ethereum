package util

import (
	"errors"
	"ethernum/src/types"
	"ethernum/src/utils/sky_logger"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
	"time"
)

//ETH连接池结构
type ETHConnectPool struct {
	cilentFactory chan *ethclient.Client
	isClose       bool
	mutex         sync.Mutex
}

func ethClientConnect() (*ethclient.Client, error) {
	clinet := new(ethclient.Client)
	for i := 0; i <= 9; i++ {
		var err error
		clinet, err = ethclient.Dial(types.EthServer())
		if err != nil {
			time.Sleep(5 * time.Second)
			sky_logger.Errorf("client_err: %s", err)
			if i == 9 {
				err := errors.New("can not connect eth_server")
				sky_logger.Errorf("%s", err.Error())
				return nil, err
			}
			continue
		}
		break
	}

	return clinet, nil
}

func NewETHConnectPool(poolSize int) (*ETHConnectPool, error) {
	if poolSize <= 0 {
		return nil, errors.New("pool size  < 0")
	}
	return &ETHConnectPool{
		cilentFactory: make(chan *ethclient.Client, poolSize),
		isClose:       false,
		mutex:         sync.Mutex{},
	}, nil
}

func (ecp *ETHConnectPool) GetConnect() (*ethclient.Client, error) {

	select {
	case client, ok := <-ecp.cilentFactory:
		if ok == false {
			return nil, errors.New("clientFactory is closed")
		}
		return client, nil
	default:
		return ethClientConnect()
	}
}

func (ecp *ETHConnectPool) PutConnect(client *ethclient.Client) {
	ecp.mutex.Lock()
	defer ecp.mutex.Unlock()

	if ecp.isClose == true {
		return
	}

	select {
	case ecp.cilentFactory <- client:
		sky_logger.Infof("Put connect into the closer")
	default:
		sky_logger.Infof("Put failed, close connection")
		client.Close()
	}
}

func (ecp *ETHConnectPool) Close() {
	ecp.mutex.Lock()
	defer ecp.mutex.Unlock()

	if ecp.isClose == false {
		ecp.isClose = true
	}
	close(ecp.cilentFactory)
	for client := range ecp.cilentFactory {
		client.Close()
	}
}
