package main

import (
	"bufio"
	"fmt"
	"github.com/jlgerber/packwrap"
	"os"
	"sort"
	"strings"
)

// processCommonArgs - modify the logging level based on args argument, which has been
// generated with docopt.
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

// stringInSlice - convenience function to determine if the supplied string is
// one of the elements in the supplied slice.
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

// createSubcmdRunner - Convenience function to create the subcmdRunner instance.
// This function does the following:
// * instantiates a new SubcmdRunner
// * registers subcommands
// * instantiates new ManifestReaderFactory
// * registers readers with factory
// * instantiates new manifestLocator passing in manifestReaderFactory.keys()
// * registers the manifestReaderFactory instance with the subcmdRunner
// * registers the manifestLocator with the subcmdRunner
func createSubcmdRunner() *SubcmdRunner {
	// create subcommand runner and register subcommands
	subcmdRunner := NewSubcmdRunner()
	subcmdRunner.RegisterSubcmd("list", "List available packages.", pawList)
	subcmdRunner.RegisterSubcmd("versions", "List available versions for a package.", pawVersions)
	subcmdRunner.RegisterSubcmd("run", "Run a package.", pawRun)
	subcmdRunner.RegisterSubcmd("env", "Print the environment for a package.", pawEnv)
	subcmdRunner.RegisterSubcmd("print", "Print the environment for a package.", pawPrint)
	subcmdRunner.RegisterSubcmd("shell", "Drop down into a subshell with appropriate environment.", pawShell)

	// create manifest reader factory and register reader instances
	manifestReaderFactory := packwrap.NewManifestReaderFactory()
	manifestReaderFactory.AddReader("json", &packwrap.JsonManifestReader{})

	// create manifest locator with appropriate extensions
	manifestLocator := packwrap.NewManifestLocator(manifestReaderFactory.Keys())

	// register manifest locator and manifest reader factory with subcommand runner.
	subcmdRunner.RegisterManifestLocator(manifestLocator)
	subcmdRunner.RegisterManifestReaderFactory(manifestReaderFactory)

	return subcmdRunner
}

// generateSubcmdString - helper function to generate the padded subcommand string used in the usage docs.
func generateSubcmdString(s *SubcmdRunner) string {
	maxlen := s.MaxNameLength()
	retstr := ""

	sortedKeys := s.Keys()
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		retstr += fmt.Sprintf("  %s%s    %s\n", key, strings.Repeat(" ", maxlen-len(key)), s.GetDesc(key))
	}
	return retstr
}
