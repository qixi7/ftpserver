package main

import (
	"ftpServer/server"
	"ftpServer/xcore/xlog"
)

func main() {
	s := server.NewFtpServer()
	if err := s.StartServer(); err != nil {
		xlog.Errorf("main err=%v", err)
		return
	}
}
