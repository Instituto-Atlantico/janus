## How to use CLI in Windows environment

Having an ssh key par configured and already passed as authorized_keys on remote device is required. Need help with this? Click [here](https://phoenixnap.com/kb/ssh-with-key).

To use the CLI in Windows is necessary to copy the bin folder generated on Ubuntu via terminal and paste it in a Windows directory previously created, like for example Janus and access it via PowerShell and/or Command Prompt

#### How to deploy an agent on a remote device via PowerShell

```
.\bin\janus-cli_windows_amd64 deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
C:\Users\username\Desktop\Janus> .\bin\janus-cli_windows_amd64 deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

#### How to deploy an agent locally via PowerShell

```
.\bin\janus-cli_windows_amd64 deploy local --agent-port <port-number> --agent-name <agent-name>
```

#### How to deploy an agent on a remote device via Command Prompt

```
janus-cli_windows_amd64.exe deploy remote --agent-port <port-number> --agent-name <agent-name> -H <device-user>@<device-ip>
```

For example:

```
C:\Users\username\Desktop\Janus\bin> janus-cli_windows_amd64.exe deploy remote --agent-port 8001 --agent-name demo -H pi@192.168.0.2
```

#### How to deploy an agent locally via Command Prompt

```
janus-cli_windows_amd64.exe deploy local --agent-port <port-number> --agent-name <agent-name>
```

## Features docs

Read more about the proposed features [here](./docs/readme.md)
