package main

import (
    "github.com/urfave/cli"
    "os"
    "runtime"
)

const VERSION string = "1.0.0"

func main() {
    initCmdArgs()
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
func executeCmd(cmd string) {
    print("execute " + cmd)
}

/**
输出帮助信息
*/
func help() {

}
