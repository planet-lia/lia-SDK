#!/usr/bin/env bash

cd "$1"
pip3 install virtualenv
if ! [[ -d "venv" ]]; then
	python -m venv venv
    source venv/Scripts/activate
fi
venv/Scripts/pip install -r requirements.txt