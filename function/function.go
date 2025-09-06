package function

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// // just for testing perposes anyways works?.
// func Test() string {
// 	return "works"
// }

type Setting struct {
	Path     string `json:"path"`
	DateFlag bool   `json:"dateflag"`
	TimeFlag bool   `json:"timeflag"`
}

type Data struct {
	TimeHour   int       `json:"hours"`
	TimeMin    int       `json:"minutes"`
	Date       string    `json:"date"`
	ActualTime time.Time `json:"time"`
}

const (

	// this is just a place holder for now ill add the real shit real
	VersionValue = `
Version of code
	version : v0.1
`

	// this is important
	HelpValue = `
uptimeJson:
	-- A simple uptime claculating CLI tool tht creats a json file so that user can then extract the .json content and use it as he want.

Usage:

	uptimeJson [OPTIONS]

	uptimeJson --help
		-- at this point you know this

	uptimeJson setPath "~/path_url"
		-- used to set new "path" for the *.json file 
		-- by deffult its "~/.local/share/uptimelogger/uptime.json"
		-- always give absalute path althought if u still dont code handels that 

	uptimeJson setDate true
		-- used to set the incuding of "date" while the code is recording the date
		-- by defulte its true
	
	uptimeJson setTime true
		-- used to set the including of time while the code is recording the time 
		-- by defult its true

Options/Flags:

	-v, --version
		-- Shows the code verison 

	-h, --help
		-- shows this page to know about tool 

Examples:

	uptimeJson setPath "~/new_path_url"

		uptimeJson setDate false

	uptimeJson setTime false
`
	// this is used when no command is given
	NoCommand = `
Error : no command given

Use "uptimejson --help" for more information about that topic.`
)

// CheckFields takes the raw uptime in seconds and the current config settings,
// and returns a Data struct that respects those settings.
//
// Logic:
//  1. Converts uptime seconds into hours and minutes using HourMin().
//  2. Builds a new Data struct and always sets TimeHour and TimeMin.
//  3. If DateFlag is true, adds the current date (YYYY-MM-DD).
//     Otherwise, leaves Date as an empty string ("").
//  4. If TimeFlag is true, adds the exact current timestamp.
//     Otherwise, leaves ActualTime as Go's zero-value (0001-01-01T00:00:00Z).
//  5. Returns the Data struct, or an error if the input seconds is invalid.
//
// This design guarantees the JSON always has consistent fields
// (hours, minutes, date, time), but the date/time fields may be "empty"
// when disabled by user config. This makes the output predictable
// and easy for tools/scripts to consume.
func CheckFields(seconds float64, conf Setting) (Data, error) {

	if seconds < 0 {
		return Data{}, fmt.Errorf("invalid uptime value: %f", seconds)
	}

	hour, min := HourMin(seconds)
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	currentDate := fmt.Sprintf("%d-%d-%d", year, int(month), day)

	newStruct := Data{
		TimeHour: hour,
		TimeMin:  min,
	}
	// // FOR DEBUG
	// fmt.Println("we don here on line")

	if conf.DateFlag {
		newStruct.Date = currentDate
	} else {
		newStruct.Date = " "
	}
	if conf.TimeFlag {
		newStruct.ActualTime = time.Now()
	} else {
		newStruct.ActualTime = time.Time{}
	}

	// // FOR DEBUG
	// fmt.Println(conf)

	return newStruct, nil
}

// ExpandPath takes a filesystem path and resolves it to an absolute path.
// - Expands '~' to the current user's home directory.
// - Returns the path unchanged if it's already absolute (starts with '/').
// - Converts relative paths into absolute ones based on the current working directory.
// This ensures all paths used in the program are safe and normalized
func ExpandPath(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, path[1:])
	}
	// if it's already absolute (starts with "/"), just return it
	if filepath.IsAbs(path) {
		return path
	}
	// otherwise, make it absolute relative to current dir
	abs, _ := filepath.Abs(path)
	return abs
}

// HourMin converts uptime in seconds into hours and minutes.
// Example: 3700 seconds -> 1 hour, 1 minute.
func HourMin(seconds float64) (int, int) {
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	return h, m
}

// LoadConfig loads user settings from "~/.config/uptimejson/config.json".
// If the config file does not exist, default settings are returned.
// Automatically expands "~" in paths to the user's home directory.
func LoadConfig() Setting {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "uptimejson", "config.json")

	config := Setting{
		Path:     filepath.Join(os.Getenv("HOME"), ".local", "share", "uptimejson"),
		DateFlag: false,
		TimeFlag: false,
	}

	if data, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(data, &config)
	} else {
		// config.json missing → create one with defaults
		_ = os.MkdirAll(filepath.Dir(configPath), 0755)
		b, _ := json.MarshalIndent(config, "", "\t")
		_ = os.WriteFile(configPath, b, 0644)
	}
	config.Path = ExpandPath(config.Path)
	return config
}

// SaveConfig saves the provided Setting struct to ~/.config/uptimejson/config.json.
// Creates the config directory if it does not exist.
// Pretty-prints JSON for easier manual editing.
func SaveConfig(C Setting) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "uptimejson", "config.json")

	// ensure parent directory exists
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Fatal(err)
		PrintCurrentLine()
	}

	data, err := json.MarshalIndent(C, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// PrintCurrentLine prints the current source file and line number to stdout.
// Useful for debugging to trace exactly where the code is running.
func PrintCurrentLine() {
	_, file, line, ok := runtime.Caller(0) // 0 indicates the current function's caller
	if ok {
		fmt.Printf("Current line in code: %s:%d\n", file, line)
	}
}
