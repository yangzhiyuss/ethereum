package executor

import (
	"context"
	"errors"
	"ethernum/src/app/executor/history"
	"ethernum/src/third_lib/src/golang.org/semaphore"
	"ethernum/src/types"
	"ethernum/src/utils/sky_logger"
	"strconv"
	"sync"
)

func Executor(args []string) error {
	//拉取历史数据
	if len(args) == 3 && types.DataFetchType() == 1 {
		start, end, err := StratAndEndNum(args[1], args[2])
		if err != nil {
			sky_logger.Errorf("Fail to get start and end error: %s", err.Error())
			return err
		}
		//抓取历史数据
		err = ETHHistoryData(start, end)
		if err != nil {
			sky_logger.Errorf("Fetch History Data happen error: %s", err.Error())
			return err
		}

	} else if len(args) != 3 && types.DataFetchType() == 1 {
		err := errors.New("fetch history data need two params, the param's count is error, " +
			"please input tow params")
		sky_logger.Errorf(err.Error())
		return err
	}

	//拉取实时数据
	if types.DataFetchType() == 2 {
		//待开发
	}

	return nil
}

//获取参数中的开始和结束块
func StratAndEndNum(param1 string, param2 string) (uint64, uint64, error) {
	startBlock, err := strconv.ParseUint(param1, 10, 64)
	if err != nil {
		sky_logger.Errorf("param1(%s) Fail of sting to uint64 "+
			"error: %s", param1, err.Error())
		return 0, 0, err
	}

	endBlock, err := strconv.ParseUint(param2, 10, 64)
	if err != nil {
		sky_logger.Errorf("param2(%s) Fail of sting to uint64 "+"error: %s",
			param2, err.Error())
		return 0, 0, err
	}

	if startBlock > endBlock {
		tmp := startBlock
		startBlock = endBlock
		endBlock = tmp
	}
	return startBlock, endBlock, nil
}

func ETHHistoryData(start uint64, end uint64) error {
	//建立队列，初始化
	queue := new(types.BlockNodeQeueu)
	queue.InitQueue(types.QueueLen())

	ctx := context.TODO()
	//初始化解解析ETH的资源量
	parsePv := semaphore.NewWeighted(int64(types.Parallize()))
	//初始化写入ETH的资源量量
	writePv := semaphore.NewWeighted(int64(types.Parallize()))

	//设置wait 等待子程序执行结束
	wait := new(sync.WaitGroup)
	//ETH数据获取解析
	go func() {
		wait.Add(1)
		defer wait.Done()
		//ETH块解析
		err := history.ParseHistoryData(ctx, parsePv, start, end, queue)
		if err != nil {
			sky_logger.Errorf("Fail to parse history data error: %s", err.Error())
		}
	}()

	//ETH块写入
	go func() {
		wait.Add(1)
		defer wait.Done()
		//ETH块写入
		err := history.WriteHistoryData(ctx, writePv, start, end, queue)
		if err != nil {
			sky_logger.Errorf("Fail to writer history data error: %s", err.Error())
			return
		}

	}()

	wait.Wait()
	return nil
}
