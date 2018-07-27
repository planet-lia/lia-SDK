#!/usr/bin/env bash

cd $1
pip3 install virtualenv
if ! [[ -d "env" ]]; then
    virtualenv -p python3 venv
    virtualenv --relocatable venv
    source venv/bin/activate
fi
venv/bin/pip install -r requirements.txt