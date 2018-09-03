package main

import (
    "fmt"
    "log"
    "os"
    "runtime"

    "github.com/urfave/cli"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/yancai/sss/util"
)

const VERSION string = "1.0.0"

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
初始化命令行参数
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
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "host, H",
                    Value: "localhost",
                    Usage: "host",
                },
                cli.IntFlag{
                    Name:  "port, p",
                    Value: 22,
                    Usage: "port",
                },
                cli.StringFlag{
                    Name:  "user, u",
                    Value: "root",
                    Usage: "user",
                },
            },
        },
    }

    app.Run(os.Args)
}

func inputPassword() string {
    fmt.Printf("input password: \n")
    pass, _ := terminal.ReadPassword(0)
    return string(pass)
}

/**
执行ssh
*/
func runSSH(c *cli.Context) error {
    // TODO: 修复进去bash后输入内容重复输出的问题
    // TODO: 修复ctrl+c退出的问题
    //var hostKey ssh.PublicKey
    // Create client config

    sshConfig := util.SSHConfig{
        Host:     c.String("host"),
        Port:     c.Int("port"),
        User:     c.String("user"),
        Password: inputPassword(),
    }

    address := fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port)

    config := &ssh.ClientConfig{
        User: sshConfig.User,
        Auth: []ssh.AuthMethod{
            ssh.Password(sshConfig.Password),
        },
        //HostKeyCallback: ssh.FixedHostKey(hostKey),
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    // TODO: 测试端口是否通

    // TODO: 密码重试三次
    // Connect to ssh server
    conn, err := ssh.Dial("tcp", address, config)
    if err != nil {
        log.Fatal("unable to connect: "+sshConfig.User+"@"+address+" error: ", err)
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

func main() {
    initCmdArgs()
}
