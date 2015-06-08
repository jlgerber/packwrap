package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
	"os"
)

// retrieve the looger and make it available in the main namespace.
var log = packwrap.GetLogger()

func main() {
	usage := `Usage: paw [options] <command> [<args>...]

Paw - PAckage Wrapper system, which provides a mechanism for defining clean, controlled
package environment upon launching an application.The system uses package manifest 
files which provide, among other things, a list of environmentvariable settings 
which get evaluated before executing the wrapped application in its own subprocess.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   
paw subcommands:
   list       List available packages.
   versions   List available versions for a package.
   run        Run a package.
   env        Print the environment for a package.
   print      Prints the manifest for a package version.
   template   Prints the manifest template.
   shell      Drop down into a subshell with appropriate environment.
   `

	args, err := docopt.Parse(usage, nil, true, "", true)

	if err != nil {
		fmt.Println("paw - problem with docopt")
		fmt.Println(err)
		os.Exit(1)
	}

	if args["--help"].(bool) == true {
		fmt.Println(usage)
		os.Exit(0)
	}

	cmd := args["<command>"].(string)
	cmdArgs := args["<args>"].([]string)

	// set the logging level if passed in
	processCommonArgs(args)

	log.Debug("paw - Arguments  ", cmdArgs)

	if err := runCommand(cmd, cmdArgs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runCommand - this function routes calls to the appropriate subcommand.
func runCommand(cmd string, args []string) (err error) {
	// argv := make([]string, 1)
	// argv[0] = cmd
	// argv = append(argv, args...)
	switch cmd {
	case "list":
		return pawList()
	case "versions":
		return pawVersions()
	case "run":
		return pawRun()
	case "env":
		return pawEnv()
	case "print":
		return pawPrint()
	case "shell":
		return pawShell()
	}
	return errors.New(fmt.Sprintf("%s is not a paw command. See 'paw help'", cmd))
}
