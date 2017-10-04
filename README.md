# <Name>

<Name> is a basic web-based file server written in Go. 

## Get it

To download <Name> to your `~/Downloads` folder, build, and run it:

```
cd ~/Downloads
git clone <link>
cd <link>/src
go build && ./src
```



For building, golang is required. 

## Features


## Changing settings

When first starting the program and accessing the web interface, the setup is triggered and the settings are stored.

After this, the settings cannot be changed from within the program. To still change them, delete `./json/settings.json`. Once it is deleted, restart the program. 

## Purpose

