#!/bin/bash

source ../../utils/PrintMessage.sh


isOpenShiftInstalled=`docker --version | grep -c ^"Docker version"`
if [ "$isOpenShiftInstalled" -eq 1 ];
then 
	dockerIsRunnung=`docker ps | grep -c "CONTAINER ID"`
	if [ "$isOpenShiftInstalled" -eq 1 ];
	then 
		success_message "Docker is up and running"
	else
		fail_message "Docker is not running"
	fi
else
	fail_message "Docker is not installed"
fi