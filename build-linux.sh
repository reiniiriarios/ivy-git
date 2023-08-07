#!/bin/sh

GOGC=100 GOMEMLIMIT=1000MiB wails build -platform linux/amd64
# todo: https://ubuntu.com/core/services/guide/snap-crafting
# todo: https://snapcraft.io/snapcraft
