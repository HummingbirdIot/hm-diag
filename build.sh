# !/bin/sh

set -e

# Set env
if [ -z "$GOOS" ] 
then
	export GOOS=linux
	export GOARCH=arm64
fi

# Set version
sed -i "s/{version}/${version}/g" ./build-deb/DEBIAN/control


# Build and Package

go build

rm -fr ./build-deb/usr/local/bin/*
cp hm-diag ./build-deb/usr/local/bin

dpkg -b ./build-deb hm-diag_${GOOS}_${GOARCH}.deb
mv hm-diag hm-diag_${GOOS}_${GOARCH}