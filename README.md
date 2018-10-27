# unicorn

unicorn  is a distributed id generator, it has feature:

- 1 second can produce 1000000 id

- support grpc restful-api client

- CAS free lock


# Usage

## deploy to local

**step1:**
git clone git@github.com:Waitfantasy/unicorn.git

**step2:**
make && make install

**step3:**
now, you need to configure the configuration file

unicorn will read configure file, the default path is `/etc/unicorn/unicorn.yaml`

you can `cp example.yaml  /etc/unicorn/unicorn.yaml` or
specify the config file path when starting the unicorn: `unicorn -config path/unicorn.yaml`


## deploy docker

make alpine

docker build -t unicorn:latest ./

docker run -d -p 9001:9001 -p 9002:9002 -v /path/unicorn.yaml:/etc/unicorn/unicorn.yaml unicorn


# config file

## etcd config

**cluster**

cluster of etcd

**enableTls**

enable etcd tls authentication, when enableTls, you need configure client c**cert file** and **key file**

**certFile**

when enableTls, you configure it.

**keyFile**

when enableTls, you configure it.

**clientAuth**
enable authentication for the client, when enable clientAuth, you need configure **ca file**, client **cert file** and **key file**

**caFile**
when enable clientAuth, you configure it.

**report**
how many second report the timestamp of the machine to etcd, default 30 seconds

**timeout**
etcd client timeout, default 5 seconds.

## id config

**epoch**
epoch is  millisecond timestamp, it is recommended to be slower than the current time

**machineId**
not use etcd to get the allocation id. this value is used when machineIdType = 0

**machineIp**
according to machine ip, etcd service will assign the machine id

**machineIdType**
- 0 use machineId value

- 1 use etcd assign by machine ip

**idType**

- 0 using a second timestamp, the sequence is 20 bits and can generate 1 048 576 id per second.

- 1 using a million second timestamp, the sequence is 10 bits and can generate 1024 id per second.

**version**
0 or 1


# next plan
[x] support etcd configure center
