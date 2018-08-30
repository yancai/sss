package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/urfave/cli"
)

const VERSION string = "1.0.0"

func main() {
	initCmdArgs()
	executeCmd("ping", "127.0.0.1")
}

/**
获取当前系统类型，是windows还是linux
*/
func getSysType() string {
	return runtime.GOOS
}

/**
读取~/.sss/sss.json文件

Linux下为~/.sss/sss.json
Windows下为%USERPROFILE%/.sss/sss.json
*/
func readConf() {

}

/**
显示session列表
*/
func showSessionList(c *cli.Context) error {
	print("1. root@127.0.0.1")
	return nil
}

/**
保存~/.sss/sss.json文件
*/
func saveConf() {

}

/**
读取命令行参数
*/
func initCmdArgs() {
	app := cli.NewApp()
	app.Name = "sss"
	app.Usage = "super ssh"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "show session list",
			Action:  showSessionList,
		},
	}

	app.Run(os.Args)
}

/**
加密
*/
func encrypt() {

}

/**
解密
*/
func decrypt() {

}

/**
执行命令
*/
func executeCmd(name string, arg ...string) {
	log.Println("executeCmd name: " + name + " arg: " + strings.Join(arg, ","))
	bin, lookErr := exec.LookPath(name)

	if lookErr != nil {
		print(lookErr)
	}

	// TODO: 此方法在linux上可用(替换当前进程)，Windows上需使用其他方式
	env := os.Environ()
	args := []string{name}
	args = append(args, arg...)
	exeErr := syscall.Exec(bin, args, env)
	if exeErr != nil {
		log.Println("execute error")
		fmt.Printf("execute error: %v", exeErr)
	}
}
