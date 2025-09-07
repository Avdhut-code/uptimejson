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
sudo ./install.sh
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
---

# Configuration

tool saves settings in config.json : 
```bash 
~/.config/uptimejson/config.json
```

Example :
```json
{
    "path": "/home/user/.local/share/uptimejson",
    "dateflag": true,
    "timeflag": true
}

```
---

# Log files

Logs are stored in :
```bash 
~/.local/share/uptimejson/YYYY-MM.json
```

Example :
```json
{
    "hours": 12,
    "minutes": 45,
    "date": "2025-09-06",
    "time": "2025-09-06T18:42:04+05:30"
}
```

---

# Services (Optional but not recommended)

- If you want to run uptimejson automatically at startup or shutdown,
- use the provided make_service.sh along with the unit files in service_related/.

command to set-up it :
```bash 
./make_service.sh --user
systemctl --user enable --now uptimejson-startup.service
systemctl --user enable --now uptimejson-shutdown.service
```