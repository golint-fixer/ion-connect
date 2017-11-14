// config.go
//
// Copyright [2016] [Selection Pressure]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ionconnect

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

//Command represents a principle function of the command line tool
type Command struct {
	Name        string
	Usage       string
	Method      string
	URL         string
	Flags       []Flag
	Args        Args
	Subcommands []Command
}

//GetArgsUsage returns documentation for all command arguments
func (command Command) GetArgsUsage() string {
	var buffer bytes.Buffer
	for _, arg := range command.Args {
		if len(arg.Usage) > 0 {
			if !arg.Required {
				buffer.WriteString("[")
			}
			buffer.WriteString(arg.Usage)
			if !arg.Required {
				buffer.WriteString("]")
			}
			buffer.WriteString(" ")

		}
	}

	return buffer.String()
}

//GetArgsForFlags returns arguments for a given command flag
func (command Command) GetArgsForFlags(flagName string) Args {
	for _, flag := range command.Flags {
		if flag.Name == flagName {
			return flag.Args
		}
	}

	return []Arg{}
}

//GetFlagsWithArgs returns flags that take arguments
func (command Command) GetFlagsWithArgs() []Flag {
	flags := []Flag{}
	for _, flag := range command.Flags {
		if len(flag.Args) > 0 {
			flags = append(flags, flag)
		}
	}

	return flags
}

//GetArgsUsageWithFlags writes documentation for a given flag
func (command Command) GetArgsUsageWithFlags(flagName string) string {
	var buffer bytes.Buffer

	for _, flag := range command.Flags {
		if flag.Name == flagName {
			for _, arg := range flag.Args {
				if len(arg.Usage) > 0 {
					if !arg.Required {
						buffer.WriteString("[")
					}
					buffer.WriteString(arg.Usage)
					if !arg.Required {
						buffer.WriteString("]")
					}
					buffer.WriteString(" ")
				}
			}
		}
	}

	return buffer.String()
}

//GetRequiredArgsCount returns the number of args required
func (args Args) GetRequiredArgsCount() int {
	var count int
	for _, arg := range args {
		Debugf("Arg (%s) is required (%b)", arg.Name, arg.Required)
		if len(arg.Usage) > 0 && arg.Required {
			count++
		}
	}

	return count
}

//GetDefaultRequiredArgsCount returns the required arguments for a command
func (command Command) GetDefaultRequiredArgsCount() int {
	return command.Args.GetRequiredArgsCount()
}

//Arg represents a field of a flag
type Arg struct {
	Name     string
	Value    string
	Usage    string
	Required bool
	Type     string
}

//Args is a collection of command arguments
type Args []Arg

//Flag represents a parameter of a command
type Flag struct {
	Name     string
	Value    string
	Usage    string
	Type     string
	Required bool
	Args     Args
}

//Config indicates the API interface and supported commands
type Config struct {
	Version  string
	Endpoint string
	Token    string
	Commands []Command
}

//GetConfig returns API config available on the filesystem
func GetConfig() Config {
	configBox, err := rice.FindBox("../config")
	if err != nil {
		log.Fatal(err)
	}

	// get file contents as string
	configString, err := configBox.String("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	err = yaml.Unmarshal([]byte(configString), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if !test {
		config.Commands = config.Commands[:len(config.Commands)-1]
	}
	return config
}

//FindCommandConfig returns a Command based on name or an error
func (config Config) FindCommandConfig(commandName string) (Command, error) {
	for _, command := range config.Commands {
		if command.Name == commandName {
			return command, nil
		}
	}

	return Command{}, errors.New("Command not found")
}

//ProcessURLFromConfig returns the command URL
func (config Config) ProcessURLFromConfig(commandName string, subcommandName string, params interface{}) (string, error) {
	subCommandConfig, err := config.FindSubCommandConfig(commandName, subcommandName)
	if err != nil {
		return "", err
	}
	return subCommandConfig.URL, nil
}

//FindSubCommandConfig returns a secondary Command based on a primary and secondary key or an error
func (config Config) FindSubCommandConfig(commandName string, subcommandName string) (Command, error) {
	command, err := config.FindCommandConfig(commandName)
	if err != nil {
		return Command{}, err
	}

	for _, subcommand := range command.Subcommands {
		if subcommand.Name == subcommandName {
			return subcommand, nil
		}
	}

	return Command{}, errors.New("Subcommand not found")
}

//LoadEndpoint returns a string from configuration
func (config Config) LoadEndpoint() string {
	endpoint := os.Getenv(endpointEnvironmentVariable)
	if endpoint == "" {
		Debugf("Endpoint env var not found returning from config file (%s)", config.Endpoint)

		return config.Endpoint
	}

	Debugf("Credential env var found (%s)", endpoint)
	return endpoint
}

//HandleConfigure loads API credentials
func HandleConfigure(context *cli.Context) {
	currentSecretKey := LoadCredential()
	truncatedSecretKey := currentSecretKey
	if len(currentSecretKey) > 4 {
		truncatedSecretKey = currentSecretKey[len(currentSecretKey)-4 : len(currentSecretKey)]
	}

	fmt.Printf("Ion Channel Api Key [%s]: ", truncatedSecretKey)
	secretKey, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

	Debugf("All you keys are belong to us! (%s)", secretKey)

	if len(secretKey) != 0 {
		saveCredentials(string(secretKey))
	}
}

//LoadCredential loads API credentials
func LoadCredential() string {
	credential := os.Getenv(credentialsEnvironmentVariable)
	if credential == "" {
		Debugln("Credential env var not found looking in file")
		exists, _ := PathExists(ionHome)
		if exists {
			bytes, _ := ReadBytesFromFile(credentialsFile)
			credentials := make(map[string]string)
			yaml.Unmarshal([]byte(bytes), &credentials)
			return credentials[credentialsKeyField]
		}

		MkdirAll(ionHome, 0775)
		return ""
	}

	Debugln("Credential env var found")
	return credential
}

func saveCredentials(secretKey string) {
	credentials := make(map[string]string)
	credentials[credentialsKeyField] = secretKey
	yamlCredentials, _ := yaml.Marshal(&credentials)
	WriteLinesToFile(credentialsFile, []string{string(yamlCredentials)}, 0600)
}
