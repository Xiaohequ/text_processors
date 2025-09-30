@echo off
setlocal enabledelayedexpansion

REM Fichier batch pour compiler l'application text_processors
REM Usage: build.bat [release]

REM Dossier de sortie "build" à la racine du projet
set OUTPUT_DIR=build
set OUTPUT_NAME=text_processors.exe

REM Créer le dossier build s'il n'existe pas
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

REM Vérifier le paramètre release
set GO_ARGS=build
if /i "%1"=="release" (
    echo Mode release activé - binaire optimisé
    set GO_ARGS=build -ldflags=-s -w
) else (
    echo Mode debug
)

REM Construire la commande complète
set FULL_COMMAND=go %GO_ARGS% -o %OUTPUT_DIR%\%OUTPUT_NAME% .

echo Execution: %FULL_COMMAND%
echo.

REM Lancer la compilation
%FULL_COMMAND%

REM Vérifier le résultat
if %ERRORLEVEL% neq 0 (
    echo ERREUR: La compilation a echoue avec le code %ERRORLEVEL%

	pause
    exit /b %ERRORLEVEL%
)

echo.
echo SUCCESS: Artifact genere: %OUTPUT_DIR%\%OUTPUT_NAME%
echo.

REM Afficher la taille du fichier
if exist "%OUTPUT_DIR%\%OUTPUT_NAME%" (
    for %%I in ("%OUTPUT_DIR%\%OUTPUT_NAME%") do echo Taille du fichier: %%~zI octets
)

cd /d "%~dp0"
start "" ".\build\text_processors.exe"
