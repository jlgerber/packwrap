package main

/*
pawShell - subcommand to invoke a subshell with an executable's environment
*/
import (
	"errors"
	//"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
	"os"
	"os/exec"
)

func pawShell(manifestLocator *packwrap.ManifestLocator,
	manifestReaderFactory *packwrap.ManifestReaderFactory) error {
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

	// find manifest
	location, err := manifestLocator.GetManifestLocationFor(command, version)
	if err != nil {
		return err
	}
	// get instance of manifest given a valid location
	manifest, err := manifestReaderFactory.NewManifestFor(location)
	if err != nil {
		return err
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
