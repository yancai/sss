package util

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
    "strconv"

    "github.com/mitchellh/go-homedir"
)

const (
    SSSConfigName = "sss.json"
)

func SSSConfigDir() string {
    home, err := homedir.Dir()
    if err != nil {
        log.Fatal("get home dir error: ", err)
        return "./"
    }
    return home + "/.sss"
}

func SSSConfig() string {
    return SSSConfigDir() + "/" + SSSConfigName
}

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
SSH缓存配置
 */
type SSHCache struct {
    ID             int    `json:"id"`
    Address        string `json:"address"`
    CipherPassword string `json:"cipher_password"`
}

type NotFoundError struct {
    error
}

func (this *NotFoundError) Error() string {
    return fmt.Sprintf("未找到cache")
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

        plain, err := decrypt(cache.CipherPassword)
        if err != nil {
            log.Fatal("decrypt password error: ", err)
        }
        return SSHConfig{
            User:     params[1],
            Host:     params[2],
            Port:     port,
            Password: plain,
        }
    } else {
        return SSHConfig{}
    }
}

func checkCacheDir() {
    if _, err := os.Stat(SSSConfigDir()); err != nil {
        if os.IsNotExist(err) {
            err = os.Mkdir(SSSConfigDir(), 0755)
            if err != nil {
                log.Fatal("mkdir `"+SSSConfigDir()+"` error: ", err)
            }
        }
    }
}

func checkCacheFile() {
    checkCacheDir()
    if _, err := os.Stat(SSSConfig()); err != nil {
        if os.IsNotExist(err) {
            ioutil.WriteFile(SSSConfig(), nil, 0644)
        }
    }
}

type SSHCacheManager struct {
    cacheMap map[string]SSHCache
    maxId    int
}

func NewSSHCacheManager() *SSHCacheManager {
    caches := ReadCache()
    cacheMap := map[string]SSHCache{}
    maxId := 0
    for _, v := range caches {
        cacheMap[v.Address] = v
        if maxId < v.ID {
            maxId = v.ID
        }
    }
    return &SSHCacheManager{
        cacheMap: cacheMap,
        maxId:    maxId,
    }
}

/**
读取cache
 */
func ReadCache() []SSHCache {
    checkCacheFile()

    b, err := ioutil.ReadFile(SSSConfig())
    if err != nil {
        log.Fatal("read json `"+SSSConfig()+"` error: ", err)
    }

    var caches []SSHCache
    if len(b) != 0 {
        err = json.Unmarshal(b, &caches)
        if err != nil {
            log.Fatal("unmarshal json `"+SSSConfig()+"` error: ", err)
            return caches
        }
    }
    return caches
}

/**
保存cache
 */
func (manager *SSHCacheManager) SaveCache() {
    checkCacheDir()

    var caches []SSHCache

    for _, v := range manager.cacheMap {
        caches = append(caches, v)
    }

    if jsonStr, err := json.Marshal(caches); err != nil {
        log.Fatal("Caches to json: `"+SSSConfig()+"` error: ", err)
    } else {
        if err = ioutil.WriteFile(SSSConfig(), []byte(jsonStr), 0644); err != nil {
            log.Fatal("save to json `"+SSSConfig()+"` error: ", err)
        }
    }
}

/**
增加cache
 */
func (manager *SSHCacheManager) AddCache(cache SSHCache) {
    if _, exist := manager.cacheMap[cache.Address]; !exist {
        manager.cacheMap[cache.Address] = cache
        manager.maxId += 1
    }
}

/**
将配置添加至cache
 */
func (manager *SSHCacheManager) AddToCache(config SSHConfig) {
    cache := config2cache(config)
    cache.ID = manager.maxId + 1
    manager.AddCache(cache)
    manager.SaveCache()
}

/**
清除cache
 */
func (manager *SSHCacheManager) PurgeCache() {
    err := os.Remove(SSSConfig())
    if err != nil {
        log.Fatal("delete `"+SSSConfig()+"` error: ", err)
    }
}

/**
通过id获取配置
 */
func (manager *SSHCacheManager) GetConfig(id int) (SSHConfig, error) {
    caches := ReadCache()
    for _, cache := range caches {
        if cache.ID == id {
            return cache2config(cache), nil
        }
    }
    return SSHConfig{}, &NotFoundError{}
}

/**
从缓存中删除指定id的记录
 */
func (manager *SSHCacheManager) DelCache(id int) {
    for k, v := range manager.cacheMap {
        if v.ID == id {
            delete(manager.cacheMap, k)
        }
    }
    manager.SaveCache()
}
