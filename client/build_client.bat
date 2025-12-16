@echo off
echo Installing PyInstaller...
pip install pyinstaller

echo Building Client...
pyinstaller --noconsole --onefile --name "ModelManagerClient" client_tray.py

echo Build Complete! Executable is in dist/ModelManagerClient.exe
pause
