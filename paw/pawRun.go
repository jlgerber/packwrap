package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
	"os"
	"os/exec"
)

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
