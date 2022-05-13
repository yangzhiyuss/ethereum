package history

import (
	"context"
	"ethernum/src/third_lib/src/golang.org/semaphore"
	"ethernum/src/types"
	"ethernum/src/utils/sky_eth"
	"ethernum/src/utils/sky_logger"
	"sync"
)

//解析ETH历史数据,将解析的数据放入队列中
func ParseHistoryData(ctx context.Context, parsePv *semaphore.Weighted,
	start uint64, end uint64, queue *types.BlockNodeQeueu) error {

	//创建eth连接池
	ethConnectPool, err := util.NewETHConnectPool(types.Parallize())
	if err != nil {
		sky_logger.Errorf("Fail to create eth pool error: %s", err)
		return err
	}

	//设置wait
	wait := sync.WaitGroup{}
	//将解析的块数据写入
	for i := start; i <= end; i += types.NUM_PER_FETCH {
		first := i
		last := i + types.NUM_PER_FETCH - 1
		//防止最后块越界
		if last > end {
			last = end
		}

		//解析线程信号量减一，减到0进行堵塞
		err := parsePv.Acquire(ctx, 1)
		if err != nil {
			sky_logger.Errorf("parse semaphore happen error: %s", err.Error())
			return err
		}

		//抓取数据线程
		go func() {
			wait.Add(1)
			defer wait.Done()
			//唤醒解析信号量
			defer parsePv.Release(1)
			err := FetchETHBlockData(ethConnectPool, queue, first, last)
			if err != nil {
				sky_logger.Errorf("Fail to prase block data error:%s ", err.Error())
			}
		}()

	}

	//等到程序结束
	wait.Wait()
	//关闭连接池
	ethConnectPool.Close()

	return nil
}

func FetchETHBlockData(ethConnectPool *util.ETHConnectPool, queue *types.BlockNodeQeueu, start uint64, end uint64) error {

	return nil
}
