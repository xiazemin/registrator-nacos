golang实现nacos 第三方自动注册

1，启动nacos
```
docker run --name nacos-quick -e MODE=standalone -p 8848:8848 -p 9848:9848 -d nacos/nacos-server:2.0.2
```

http://127.0.0.1:8848/nacos/#/login

用户名密码
nacos

2,本地运行
```
go run .  nacos://localhost
```
 % go run .  nacos://127.0.0.1
2022/11/20 22:58:12 Starting registrator  ...
2022/11/20 22:58:12 Using nacos adapter: nacos://127.0.0.1
2022/11/20 22:58:12 Connecting to backend (0/0)
2022/11/20 22:58:12 Listening for Docker events ...
2022/11/20 22:58:12 Get "http://unix.sock/containers/json?": dial unix /tmp/docker.sock: connect: no such file or directory
exit status 1

```
 ln -s /var/run/docker.sock /tmp/docker.sock
```

% go run .  nacos://127.0.0.1
2022/11/20 23:03:38 ignored: ad38892f8665 no published ports
2022/11/20 23:03:38 ignored: 6ec10caa82cc no published ports
2022/11/20 23:03:38 ignored: 989fca3fcebd no published ports
2022/11/20 23:03:38 ignored: 28a5ba40a31c no published ports
2022/11/20 23:03:38 ignored: 2a4854e4b395 no published ports
2022/11/20 23:03:38 added: 7636a753ca81 xiazemindeMacBook-Pro.local:prometheus:9090

3,docker运行
A, 方式一：自己build docker
```
docker build -t xiazemin/registrator-nacos:v0.0.1 .
```
naming to docker.io/xiazemin/registrator-nacos:v0.0.1   

```
docker login
docker push docker.io/xiazemin/registrator-nacos:v0.0.1  
```

B, 方式二：拉取已经build完成的镜像
https://hub.docker.com/repository/docker/xiazemin/registrator-nacos
```
docker pull docker.io/xiazemin/registrator-nacos:v0.0.1  
```

4,启动服务
```
docker run -d \
    --name=registrator-nacos \
    --net=host \
    --volume=/var/run/docker.sock:/tmp/docker.sock \
    xiazemin/registrator-nacos:v0.0.1 \
      nacos://127.0.0.1
```
查看下日志：
docker logs da24cd7735e8a92aa00e978247c0a7abc1ac3c3431e53a134d5a2952ca653add
2022/11/20 15:45:42 ignored: ad38892f8665 no published ports
2022/11/20 15:45:42 ignored: 6ec10caa82cc no published ports
2022/11/20 15:45:42 ignored: 989fca3fcebd no published ports
2022/11/20 15:45:42 ignored: 28a5ba40a31c no published ports
2022/11/20 15:45:42 ignored: 2a4854e4b395 no published ports
2022/11/20 15:45:42 added: 7636a753ca81 docker-desktop:prometheus:9090


5,启动service
```
docker run  -p5678:5678 apple:5678
```
curl http://127.0.0.1:5678/apple
/apple

发现我们的服务已经起来了
![nacos sidecar服务自动注册](https://github.com/xiazemin/registrator-nacos/blob/main/nacos.jpeg)