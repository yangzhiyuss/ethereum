package history

import (
	"context"
	"ethernum/src/third_lib/src/golang.org/semaphore"
	"ethernum/src/types"
	"ethernum/src/utils/sky_logger"
	"sync"
)

func WriteHistoryData(ctx context.Context, writePv *semaphore.Weighted,
	start uint64, end uint64, queue *types.BlockNodeQeueu) error {

	//设置wait
	wait := sync.WaitGroup{}
	//写入数据到文件
	for true {
		//读取线程信号量减一，减到0进行堵塞
		err := writePv.Acquire(ctx, 1)
		if err != nil {
			if err != nil {
				sky_logger.Errorf("write semaphore happen error: %s", err.Error())
				return err
			}
		}

		//多线程写入数据到文件
		go func() {
			wait.Add(1)
			defer wait.Done()
			//执行结束，唤醒写入信号量
			defer writePv.Release(1)
			err := WriteETHBlockData(queue, start, end)
			if err != nil {
				sky_logger.Errorf("Fail to write eth block data error :%s", err.Error())
				return
			}
		}()

		//全部读取完毕,跳出循环,待开发
		if end == 1 {
			break
		}
	}

	wait.Wait()
	return nil
}

func WriteETHBlockData(queue *types.BlockNodeQeueu, start uint64, end uint64) error {

	return nil
}
