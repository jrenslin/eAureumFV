# eAureumFV

eAureumFV offers a web-based interface for searching and previewing files. It's written in Go and (Vanilla) JavaScript.

## Features

- Keyboard-driven navigation
- Navigation to files using the given folder structure or their use cases (images, videos, plain text files)
- Search through files using regular expressions
- Previews using HTML5
- Built-in CBZ viewer

## Setup

To download eAureumFV to your `~/Downloads` folder, build, and run it:

```
cd ~/Downloads
git clone https://github.com/jrenslin/eAureumFV.git
cd eAureumFV
go build && ./eAureumFV
```

Next, connect to `localhost` on port `9090`, the default port. 

A setup page is triggered. Here you can set the port eAureumFV listens on and the folders it serves. 

## Changing settings

After the initial setup, the settings cannot be changed from within the program. To change them anyway, delete `./json/settings.json`. Once it is deleted, restart the program. 



