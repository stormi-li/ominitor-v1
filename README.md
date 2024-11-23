# Ominitor 节点监控框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  
## 简介

**Ominitor** 是一个用于监控和管理由 **omi** 系列框架启动的所有节点的工具。通过 **Ominitor**，您可以实时查看各节点的运行状态，并进行远程管理，确保系统的健康和稳定。

## 功能

- **支持实时监控节点信息**：实时展示节点的运行状态、资源消耗等重要信息，便于运维人员及时发现问题。
- **支持远程管理注册服务**：能够远程控制和管理注册的服务，便于动态调整和优化系统。
## 教程
### 安装
```shell
go get github.com/stormi-li/ominitor-v1
```
### 启动节点监控
```go
package main

import (
	"github.com/go-redis/redis/v8"
	ominitor "github.com/stormi-li/ominitor-v1"
)

func main() {
	c := ominitor.NewClient(&redis.Options{Addr: "localhost:6379"})
	c.Start("localhost:9013")
}
```
在浏览区输入[http://localhost:9013](http://localhost:9013),如果跳转到节点管理界面则表示启动成功。