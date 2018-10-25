:: qclauncher, syncore <syncore@syncore.org> 2017-2018
cd %GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src
if exist .\blffDistPkg rd .\blffDistPkg /S /Q
if exist .\blff-master rd .\blff-master /S /Q
if exist .\blff_src rd .\blff_src /S /Q
if not exist .\blffDistPkg md .\blffDistPkg
if exist .\blff_src\Blff.Console\bin rd .\blff_src\Blff.Console\bin /S /Q
if exist .\blff_src\Blff.Console\obj rd .\blff_src\Blff.Console\obj /S /Q
if exist .\blff_src\Blff.Win\bin rd .\blff_src\Blff.Win\bin /S /Q
if exist .\blff_src\Blff.Win\obj rd .\blff_src\Blff.Win\obj /S /Q
if exist .\blff_src\Blff.Lib\bin rd .\blff_src\Blff.Lib\bin /S /Q
if exist .\blff_src\Blff.Lib\obj rd .\blff_src\Blff.Lib\obj /S /Q
go build -o get_blff_src.exe get_blff_src.go
call "%GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src\get_blff_src.exe"
set blffSrcDir=%GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src\blff_src

:: NOTE: Building this will require Visual Studio or MSBuild
:: If you don't have MSBuild (or Visual Studio), download Build Tools For Visual Studio 2017 at:
:: https://visualstudio.microsoft.com/downloads/#build-tools-for-visual-studio-2017

:: Afterwards, uncomment the following line:
:: set msBuildDir=%programfiles(x86)%\Microsoft Visual Studio\2017\BuildTools\MSBuild\15.0\Bin

:: and comment out the following line:
set msBuildDir=%programfiles(x86)%\Microsoft Visual Studio\2017\Enterprise\MSBuild\15.0\Bin

call "%GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src\nuget.exe" restore %blffSrcDir%\blff.sln
call "%msBuildDir%\msbuild.exe" %blffSrcDir%\blff.sln /p:Configuration=Release /p:Platform=x86 /l:FileLogger,Microsoft.Build.Engine;logfile=blff_src_build.log
set msBuildDir=
set blffSrcDir=
XCOPY .\blff_src\Blff.Console\bin\DistPackage\*.zip .\blffDistPkg\
XCOPY .\blff_src\Blff.Win\bin\DistPackage\*.zip .\blffDistPkg\
call "%GOPATH%\src\github.com\syncore\qclauncher\resources\bin_src\get_blff_src.exe" -extractDistPkgs
if exist get_blff_src.exe del /Q get_blff_src.exe