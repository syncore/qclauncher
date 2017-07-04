#!/usr/bin/env bash
export GOPATH=$(go env GOPATH)
cd $GOPATH/src/github.com/syncore/qclauncher
rm -rf $GOPATH/src/github.com/syncore/qclauncher/bin 2>/dev/null
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get github.com/jteeuwen/go-bindata/...
rm $GOPATH/src/github.com/syncore/qclauncher/qclauncher.exe 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/qclauncher.log 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/data.qcl 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/qclaunchargs.txt 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/data.qcl.lock 2>/dev/null
cp $GOPATH/bin/goversioninfo cmd/qclauncher/
cp $GOPATH/bin/go-bindata cmd/qclauncher/
chmod +x $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher/goversioninfo
chmod +x $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher/go-bindata
./cmd/qclauncher/goversioninfo -manifest "resources/qclauncher.manifest" -icon "resources/qclauncher.ico" -o "qclauncher_amd64.syso" -64 resources/versioninfo.json
cd $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher
#chdir because ../../resources/img must match exactly
./go-bindata -pkg "resources" -o ../../resources/logo.go ../../resources/img
GOOS=windows GOARCH=amd64 go build -o qclauncher.exe -ldflags="-H windowsgui -s -w" main.go
rm $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher/qclauncher_amd64.syso 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher/goversioninfo 2>/dev/null
rm $GOPATH/src/github.com/syncore/qclauncher/cmd/qclauncher/go-bindata 2>/dev/null
mkdir -p $GOPATH/src/github.com/syncore/qclauncher/bin
cd ../../
mv cmd/qclauncher/qclauncher.exe bin/ 2>/dev/null