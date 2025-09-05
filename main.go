package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Avdhut-code/function"
)

func main() {

	// FOR DEBUG
	// function.Test()
	// function.PrintCurrentLine()

	// file_path_To_data := filepath.Join(os.Getenv("HOME"), ".local", "share", "uptimejson", "timestamp.json")

	// defaultPath := function.ExpandPath(file_path_To_data)

	conf := function.LoadConfig()

	var logs []function.Data

	secondValue, _ := os.ReadFile("/proc/uptime")

	data, err := os.ReadFile(conf.Path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &logs)
	}

	fields := strings.Fields(strings.TrimSpace(string(secondValue)))
	if len(fields) < 1 {
		fmt.Println("invalid /proc/uptime format")
		return
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		log.Fatal(err)
		function.PrintCurrentLine()
	}

	currentTime := time.Now()

	hour, min := function.HourMin(seconds)

	year, month, day := currentTime.Date()

	variable := function.Data{
		TimeHour:   hour,
		TimeMin:    min,
		Date:       fmt.Sprintf("%d-%d-%d", year, int(month), day),
		ActualTime: currentTime,
	}

	setpath := flag.String("set-path", conf.Path, "path to the json file u can change this as well")

	setdate := flag.Bool("set-date", true, "include date in log entries")
	settime := flag.Bool("set-time", true, "include time in log entries")

	run_logs := flag.Bool("log", true, "initiats the logging.")
	version := flag.Bool("version", false, "shows the verion of code.")
	help := flag.Bool("help", false, "show help menu.")

	flag.Parse()

	// FOR DEBUG
	// args := flag.Args()
	// if len(args) < 1 {
	// 	fmt.Print(function.NoCommand, "\n")
	// 	os.Exit(1)
	// }

	if *version {
		fmt.Print(function.VersionValue)
	}

	if *help {
		fmt.Print(function.HelpValue)
	}

	if *setpath != conf.Path {
		conf.Path = function.ExpandPath(*setpath)
		function.SaveConfig(conf)
		fmt.Println("New Path Saved:", conf.Path)
		return

		// FOR DEBUG
		// new_settings := function.Setting{Path: *logPath}
		// function.SaveConfig(new_settings)
		// fmt.Println("New Path Saved.", conf.Path)
	}

	if *setdate != conf.DateFlag {

		conf.DateFlag = *setdate
		function.SaveConfig(conf)
		return

		// FOR DEBUG
		// fmt.Println(*setdate)
	}

	if *settime != conf.TimeFlag {

		conf.TimeFlag = *settime
		function.SaveConfig(conf)
		return

		// FOR DEBUG
		// fmt.Println(*settime)
	}

	if *run_logs {

		//  HERE BEFOR EDOING THE APPENDIG WE DO CHACKING OF WHAT TO ADD

		logs = append(logs, variable)

		jsonData, err := json.MarshalIndent(logs, "", " ")
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}

		err = os.MkdirAll(filepath.Dir(conf.Path), 0755)
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}

		err = os.WriteFile(conf.Path, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}

	}

	// FOR DEBUG
	fmt.Println("End for test run...")
}
