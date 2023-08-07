#!/usr/bin/env pwsh

$Env:GOGC = '100'
$Env:GOMEMLIMIT = '1000MiB'
wails build -platform windows/amd64 -nsis
Compress-Archive -Path "$PWD\build\bin\Ivy Git.exe" -DestinationPath "$PWD\IvyGit_dev_Windows_amd64.zip"
Move-Item -Path "$PWD\build\bin\Ivy Git-amd64-installer.exe" -Destination "$PWD\IvyGit_dev_Windows_amd64_installer.exe"
