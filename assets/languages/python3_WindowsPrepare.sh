#!/usr/bin/env bash

cd "$1"
pip3 install virtualenv
if ! [[ -d "env" ]]; then
	python -m venv venv
    source venv/Scripts/activate
fi
venv/Scripts/pip install -r requirements.txt