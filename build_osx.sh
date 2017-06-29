#!/usr/bin/env bash
#export GOPATH=$(go env GOPATH)
#GOOS=windows GOARCH=amd64 go get -d ./...
GOOS=windows GOARCH=amd64 go install -ldflags="-H windowsgui -s -w" github.com/syncore/qclauncher
go get -u github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get -u github.com/jteeuwen/go-bindata/...
GOOS=windows GOARCH=amd64 go get -u github.com/lxn/walk
rm ../../bin/qclauncher.exe 2>/dev/null
rm ../../bin/qclauncher.log 2>/dev/null
rm qclauncher.exe 2>/dev/null
rm qclauncher.log 2>/dev/null
rm data.qcl 2>/dev/null
rm qclaunchargs.txt 2>/dev/null
rm data.qcl.lock 2>/dev/null
cp $GOPATH/bin/goversioninfo cmd/qclauncher/
cp $GOPATH/bin/go-bindata cmd/qclauncher/
chmod +x cmd/qclauncher/goversioninfo
chmod +x cmd/qclauncher/go-bindata
./cmd/qclauncher/goversioninfo -manifest "resources/qclauncher.manifest" -icon "resources/qclauncher.ico" -o "qclauncher_amd64.syso" -64 resources/versioninfo.json
cd cmd/qclauncher/
#chdir because ../../resources/img must match exactly
./go-bindata -pkg "resources" -o ../../resources/logo.go ../../resources/img
GOOS=windows GOARCH=amd64 go build -o qclauncher.exe -ldflags="-H windowsgui -s -w" main.go
rm qclauncher_amd64.syso 2>/dev/null
rm goversioninfo 2>/dev/null
rm go-bindata 2>/dev/null
mkdir -p ../../bin
cd ../../
mv cmd/qclauncher/qclauncher.exe bin/ 2>/dev/null
