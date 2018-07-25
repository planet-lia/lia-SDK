#!/usr/bin/env bash

cd $1
pip3 install virtualenv
if ! [[ -d "env" ]]; then
    virtualenv env
    virtualenv --relocatable env
fi
env/bin/pip3 install -r requirements.txt