@echo off
cd %GOPATH%\src\github.com\syncore\qclauncher
if exist %GOPATH%\src\github.com\syncore\qclauncher\bin rmdir /S /Q %GOPATH%\src\github.com\syncore\qclauncher\bin
go get github.com/kardianos/govendor
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get github.com/kevinburke/go-bindata/...
if exist %GOPATH%\bin\govendor.exe copy /Y %GOPATH%\bin\govendor.exe cmd\qclauncher\ >NUL
if exist %GOPATH%\bin\goversioninfo.exe copy /Y %GOPATH%\bin\goversioninfo.exe cmd\qclauncher\ >NUL
if exist %GOPATH%\bin\go-bindata.exe copy /Y %GOPATH%\bin\go-bindata.exe cmd\qclauncher\ >NUL
if exist qclauncher.exe del /Q qclauncher.exe
if exist qclauncher.log del /Q qclauncher.log
if exist data.qclr del /Q data.qclr
if exist qclaunchargs.txt del /Q qclaunchargs.txt
if exist data.qcl.lock del /Q data.qcl.lock
cd cmd\qclauncher
go generate
govendor build -o qclauncher.exe -ldflags="-H windowsgui -s -w"
::govendor build -o qclauncher.exe
if exist qclauncher_amd64.syso del /Q qclauncher_amd64.syso
if exist govendor.exe del /Q govendor.exe
if exist goversioninfo.exe del /Q goversioninfo.exe
if exist go-bindata.exe del /Q go-bindata.exe
if not exist ..\..\bin\ mkdir ..\..\bin
if exist qclauncher.exe move /Y qclauncher.exe ..\..\bin\ >NUL
cd ..\..\