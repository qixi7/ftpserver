package server

import (
	"fmt"
	"ftpServer/global"
	"ftpServer/xcore/xlog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var regNet *regexp.Regexp       // 网络统计正则
var regClientLog *regexp.Regexp // 客户端日志正则

var netLogString = `网络统计`
var clientLogString = `客户端日志`

func init() {
	regNet = regexp.MustCompile(netLogString)
	regClientLog = regexp.MustCompile(clientLogString)
}

func fillFilePath(fileName string) (bool, string) {
	if regNet.MatchString(fileName) {
		return true, fmt.Sprintf("网络统计/%s", fileName)
	}
	if regClientLog.MatchString(fileName) {
		return true, fmt.Sprintf("客户端日志/%s", fileName)
	}
	return false, fileName
}

// 建文件夹
func makeDir(dirPath, today string) bool {
	var err error
	// 创建文件夹
	if err = os.MkdirAll(filepath.Join(dirPath), 0777); err != nil {
		xlog.Errorf("Mkdir err=%v, today=%s, newDir=%s", err, today, dirPath)
		return false
	}
	// 修改权限
	if err = os.Chmod(dirPath, 0777); err != nil {
		xlog.Errorf("Chmod err=%v, today=%s, newDir=%s", err, today, dirPath)
		return false
	}
	return true
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//parse the multipart form in the request
		err := r.ParseMultipartForm(2 * 1024 * 1024 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		today := time.Now().Format(global.TimeFormat)
		folderPath := filepath.Join(global.FileDir, "/", today)
		if _, err = os.Stat(folderPath); os.IsNotExist(err) {
			if !makeDir(fmt.Sprintf("%s/%s", folderPath, netLogString), today) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !makeDir(fmt.Sprintf("%s/%s", folderPath, clientLogString), today) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// get the pure file text
		for key, sliceStr := range m.Value {
			ok, filePath := fillFilePath(key)
			if !ok {
				xlog.Errorf("fileName=%s, can not parse dir.", key)
				continue
			}
			dst, err := os.Create(folderPath + "/" + filePath)
			if err != nil {
				xlog.Errorf("Create filePath=%s, file=%s, err=%v", filePath, key, err)
				dst.Close()
				continue
			}
			for i := range sliceStr {
				_, err = dst.WriteString(sliceStr[i])
				if err != nil {
					xlog.Errorf("filePath=%s, file=%s, WriteString err=%v", filePath, key, err)
					continue
				}
			}
			dst.Close()
		}

		// todo.暂不需要去支持传文件
		//get the *fileheaders
		files := m.File["uploadfile"]
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			if err != nil {
				file.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				xlog.Errorf("Open index=%d File err=%v", i, err)
				return
			}
			//create destination file making sure the path is writeable.
			dst, err := os.Create(folderPath + "/" + files[i].Filename)
			if err != nil {
				file.Close()
				dst.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				xlog.Errorf("Create file index=%d, folderPath=%s, Filename=%s File err=%v",
					i, folderPath, files[i].Filename, err)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				file.Close()
				dst.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				xlog.Errorf("Copy file index=%d, folderPath=%s, Filename=%s File err=%v",
					i, folderPath, files[i].Filename, err)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
