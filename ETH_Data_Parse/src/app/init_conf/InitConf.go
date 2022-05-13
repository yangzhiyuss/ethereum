package initconf

import (
	"ethernum/src/third_lib/src/github.com/ini"
	"ethernum/src/types"
	"ethernum/src/utils/sky_file"
	"ethernum/src/utils/sky_logger"
)

func InitConf() error {
	//初始化日志
	err := initLog()
	if err != nil {
		sky_logger.Errorf("Fail to init log error: %s", err.Error())
		return err
	}

	//系统参数
	err = initConfParam()
	if err != nil {
		sky_logger.Errorf("Fail to init system param error:%s", err.Error())
		return err
	}

	return nil
}

func initConfParam() error {
	confPath, err := sky_file.GetCurrentDirectory()
	if err != nil {
		sky_logger.Errorf("Fail to get current directory:%s error: %s",
			confPath, err.Error())
		return err
	}

	//confPath = confPath + "/conf/eth_conf.ini"
	confPath = "./src/conf/eth_conf.ini"
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, confPath)
	if err != nil {
		sky_logger.Errorf("Fail to read the configuration file %s error %s",
			confPath, err.Error())
		return err
	}

	ethConf, err := cfg.GetSection("eth_conf")
	if err != nil {
		sky_logger.Errorf("Fail to get eth_conf error: %s", err.Error())
		return err
	}

	key, err := ethConf.GetKey("history_data_dir")
	if err != nil {
		sky_logger.Errorf("Fail to get history data dir error: %s", err.Error())
		return err
	}
	hDataDir := key.String()

	key, err = ethConf.GetKey("realtime_data_dir")
	if err != nil {
		sky_logger.Errorf("Fail to get realtime data dir error: %s", err.Error())
		return err
	}
	rDataDir := key.String()

	key, err = ethConf.GetKey("data_fetch_type")
	if err != nil {
		sky_logger.Errorf("Fail to get data_fetch_type error: %s", err.Error())
		return err
	}
	dataFetchType, err := key.Int()
	if err != nil {
		dataFetchType = 0
	}

	key, err = ethConf.GetKey("core_num")
	if err != nil {
		sky_logger.Errorf("Fail to get core_num error: %s", err.Error())
		return err
	}
	coreNum, err := key.Int()
	if err != nil {
		coreNum = 0
	}

	key, err = ethConf.GetKey("sys_mem")
	if err != nil {
		sky_logger.Errorf("Fail to get data_fetch_type error: %s", err.Error())
		return err
	}
	sysMen, err := key.Int()
	if err != nil {
		sysMen = 0
	}

	key, err = ethConf.GetKey("eth_server")
	if err != nil {
		sky_logger.Errorf("Fail to get data_fetch_type error: %s", err.Error())
		return err
	}
	ethServer := key.String()

	types.InitSysParm(hDataDir, rDataDir, dataFetchType, coreNum, sysMen, ethServer)
	return nil
}

func initLog() error {
	err := sky_logger.InitLogger(types.LOG_DIR, -1)
	if err != nil {
		sky_logger.Errorf("Fail to create log directory error: %s", err.Error())
		return err
	}
	return nil
}
