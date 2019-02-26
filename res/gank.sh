#!/bin/bash

HOME=`pwd`
PID_HOME=${HOME}/pid
LOG_HOME=${HOME}/log

prepare(){
    mkdir -p ${PID_HOME}
    mkdir -p ${LOG_HOME}
}

KILL(){
    pidfile=$1   
    pid=$(cat $PID_HOME/$pidfile)

    echo "stopping $pidfile:$pid"
    kill -9 $pid
    if [ $? == 0 ]; then        
	while [ 1 == `ps --no-header $pid | wc -l` ]; do
            echo "sleeping ..."
            sleep 0.1
        done
    fi
    
    rm -rf $PID_HOME/$pidfile
}

start_gs(){
    if [ $# == 0 ]; then
        echo "game id or port not specified."
        return
    fi

    cfg=${4:-'gs.json'}
    jp=$3
    echo "starting=> gmid:$1 port:$2 cfg:$cfg"
    nohup dotnet ./gs.dll -entry launcher -c $cfg -db db.json -http http://+:$2/ -gid $1 $jp 2>&1 >${LOG_HOME}/gs_$1.log&
    echo $! > $PID_HOME/gs_$1.pid
}

stop_gs(){
    if [ $# == 0 ]; then
        echo "game id not specified."
        return
    fi

    KILL "gs_$1.pid"
}

start_all_gs(){
    start_gs 1001 9010 -jp
    start_gs 1002 9000 -jp
    start_gs 1003 9030 -jp
    start_gs 1004 9040 -jp
    start_gs 1005 9050 -jp

    start_gs 2001 9020
    start_gs 2002 9060    
    start_gs 2003 9070

    start_gs 3001 9080 -jp
    start_gs 3002 9090 -jp

    echo "gameservice started."
}

stop_all_gs(){
    if [ ! -d ${PID_HOME}  ]; then 
        echo "${PID_HOME} not exist."
        return
    fi

    for element in `ls ${PID_HOME}`
    do  
        KILL ${element}
    done
}

start_hall(){
    echo "starting hall..."
    nohup dotnet ./hall.dll -entry launcher -c hall.json -db db.json 2>&1 >${LOG_HOME}/hall.log&
    echo $! > $PID_HOME/hall.pid
    
    echo "hall started."
}

stop_hall(){    
    KILL "hall.pid"
}

start_client(){
    echo "starting client..."
    dotnet ./client.dll -entry launcher -c client.json
}

clear_log(){
    if [ $1 ]; then
        echo "removing ${LOG_HOME}/$1"
        rm -f ${LOG_HOME}/$1
    else
        echo "removing ${LOG_HOME}"
        rm -rf ${LOG_HOME}
    fi
}

logs(){
    if [ $1 ]; then
        tail ${LOG_HOME}/$1.log -f -n 1000
    else
        tail ${LOG_HOME}/*.log -f -n 1000
    fi
}

prepare
case "$1" in
    'start')
        start_hall
        start_all_gs
        ;;
    'stop')
        stop_hall
        stop_all_gs
        ;;
    'restart')
        stop_hall
        stop_all_gs
        sleep 5
        start_hall
        start_all_gs     
        ;;
    'start_hall')
        start_hall
        ;;
    'stop_hall')
        stop_hall
        ;;
    'restart_hall')
        stop_hall
        start_hall
        ;;
    'start_client')
        start_client
        ;;
    'start_gs')
        start_gs $2 $3 $4
        ;;
    'stop_gs')
        stop_gs $2
        ;;
    'start_all_gs')
        start_all_gs
        ;;
    'stop_all_gs')
        stop_all_gs
        ;;
    'restart_all_gs')
        stop_all_gs
        start_all_gs
        ;;
    'clear_log')
        clear_log $2
        ;;
    'logs')
        logs $2
        ;;
    *)
        echo "usage: $0 {start|stop|start_hall|stop_hall|restart_hall|start_gs|stop_gs|start_all_gs|stop_all_gs|restart_all_gs|clear_log|logs}"
        exit 1
        ;;
    esac

exit