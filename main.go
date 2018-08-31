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
	"golang.org/x/crypto/ssh"
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
	app.Usage = "super ssh"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "show session list",
			Action:  showSessionList,
		},
		{
			Name:    "connect",
			Aliases: []string{"c"},
			Usage:   "connect to ssh",
			Action:  runSSH,
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

/**
执行ssh
*/
func runSSH(c *cli.Context) error {
	// TODO: 修复进去bash后输入内容重复输出的问题
	// TODO: 修复ctrl+c退出的问题
	//var hostKey ssh.PublicKey
	// Create client config
	config := &ssh.ClientConfig{
		User: "yancai",
		Auth: []ssh.AuthMethod{
			ssh.Password("`123qwe"),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "192.168.74.129:22", config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 30, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	session.Run("bash")
	return nil
}
