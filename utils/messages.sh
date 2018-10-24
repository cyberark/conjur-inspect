#!/bin/bash

function success_message(){
	GREEN='\33[0;32m'
	NC='\033[0m'
	printf "${GREEN} Success - $1${NC}\n"
}

function fail_message(){
	RED='\033[0;31m'
	NC='\033[0m'
	printf "${RED} Failed - $1${NC}\n"
}

function info_message(){
	NC='\033[0m'
	printf " Info - $1${NC}\n"
}