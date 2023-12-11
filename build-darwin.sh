#!/bin/sh

GOGC=100 GOMEMLIMIT=1000MiB wails build -platform darwin/universal
pkgbuild --root build/bin --component-plist build/darwin/components.plist --identifier "me.reinii.ivy-git.pkg" --install-location /Applications ivy-git.pkg
productbuild --package ivy-git.pkg "IvyGit_dev_Darwin_Universal.pkg"
