package main

import (
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
`

	subcmdRunner := createSubcmdRunner()
	usage += generateSubcmdString(subcmdRunner)

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

	//if err := runCommand(cmd, cmdArgs); err != nil {
	if err := subcmdRunner.Run(cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
