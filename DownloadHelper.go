package main

import (
	"os"
	"os/exec"
	"runtime"
)

func OneDownload(nowDownloadRoot, nowFileName, downloadLink string, downloadInfo DownloadInfo) error {
	// 设置代理有差异
	var err error
	var tmpCommand = ""
	var setProxy = ""
	var unSetProxy = ""
	switch runtime.GOOS {
	case "windows":
		setProxy = "$env:http_proxy=" + "\"" +configs.DownloadHttpProxy + "\""
		unSetProxy = "$env:http_proxy=\"\""
	default:
		setProxy = "http_proxy=" + configs.DownloadHttpProxy
		unSetProxy = "http_proxy="
	}
	if downloadInfo.UseProxy == true {
		tmpCommand = setProxy
	} else {
		tmpCommand = unSetProxy
	}
	// 设置代理
	cmd := exec.Command(tmpCommand)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", tmpCommand)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	// 开始下载
	tmpCommand = "annie -o " + nowDownloadRoot + " -O \"" + nowFileName +  "\" " + "\"" + downloadLink + "\""
	cmd = exec.Command(tmpCommand)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", tmpCommand)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}