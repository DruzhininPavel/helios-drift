#!/bin/sh
# Cross-compiles the release builds. Ebiten needs no cgo for these targets,
# so one machine can produce every binary.
set -e
cd "$(dirname "$0")/go"
mkdir -p ../dist

build() {
	os=$1 arch=$2 out=$3
	echo "building $os/$arch -> $out"
	GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
		-buildmode=exe -trimpath -ldflags="-s -w" -o "../dist/$out" .
}

build linux   amd64 heliosdrift-linux-amd64
build linux   arm64 heliosdrift-linux-arm64
build windows amd64 heliosdrift-windows-amd64.exe
build darwin  arm64 heliosdrift-darwin-arm64

# the shipped Windows build is a GUI binary: no console window behind the game
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
	-buildmode=exe -trimpath -ldflags="-s -w -H windowsgui" \
	-o ../dist/heliosdrift.exe .

echo "packaging"
cd ../dist
rm -rf pkg-linux pkg-win
mkdir -p pkg-linux pkg-win
cp heliosdrift-linux-amd64 pkg-linux/heliosdrift
cp heliosdrift.exe pkg-win/heliosdrift.exe
cp heliosdrift-windows-amd64.exe pkg-win/heliosdrift-console.exe
cp ../PACKAGE-README.txt pkg-linux/README.txt
cp ../PACKAGE-README.txt pkg-win/README.txt
if [ "$(uname -s)" = "Darwin" ]; then
	# bsdtar on macOS otherwise writes xattr headers GNU tar warns about
	COPYFILE_DISABLE=1 tar --no-mac-metadata -czf HeliosDrift-linux-x86_64.tar.gz -C pkg-linux .
else
	tar czf HeliosDrift-linux-x86_64.tar.gz -C pkg-linux .
fi
(cd pkg-win && zip -q ../HeliosDrift-windows-x86_64.zip heliosdrift.exe heliosdrift-console.exe README.txt)
ls -lh HeliosDrift-*.tar.gz HeliosDrift-*.zip
