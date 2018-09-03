package util

import (
    "encoding/json"
    "log"
    "io/ioutil"
    "regexp"
    "strconv"
)

type SSHCache struct {
    ID             int    `json:"id"`
    Address        string `json:"address"`
    CipherPassword string `json:"cipher_password"`
}

type SSHCacheManager struct {
    cacheMap map[string]SSHCache
    maxId    int
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

func NewSSHCacheManager(caches []SSHCache) *SSHCacheManager {
    cacheMap := map[string]SSHCache{}
    maxId := 0
    for k, v := range caches {
        cacheMap[v.Address] = v
        maxId = k
    }
    return &SSHCacheManager{
        cacheMap: cacheMap,
        maxId:    maxId,
    }
}

/**
读取cache
 */
func (manager SSHCacheManager) ReadCache() {

}

/**
保存cache
 */
func (manager *SSHCacheManager) SaveCache() {
    var caches []SSHCache

    for _, v := range manager.cacheMap {
        caches = append(caches, v)
    }

    if jsonStr, err := json.Marshal(caches); err != nil {
        log.Fatal("Caches to json error: ", err)
    } else {
        if err = ioutil.WriteFile("./sss.json", []byte(jsonStr), 0644); err != nil {
            log.Fatal("save to sss.json error: ", err)
        }
    }
}

/**
增加cache
 */
func (manager *SSHCacheManager) AddCache(cache SSHCache) {
    manager.cacheMap[cache.Address] = cache
    manager.maxId += 1
}

/**
将配置添加至cache
 */
func (manager *SSHCacheManager) AddToCache(config SSHConfig) {
    cache := config2cache(config)
    cache.ID = manager.maxId
    manager.AddCache(cache)
}

/**
删除cache
 */
func (manager *SSHCacheManager) DelCache() {

}
