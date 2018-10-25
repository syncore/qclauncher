@echo off
cd %GOPATH%\src\github.com\syncore\qclauncher
if exist %GOPATH%\src\github.com\syncore\qclauncher\bin rmdir /S /Q %GOPATH%\src\github.com\syncore\qclauncher\bin

:: NOTE: building the blff tool that is packaged with qclauncher will require Visual Studio or MSBuild (see resources\bin_src\README.md)
:: If you don't have MSBuild or Visual Studio, download Build Tools For Visual Studio 2017 at:
:: https://visualstudio.microsoft.com/downloads/#build-tools-for-visual-studio-2017
:: Afterwards, adjust the `msBuildDir` inside of resources\bin_src\build_blff_src.bat to: %programfiles(x86)%\Microsoft Visual Studio\2017\BuildTools\MSBuild\15.0\Bin
call "%GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src\build_blff_src.bat"

cd %GOPATH%\src\github.com\syncore\qclauncher
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
:: Comment the previous line and uncomment the following line to display errors/msg on stdout
::govendor build -o qclauncher.exe
if exist qclauncher_amd64.syso del /Q qclauncher_amd64.syso
if exist govendor.exe del /Q govendor.exe
if exist goversioninfo.exe del /Q goversioninfo.exe
if exist go-bindata.exe del /Q go-bindata.exe
if not exist ..\..\bin\ mkdir ..\..\bin
if exist qclauncher.exe move /Y qclauncher.exe ..\..\bin\ >NUL
cd ..\..\