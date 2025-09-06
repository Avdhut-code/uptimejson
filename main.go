package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Avdhut-code/function"
)

func main() {

	// FOR DEBUG
	// function.Test()
	// function.PrintCurrentLine()

	conf := function.LoadConfig()

	seconds := function.GiveSeconds()

	// FOR DEBUG
	// fmt.Println(conf.Path)
	// fmt.Println(seconds)

	var logs []function.Data

	Change := false

	setpath := flag.String("set-path", conf.Path, "path to the json file u can change this as well")

	run_logs := flag.Bool("log", false, "initiats the logging.")

	setdate := flag.Bool("set-date", false, "include date in log entries")
	settime := flag.Bool("set-time", false, "include time in log entries")

	version := flag.Bool("version", false, "shows the verion of code.")
	help := flag.Bool("help", false, "show help menu.")

	flag.Parse()

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
	// 	function.SaveConfig(conf)
	// 	fmt.Println("New Path Saved:", conf.Path)
	// 	return
	// }

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
		fullPath := filepath.Join(conf.Path, filename) //  FOR DEBUG "2025-10.json")

		// FOR DEBUG
		// fmt.Println(fullPath)

		//checkign the filepath joineds existance
		_, err := os.Stat(fullPath)
		if err != nil {

			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			if err != nil {
				log.Fatal(err)
				function.CurrentLine()
			}

			// makign sure that the file is gettin created after knowing tht it dosent existes
			os.WriteFile(fullPath, []byte(""), 0644)
		}

		// // Creating a empty file i k its risky to do it like this
		// os.WriteFile(fullPath, []byte(""), 0644)

		Struct_after_checking, err := function.CheckFields(seconds, conf)

		if err != nil {
			log.Fatal(err)
			function.CurrentLine()
		}

		// // FOR DEBUG
		// fmt.Println("passs the function ")

		data, err := os.ReadFile(fullPath)
		if err == nil && len(data) > 0 {
			// currentl empty data [-] so it formas a arry in file
			_ = json.Unmarshal(data, &logs)
		}

		// then we append it to the emty array then we marshal it and then write
		logs = append(logs, Struct_after_checking)

		jsonData, err := json.MarshalIndent(logs, "", "\t")
		if err != nil {
			log.Fatal(err)
			function.CurrentLine()
		}

		// err = os.MkdirAll(filepath.Dir(fullPath), 0755)
		// if err != nil {``
		// 	log.Fatal(err)
		// 	function.PrintCurrentLine()
		// }

		err = os.WriteFile(fullPath, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
			function.CurrentLine()
		}
	}

	// FOR DEBUG
	fmt.Println("log Registred...")
	// fmt.Println("End...")
}
