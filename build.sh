# !/bin/sh

set -e

# Set env
if [ -z "$GOOS" ] 
then
	export GOOS=linux
	export GOARCH=arm64
fi

function buildDeb() {
	if [ -z "$version" ]
	then
		echo "environment varialble \"version\" is not set"
		exit 1
	fi

	# Set version
	sed -i "s/{version}/${version}/g" ./build-deb/DEBIAN/control


	# Build and Package

	go build


	rm -fr ./build-deb/usr/local/bin/.gitkeep
	rm -fr ./build-deb/usr/local/bin/*
	cp hm-diag ./build-deb/usr/local/bin

	dpkg-deb --root-owner-group --build ./build-deb hm-diag_${GOOS}_${GOARCH}.deb
	mv hm-diag hm-diag_${GOOS}_${GOARCH}
}

function pack() {
	cd ./web \
		&& yarn install \
		&& yarn run release \
		&& cd .. \
		&& go build \
		&& upx ./hm-diag	
}

case $1 in
	'' )
		buildDeb ;;
	pack )
		pack ;;
	* )
		echo "unknown subcommand"
		exit 1 
		;;
esac
