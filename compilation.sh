#!/bin/bash
#https://golang.org/doc/install/source
ALL_GOOS_GOARCH="\
darwin_386 \
darwin_amd64 \
dragonfly_386 \
dragonfly_amd64 \
freebsd_386 \
freebsd_amd64 \
freebsd_arm \
linux_386 \
linux_amd64 \
linux_arm \
netbsd_386 \
netbsd_amd64 \
netbsd_arm \
openbsd_386 \
openbsd_amd64 \
plan9_386 \
plan9_amd64 \
solaris_amd64 \
windows_386 \
windows_amd64"


#firt parameter GOROOT
#second parameter OS_ARCH
function compile(){
	GOOS=${2%_*}
	GOARCH=${2#*_}
	#build env
	CURRENT_DIR=${PWD}
	
	if [ ! -d $1/pkg/$2 ]; then
		#make.bash
		echo make
		cd $1/src
		GOOS=${GOOS} GOARCH=${GOARCH} ./make.bash -v --no-clean > /dev/null 2>&1
		if [ $? -eq 0 ]; then
			echo Yes we can $2
    		cd $CURRENT_DIR
		else
			echo No try again $2
    		cd $CURRENT_DIR
    		return 0	
			
		fi
		
	else 
		echo $2 exists
	fi


}

for GOOS_GOARCH in $ALL_GOOS_GOARCH; do
	compile $(go env GOROOT) $GOOS_GOARCH
done