#!/bin/sh

APP="ftpserver"
APPPath="./"
APP_PID=`pgrep ${APP}`
RunArg="-log.dir=./logs/"

case "$1" in
    -h|-?|h|help)
		printUsage
		;;
	start)
		if [ $# -lt 1 ]; then
			echo "Invalid number of arg"
			exit 1
	    fi
	    if [ -e ${APPPath}${APP} ]; then
		    ${APPPath}${APP} ${RunArg} 1>>/dev/null 2>&1 &
            echo "start ${APP} ok..."
		else
		    echo "start error: ${APPPath}${APP} not exist"
		    exit 3
		fi
		;;
	stop)
		if [ $# -lt 1 ]; then
			echo "Invalid number of arg"
			exit 1
		fi
	    if [ -e /proc/${APP_PID} ]; then
	        kill ${APP_PID}
            echo "stop ${APP} ok..."
		else
		    echo "stop error: ${APP} alread stopped..."
		    exit 2
		fi
		;;
	*)
		echo "unknown Operation: $1"
		printUsage
		;;
esac

exit 0