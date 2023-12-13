# Determine OS/Processor
USER_OS = 
USER_PROCESSOR = 
ifeq ($(OS),Windows_NT)
USER_OS = windows
ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
USER_PROCESSOR = amd64
else
ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
USER_PROCESSOR = amd64
endif
ifeq ($(PROCESSOR_ARCHITECTURE),x86)
USER_PROCESSOR = x86
endif
endif
else
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
USER_OS = linux
endif
ifeq ($(UNAME_S),Darwin)
USER_OS = darwin
endif
UNAME_P := $(shell uname -p)
ifeq ($(UNAME_P),x86_64)
USER_PROCESSOR = amd64
endif
ifneq ($(filter %86,$(UNAME_P)),)
USER_PROCESSOR = x86
endif
ifneq ($(filter arm%,$(UNAME_P)),)
USER_PROCESSOR = arm
endif
endif

THIS_FILE := $(lastword $(MAKEFILE_LIST))

all: deps build package

deps:
	go get -d ./...
	go install github.com/wailsapp/wails/v2/cmd/wails@v2.6.0
	(cd frontend && npm ci)
	wails doctor

run:
	IVY_GIT_DEBUG=true GODEBUG=gctrace=1 GOGC=100 GOMEMLIMIT=1000MiB wails dev

build:
	@mkdir -p build
	ifeq ($(USER_OS), darwin)
		GOGC=100 GOMEMLIMIT=1000MiB wails build -platform darwin/universal
	endif
	ifeq ($(USER_OS), linux)
		GOGC=100 GOMEMLIMIT=1000MiB wails build -platform linux/amd64 -o ivy-git
	endif
	ifeq ($(USER_OS), windows)
		GOGC=100 GOMEMLIMIT=1000MiB wails build -platform windows/amd64 -nsis
	endif
	ifeq ($(USER_OS),)
		echo "Build error, unrecognized OS"
	endif

linux:
	@mkdir -p build
	GOGC=100 GOMEMLIMIT=1000MiB wails build -platform linux/amd64 -o ivy-git

darwin:
	@mkdir -p build
	GOGC=100 GOMEMLIMIT=1000MiB wails build -platform darwin/universal

windows:
	@mkdir -p build
	GOGC=100 GOMEMLIMIT=1000MiB wails build -platform windows/amd64 -nsis

package:
	ifeq ($(USER_OS), darwin)
		pkgbuild --root dist/bin --component-plist dist/darwin/components.plist --identifier "me.reinii.ivy-git.pkg" --install-location /Applications ivy-git.pkg
		productbuild --package ivy-git.pkg "IvyGit_dev_Darwin_Universal.pkg"
		mv ivy-git.pkg dist/bin/ivy-git.pkg
		mv IvyGit_dev_Darwin_Universal.pkg build/IvyGit_dev_Darwin_Universal.pkg
	endif
	ifeq ($(USER_OS), linux)
		cp dist/bin/ivy-git build/
		cp dist/linux/ivy-git.desktop build/
		cp dist/appicon.png build/
	endif
	ifeq ($(USER_OS), windows)
		pwsh -noprofile -command Compress-Archive -Path "$PWD\dist\bin\Ivy Git.exe" -DestinationPath "$PWD\build\IvyGit_dev_Windows_amd64.zip"
		pwsh -noprofile -command Move-Item -Path "$PWD\dist\bin\Ivy Git-amd64-installer.exe" -Destination "$PWD\build\IvyGit_dev_Windows_amd64_installer.exe"
	endif
	ifeq ($(USER_OS),)
		echo "Package error, unrecognized OS"
	endif

install:
	ifeq ($(USER_OS), darwin)
		build/IvyGit_dev_Darwin_Universal.pkg
	endif
	ifeq ($(USER_OS), linux)
		install -v build/ivy-git /usr/bin/ivy-git
		mkdir -p /etc/ivy-git
		install -v build/ivy-git.desktop /usr/share/applications/ivy-git.desktop
		install -v build/appicon.png /etc/ivy-git/icon.png
	endif
	ifeq ($(USER_OS), windows)
		build/IvyGit_dev_Windows_amd64_installer.exe
	endif
	ifeq ($(USER_OS),)
		echo "Install error, unrecognized OS"
	endif

clean:
	rm -rf dist/bin/*
	rm -rf build
