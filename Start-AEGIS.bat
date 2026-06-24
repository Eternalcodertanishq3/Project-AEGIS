@echo off
title AEGIS Bootloader
color 0A

echo ===================================================
echo               PROJECT AEGIS 
echo         Offline Survival Computer
echo ===================================================
echo.

:: Check if aegis.exe exists in the current folder
if not exist "aegis.exe" (
    color 0C
    echo [FATAL ERROR] aegis.exe not found!
    echo Ensure this script is placed in the exact same folder as aegis.exe on your pendrive.
    echo.
    pause
    exit /b
)

echo [SYSTEM] Booting AEGIS Core...
echo [SYSTEM] Launching Command Center in default browser...
echo.

:: Open the default web browser to the local server
start "" http://localhost:8080

:: Start the actual AEGIS binary
aegis.exe

:: Keep window open if aegis crashes or closes
pause
