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

	// fmt.Println(conf.Path)

	var logs []function.Data

	// if data, err := os.ReadFile(conf.Path); err == nil && len(data) > 0 {
	// 	_ = json.Unmarshal(data, &logs)
	// }

	secondValue, _ := os.ReadFile("/proc/uptime")

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

	setpath := flag.String("set-path", conf.Path, "path to the json file u can change this as well")

	setdate := flag.Bool("set-date", false, "include date in log entries")
	settime := flag.Bool("set-time", false, "include time in log entries")

	run_logs := flag.Bool("log", false, "initiats the logging.")

	version := flag.Bool("version", false, "shows the verion of code.")
	help := flag.Bool("help", false, "show help menu.")

	flag.Parse()

	Change := false

	// FOR DEBUG
	// args := flag.Args()
	// if len(args) < 1 {
	// 	fmt.Println(function.NoCommand)
	// 	os.Exit(1)
	// }

	if *version {
		fmt.Print(function.VersionValue)
		return
	}

	if *help {
		fmt.Print(function.HelpValue)
		return
	}

	// if *setpath != conf.Path {
	// 	conf.Path = function.ExpandPath(*setpath)
	// 	fmt.Println("New Path Saved:", conf.Path)
	// 	Change = true

	// 	// FOR DEBUG
	// 	// new_settings := function.Setting{Path: *logPath}
	// 	// function.SaveConfig(new_settings)
	// 	// fmt.Println("New Path Saved.", conf.Path)
	// }

	// if *setdate != conf.DateFlag {
	// 	conf.DateFlag = *setdate
	// 	Change = true

	// 	// FOR DEBUG
	// 	// fmt.Println(*setdate)
	// }

	// if *settime != conf.TimeFlag {
	// 	conf.TimeFlag = *settime
	// 	Change = true
	// 	// FOR DEBUG
	// 	// fmt.Println(*settime)
	// }

	if *setpath != conf.Path {
		conf.Path = function.ExpandPath(*setpath)
		function.SaveConfig(conf)
		fmt.Println("New Path Saved:", conf.Path)
		return
	}

	seen := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	if seen["set-date"] && *setdate != conf.DateFlag {
		conf.DateFlag = *setdate
		Change = true
	}

	if seen["set-time"] && *settime != conf.TimeFlag {
		conf.TimeFlag = *settime
		Change = true
	}

	if seen["set-path"] {
		newPath := function.ExpandPath(*setpath)
		if newPath != conf.Path {
			conf.Path = newPath
			Change = true
			fmt.Println("New Path set to:", conf.Path)
		}
	}
	if Change {
		function.SaveConfig(conf)
		fmt.Println("Saved config:", conf)
	}

	if *run_logs {

		year, month, _ := time.Now().Date()

		// filename creation for path the use
		filename := fmt.Sprintf("%d-%02d.json", year, int(month))

		// joining the path to the file tht way we go and make/check the file
		fullPath := filepath.Join(conf.Path, filename) // "2025-10.json")

		// FOR DEBUG
		// fmt.Println(fullPath)

		// if data, err := os.ReadFile(fullPath); err == nil && len(data) > 0 {
		// 	fmt.Println("we Unmarshaling the data")
		// 	_ = json.Unmarshal(data, &logs)
		// }

		//checkign the filepath joineds existance
		_, err = os.Stat(fullPath)
		if err != nil {
			// // makign sure that the file is gettin created after knowing tht it dosent existes
			// os.WriteFile(fullPath, []byte(""), 0644)
			// log.Fatal(err)
			// function.PrintCurrentLine()

			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			if err != nil {
				log.Fatal(err)
				function.PrintCurrentLine()
			}
		}

		// // Creating a empty file i k its risky to do it like this
		// os.WriteFile(fullPath, []byte(""), 0644)

		Struct_after_checking, err := function.CheckFields(seconds, conf)
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}

		// // FOR DEBUG
		// fmt.Println("passs the function ")

		data, err := os.ReadFile(fullPath)
		if err == nil && len(data) > 0 {
			_ = json.Unmarshal(data, &logs)
		}

		logs = append(logs, Struct_after_checking)

		jsonData, err := json.MarshalIndent(logs, "", "\t")
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}

		// err = os.MkdirAll(filepath.Dir(fullPath), 0755)
		// if err != nil {
		// 	log.Fatal(err)
		// 	function.PrintCurrentLine()
		// }

		err = os.WriteFile(fullPath, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
			function.PrintCurrentLine()
		}
	}

	// FOR DEBUG
	fmt.Println("End for test run...")
}
