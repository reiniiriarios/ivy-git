#!/usr/bin/env pwsh

wails build -platform windows/amd64 -nsis
Compress-Archive -Path "$PWD\build\bin\Ivy Git.exe" -DestinationPath "$PWD\IvyGit_dev_Windows_amd64.zip"
Move-Item -Path "$PWD\build\bin\Ivy Git-amd64-installer.exe" -Destination "$PWD\IvyGit_dev_Windows_amd64_installer.exe"
