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
	"golang.org/x/crypto/ssh/terminal"
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

func runSSH(c *cli.Context) error {
	//var hostKey ssh.PublicKey
	// Create client config
	config := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			ssh.Password("123456"),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "127.0.0.1:22", config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal("make raw error: ", err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	tw, th, err := terminal.GetSize(fd)
	if err != nil {
		log.Fatal("get size error: ", err)
	}

	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", th, tw, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	session.Run("bash")
	return nil
}
