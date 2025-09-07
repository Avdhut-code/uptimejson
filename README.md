# uptimejson

A simple CLI tool written in Go that logs your system uptime into JSON files.  
Each run records the current uptime and appends it into a monthly JSON log.  
This makes it easy to track system usage, analyze uptime history, or integrate with other tools.

---

## Features

- Logs system uptime (hours + minutes).
- Optionally include the current date and/or timestamp in each entry.
- Stores logs in `~/.local/share/uptimejson/YYYY-MM.json` (one file per month).
- Keeps user settings in `~/.config/uptimejson/config.json`.
- Lightweight and fast — reads directly from `/proc/uptime`.

---

## Installation

Clone the repo and run the installer:

```bash
git clone https://github.com/Avdhut-code/uptimejson.git
cd uptimejson
./install.sh
```

---

## Usage 

```bash 
# Run witout no flags to see help :
uptimejson --help

# Log time : 
uptimejson --log

# Enable date and time in logs :
uptimejson --set-date=true --set-time=true

# Change log storage path :
uptimejson --set-path "$HOME/custom_path"

# Show Version
uptimejson --version

```

# Configuration

tool saves settings in config.json : 

```bash 
~/.config/uptimejson/config.json
````
