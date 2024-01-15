// main.go

package main

import (
	"fmt"
	"github.com/git719/maz"
	"github.com/git719/utl"
	"os"
	"path/filepath"
)

const (
	prgname = "pwrep"
	prgver  = "1.0.0"
)

func printUsage() {
	fmt.Printf(prgname + " Azure IDs password expiry report utility v" + prgver + "\n" +
		"    -ap  [DAYS]                     Password expiry report for all apps in tenant; optional within DAYS\n" +
		"    -sp  [DAYS]                     Password expiry report for all SPs in tenant; optional within DAYS\n" +
		"    -csv [DAYS]                     Password expiry report for all apps and SPs in CSV format; optional within DAYS\n" +
		"\n" +
		"    -id                             Display configured login values\n" +
		"    -id TenantId Username           Set up user for interactive login\n" +
		"    -id TenantId ClientId Secret    Set up ID for automated login\n" +
		"    -tx                             Delete current configured login values and token\n" +
		"    -v                              Print this usage page\n")
	os.Exit(0)
}

func setupVariables(z *maz.Bundle) maz.Bundle {
	// Set up variable object struct
	*z = maz.Bundle{
		ConfDir:      "",
		CredsFile:    "credentials.yaml",
		TokenFile:    "accessTokens.json",
		TenantId:     "",
		ClientId:     "",
		ClientSecret: "",
		Interactive:  false,
		Username:     "",
		AuthorityUrl: "",
		MgToken:      "",
		MgHeaders:    map[string]string{},
		AzToken:      "",
		AzHeaders:    map[string]string{},
	}
	// Set up configuration directory
	z.ConfDir = filepath.Join(os.Getenv("HOME"), ".maz") // IMPORTANT: Setting config dir = "~/.maz"

	if utl.FileNotExist(z.ConfDir) {
		if err := os.Mkdir(z.ConfDir, 0700); err != nil {
			panic(err.Error())
		}
	}
	return *z
}

func main() {
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 4 {
		printUsage() // Don't accept less than 1 or more than 4 arguments
	}

	var z maz.Bundle
	z = setupVariables(&z)

	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		// This first set of 1-arg requests do not require API tokens to be set up
		switch arg1 {
		case "-v":
			printUsage()
		case "-id":
			maz.DumpLoginValues(z)
		case "-tx":
			maz.RemoveCacheFile("t", z)
			maz.RemoveCacheFile("id", z)
		}
		z = maz.SetupApiTokens(&z) // The remaining 1-arg requests do required API tokens to be set up
		switch arg1 {
		case "-ap":
			maz.PrintExpiringSecretsReport("ap", "-1", z)
		case "-sp":
			maz.PrintExpiringSecretsReport("sp", "-1", z)
		case "-csv":
			maz.PrintExpiringSecretsReport("csv", "-1", z)
		default:
			printUsage()
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		z = maz.SetupApiTokens(&z) // The remaining 1-arg requests do required API tokens to be set up
		switch arg1 {
		case "-ap":
			maz.PrintExpiringSecretsReport("ap", arg2, z)
		case "-sp":
			maz.PrintExpiringSecretsReport("sp", arg2, z)
		case "-csv":
			maz.PrintExpiringSecretsReport("csv", arg2, z)
		default:
			printUsage()
		}
	case 3: // Process 3-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		switch arg1 {
		case "-id":
			z.TenantId = arg2
			z.Username = arg3
			maz.SetupInterativeLogin(z)
		default:
			printUsage()
		}
	case 4: // Process 4-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		arg4 := os.Args[4]
		switch arg1 {
		case "-id":
			z.TenantId = arg2
			z.ClientId = arg3
			z.ClientSecret = arg4
			maz.SetupAutomatedLogin(z)
		default:
			printUsage()
		}
	}
	os.Exit(0)
}
