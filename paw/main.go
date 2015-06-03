package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
	"os"
	"os/exec"
	"strings"
)

var log = packwrap.GetLogger()

func main() {
	usage := `Usage: paw [-h | --help] [-d | --debug] [-q | --quiet ] [-l | --loglevel=<level>] 
           <command> [<args>...]

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
   `

	args, err := docopt.Parse(usage, nil, false, "", false)

	if err != nil {
		fmt.Println("problem with docopt")
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := args["<command>"].(string)
	cmdArgs := args["<args>"].([]string)

	// set the logging level if passed in
	if level := args["--loglevel"].([]string); len(level) == 0 {
		log.SetLevel("info")
	} else {
		log.SetLevel(level[0])
	}

	if args["--debug"].(bool) == true {
		log.SetLevel("debug")
	}

	if args["--quiet"].(bool) == true {
		log.SetLevel("error")
	}

	log.Debug("Args: ", args)
	log.Info("SubCommand ", cmd)
	log.Info("Arguments  ", cmdArgs)

	if err := runCommand(cmd, cmdArgs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func runCommand(cmd string, args []string) (err error) {
	//argv := make([]string, 1)
	//argv[0] = cmd
	//argv = append(argv, args...)
	switch cmd {
	case "list":
		// subcommand is a function call
		return pawList(args)
	case "versions":
		// subcommand is a script
		return pawVersions(args)
	case "run":
		// subcommand is a script
		return pawRun(args)
	case "env":
		return printEnv(args)
	case "print":
		return printManifest(args)

	}
	return errors.New(fmt.Sprintf("%s is not a paw command. See 'paw help'", cmd))
}

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

func printManifest(args []string) error {
	if len(args) < 2 {
		err := errors.New("wrong number of arguements. paw run <package> <version>")
		return err
	}
	manifest, err := packwrap.GetManifestLocationFor(args[0], args[1])
	if err != nil {
		return err
	}
	contents, err := readLines(manifest)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println(manifest)
	fmt.Println("")

	for _, ln := range contents {
		fmt.Println(ln)
	}
	return nil
}

func pawList(args []string) error {
	lst := packwrap.GetPackageList()
	for _, pack := range lst {
		fmt.Println(pack)
	}
	return nil
}

func pawVersions(args []string) error {
	versions := packwrap.GetPackageVersions(args[0])
	if versions == nil {
		log.Info("No Package Versions Found for ", args[0])
		return nil
	}

	for _, version := range versions {
		fmt.Println(version)
	}

	return nil
}

func pawRun(args []string) error {
	if len(args) < 2 {
		err := errors.New("wrong number of arguements. paw run <package> <version>")
		return err
	}
	manifest, err := packwrap.NewManifestFor(args[0], args[1])
	if err != nil {
		return err
	}

	//err = manifest.Setenv()
	if err = manifest.Setenv(); err != nil {
		log.Fatal(err)

	}
	//_ = sp
	remainingArgs := []string{}

	if len(args) > 2 {
		remainingArgs = args[2:]
	}

	cmd := exec.Command(manifest.Name, remainingArgs...)

	log.Info("Running", manifest.Name, " Version: ", manifest.Version())
	log.Info(manifest.Name, remainingArgs)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)

	}

	outlines := strings.Split(out.String(), "\n")
	for ln := range outlines {
		log.Info(outlines[ln])
	}
	return nil
}

// printEnv - given a string slice of arguments to the env command, lad
// the package and print it out
func printEnv(args []string) error {
	if len(args) < 2 {
		err := errors.New("wrong number of arguements. paw env <package> <version>")
		return err
	}

	manifest, err := packwrap.NewManifestFor(args[0], args[1])

	if err != nil {
		return err
	}

	for _, value := range manifest.Getenv() {
		fmt.Println(value)
	}
	return nil

}
