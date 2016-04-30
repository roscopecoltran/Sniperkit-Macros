#!/bin/sh
# symlinks workaround for vendor directory
rm -f ./vendor/github.com/matthieudelaro/nut/persist;
rm -f ./vendor/github.com/matthieudelaro/nut/nvidia;
ln -s `pwd`/persist `pwd`/vendor/github.com/matthieudelaro/nut/persist;
ln -s `pwd`/nvidia `pwd`/vendor/github.com/matthieudelaro/nut/nvidia;
export GOPATH=`pwd`;
code -w .