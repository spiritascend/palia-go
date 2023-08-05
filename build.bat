@echo off
setlocal

set "go_executable=go"
set "go_executable_name=palia-go.exe"
set "destination_dir=build"

if not exist %destination_dir% (
    echo Creating the "%destination_dir%" directory...
    mkdir %destination_dir%
)

set /p use_upx="Do you want to use UPX packing? (y/n): "
if /i "%use_upx%"=="y" (
    echo Using UPX packing...
    %go_executable% build -ldflags="-s" -o %destination_dir%\%go_executable_name% .
    upx --ultra-brute -9 -k %destination_dir%\%go_executable_name%
) else (
    echo Building without UPX packing...
    %go_executable% build -ldflags="-s" -o %destination_dir%\%go_executable_name% .
)

if %errorlevel% neq 0 (
    echo Build failed.
    exit /b
)

cd %destination_dir%
.\%go_executable_name%

pause
endlocal