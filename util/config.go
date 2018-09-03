package util

import (
    "fmt"
)

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
SSH配置 转 SSH缓存
 */
func config2cache(config SSHConfig) SSHCache {
    return SSHCache{
        Address:        fmt.Sprintf("%s@%s:%d", config.User, config.Host, config.Port),
        CipherPassword: encrypt(config.Password),
    }
}

