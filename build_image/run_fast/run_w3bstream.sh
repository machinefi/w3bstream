#!/bin/bash
#
if command -v docker-compose >/dev/null 2>&1; then 
  echo 'docker-compose command exists!' 
  if [ ! -f ./docker-compose.yaml ];then
     echo 'download docker-compose.yaml'
     curl -o docker-compose.yaml http://35.227.150.243/docker-compose.yaml
  fi
  if [ ! -f ./build_image.tgz ];then
     echo 'download build_image.tgz'
     curl -o build_image.tgz http://35.227.150.243/build_image.tgz
     tar -xzvf ./build_image.tgz 
  fi
else 
  echo 'docker-compose command does not exist, it needs to be installed!' 
  exit
fi

start(){
	echo "start w3bstream server"
        docker-compose -f ./docker-compose.yaml up -d
	echo "----------------"
}

stop() {
	echo "stop w3bstream server"
	docker-compose -f ./docker-compose.yaml down
	echo "----------------"
}

restart() {
	echo "restart w3bstream server"
	echo "----------------"
	stop
	start
}

case "$1" in
	start )
		echo "****************"
		start
		echo "****************"
		;;
	stop )
		echo "****************"
		stop
		echo "****************"
		;;
	restart )
		echo "****************"
		restart
		echo "****************"
		;;
	* )
		echo "****************"
		echo "no command"
		echo "****************"
		;;
esac
