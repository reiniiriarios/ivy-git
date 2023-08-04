#!/bin/sh

IVY_GIT_DEBUG=true GODEBUG=gctrace=1 GOGC=100 GOMEMLIMIT=1000MiB wails dev
