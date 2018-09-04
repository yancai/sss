package main

import (
    "fmt"
    "log"
    "os"
    "github.com/urfave/cli"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/yancai/sss/util"
    "strconv"
)

const VERSION string = "1.0.0"

/**
显示cache列表
*/
func showCacheList(c *cli.Context) {
    caches := util.ReadCache()
    fmt.Printf("%2v\t%v\n", "ID", "Address")
    for _, v := range caches {
        fmt.Printf("%2v\t%v\n", v.ID, v.Address)
    }
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
            Action:  showCacheList,
        },
        {
            Name:    "connect",
            Aliases: []string{"c"},
            Usage:   "connect to ssh",
            Action:  connectTo,
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "host, H",
                    Value: "localhost",
                    Usage: "Host",
                },
                cli.IntFlag{
                    Name:  "port, p",
                    Value: 22,
                    Usage: "Port",
                },
                cli.StringFlag{
                    Name:  "user, u",
                    Value: "root",
                    Usage: "User",
                },
                cli.IntFlag{
                    Name:  "id, i",
                    Value: -1,
                    Usage: "cache Id",
                },
            },
        },
        {
            Name:    "id",
            Aliases: []string{"i"},
            Usage:   "connect to ssh by id",
            Action:  sshById,
        },
        {
            Name:    "delete",
            Aliases: []string{"d"},
            Usage:   "delete session by id",
            Action:  delSSSCache,
        },
        {
            Name:    "purge",
            Aliases: []string{"p"},
            Usage:   "purge sss cache",
            Action:  purgeSSSCache,
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
通过id连接ssh
 */
func sshById(c *cli.Context) error {
    idStr := c.Args().Get(0)
    id, _ := strconv.Atoi(idStr)
    manager := util.NewSSHCacheManager()
    sshConfig, err := manager.GetConfig(id)
    if err != nil {
        return err
    }
    runSSH(sshConfig)
    return nil
}

/**
通过host, port, user, password连接
 */
func connectTo(c *cli.Context) error {
    manager := util.NewSSHCacheManager()

    id := c.Int("id")
    var sshConfig util.SSHConfig
    if id == -1 {
        sshConfig = util.SSHConfig{
            Host:     c.String("host"),
            Port:     c.Int("port"),
            User:     c.String("user"),
            Password: inputPassword(),
        }
    } else {
        var err error
        sshConfig, err = manager.GetConfig(id)
        if err != nil {
            return err
        }
    }
    runSSH(sshConfig)
    return nil
}

/**
执行ssh
*/
func runSSH(sshConfig util.SSHConfig) {
    // TODO: 修复进去bash后输入内容重复输出的问题
    // TODO: 修复ctrl+c退出的问题
    //var hostKey ssh.PublicKey
    // Create client config

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

    manager := util.NewSSHCacheManager()
    // 保存cache
    manager.AddToCache(sshConfig)
    manager.SaveCache()

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
}

/**
清理ssh连接缓存
 */
func purgeSSSCache(c *cli.Context) error {
    manager := util.NewSSHCacheManager()
    manager.PurgeCache()
    return nil
}

func delSSSCache(c *cli.Context) error {
    manager := util.NewSSHCacheManager()
    idStr := c.Args().Get(0)
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Fatal("parse id error: ", err)
    }

    manager.DelCache(id)
    return nil
}

func main() {
    initCmdArgs()
}
