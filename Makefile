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

all: deps build

deps:
	go get -d ./...
	(cd frontend && npm ci)

run:
	IVY_GIT_DEBUG=true GODEBUG=gctrace=1 GOGC=100 GOMEMLIMIT=1000MiB wails dev

build:
  ifeq ($(USER_OS), darwin)
		GOGC=100 GOMEMLIMIT=1000MiB wails build -platform darwin/universal
		pkgbuild --root build/bin --component-plist build/darwin/components.plist --identifier "me.reinii.ivy-git.pkg" --install-location /Applications ivy-git.pkg
		productbuild --package ivy-git.pkg "IvyGit_dev_Darwin_Universal.pkg"
  endif
  ifeq ($(USER_OS), linux)
		GOGC=100 GOMEMLIMIT=1000MiB wails build -platform linux/amd64 -o ivy-git
  endif
  ifeq ($(USER_OS), windows)
		wails build -platform windows/amd64 -nsis
		pwsh -noprofile -command Compress-Archive -Path "$PWD\build\bin\Ivy Git.exe" -DestinationPath "$PWD\IvyGit_dev_Windows_amd64.zip"
		pwsh -noprofile -command Move-Item -Path "$PWD\build\bin\Ivy Git-amd64-installer.exe" -Destination "$PWD\IvyGit_dev_Windows_amd64_installer.exe"
  endif
  ifeq ($(USER_OS),)
		echo "Unrecognized OS"
  endif

install:
  ifeq ($(USER_OS), darwin)
		build/bin/ivy-git.pkg
  endif
  ifeq ($(USER_OS), linux)
		install -v build/bin/ivy-git /usr/bin/ivy-git
		mkdir -p /etc/ivy-git
		install -v build/linux/ivy-git.desktop /usr/share/applications/ivy-git.desktop
		install -v build/appicon.png /etc/ivy-git/icon.png
  endif
  ifeq ($(USER_OS), windows)
		build/bin/IvyGit_dev_Windows_amd64_installer.exe
  endif
