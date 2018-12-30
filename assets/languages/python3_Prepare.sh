#!/usr/bin/env bash

cd "$1"
if ! [[ -d "venv" ]]; then
    echo "Setting up environment. This may take some time but is only done once."
    pip3 install virtualenv
    echo "Installing virtual environment."
    virtualenv -p python3 venv
    venv/bin/python -m pip install --upgrade pip
fi
venv/bin/pip install -r requirements.txt