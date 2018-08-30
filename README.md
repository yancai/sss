# sss
增强ssh工具

sss是增强的ssh工具，支持查询ssh连接过的列表，支持记住密码，方便使用ssh

# 主要文件说明

`sss.go`为主程序入口

# 路线图
 - [ ] 使用Go开发
 - [ ] 完成基本功能: 保存查询ssh连接列表；记住密码
 - [ ] 使用安全的密码保存方式
 - [ ] 使用标准的包构建管理方式

# 代码开发使用方式

```sh
# 下载代码
go get github.com/yancai/sss

# IDE打开 $GOTPATH/src/github.com/yancai/sss 编辑代码

# 构建
go build github.com/yancai/sss

# 运行
sss
```
