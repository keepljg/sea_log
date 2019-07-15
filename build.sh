#!/usr/bin/env bash

#IP
log147="root@192.168.182.147"
log148="root@192.168.182.148"
log188="root@192.168.182.188"
log189="root@192.168.182.189"
basepath=$(cd `dirname $0`; pwd)


if [ $1 == "log147" ]
then
    echo "编译中..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sea_log slaver/cmd/main.go
    echo "rm -rf sea_log..."
    ssh -p 58422 $log147 << log
        mkdir -p /www/sea_log/slaver/cmd/
        mkdir -p /www/sea_log/slaver/conf/
        rm -rf /www/sea_log/slaver/cmd/sea_log
#        kill -9 `ps -ef | grep sea_log | grep -v "grep" | awk '{print $2}'`
        exit
log
    echo "上传sea_log..."
    scp -P 58422 $basepath/sea_log $log147:/www/sea_log/slaver/cmd/
    scp -P 58422 $basepath/slaver/conf/conf.ini $log147:/www/sea_log/slaver/conf/


elif [ $1 == "log148" ]
then
    echo "编译中..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sea_log slaver/cmd/main.go
    echo "rm -rf sea_log..."
    ssh -p 58422 $log148 << log
        mkdir -p /www/sea_log/slaver/cmd/
        mkdir -p /www/sea_log/slaver/conf/
        rm -rf /www/sea_log/slaver/cmd/sea_log
        exit
log
    echo "上传sea_log..."
    scp -P 58422 $basepath/sea_log $log148:/www/sea_log/slaver/cmd/
    scp -P 58422 $basepath/slaver/conf/conf.ini $log148:/www/sea_log/slaver/conf/


elif [ $1 == "log188" ]
then
    echo "编译中..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sea_log slaver/cmd/main.go
    echo "rm -rf sea_log..."
    ssh -p 58422 $log188 << log
        mkdir -p /www/sea_log/slaver/cmd/
        mkdir -p /www/sea_log/slaver/conf/
        rm -rf /www/sea_log/slaver/cmd/sea_log
        exit
log
    echo "上传sea_log..."
    scp -P 58422 $basepath/sea_log $log188:/www/sea_log/slaver/cmd/
    scp -P 58422 $basepath/slaver/conf/conf.ini $log188:/www/sea_log/slaver/conf/


elif [ $1 == "log189" ]
then
    echo "编译中..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sea_log slaver/cmd/main.go
    echo "rm -rf sea_log..."
    ssh -p 58422 $log189 << log
        mkdir -p /www/sea_log/slaver/cmd/
        mkdir -p /www/sea_log/slaver/conf/
        rm -rf /www/sea_log/slaver/cmd/sea_log
        exit
log
    echo "上传sea_log..."
    scp -P 58422 $basepath/sea_log $log189:/www/sea_log/slaver/cmd/
    scp -P 58422 $basepath/slaver/conf//conf.ini $log189:/www/sea_log/slaver/conf/

else
    echo "缺失部署环境参数！！！"
fi