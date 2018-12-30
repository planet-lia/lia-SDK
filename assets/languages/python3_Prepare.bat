@echo off

cd "%1"
IF exist "venv" (
    echo virtualenv exists in venv
) ELSE (
    echo Setting up environment. This may take some time but is only done once.
    python -m pip install virtualenv
    echo Installing virtual environment.
    python -m venv venv
    .\venv\Scripts\python -m pip install --upgrade pip
)
.\venv\Scripts\pip install -r requirements.txt
