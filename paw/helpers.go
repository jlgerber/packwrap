package main

import (
	"bufio"
	"os"
)

func processCommonArgs(args map[string]interface{}) {
	log.Debug("paw.processCommonArgs - ", args)
	// set the logging level if passed in
	if args["--loglevel"] == nil {
		log.SetLevel("info")
	} else {
		log.SetLevel(args["--loglevel"].(string))
	}

	if args["--debug"].(bool) == true {
		log.SetLevel("debug")
	}

	if args["--quiet"].(bool) == true {
		log.SetLevel("error")
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// readLines - helper function to slurp in a text file and return a list of
// lines, a la python.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
