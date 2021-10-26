#!/bin/sh

go build -o archinstall2 *.go &&
git add . &&
git commit -m "push command" &&
git push
