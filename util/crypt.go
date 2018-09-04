package util

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/des"
    "encoding/base64"
    "log"
    "math/rand"
    "strings"
    "time"
)

type DecodeError struct {
    error
}

/**
 * 实现明文的补全
 * 如果ciphertext的长度为blockSize的整数倍，则不需要补全
 * 否则差几个则被几个，例：差5个则补5个5
 */
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

/**
 * 实现去补码，PKCS5Padding的反函数
 */
func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

/**
 * DES加密方法
 */
func DesEncrypt(orig, key string) string {
    // 将加密内容和秘钥转成字节数组
    origData := []byte(orig)
    k := []byte(key)
    // 秘钥分组
    block, _ := des.NewCipher(k)
    //将明文按秘钥的长度做补全操作
    origData = PKCS5Padding(origData, block.BlockSize())
    //设置加密方式－CBC
    blockMode := cipher.NewCBCDecrypter(block, k)
    //创建明文长度的字节数组
    crypted := make([]byte, len(origData))
    //加密明文
    blockMode.CryptBlocks(crypted, origData)
    //将字节数组转换成字符串，base64编码
    return base64.StdEncoding.EncodeToString(crypted)
}

/**
 * DES解密方法
 */
func DESDecrypt(data string, key string) string {
    k := []byte(key)
    //将加密字符串用base64转换成字节数组
    crypted, _ := base64.StdEncoding.DecodeString(data)
    //将字节秘钥转换成block快
    block, _ := des.NewCipher(k)
    //设置解密方式－CBC
    blockMode := cipher.NewCBCEncrypter(block, k)
    //创建密文大小的数组变量
    origData := make([]byte, len(crypted))
    //解密密文到数组origData中
    blockMode.CryptBlocks(origData, crypted)
    //去掉加密时补全的部分
    origData = PKCS5UnPadding(origData)

    return string(origData)
}

func AesEncrypt(orig string, key string) string {
    // 转成字节数组
    origData := []byte(orig)
    k := []byte(key)

    // 分组秘钥
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 补全码
    origData = PKCS7Padding(origData, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
    // 创建数组
    cryted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(cryted, origData)

    return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
    // 转成字节数组
    crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
    k := []byte(key)

    // 分组秘钥
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
    // 创建数组
    orig := make([]byte, len(crytedByte))
    // 解密
    blockMode.CryptBlocks(orig, crytedByte)
    // 去补全码
    orig = PKCS7UnPadding(orig)
    return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
    padding := blocksize - len(ciphertext)%blocksize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func GetRandomString(l int) string {
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < l; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

/**
加密
des(k1): plain -> c1
aes(k2): c1-k1 -> c2
b64(): b64(b64(c2)-b64(k2))
*/
func encrypt(plain string) string {
    key1 := GetRandomString(8)
    c1 := DesEncrypt(plain, key1)
    key2 := GetRandomString(32)
    c2 := AesEncrypt(c1+"-"+key1, key2)
    return base64.StdEncoding.EncodeToString([]byte(c2 + "-" + base64.StdEncoding.EncodeToString([]byte(key2))))
}

/**
解密
*/
func decrypt(cipher string) (string, error) {
    decodeBytes, err := base64.StdEncoding.DecodeString(cipher)
    if err != nil {
        log.Fatal("decode base64 error: ", err)
        return "", DecodeError{}
    }
    decodeStr := string(decodeBytes)
    splitStrs := strings.Split(decodeStr, "-")
    c2 := splitStrs[0]
    key2Byte, err := base64.StdEncoding.DecodeString(splitStrs[1])
    if err != nil {
        log.Fatal("decode base64 error: ", err)
        return "", DecodeError{}
    }
    key2 := string(key2Byte)
    p2 := AesDecrypt(c2, key2)
    c1key1 := strings.Split(p2, "-")
    c1 := c1key1[0]
    key1 := c1key1[1]
    p := DESDecrypt(c1, key1)
    return p, nil
}
