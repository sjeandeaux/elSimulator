#!/bin/bash
####################################################################################################
#
#TODO refactoring code
#
####################################################################################################
#https://golang.org/doc/install/source
ALL_GOOS_GOARCH="\
darwin_amd64 \
darwin_386 \
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

#execute GOROOT/src/make.bash
function prepare(){
	WHERE_GO=$(go env GOROOT)
	GOOS=${1%_*}
	GOARCH=${1#*_}
	#we have not package, we create
	if [ ! -d $WHERE_GO/pkg/$1 ]; then
		#make.bash
		echo make
		cd $WHERE_GO/src
		GOOS=${GOOS} GOARCH=${GOARCH} ./make.bash -v --no-clean 
		if [ $? -eq 0 ]; then
			echo Yes we can $1
    			cd $CURRENT_DIR
		else
			echo No try again $1
    			cd $CURRENT_DIR 	
    			return 1
		fi
	else 
		echo $1 exists
	fi
	return 0
}

#firt parameter GOROOT
#second parameter OS_ARCH
function compile(){
	#build env
	CURRENT_DIR=${PWD}
	prepare $1
	if [ $? -ne 0 ]; then
		return 1
	fi
	if [ ! -d bin/$2/ ]; then
		mkdir -p bin/$2/
	fi
	make build-output output=bin/$1/${CURRENT_DIR##*/}
	return 0
}

for GOOS_GOARCH in $ALL_GOOS_GOARCH; do
	compile $GOOS_GOARCH
	if [ $? -ne 0 ]; then
		echo Sh***t
		return 1
	fi
done
