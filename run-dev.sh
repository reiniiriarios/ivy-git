#!/bin/sh

GODEBUG=gctrace=1 GOGC=100 GOMEMLIMIT=1000MiB wails dev
