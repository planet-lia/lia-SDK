#!/usr/bin/env bash

cd "$1"
if ! [[ -d "venv" ]]; then
    pip3 install virtualenv
	python -m venv venv
fi
venv/Scripts/pip install -r requirements.txt