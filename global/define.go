/*
	This define.go: some global define.
*/
package global

import (
	"fmt"
	"net"
	"os"
)

const TimeFormat = "20060102"  // time格式化字符串
const DelDayOffset = 30        // 过去多少天的日志文件夹删除
const FileDir = "./filedir"    // 对外文件系统文件夹
const FileDirPre = "/filedir/" // 对外文件系统文件夹
const ResDir = "./resdir/"     // 资源文件夹
const ResDirPre = "/resdir/"   // 资源文件夹

// 获取ip地址最后一组数字
// 返回本机IP
func GetIP() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Errorf("getIPFinalNum err=%v\n", err.Error())
		os.Exit(1)
	}
	for _, address := range addr {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	fmt.Errorf("getIP failed.\n")
	os.Exit(2)
	return ""
}
