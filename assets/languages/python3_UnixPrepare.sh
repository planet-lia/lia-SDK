#!/usr/bin/env bash

cd "$1"
if ! [[ -d "venv" ]]; then
    pip3 install virtualenv
    virtualenv -p python3 venv
fi
venv/bin/pip install -r requirements.txt