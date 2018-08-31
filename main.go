package main

import (
    "fmt"
    "log"
    "os"
    "regexp"
    "runtime"
    "strconv"

    "github.com/urfave/cli"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
)

const VERSION string = "1.0.0"

/**
SSH连接配置
 */
type SSHConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
}

/**
SSH配置缓存
 */
type SSHCache struct {
    Address        string `json:"address"`
    CipherPassword string `json:"cipher_password"`
}

/**
SSH配置 转 SSH缓存
 */
func config2cache(config SSHConfig) SSHCache {
    return SSHCache{
        Address:        fmt.Sprintf("%s@%s:%d", config.User, config.Host, config.Port),
        CipherPassword: encrypt(config.Password),
    }
}

/**
SSH缓存 转 SSH配置
 */
func cache2config(cache SSHCache) SSHConfig {
    pattern := regexp.MustCompile(`([\w]+)@([\w1-9.]+):([\d{2,5}]+)`)
    params := pattern.FindStringSubmatch(cache.Address)

    if params != nil {
        port, _ := strconv.Atoi(params[3])

        return SSHConfig{
            User:     params[1],
            Host:     params[2],
            Port:     port,
            Password: decrypt(cache.CipherPassword),
        }
    } else {
        return SSHConfig{}
    }

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

/**
加密
*/
func encrypt(plain string) string {
    // TODO: to finish
    return plain
}

/**
解密
*/
func decrypt(cipher string) string {
    // TODO: to finish
    return cipher
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

    sshConfig := SSHConfig{
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
