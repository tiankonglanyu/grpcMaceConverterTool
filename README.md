# grpcMaceConverterTool
 The mace docker  container  rpc  transfer method ， you can  call it outside the  container and get the transform  log and the result
 
 
### 关于工具安装
工具安装  mace：

参考： https://mace.readthedocs.io/en/latest/installation/using_docker.html 

这里直接安装 mace 的 docker 环境：

```shell
docker pull registry.cn-hangzhou.aliyuncs.com/xiaomimace/mace-dev-lite

docker run -it --privileged -d --name mace-dev \
           -v /dev/bus/usb:/dev/bus/usb --net=host \
           -v /local/path:/container/path \
           -v /usr/bin/docker:/usr/bin/docker \
           -v /var/run/docker.sock:/var/run/docker.sock \
           registry.cn-hangzhou.aliyuncs.com/xiaomimace/mace-dev-lite
```



编译：mace  grpc 工具：

- 安装 go module 环境（略） 

- 进入 grpcMaceConverterTool 执行 

```
# 设置运行平台为  linux amd64 
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
```



```
go build -o mace-server   server.go   convert_zip_file.go    生成 mace-server  可运行程序
```

```shell
go build -o mace-clinet   client.go    生成 mace-clinet  可运行程序
```

把  mace-server 放在 mace 容器的 /mace 下 后台运行即可

mace-client 放在 其他容器或者机器 /usr/src/ 目录下 加可执行权限

