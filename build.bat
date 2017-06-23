@echo off
if exist bin rmdir /S /Q bin
go install -ldflags="-H windowsgui -s -w" github.com/syncore/qclauncher
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go install github.com/jteeuwen/go-bindata/...
if exist %GOPATH%\bin\goversioninfo.exe copy /Y %GOPATH%\bin\goversioninfo.exe cmd\qclauncher\ >NUL
if exist %GOPATH%\bin\go-bindata.exe copy /Y %GOPATH%\bin\go-bindata.exe cmd\qclauncher\ >NUL
if exist qclauncher.exe del /Q qclauncher.exe
if exist qclauncher.log del /Q qclauncher.log
if exist data.qclr del /Q data.qclr
if exist qclaunchargs.txt del /Q qclaunchargs.txt
if exist data.qcl.lock del /Q data.qcl.lock
cd cmd\qclauncher
go generate
go build -o qclauncher.exe -ldflags="-H windowsgui -s -w"
if exist qclauncher_amd64.syso del /Q qclauncher_amd64.syso
if exist goversioninfo.exe del /Q goversioninfo.exe
if exist go-bindata.exe del /Q go-bindata.exe
if not exist ..\..\bin\ mkdir ..\..\bin
if exist qclauncher.exe move /Y qclauncher.exe ..\..\bin\ >NUL
cd ..\..\