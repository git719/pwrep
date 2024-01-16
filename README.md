# pwrep
`pwrep` is a [CLI](https://en.wikipedia.org/wiki/Command-line_interface) utility that reports the **expiry** status of every secret on every application and service principal in the current Azure tenant. It is written in [Go](https://go.dev/) and compiled into a binary executable, to make it a very quick little _Swiss Army knife_ tool for reporting on Azure ID secrets. By default it reports in regular text, but has an option to produce [Comma-separated_values (CSV)](https://en.wikipedia.org/wiki/Comma-separated_values) formatted output that can be piped into a file and then view in a spreadsheet.

## Quick Example
For a quick example, to list the secret expiry dates of all applications in the tenant: 

```
$ pwrep -ap
API call 1: 0 objects
PASSWORD EXPIRY REPORT: Apps
DISPLAY_NAME                             APP_ID                                 SECRET_ID                              SECRET_NAME      EXPIRY_DATE_TIME
localtest-sp                             f706dd63-a9d5-4ba2-b57a-476225d5f23b   42f0558c-bf40-4bc3-bc60-03003932d07d   MyName2          2026-01-01 00:00
sp-validator                             20726181-d443-426e-a07d-6e13f592cc57   a9a775c7-aaa1-47cc-ac0f-edad9749a4d9   Initial          2024-11-24 16:27
tf-az-sp00                               43e6a637-587f-49bf-b4f1-ae5473d2b9b4   bb5b4b41-4cea-4949-bbf6-086cf8e1605b   today            2023-01-23 04:59
tf-az-sp00                               43e6a637-587f-49bf-b4f1-ae5473d2b9b4   5e77300b-0d36-4f81-889d-30ec4818423c   Joe's Test       2023-04-23 00:30
tf-az-sp00                               43e6a637-587f-49bf-b4f1-ae5473d2b9b4   673cb88c-4845-4798-b980-5dc3480b7feb   Initial          2024-10-27 22:43
sp_site_extension                        5c6daa9d-27c6-4b5b-9f76-c1ef09af406e   a7775243-61e7-452b-9f74-3a236d2f2625                    2025-09-25 02:06
sp_site_reader                           ce882285-1954-4b07-a38e-615bd0e931f1   07db8e95-0375-4e1b-8a99-f5aec348fc8e   2nd_secret       2024-05-19 00:00
sp_site_reader                           ce882285-1954-4b07-a38e-615bd0e931f1   bfc7f9fa-57ba-4785-b08f-1bec4f9aef98   new-secret       2024-10-01 15:45
```

- Already expired secrets are highligted in <span style="color:red">red color text</span> .

## Introduction
The utility was primarily developed as a __proof-of-concept__ to:

- Learn to develop Azure utilities in the Go language
- Develop a small framework library for acquiring Azure [JWT](https://jwt.io/) token using the [MSAL library for Go](https://github.com/AzureAD/microsoft-authentication-library-for-go)
- Get a token for either an Azure user or a Service Principal (SP)
- Get the token to access the tenant's **Security** Services API via endpoint <https://graph.microsoft.com> ([MS Graph](https://learn.microsoft.com/en-us/graph/overview))
- Do quick and dirty reporting on password expiry for Azure IDs
- Develop small CLI utilities that call and use other Go library packages

## Getting Started
To compile `pwrep`, first make sure you have installed and set up the Go language on your system. You can do that by following [these instructions here](https://que.tips/golang/#install-go-on-macos) or by following other similar recommendations found across the web.

- Also ensure that `$GOPATH/bin/` in your `$PATH`, since that's where the executable binary will be placed.
- Open a `bash` shell, clone this repo, then switch to the `pwrep` working directory
- Type `./build` to build the binary executable
- To build from a regular Windows Command Prompt, just run the corresponding line in the `build` file (`go build ...`)
- If no errors, you should now be able to type `pwrep` and see the usage screen for this utility.

This utility has been successfully tested on macOS, Ubuntu Linux, and Windows. In Windows it works from a regular CMD.EXE, or PowerShell prompts, as well as from a GitBASH prompt.

Below other sections in this README explain how to set up access and use the utility in your own Azure tenant. 

## Access Requirements
First and foremost you need to know the special **Tenant ID** for your tenant. This is a UUID that uniquely identifies your Microsoft Azure tenant.

Then, you need a User ID or a Service Principal with the appropriate access rights. Either one will need the necessary _Global Reader_ role access to read security objects.

When you run `pwrep` without any arguments you will see the **usage** screen listed at the bottom of this README. As you can probably surmise, the `-id` arguments will allow you to set up these 2 optional ways to connect to your tenant. You can either set up to use a User ID with an interactive browser popup login, also known as a [User Principal Name (UPN)](https://learn.microsoft.com/en-us/entra/identity/hybrid/connect/plan-connect-userprincipalname), or you can set up to use special Service Principal or SP with a password.

## User ID Logon
For example, if your Tenant ID was **c44154ad-6b37-4972-8067-0ef1068079b2**, and your User ID was __bob@contoso.com__, you would type:

```
$ pwrep -id c44154ad-6b37-4972-8067-0ef1068079b2 bob@contoso.com
Updated /Users/myuser/.maz/credentials.yaml file
```
`pwrep` responds that the special `credentials.yaml` file has been updated accordingly.

To view, dump all configured logon values type the following:

```
$ pwrep -id
config_dir: /Users/myuser/.maz  # Config and cache directory
config_env_variables:
  # 1. Credentials supplied via environment variables override values provided via credentials file
  # 2. MAZ_USERNAME+MAZ_INTERACTIVE login have priority over MAZ_CLIENT_ID+MAZ_CLIENT_SECRET login
  MAZ_TENANT_ID:
  MAZ_USERNAME:
  MAZ_INTERACTIVE:
  MAZ_CLIENT_ID:
  MAZ_CLIENT_SECRET:
config_creds_file:
  file_path: /Users/myuser/.maz/credentials.yaml
  tenant_id: c44154ad-6b37-4972-8067-0ef1068079b2
  username: bob@contoso.com
  interactive: true
```

Above tells you that the utility has been configured to use Bob's UPN for access via the special credentials file. Note that above is only a configuration setup, it actually hasn't logged Bob onto the tenant yet. To logon as Bob you have have to run any command, and the logon will happen automatically, in this case it will be an interactively browser popup logon.

Note also, that instead of setting up Bob's login via the `-id` argument, you could have setup the special 3 operating system environment variables to achieve the same. Had you done that, running `pwrep -id` would have displayed below instead:

```
$ pwrep -id
config_dir: /Users/myuser/.maz  # Config and cache directory
config_env_variables:
  # 1. Credentials supplied via environment variables override values provided via credentials file
  # 2. MAZ_USERNAME+MAZ_INTERACTIVE login have priority over MAZ_CLIENT_ID+MAZ_CLIENT_SECRET login
  MAZ_TENANT_ID: c44154ad-6b37-4972-8067-0ef1068079b2
  MAZ_USERNAME: bob@contoso.com
  MAZ_INTERACTIVE: true
  MAZ_CLIENT_ID:
  MAZ_CLIENT_SECRET:
config_creds_file:
  file_path: /Users/myuser/.maz/credentials.yaml
  tenant_id: c44154ad-6b37-4972-8067-0ef1068079b2
  username: bob@contoso.com
  interactive: true
```

## SP Logon
To use an SP logon it means you first have to set up a dedicated App Registrations, grant it the same Reader resource and Global Reader security access roles mentioned above. Please reference other sources on the Web for how to do an Azure App Registration.

Once above is setup, you then follow the same logic as for User ID logon above, but using `-id` instead; or use the other environment variables (MAZ_CLIENT_ID andMAZ_CLIENT_SECRET ). 

The utility ensures that the permissions for configuration directory where the `credentials.yaml` file is placed is only accessible by the owning user. However, storing a secrets in a clear-text file is a very poor security practice and should __never__ be use other than for quick tests, etc. The environment variable options was developed pricisely for this SP logon pattern, where the `zls` utility could be setup to run from say a [Docker container](https://en.wikipedia.org/wiki/Docker_(software)) and the secret injected as an environment variable, and that would be a much more secure way to run the utility.

These login methods and the environment variables are described in more length in the [maz](https://github.com/git719/maz) package README.

## To-Do and Known Issues
The program is stable enough to be relied on as a reading/listing utility, but there are a number of little niggly things that could be improved. Will put a list together at some point.

At any rate, no matter how stable any code is, it is always worth remembering computer scientist [Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)'s famous quote:
> "Inside every large program is a small program struggling to get out."

## Coding Philosophy
The primary goal of this utility is to serve as a study aid for coding Azure utilities in Go, as well as to serve as a quick, _Swiss Army knife* utility for working with Azure.

If you look through the code you will note that it is very straightforward. Keeping the code clear, and simple to understand and maintain is another coding goal.

Note that the bulk of the code is actually in the [maz](https://github.com/git719/maz) library, and other packages.

## Usage
```
pwrep Azure IDs password expiry report utility v1.0.0
    -ap  [DAYS]                     Password expiry report for all apps in tenant; optional within DAYS
    -sp  [DAYS]                     Password expiry report for all SPs in tenant; optional within DAYS
    -csv [DAYS]                     Password expiry report for all apps and SPs in CSV format; optional within DAYS

    -id                             Display configured login values
    -id TenantId Username           Set up user for interactive login
    -id TenantId ClientId Secret    Set up ID for automated login
    -tx                             Delete current configured login values and token
    -v                              Print this usage page```
```

Instead of documenting individual examples of all of the above switches, it is best for you to play around with the utility to see the different listing functionality that it offers.

## Feedback
This utility along with the required libraries are obviously very useful to me, which is why I wrote them. However, I don't know if anyone else feels the same, which is why I have not yet thought about formalizing a proper feedback process.

The licensing is fairly open, so if you find it useful, feel free to clone and use on your own, with proper attributions. However, if you do see anything that could help improve any of this please let me know.
