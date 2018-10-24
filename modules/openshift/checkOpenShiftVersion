#!/bin/bash

source ../../utils/PrintMessage.sh


isOpenShiftInstalled=`oc version | grep -c ^"oc"`
if [ "$isOpenShiftInstalled" -eq 1 ];
then 
	k8sVersion=`oc version | grep "kubernetes" | cut -d ' ' -f2`
	info_message "Kubernetes version is $k8sVersion"
	ocVersion=`oc version | grep "oc" | cut -d ' ' -f2`
	info_message "OpenShift version is $ocVersion"
else
	fail_message "OpenShift is not installed"
fi