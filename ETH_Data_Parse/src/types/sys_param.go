package types

const BLOCK_MAX_SIZE = 50            //默认100个块的全部大小为xxxG
const H_DATA_DIR = "./history_data"  //历史数据目录
const R_DATA_DIR = "./realtime_data" //实时数据目录
const FETCH_TYPE = 3                 //1:采集历史数据  2:采集实时 3:实时历史数据一起采
const NUM_PER_FETCH = 100            //每次从ETH抓取的块数目
const CORE_NUME = 64                 //系统核数目
const SYS_MEN = 64                   //内存大小xxxG
const LOG_DIR = "./log"              //日志目录
const ETH_SERVER = "ws://183.60.141.2:8546"

type SystemParam struct {
	parallize     int    //同功能子线程线程默认最大并发度
	queueLen      int    //队列长度
	hDataDir      string //历史数据目录
	rDataDir      string //实时数据目录
	dataFetchType int    //1:采集历史数据  2:采集实时 3:实时历史数据一起采
	coreNum       int    //内核数目
	sysMem        int    //内存大小
	ethServer     string
}

var SP *SystemParam

func InitSysParm(hDataDir string, rDataDir string, dataFetchType int,
	coreNum int, sysMem int, ethServer string) {

	SP = new(SystemParam)
	if hDataDir == "" {
		SP.hDataDir = H_DATA_DIR
	} else {
		SP.hDataDir = hDataDir
	}

	if rDataDir == "" {
		SP.rDataDir = R_DATA_DIR
	} else {
		SP.rDataDir = rDataDir
	}

	if dataFetchType <= 0 {
		SP.dataFetchType = FETCH_TYPE
	} else {
		SP.dataFetchType = dataFetchType
	}

	if coreNum <= 0 {
		SP.coreNum = CORE_NUME
	} else {
		SP.coreNum = coreNum
	}

	if sysMem <= 0 {
		SP.sysMem = SYS_MEN
	} else {
		SP.sysMem = sysMem
	}

	if ethServer == "" {
		SP.ethServer = ETH_SERVER
	} else {
		SP.ethServer = ethServer
	}

	SP.parallize = computeParallize()
	//队列长度=默认最大并发度 * 2
	SP.queueLen = SP.parallize << 1
}

func Parallize() int {
	return SP.parallize
}

func QueueLen() int {
	return SP.queueLen
}

func HDataDir() string {
	return SP.hDataDir
}

func RDataDir() string {
	return SP.rDataDir
}

func DataFetchType() int {
	return SP.dataFetchType
}

func EthServer() string {
	return SP.ethServer
}

//读或者写的最大并行度
func computeParallize() int {
	//系统留四核,注意：count=block块的读或者写并发度
	count := (SP.coreNum - 4) / 6
	//系统内存小于队列能够存储的数据的最大时，并发度减少
	for SP.sysMem < count*BLOCK_MAX_SIZE {
		count = count >> 1
	}
	if count < 1 {
		count = 1
	}
	return count
}
