package writer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type centerWriter struct {
	url               string
	logDir            string // 当写入日志中心失败时，写入本地日志目录
	wg                *sync.WaitGroup
	isSync            bool // 同步写入日志
	storage           chan []byte
	removeFileLogLock sync.Mutex // 为删除历史记录加锁
}

// newCenterWriter
//
// storageChanLength 日志通道长度，根据日志写入频率适当修改，当为0时用于同步接收Fatal，Panic日志，大于0用于异步接收其他日志
func newCenterWriter(wg *sync.WaitGroup, url, logDir string, storageChanLength uint) centerWriter {
	cent := centerWriter{wg: wg, url: url, logDir: logDir}
	if storageChanLength != 0 {
		storage := make(chan []byte, storageChanLength)
		cent.storage = storage
		// 启动写入日志中心的goroutine
		cent.wg.Add(1)
		go cent.work()
	} else {
		cent.isSync = true
	}
	return cent
}

func (w centerWriter) Write(p []byte) (n int, err error) {
	if w.isSync {
		_, err = w.write(p)
		// 如果写入日志中心失败，则把日志写入到本地文件
		if err != nil {
			traceid, _ := w.writeToFile(w.logDir, "", p)
			_, err = w.writeToFile(w.logDir, traceid, []byte(err.Error()))
		}
	} else {
		w.storage <- p
	}
	return len(p), nil
}

// work 循环写入日志中心
func (w centerWriter) work() {
	defer w.wg.Done()
	// 创建信号通道
	sigCh := make(chan os.Signal, 1)
	// 注册要捕获的信号
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case p, ok := <-w.storage:
			if !ok {
				return
			}
			_, err := w.write(p)
			// 如果写入日志中心失败，则把日志写入到本地文件
			if err != nil {
				traceid, _ := w.writeToFile(w.logDir, "", p)
				w.writeToFile(w.logDir, traceid, []byte(err.Error()))
			}
		case _ = <-sigCh:
			return
		}
	}
}
func (w centerWriter) write(p []byte) (n int, err error) {
	payload := bytes.NewReader(p)
	if Debug {
		fmt.Println(string(p))
	}
	client := &http.Client{
		Transport: &http.Transport{ // 直接跳过 SSL 认证
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("POST", w.url, payload)
	if err != nil {
		return
	}

	req.Header.Add("token", "zhy")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	reBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// 如果写入日志中心失败，则把日志写入到本地文件
	var rsp responseModel
	err = json.Unmarshal(reBody, &rsp)
	if err != nil {
		return
	}
	if rsp.Code != 200 {
		err = fmt.Errorf("%s", reBody)
		return
	}
	return len(p), nil
}

// writeToFile 写到文件
func (w centerWriter) writeToFile(logDir, inTraceid string, msg []byte) (traceid string, err error) {
	if inTraceid == "" {
		var evt map[string]interface{}
		d := json.NewDecoder(bytes.NewReader(msg))
		d.UseNumber()
		err = d.Decode(&evt)
		if evt[TraceidFieldName] != nil {
			traceid = evt[TraceidFieldName].(string)
		}
	} else {
		traceid = inTraceid
	}
	// 创建目录
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return
	}

	path := filepath.Join(logDir, time.Now().Format("20060102_15"))
	// 创建目录
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return
	}
	// 创建日志文件
	logFilePath := filepath.Join(path, fmt.Sprintf("%v.log", traceid))
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	// 写入日志
	_, err = file.WriteString(fmt.Sprintf("%s\n", msg))
	if err != nil {
		return
	}
	// 删除过期日志
	// 同时只允许一个goroutine
	if w.removeFileLogLock.TryLock() {
		go w.deleteFileLog(logDir)
	}
	return
}

// deleteFileLog 删除过期日志文件
func (w centerWriter) deleteFileLog(logDir string) (err error) {
	defer w.removeFileLogLock.Unlock()
	dir, e := os.ReadDir(logDir)
	if e != nil {
		return
	}

	logAge := 7 * 24 * time.Hour
	cutoff := time.Now().Add(-1 * logAge)
	for _, v := range dir {
		if v.IsDir() {
			fl, e := v.Info()
			if e != nil {
				continue
			}
			if fl.ModTime().After(cutoff) {
				continue
			}
			err = os.RemoveAll(filepath.Join(logDir, v.Name()))
			if err != nil {
				return
			}
		}
	}
	return
}
