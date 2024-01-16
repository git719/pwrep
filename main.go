// main.go

package main

import (
	"fmt"
	"github.com/git719/maz"
	"github.com/git719/utl"
	"os"
	"path/filepath"
	"time"
)

const (
	prgname = "pwrep"
	prgver  = "1.1.0"
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

func PrintSecretsExpirations(t, days string, z maz.Bundle) {
	var list []interface{} = nil
	var format = "text"
	switch t {
	case "ap":
		objList := maz.GetMatchingApps("", true, z) // true = force call to Azure to get latest
		for _, i := range objList {
			x := i.(map[string]interface{}) // Assert as JSON object
			x["oType"] = "App"              // Extend object with mazType as an ADDITIONAL field
			list = append(list, x)
		}
	case "sp":
		objList := maz.GetMatchingSps("", true, z)
		for _, i := range objList {
			x := i.(map[string]interface{})
			x["oType"] = "SP"
			list = append(list, x)
		}
	case "csv":
		objList := maz.GetMatchingApps("", true, z)
		for _, i := range objList {
			x := i.(map[string]interface{})
			x["oType"] = "App"
			list = append(list, x)
		}
		objList = maz.GetMatchingSps("", true, z)
		for _, i := range objList {
			x := i.(map[string]interface{})
			x["oType"] = "SP"
			list = append(list, x)
		}
		format = "csv"
	}

	if format == "csv" {
		fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", "OBJ", "DISPLAY_NAME", "APP_ID", "SECRET_ID", "EXPIRY_DATE_TIME")
	} else {
		fmt.Printf("%-6s %-40s %-38s %-38s %s\n", "OBJ", "DISPLAY_NAME", "APP_ID", "SECRET_ID", "EXPIRY_DATE_TIME")
	}

	for _, i := range list {
		x := i.(map[string]interface{}) // Assert as JSON object
		if x["passwordCredentials"] != nil {
			secretsList := x["passwordCredentials"].([]interface{})
			if len(secretsList) > 0 {
				oType := utl.Str(x["oType"])
				displayName := utl.Str(x["displayName"])
				appId := utl.Str(x["appId"])
				PrintExpiringSecrets(oType, displayName, appId, secretsList, days, format)
			}
		}
	}
}

func PrintExpiringSecrets(oType, displayName, appId string, secretsList []interface{}, days, format string) {
	// Print expiring secrets within 'days'; if days == -1 print regular expiry date
	if len(secretsList) < 1 {
		return
	}
	daysInt, err := utl.StringToInt64(days)
	if err != nil {
		utl.Die("Error converting 'days' to valid integer number.\n")
	}
	for _, i := range secretsList {
		pw := i.(map[string]interface{}) // Assert as JSON object
		secretId := utl.Str(pw["keyId"])
		expiry := utl.Str(pw["endDateTime"])

		// Convert expiry date to string and int64 epoch formats
		expiryStr, err := utl.ConvertDateFormat(expiry, time.RFC3339Nano, "2006-01-02 15:04")
		if err != nil {
			utl.Die(utl.Trace() + err.Error() + "\n")
		}
		expiryInt, err := utl.DateStringToEpocInt64(expiry, time.RFC3339Nano)
		if err != nil {
			utl.Die(utl.Trace() + err.Error() + "\n")
		}

		now := time.Now().Unix()
		daysDiff := (expiryInt - now) / 86400

		cExpiryStr := expiryStr
		if daysDiff <= 0 {
			cExpiryStr = utl.Red(expiryStr) // If it's expired print in red
		}

		if daysInt == -1 || daysDiff <= daysInt {
			if format == "csv" {
				fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", oType, displayName, appId, secretId, expiryStr)
			} else {
				fmt.Printf("%-6s %-40s %-38s %-38s %s\n", oType, displayName, appId, secretId, cExpiryStr)
			}
		}
	}
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
			PrintSecretsExpirations("ap", "-1", z)
		case "-sp":
			PrintSecretsExpirations("sp", "-1", z)
		case "-csv":
			PrintSecretsExpirations("csv", "-1", z)
		default:
			printUsage()
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		z = maz.SetupApiTokens(&z) // The remaining 1-arg requests do required API tokens to be set up
		switch arg1 {
		case "-ap":
			PrintSecretsExpirations("ap", arg2, z)
		case "-sp":
			PrintSecretsExpirations("sp", arg2, z)
		case "-csv":
			PrintSecretsExpirations("csv", arg2, z)
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
