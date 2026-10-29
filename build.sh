#!/bin/sh
# Cross-compiles the release builds. Ebiten needs no cgo for these targets,
# so one machine can produce every binary.
set -e
cd "$(dirname "$0")"
ROOT=$(pwd)
cd go
mkdir -p ../dist

# version stamp: the tag when the tree is tagged, else the commit
if [ -z "$VERSION" ]; then
	VERSION=$(git -C "$ROOT" describe --tags --always --dirty 2>/dev/null || echo dev)
fi
echo "version $VERSION"
LD="-s -w -X heliosdrift/game.Version=$VERSION"

build() {
	os=$1 arch=$2 out=$3
	echo "building $os/$arch -> $out"
	GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
		-buildmode=exe -trimpath -ldflags="$LD" -o "../dist/$out" .
}

build linux   amd64 heliosdrift-linux-amd64
build linux   arm64 heliosdrift-linux-arm64
build windows amd64 heliosdrift-windows-amd64.exe
build darwin  arm64 heliosdrift-darwin-arm64

# the shipped Windows build is a GUI binary: no console window behind the game
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
	-buildmode=exe -trimpath -ldflags="$LD -H windowsgui" \
	-o ../dist/heliosdrift.exe .

echo "packaging"
cd ../dist
rm -rf pkg-linux pkg-win HeliosDrift.app
mkdir -p pkg-linux pkg-win
cp heliosdrift-linux-amd64 pkg-linux/heliosdrift
cp heliosdrift.exe pkg-win/heliosdrift.exe
cp heliosdrift-windows-amd64.exe pkg-win/heliosdrift-console.exe
cp ../PACKAGE-README.txt pkg-linux/README.txt
cp ../PACKAGE-README.txt pkg-win/README.txt

# macOS ships as a real .app: a bare Mach-O binary cannot open a window, and
# double-clicking one only opens a terminal
mkdir -p HeliosDrift.app/Contents/MacOS HeliosDrift.app/Contents/Resources
cp heliosdrift-darwin-arm64 HeliosDrift.app/Contents/MacOS/heliosdrift
cp ../packaging/Info.plist HeliosDrift.app/Contents/Info.plist
cp ../packaging/HeliosDrift.icns HeliosDrift.app/Contents/Resources/heliosdrift.icns
chmod +x HeliosDrift.app/Contents/MacOS/heliosdrift
if [ "$(uname -s)" = "Darwin" ]; then
	codesign --force --sign - HeliosDrift.app >/dev/null 2>&1 || true
fi

if [ "$(uname -s)" = "Darwin" ]; then
	# bsdtar on macOS otherwise writes xattr headers GNU tar warns about
	COPYFILE_DISABLE=1 tar --no-mac-metadata -czf HeliosDrift-linux-x86_64.tar.gz -C pkg-linux .
	COPYFILE_DISABLE=1 zip -qry HeliosDrift-macos-arm64.zip HeliosDrift.app
else
	tar czf HeliosDrift-linux-x86_64.tar.gz -C pkg-linux .
	find HeliosDrift.app -type f | xargs chmod 644
	chmod +x HeliosDrift.app/Contents/MacOS/heliosdrift
	zip -qry HeliosDrift-macos-arm64.zip HeliosDrift.app
fi
(cd pkg-win && zip -q ../HeliosDrift-windows-x86_64.zip heliosdrift.exe heliosdrift-console.exe README.txt)
ls -lh HeliosDrift-*.tar.gz HeliosDrift-*.zip
