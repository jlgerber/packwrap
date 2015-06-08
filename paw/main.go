package main

import (
	"bufio"
	//"bytes"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
	"os"
	"os/exec"
	//"strings"
)

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
		fmt.Println("problem with docopt")
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

	log.Debug("Arguments  ", cmdArgs)

	if err := runCommand(cmd, cmdArgs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func processCommonArgs(args map[string]interface{}) {
	log.Debug("ARGs", args)
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

// runCommand - this function routes to the appropriate function
func runCommand(cmd string, args []string) (err error) {
	argv := make([]string, 1)
	argv[0] = cmd
	argv = append(argv, args...)
	switch cmd {
	case "list":
		// subcommand is a function call
		return pawList()
	case "versions":
		// subcommand is a script
		return pawVersions()
	case "run":
		// subcommand is a script
		return pawRun()
	case "env":
		return printEnv()
	case "print":
		return printManifest(argv)
	case "shell":
		return pawShell()

	}
	return errors.New(fmt.Sprintf("%s is not a paw command. See 'paw help'", cmd))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func pawShell() error {
	usage := `Usage: paw shell [options] <command> <version> [<args>...]

paw run - execute the supplied command and version. Alternatively, you may supply an auxiliary
command to run (cmd) in the versioned command's environment, which it gets from the manifest.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   -s, --shell=<shell>
   `

	args, _ := docopt.Parse(usage, nil, true, "", false)
	processCommonArgs(args)

	log.Debug("pawShell - args", args)

	if args["<command>"] == nil || args["<version>"] == nil {
		err := errors.New("pawShell - wrong number of arguements. paw run <package> <version>")
		return err
	}

	command := args["<command>"].(string)
	version := args["<version>"].(string)

	manifest, err := packwrap.NewManifestFor(command, version)
	if err != nil {

		return errors.New(fmt.Sprint(err.Error(), " args: ", command, " ", version))
	}

	if err = manifest.Setenv(); err != nil {
		log.Fatal(err)

	}

	var shell string
	if args["--shell"] == nil {
		log.Debug("pawShell - setting shell default to bash")
		shell = "bash"
	} else {
		shell = args["--shell"].(string)

		if stringInSlice(shell, packwrap.VALID_SHELLS) == false {
			log.Warningf("pawShell - %s is not a valid shell. invoking %s", shell, packwrap.DEFAULT_SHELL)
			shell = packwrap.DEFAULT_SHELL
		} else {
			log.Debugf("pawShell - shell set to: %s.",
				shell)
		}

	}
	// all supported shels take a -i flag to make them interactive
	callingargs := []string{"-i"}
	callingargs = append(callingargs, args["<args>"].([]string)...)

	cmd := exec.Command(shell, callingargs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)

	}

	return nil
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

// printManifest - Given the name and version of a particular executable, find its
// manifest and print its contents.
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

// pawList - List the packages in the system.
func pawList() error {
	usage := `Usage: paw list [options] 

paw list - list the packages and the paths to their respective manifests.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `
	args, _ := docopt.Parse(usage, nil, true, "", false)

	processCommonArgs(args)

	lst := packwrap.GetPackageList()
	fmt.Println()
	for _, pack := range lst {
		fmt.Println(pack)
	}
	return nil
}

// pawVersions - Lists package versions for a named package supplied as
// the first arugment.
func pawVersions() error {

	usage := `Usage: paw versions [options] <command> [<args>...]

paw versions - list the versions for the provided command.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `
	args, _ := docopt.Parse(usage, nil, true, "", false)

	if args["<command>"] == nil {
		log.Fatal("pawVersions - Need to provide a command to look up versions for")
	}
	command := args["<command>"].(string)
	versions := packwrap.GetPackageVersions(command)
	if versions == nil {
		log.Info("pawVersions - No Package Versions Found for ", command)
		return nil
	}
	fmt.Println()
	for _, version := range versions {
		fmt.Println(version)
	}

	return nil
}

// pawRun - Runs a version of an executable, as specified in the args input.
// Minimally, the args input consists of an executable name, and a version. This
// fucntion initializes the environment based on a manifest for the supplied package
// and version, and then executes it in a separate process.
func pawRun() error {
	usage := `Usage: paw run [options] <command> <version> [<args>...]

paw run - execute the supplied command and version. Alternatively, you may supply an auxiliary
command to run (cmd) in the versioned command's environment, which it gets from the manifest.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   -c, --cmd=<cmd>  execute a command using the environment from the supplied manifest environment. 
   `

	args, _ := docopt.Parse(usage, nil, true, "", false)
	log.Debug("pawRun args", args)

	if args["<command>"] == nil || args["<version>"] == nil {
		err := errors.New("pawRun - wrong number of arguements. paw run <package> <version>")
		return err
	}

	command := args["<command>"].(string)
	version := args["<version>"].(string)

	manifest, err := packwrap.NewManifestFor(command, version)
	if err != nil {

		return errors.New(fmt.Sprint("pawRiun - ", err.Error(), " args: ", command, " ", version))
	}

	processCommonArgs(args)

	//err = manifest.Setenv()
	if err = manifest.Setenv(); err != nil {
		log.Fatal(err)

	}
	//_ = sp
	remainingArgs := args["<args>"].([]string)

	cmd := exec.Command(manifest.Name, remainingArgs...)
	if runcmd := args["--cmd"]; runcmd != nil {
		runcmd := args["--cmd"].(string)
		log.Debugf("pawRun - exec.Command %s", runcmd)
		cmd = exec.Command(runcmd, remainingArgs...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("pawRun - executing ", manifest.Name, " Version: ", manifest.Version(), " args:", remainingArgs)

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)

	}

	return nil
}

// printEnv - given a string slice of arguments to the env command, lad
// the package and print it out
func printEnv() error {
	usage := `Usage: paw env [options] <command> <version> [<args>...]

paw env - execute the supplied command and version. Alternatively, you may supply an auxiliary
command to run (cmd) in the versioned command's environment, which it gets from the manifest.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `

	args, _ := docopt.Parse(usage, nil, true, "", false)
	log.Debug("printEnv -cargs", args)

	if args["<command>"] == nil || args["<version>"] == nil {
		err := errors.New("printEnv - wrong number of arguments. paw env <package> <version>")
		return err
	}

	command := args["<command>"].(string)
	version := args["<version>"].(string)

	manifest, err := packwrap.NewManifestFor(command, version)

	if err != nil {
		return err
	}

	for _, value := range manifest.Getenv() {
		fmt.Println(value)
	}
	return nil

}
