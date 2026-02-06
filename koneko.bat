@echo off

if exist bin\koneko.exe del bin\koneko.exe

set CGO_ENABLED=0
go build -mod=vendor -ldflags "-H=windowsgui -s -w" -o bin/koneko.exe ./cmd/koneko

bin\koneko.exe