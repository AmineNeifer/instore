package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fake.com/instore/client"
)

const (
	bCmd string = "[cmd]  "
	bUsr string = "[usr]> "
)

// if exitSign is present, the program terminates
var exitSigns = []string{"exit", "quit", "bye", "goodbye"}

// commands available to be used
var commands = []string{"use", "add", "addcsv", "get", "getall", "getk", "remove", "removecsv"}

// if a dbsign is present after command `use` the following commands are done on MongoDB otherwise on csvfile
var dbSigns = []string{"db", "database", "data-base", "mongo", "mongod", "mongodb"}

// storeType = "db" means that we are going to use mongodb, otherwise csv
var storeType = "csv"

func main() {
	// initializeing filePath to for later use
	var _, b, _, _ = runtime.Caller(0)
	// basepath = "/path/to/command" the directory in which we find command.go
	var basepath = filepath.Dir(b)
	// filePath = /path/to/command/csvFiles/ the directory in which we find csv file to be used by client
	var filePath = basepath + "/csvFiles/"

	// First screen to be shown to the user
	fmt.Println(bCmd + "Hello! I am excited to have you as a user :D Enjoy!")
	fmt.Println(bCmd + "You could use any of these commands")
	fmt.Println(bCmd + "----------------------------------------------------------------------------------------------")
	fmt.Println(bCmd + "add       <key>  <value>  ------- to add a key-value pair to the store")
	fmt.Println(bCmd + "get[v]    <key>           ------- to get existing values corresponding to a key from the store")
	fmt.Println(bCmd + "getk      <value>         ------- to get existing keys corresponding to a value from the store")
	fmt.Println(bCmd + "remove    <key>  <value>  ------- to remove a key-value pair from the store")
	fmt.Println(bCmd + "exit/quit                 ------- to quit the program")
	fmt.Println(bCmd + "----------------------------------------------------------------------------------------------")

	// infinite loop (command line)
	for {
		// Get input fromo User
		fmt.Print(bUsr)
		reader := bufio.NewReader(os.Stdin)
		cmdString, err := reader.ReadString('\n')

		// in case error happend while reading user input
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		// Trim leading and trailing spaces and new-line char
		cmdString = strings.Trim(cmdString, " \n")

		// split input into keywords
		cmdKeys := strings.Split(cmdString, " ")

		// lowercase the command
		cmdKeys[0] = strings.ToLower(cmdKeys[0])

		// similar commands
		switch cmdKeys[0] {
		case "delete":
			cmdKeys[0] = "remove"
		case "getv":
			cmdKeys[0] = "get"
		default:
		}

		// check for exiting tokens
		if contains(cmdKeys[0], exitSigns) {
			fmt.Println(bCmd + "Bye Bye! :)")
			return
		}

		// check for length of the input
		if len(cmdKeys) > 3 {
			fmt.Println(bCmd + "Usage: <command> <key> <value>")
			continue
		}

		// check for empty line
		if strings.Trim(cmdKeys[0], " ") == "" {
			continue
		}

		// Commands execution
		switch cmdKeys[0] {
			// choose to use database or csv mode
			case "use":
				if len(cmdKeys) != 2 {
					fmt.Println(bCmd + "Usage: use [db|csv]")
				} else if contains(strings.ToLower(cmdKeys[1]), dbSigns) {
					if storeType == "db" {
						fmt.Println(bCmd + "Already using MongoDB mode")
						continue
					} else {
						storeType = "db"
						client.UseDb()
					}
				} else {
					if storeType == "csv" {
						fmt.Println(bCmd + "Already using CSV mode")
						continue
					} else {
						storeType = "csv"
						client.UseCsv()
					}
				}
			// add key-value pair
			case "add":
				if len(cmdKeys) != 3 {
					fmt.Println(bCmd + "Usage: add <key> <value>")
				} else if storeType == "db" {
					client.AddDb(cmdKeys[1], cmdKeys[2])
				} else {
					client.AddCsv(cmdKeys[1], cmdKeys[2])
				}
			// add key-value pairs from csv
			case "addcsv":
				if storeType == "db" {
					fmt.Println(bCmd + "Command still not implemented for MongoDB mode...")
				} else if len(cmdKeys) != 2 {
					fmt.Println(bCmd + "Usage: addcsv <csv_filename>")
				} else {
					// get path of csv file which should be found under client/csvFiles directory
					filePath += cmdKeys[1]
					_, err := os.Stat(filePath)
					// printing error in case file doesn't exist
					if os.IsNotExist(err) {
						fmt.Printf(bCmd+"%v\n", err)
						continue
					}
					client.AddCsvFromFile(filePath)
				}
			// get values by key
			case "get":
				if storeType == "db" {
					client.GetvDb(cmdKeys[1])
				} else if len(cmdKeys) != 2 {
					fmt.Println(bCmd + "Usage: get <key>")
				} else {
					client.GetvCsv(cmdKeys[1])
				}
			// get all key-value pair
			case "getall":
				if storeType == "db" {
					fmt.Println(bCmd + "Command still not implemented for MongoDB mode...")
				} else if len(cmdKeys) != 1 {
					fmt.Println(bCmd + "Usage: getall")
				} else {
					client.GetAllCsv()
				}
			// get keys by value
			case "getk":
				if storeType == "db" {
					client.GetkDb(cmdKeys[1])
				} else if len(cmdKeys) != 2 {
					fmt.Println(bCmd + "Usage: getk <value>")
				} else {
					client.GetkCsv(cmdKeys[1])	
				}
			// remove key-value pair
			case "remove":
				if len(cmdKeys) < 2 {
					println(bCmd + "Usage: remove <key> <value>")
				} else if strings.ToLower(cmdKeys[1]) == "all" {
					client.RemoveAllCsv()
				} else if len(cmdKeys) != 3 {
					println(bCmd + "Usage: remove <key> <value>")
				} else if storeType == "db" {
					client.RemoveDb(cmdKeys[1], cmdKeys[2])
				} else {
					client.RemoveCsv(cmdKeys[1], cmdKeys[2])
				}
			// remove key-value pairs by csv
			case "removecsv":
				if storeType == "db" {
					fmt.Println(bCmd + "Command still not implemented for MongoDB mode...")
					continue
				} else if len(cmdKeys) != 2 {
					fmt.Println(bCmd + "Usage: removecsv <csv_filename>")
				} else {
					// get path of csv file which should be found under client/csvFiles directory
					filePath += cmdKeys[1]
					_, err := os.Stat(filePath)
					// printing error in case file doesn't exist
					if os.IsNotExist(err) {
						fmt.Printf(bCmd+"%v\n", err)
						continue
					}
					client.RemoveCsvFromFile(filePath)
				}
			default:
				fmt.Println(bCmd + "'" + string(cmdKeys[0]) + "' is not a command!")
				fmt.Println(bCmd + "Please use one of the commands provided above!")
		}
		
	}
}

// contains is a function that checks whether token is in s or not
func contains(token string, s []string) bool {
	for _, v := range s {
		if strings.Contains(token, v) {
			return true
		}
	}
	return false
}
