#!/bin/bash

source "${BASH_SOURCE%/*}/../../utils/PrintMessage.sh"

ensure_docker() {
	isDockerInstalled=`docker --version | grep -c ^"Docker version"`
	if [ "$isDockerInstalled" -eq 1 ];
	then 
		dockerIsRunning=`docker ps | grep -c "CONTAINER ID"`
		if [ "$dockerIsRunning" -eq 1 ];
		then 
			success_message "Docker is up and running"
		else
			fail_message "Docker is not running"
		fi
	else
		fail_message "Docker is not installed"
	fi
}