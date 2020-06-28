## 腾讯短信网关服务 

#### 使用go mod 解决依赖问题:
##### Golang 版本大于 1.12
```bash
export GO111MODULE="on"
export GOPROXY="https://goproxy.cn,direct"
```

#### Mac交叉编译方法：
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tsmsrv app/main.go
```

* WeChat：seraphico
* Email: osx1260@163.com