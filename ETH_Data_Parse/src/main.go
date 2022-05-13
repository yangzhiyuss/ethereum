package main

import (
	"ethernum/src/app/executor"
	"ethernum/src/app/init_conf"
	"ethernum/src/types"
	"ethernum/src/utils/sky_logger"
	"os"
)

func main() {
	//初始化系统配置
	err := initconf.InitConf()
	if err != nil {
		sky_logger.Errorf("main fail to init the project error: %s", err.Error())
		os.Exit(1)
	} else {
		sky_logger.Infof("success init the project: %d parallize(r | w)",
			types.Parallize())
	}

	//获取参数
	args := os.Args

	//开始执行
	err = executor.Executor(args)
	if err != nil {
		sky_logger.Errorf("Executor is fail error: %s", err.Error())
		return
	}
	sky_logger.Infof("Executor is success, end....")
}
