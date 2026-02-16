
# awslogin - Text-based user interfaces (TUI) 

  

[](https://github.com/faustobranco/awslogin#awslogin)

  

[![Go Report Card](https://goreportcard.com/badge/github.com/faustobranco/awslogin)](https://goreportcard.com/report/github.com/faustobranco/awslogin)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![GitHub License](https://img.shields.io/github/license/faustobranco/awslogin)

  

![Views](https://hits.dwyl.com/faustobranco/awslogin.svg)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/faustobranco/awslogin)

![Homebrew](https://img.shields.io/badge/homebrew-tap-orange)


### Help:
<img width="814" height="187" alt="Screenshot 2026-02-16 at 11 14 26" src="https://github.com/user-attachments/assets/8b08f65f-0d57-4c88-ba96-511da4c521e9" />


### List:
<img width="818" height="157" alt="Screenshot 2026-02-16 at 11 17 22" src="https://github.com/user-attachments/assets/a6200111-046f-4c9f-ad6d-621291607db3" />


### Direct Connect
<img width="960" height="78" alt="Screenshot 2026-02-16 at 11 18 01" src="https://github.com/user-attachments/assets/377ede93-493c-4640-8d8d-d7c1e44f9f63" />


### Select / List Connection 
<img width="786" height="149" alt="Screenshot 2026-02-16 at 11 18 39" src="https://github.com/user-attachments/assets/e086b7a8-1a19-49b0-9b8d-6adfdeb07a40" />


 

**awslogin** is a lightweight CLI connection manager for AWS Accounts / Profiles. It allows you to quickly switch between multiple aws accounts defined in a aws default configuration file, supporting both interactive selection and direct connection flags.

  



  

## Features

  

[](https://github.com/faustobranco/awslogin#features)

  

-  **Interactive Menu**: Beautiful fuzzy-search selection powered by `promptui`.

-  **Direct Connect**: Use `--connect <name>` for instant access.

-  **Smart Defaults**: Automatic fallback for ports (5432) and CLI tools (pgcli).

-  **Shell Autocomplete**: Full support for Zsh and Bash completions (including server names).

-  **Environment Aware**: Prompts for missing credentials (User/Database) if not defined in config.

  

## Installation

  

[](https://github.com/faustobranco/awslogin#installation)

  

### Using Homebrew (Recommended)

  

[](https://github.com/faustobranco/awslogin#using-homebrew-recommended)

  

```
brew tap faustobranco/devops
brew install awslogin
```

  

### Manual Installation

  

1. Ensure you have [Go](https://golang.org/doc/install) installed.

2. Clone the repository and build:

  

Bash

 

```
go build -o awslogin .
sudo mv awslogin /usr/local/bin/
```

  

## Configuration

  

The tool expects a AWS Config file at `~/.aws/config`.

  

  

## Usage

  
  

|Command| Description |
|--|--|
|awslogin |Open interactive menu|
|awslogin --connect "Name" |Connect login to a specific account|
|awslogin --list |List all configured profiles|
|awslogin --config ./alt.config |Use a different config file|

  
  

## Shell Completion

  

If installed via Homebrew, completions are handled automatically. For manual setup, refer to the `completion/` directory.

  

## Requirements

  

- aws cli.
