# unicorn

unicorn 是一个分布式uuid生成器, 具有以下特性：

- 最高1s内可以并发生成2^1024个不重复的id

- 核心组件采用了CAS无锁编程, 性能更高

- 部署简单, 无其他依赖


# 配置

**使用配置文件**

```yaml
etcd:
  # etcd 集群节点
  cluster:
  - "192.168.10.10:2379"
  - "192.168.10.11:2379"

  # 开启tls连接
  enableTls: false
 
  # 当开启tls连接时, 可以设置该值为true, 可以不配置证书文件
  insecure: false
  
  # 开启客户端认证
  clientAuth: false

  # 如果clientAuth 是true, 需要配置ca证书路径
  caFile: ""

  # 如果开启了客户端认证或者开启了tls, 需要配置cert文件路径
  certFile: ""

  # 如果开启了客户端认证或者开启了tls, 需要配置key文件路径
  keyFile: ""
  
  # etcd客户端连接超时时间, 单位是秒
  timeout: 5

  # 定期同步时间戳到etcd, 单位是秒
  report: 3

  # 同步时间戳到本地文件路径
  localReportFile: "unicorn-report.txt"

  # 定期同步时间戳到本地文件, 单位是秒
  localReport: 60

id:

  # 起始时间戳, 要比项目启动时间小
  epoch: 1539660973223
  
  # 使用机器id获取方式, 0: 使用手动配置的id, 1: 通过ip地址向etcd集群获取
  machineIdType: 0
  
  # 手动配置的机器id, 当machineIdType = 0 是,需要配置它
  machineId: 1

  # 向etcd集群获取id的ip地址
  machineIp: "127.0.0.1"

  # id生成类型, 0: 秒级时间戳 1: 毫秒级时间戳
  idType: 1

  # 版本号, 可选值0和1, 目前暂未对版本号做出限制,默认使用0,可以根据业务来自行处理
  version: 0

http:
  
  # http地址
  addr: ":6001"
  
  # 开启https
  enableTls: false
  
  # 开启客户端认证
  clientAuth: false
  
  # 如果开启了客户端认证,需要配置ca文件路径
  caFile: ""

  # 如果开启https或客户端认证需要配置cert文件路径 
  certFile: ""
  
  # 如果开启https或客户端认证需要配置key文件路径 
  keyFile: ""

grpc:

  # grpc 地址
  addr: ":6002"
  
  # 开启tls
  enableTls: false
  
  # 开启tls 需要配置cert文件路径
  certFile: ""
  
  # 开启tls 需要配置key文件路径
  keyFile: ""  
  
  # 开启tls 需要配置server name
  serverName: ""

log:
  # 日志输出级别 
  level: "info|warning|error|debug"

  # 日志输出控制: console对应终端, file为文件
  output:  "console|file"

  # 日志文件输出目录, 当output开启file时生效 默认为: /var/log/unicorn
  filePath: "./"

  # 日志默认格式为 unicorn-2006-01-02.log
  # 日志文件前缀 默认为: unicorn
  filePrefix: ""

  # 日志文件后缀 默认为: .log
  fileSuffix: ""

  # 是否开启日志分割, 分割方式按天
  split: true

```


**使用环境变量**
如果不指定配置文件, 也可以使用环境变量

下面是环境变量列表

```bash
# etcd
UNICORN_ETCD_CLUSTER # etcd集群 格式为: http://127.0.0.1:2379;http://127.0.0.1:2380;
UNICORN_ETCD_TLS     # etcd是否开启tls 格式为: true/false
UNICORN_ETCD_INSECURE # etcd是否跳过配置文件 格式为: true/false
UNICORN_ETCD_CLIENT_AUTH #etcd是否开启客户端认证 格式为: true/false
UNICORN_ETCD_CA_FILE_PATH
UNICORN_ETCD_CERT_FILE_PATH
UNICORN_ETCD_KEY_FILE_PATH
UNICORN_ETCD_TIMEOUT   #etcd客户端连接超时时间, 单位是秒
UNICORN_ETCD_REPORT
UNICORN_ETCD_LOCAL_REPORT
UNICORN_ETCD_LOCAL_REPORT_FILE

# grpc
UNICORN_GRPC_ADDR
UNICORN_GRPC_TLS
UNICORN_GRPC_CERT_FILE_PATH
UNICORN_GRPC_KEY_FILE_PATH
UNICORN_GRPC_SERVER_NAME

# http
UNICORN_HTTP_ADDR
UNICORN_HTTP_TLS
UNICORN_HTTP_CLIENT_AUTH
UNICORN_HTTP_CERT_FILE_PATH
UNICORN_HTTP_KEY_FILE_PATH
UNICORN_HTTP_CA_FILE_PATH

# id
UNICORN_EPOCH
UNICORN_MACHINE_ID_TYPE
UNICORN_MACHINE_ID
UNICORN_MACHINE_IP

# log
UNICORN_LOG_LEVEL
UNICORN_LOG_OUTPUT
UNICORN_LOG_FILE_PATH
UNICORN_LOG_FILE_PREFIX
UNICORN_LOG_FILE_SUFFIX
UNICORN_LOG_SPLIT

```

# 部署到本地

**克隆项目**

git clone git@github.com:Waitfantasy/unicorn.git

**编译项目**

windows用户:

```bash
cd cmd

go build -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -o ./bin/unicorn ./cmd
```
linux用户可以直接使用Makefile文件进行编译


**运行项目**

```bash

# 使用配置文件
unicorn -config path/unicorn.yaml

# 也可以配置好环境变量直接运行
unicorn
```


# 部署到docker

**克隆项目**

git clone git@github.com:Waitfantasy/unicorn.git

**编译项目**

> 注意编译到docker中的是二进制文件, 需要单独编译alpine分发版本

`make alpine`

**编译docker镜像**

```bash
# 进入Dockerfile所在目录, 运行
docker build -t youname/unicorn:latest .
```

**启动docker**

> 如果使用配置文件,无需将配置文件打包进docker, 可以直接使用环境变量

```bash
docker run -d -v path/unicorn.yaml:path/unicorn.yaml -e UNICORN_CONF=path/unicorn.yaml -p 6001:6001 -p 6002:6002 --name unicorn unicorn
```

> 如果不使用配置文件, 可以在运行容器时指定相关环境变量

