package sky_file

import (
	"ethernum/src/utils/sky_logger"
	"os"
	"path/filepath"
	"strings"
)

//获取程序的绝对路径
func GetCurrentDirectory() (string, error) {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		sky_logger.Errorf("Can't get the absoluted path, process=%s,err_desc=%s",
			os.Args[0], err.Error())
		return "", err
	}

	return strings.Replace(dir, "\\", "/", -1), nil //将\替换成/
}
