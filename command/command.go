package main

import (
	"bufio"
	"fmt"
	"os"
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
var commands = []string{"add", "addcsv", "get", "getv", "getk", "remove", "delete", "removecsv"}

func main() {
	fmt.Println(bCmd + "Hello! I am excited to have you as a user :D Enjoy!")
	fmt.Println(bCmd + "You could use any of these commands")
	fmt.Println(bCmd + "-----------------------------------------------------------------------------")
	fmt.Println(bCmd + "add       <key>  <value>  ------- to add a key-value pair to the store")
	fmt.Println(bCmd + "get[v]    <key>           ------- to get existing values corresponding to a key from the store")
	fmt.Println(bCmd + "getk      <value>         ------- to get existing keys corresponding to a value from the store")
	fmt.Println(bCmd + "remove    <key>  <value>  ------- to remove a key-value pair from the store")
	fmt.Println(bCmd + "exit/quit                 ------- to quit the program")
	fmt.Println(bCmd + "-----------------------------------------------------------------------------")
	for {
		// Get input fromo User
		fmt.Print(bUsr)
		reader := bufio.NewReader(os.Stdin)
		cmdString, err := reader.ReadString('\n')
		// Trim leading and trailing spaces and new-line char
		cmdString = strings.Trim(cmdString, " \n")

		// split input into keywords
		cmdKeys := strings.Split(cmdString, " ")
		
		// lowercase the command
		cmdKeys[0] = strings.ToLower(cmdKeys[0])

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

		// check for existing commands
		if !contains(cmdKeys[0], commands) {
			fmt.Println(bCmd + "'" +string(cmdKeys[0]) + "' is not a command!")
			fmt.Println(bCmd + "Please use one of the commands provided above!")
			continue
		}

		// Commands execution
		if cmdKeys[0] == "add" {
			client.Add(cmdKeys[1], cmdKeys[2])
		} else if cmdKeys[0] == "addcsv" {
			_, err := os.Stat(cmdKeys[1])
			if os.IsNotExist(err) {
				fmt.Printf(bCmd+"%v\n", err)
				continue
			}
			client.AddCsv(cmdKeys[1])
		} else if cmdKeys[0] == "get" || cmdKeys[0] == "getv" {
			if strings.ToLower(cmdKeys[1]) == "all" {
				client.GetAll()
			} else {
				client.GetV(cmdKeys[1])
			}
		} else if cmdKeys[0] == "getk" {
			client.GetK(cmdKeys[1])
		} else if cmdKeys[0] == "remove" || cmdKeys[0] == "delete" {
			if  strings.ToLower(cmdKeys[1]) == "all" {
				client.RemoveAll()
			} else {
				client.Remove(cmdKeys[1], cmdKeys[2])
			}
		} else if cmdKeys[0] == "removecsv" {
			_, err := os.Stat(cmdKeys[1])
			if os.IsNotExist(err) {
				fmt.Printf(bCmd+"%v\n", err)
				continue
			}
			client.RemoveCsv(cmdKeys[1])
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
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
