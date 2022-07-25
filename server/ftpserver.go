package server

import (
	"fmt"
	"ftpServer/global"
	"ftpServer/xcore/xlog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// ftp server
type FtpServer struct {
}

// --------------------------- 内置函数 ---------------------------

// --------------------------- 外置函数 ---------------------------

// new
func NewFtpServer() *FtpServer {
	return &FtpServer{}
}

// 开启服务器
func (s *FtpServer) StartServer() error {
	http.HandleFunc("/upload", uploadHandler)
	// for debug
	http.HandleFunc("/debug", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "golang大法好!")
	})
	// 检查资源文件夹是否存在, 不存在就新建
	checkAndCreateDir(global.FileDir)
	checkAndCreateDir(global.ResDir)
	selfIP := global.GetIP()
	http.Handle(global.FileDirPre, http.StripPrefix(global.FileDirPre, http.FileServer(http.Dir(global.FileDir))))
	http.Handle(global.ResDirPre, http.StripPrefix(global.ResDirPre, http.FileServer(http.Dir(global.ResDir))))
	xlog.InfoF("ftp server start. listen http://localhost:8080, dir=%s", global.FileDir)
	xlog.InfoF("logDir=http://%s:8080%s", selfIP, global.FileDirPre)
	xlog.InfoF("resDir=http://%s:8080%s", selfIP, global.ResDirPre)
	xlog.InfoF("DebugUrl: http://%s:8080/debug", selfIP)
	startDelFileTimer()
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		xlog.Errorf("start server err=%v", err)
		return err
	}
	return nil
}

func checkAndCreateDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			xlog.Fatalf("checkAndCreateDir path=%s, err=%s", path, err)
		}
	}
}

// 开启检测定时器, 删除5天前的日志, 只保留5天
func startDelFileTimer() {
	// 开服的时候先检测一下删除
	checkDelDir(global.FileDir)
	go func() {
		// 计算下一次凌晨5点在哪个时刻
		nowBaseTimeStr := time.Now().Format("20060102")
		tempTimeStr, err := time.ParseInLocation("20060102", nowBaseTimeStr, time.Local)
		if err != nil {
			xlog.Errorf("startDelFileTimer ParseInLocation now err=%v", err)
			return
		}
		next5Clock := tempTimeStr.AddDate(0, 0, 1).Add(time.Hour * 5)
		timer := time.NewTimer(next5Clock.Sub(time.Now()))
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				// 检测
				checkDelDir(global.FileDir)
				// 之后直接24小时后调用一次
				timer.Reset(time.Hour * 24)
			}
		}
	}()
}

func checkDelDir(pateName string) {
	rd, err := ioutil.ReadDir(pateName)
	if err != nil {
		xlog.Errorf("checkDelDir read dir fail, err=%v", err.Error())
		return
	}
	nowTimeF := time.Now().Format(global.TimeFormat)
	nowTimeFake, err := time.Parse(global.TimeFormat, nowTimeF)
	if err != nil {
		xlog.Errorf("nowTimeF=%s, Parse Time err=%v", nowTimeF, err)
		return
	}
	xlog.Debugf("invoke once")
	for _, fi := range rd {
		if fi.IsDir() {
			// 如果是合法文件夹. 判断是否超过今天3天了, 如果是就删除
			tempTime, err := time.ParseInLocation(global.TimeFormat, fi.Name(), time.Local)
			if err != nil {
				xlog.Errorf("dirName=%s, Parse Time err=%v", fi.Name(), err)
				return
			}
			if int32(nowTimeFake.Sub(tempTime).Hours()) >= global.DelDayOffset*24 {
				tmpDir := fmt.Sprintf("%s%s%s",
					global.FileDir, string(os.PathSeparator), fi.Name())
				err = os.RemoveAll(tmpDir)
				if err != nil {
					xlog.Errorf("Del dir=%s, err=%v", fi.Name(), err)
					continue
				}
				xlog.InfoF("Del dir=%s", tmpDir)
			}
		}
	}
}
