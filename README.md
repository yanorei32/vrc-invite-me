# vrc-invite-me

## Quick Start

### Install vrc-invite-me.exe
Ex: `C:\Program Files\vrc-invite-me\vrc-invite-me.exe`

### Write configure.yml and save application directory

#### Get cookie by your browser

![image](https://user-images.githubusercontent.com/11992915/88283162-df2d6680-cd25-11ea-8faa-54b680c45317.png)

Firefox: https://developer.mozilla.org/en/docs/Tools/Storage_Inspector#Cookies

Chromium: https://developers.google.com/web/tools/chrome-devtools/storage/cookies

EdgeHTML: https://docs.microsoft.com/en-us/microsoft-edge/devtools-guide/storage#cookies-manager

#### Create `configure.yml`

Create `C:\Program Files\vrc-invite-me\configure.yml`.

Example:
```yaml
apiKey: JlE5Jldo5Jibnk5O5hTx6XVqsJu4WJ26
auth: authcookie_xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### Replace launch.bat
Replace `C:\Program Files (x86)\Steam\steamapps\common\VRChat\launch.bat`.

Example:
```bat
rem cd /d %1
rem VRChat.exe %2

@echo off

"C:\Program Files\vrc-invite-me\vrc-invite-me.exe" %2

if not %ERRORLEVEL% == 0 (
    pause
)
```


