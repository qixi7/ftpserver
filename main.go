package main

import (
	"ftpserver/server"
	"ftpserver/xcore/xlog"
)

func main() {
	s := server.NewFtpServer()
	if err := s.StartServer(); err != nil {
		xlog.Errorf("main err=%v", err)
		return
	}
}
