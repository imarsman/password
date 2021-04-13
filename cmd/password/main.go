package main

import (
	_ "embed"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"unicode"

	"cuelang.org/go/pkg/strings"
	"github.com/foomo/htpasswd"
	"github.com/galdor/go-cmdline"
	"github.com/imarsman/password/cmd/password/internal/common"
	"github.com/jwalton/gchalk"
	"github.com/sethvargo/go-password/password"
)

//go:embed .appbuildts
var appBuildTS string

//go:embed .appbuildversion
var appBuildVersion string

func init() {
	// GetOptions()
	commandLine = readOpts()
}

// options options for commandline parameters
type options struct {
	File     string
	Add      bool
	Remove   bool
	Generate bool
	User     string
	Pass     string
	Once     bool
	Length   int
	Version  bool
}

// Opts instantiation of Options
var Opts options

//var start time.Duration
// var cmdLineRead bool

var commandLine = cmdline.New()

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func readOpts() *cmdline.CmdLine {
	commandLine := cmdline.New()
	if common.RunningTest() {
		var _ = func() bool {
			testing.Init()
			return true
		}()
		commandLine.AddFlag("a", "add", "add user")
		commandLine.AddFlag("r", "remove", "remove user")
		commandLine.AddFlag("g", "generate", "generate password and print it")
		commandLine.AddOption("l", "length", "length", "length of password to generate (15-30)")
		commandLine.AddOption("u", "user", "username", "username for call")
		commandLine.AddOption("p", "pass", "password", "password to set in file using bcrypt encryption")
		commandLine.AddOption("f", "file", "file", "path to password file")
		commandLine.AddFlag("v", "version", "print version information")

		commandLine.Parse(commandLine.CommandArguments)

		return commandLine
	}
	// fmt.Println("Reading options")
	// http://snowsyn.net/2016/08/11/parsing-command-line-options-in-go/

	// Library handles -h (prints usage)
	commandLine.AddFlag("a", "add", "add user")
	commandLine.AddFlag("r", "remove", "remove user")
	commandLine.AddFlag("g", "generate", "generate password and print it with default length of 15")
	commandLine.AddOption("l", "length", "length", "length of password to generate (15-30)")
	commandLine.AddOption("u", "user", "username", "username for call")
	commandLine.AddOption("p", "pass", "password", "password to set in file using bcrypt encryption")
	commandLine.AddOption("f", "file", "file", "path to password file")
	commandLine.AddFlag("v", "version", "print version information")

	commandLine.Parse(os.Args)

	Opts.Add = commandLine.IsOptionSet("add")
	Opts.Version = commandLine.IsOptionSet("version")

	Opts.Remove = commandLine.IsOptionSet("remove")
	Opts.Generate = commandLine.IsOptionSet("generate")

	if Opts.Version {
		fmt.Println()
		fmt.Println(gchalk.WithBlue().Bold("Password tool"))
		fmt.Println(gchalk.Green("Generate a password or update a password file using bcrypt encryption"))
		fmt.Println(gchalk.Green("Use -h flag to see options"))
		fmt.Println()
		fmt.Println(gchalk.BrightYellow("Version: \t") + gchalk.Green(strings.TrimSpace(appBuildVersion)))
		fmt.Println(gchalk.BrightYellow("Build: \t\t") + gchalk.Green(strings.TrimSpace(appBuildTS)))
		fmt.Println(gchalk.BrightYellow("Platform: \t") + gchalk.Green(runtime.GOOS))
		fmt.Println(gchalk.BrightYellow("Architecture: \t") + gchalk.Green(runtime.GOARCH))
		fmt.Println()
		os.Exit(0)
	}

	if commandLine.IsOptionSet("length") {
		length, err := strconv.Atoi(commandLine.OptionValue("length"))
		if err != nil {
			fmt.Println("Invalid length specified")
			os.Exit(1)
		} else {
			if length < 15 {
				length = 15
			} else if length > 30 {
				length = 30
			} else {
				Opts.Length = 15
			}
		}
		Opts.Length = length
	} else {
		Opts.Length = 15
	}

	if commandLine.IsOptionSet("length") {
		Opts.File = commandLine.OptionValue("length")

	}

	if commandLine.IsOptionSet("file") {
		Opts.File = commandLine.OptionValue("file")
	}

	if commandLine.IsOptionSet("user") {
		Opts.User = commandLine.OptionValue("user")
	}

	if commandLine.IsOptionSet("pass") {
		Opts.Pass = commandLine.OptionValue("pass")
	}

	if Opts.User != "" && Opts.Remove == true && Opts.File != "" {

	}

	if Opts.Generate == false {
		if Opts.File == "" {
			fmt.Println("No file specified")
			commandLine.PrintUsage(os.Stderr)
			os.Exit(1)
		}
	}

	return commandLine
}

func main() {
	if Opts.User != "" && Opts.Remove == true && Opts.File != "" {
		err := htpasswd.RemoveUser(Opts.File, Opts.User)
		if err != nil {
			fmt.Printf("User %s could not be removed from file %s %v\n", Opts.User, Opts.File, err)
			os.Exit(1)
		}
		fmt.Printf("User %s removed from file %s\n", Opts.User, Opts.File)
		os.Exit(0)
	}

	if Opts.User != "" && Opts.Add == true && Opts.File != "" {
		err := htpasswd.SetPassword(Opts.File, Opts.User, Opts.Pass, htpasswd.HashBCrypt)
		if err != nil {
			fmt.Printf("User %s could not be added to file %s %v\n", Opts.User, Opts.File, err)
			os.Exit(1)
		}
		fmt.Printf("User %s added to file %s\n", Opts.User, Opts.File)
		os.Exit(0)
	}

	if Opts.Add == false && Opts.Remove == false {
		if Opts.Generate == true {
			newPass, err := password.Generate(Opts.Length, 5, 2, false, false)
			if err != nil {
				fmt.Printf("Problem generating password %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s", newPass)
			os.Exit(0)
		}
	}

	fmt.Printf("Add or remove or generate new password not specified")
	commandLine.PrintUsage(os.Stdout)
}
