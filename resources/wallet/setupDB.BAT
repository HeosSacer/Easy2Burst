
REM  --> Check for permissions
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"

REM --> If error flag set, we do not have admin.
if '%errorlevel%' NEQ '0' (
    echo Requesting administrative privileges...
    goto UACPrompt
) else ( goto gotAdmin )

:UACPrompt
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
    echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs"

    "%temp%\getadmin.vbs"
    exit /B

:gotAdmin
    if exist "%temp%\getadmin.vbs" ( del "%temp%\getadmin.vbs" )
    pushd "%CD%"
    CD /D "%~dp0"
:: BatchGotAdmin (Run as Admin code ends)

%~dp0/mysqld.exe --install easy_db --defaults-file=%~dp0/my-medium.ini
net start easy_db
%~dp0/mysql.exe -uroot --port=3399 < C:\Users\Heos\AppData\Roaming\Easy2Burst\MariaDB\bin\init-easy-db.sql

:: TO Replace string C:/Users/Heos/AppData/Roaming/Easy2Burst/BurstWallet/init-mysql.sql