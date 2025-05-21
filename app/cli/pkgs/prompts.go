package pkg

import (
	"github.com/manifoldco/promptui"
)

var KubeConfigPath = promptui.Prompt{
	Label:   "Enter path to kubeconfig file, leave empty if you want to use default kubeconfig path of your home directory",
	Default: "",
}

var Namespace = promptui.Prompt{
	Label:   "Enter namespace to install vapusai",
	Default: "vapusai",
}

var VapusDevMode = promptui.Prompt{
	Label:   "Is vapusai running in development mode? (Y/N)",
	Default: "Y",
}

var HelmCharurl = promptui.Prompt{
	Label:   "Enter helm chart url",
	Default: "",
}

var DevPostgresUsername = promptui.Prompt{
	Label:    "Enter postgres username for vapusai installation",
	Default:  "vapusai-postgres",
	Validate: validateNonEmpty,
}

var DevPostgresPassword = promptui.Prompt{
	Label:    "Enter postgres password for vapusai installation",
	Default:  "vapusai@postgres123",
	Validate: validateNonEmpty,
}

var DevPostgresDBName = promptui.Prompt{
	Label:    "Enter database name for vapusai installation",
	Default:  "vapusai-dev",
	Validate: validateNonEmpty,
}

var DevRedisPassword = promptui.Prompt{
	Label:    "Enter redis password for vapusai installation",
	Default:  "vapusai@redis123",
	Validate: validateNonEmpty,
}

var VapusInstallationName = promptui.Prompt{
	Label:    "Enter name for vapusai installation",
	Default:  "vapusai",
	Validate: validateNonEmpty,
}

var VapusPlatformVersion = promptui.Prompt{
	Label:    "Vapusdata platform version to install",
	Default:  "",
	Validate: validateNonEmpty,
}

var DevMode = promptui.Prompt{
	Label:    "Is vapusai running in development mode? (Y/N)",
	Default:  "Y",
	Validate: validateBool,
}

var InstallationName = promptui.Prompt{
	Label:    "Enter name for vapusai installation",
	Default:  "vapusai",
	Validate: validateNonEmpty,
}

var AIStudioCreator = promptui.Prompt{
	Label:    "Enter AI Studio creator email",
	Default:  "",
	Validate: validateNonEmpty,
}

var EncryptionAlgorithm = promptui.Prompt{
	Label:    "Select encryption algorithm, available options are ECDSA and RSA. Default is ECDSA",
	Default:  "ECDSA",
	Validate: validateEncryptionAlgorithm,
}

var TlsBitsize = promptui.Prompt{
	Label:    "Select tls bitsize, available options are 256, 384 and 521. Default is 521",
	Default:  "521",
	Validate: validateEncryptionAlgorithmBitSize,
}

var EncryptionAlgorithmBitSize = promptui.Prompt{
	Label:    "Select encryption algorithm bit size, available options are 2048, 3072 and 4096. Default is 2048",
	Default:  "2048",
	Validate: validateEncryptionAlgorithmBitSize,
}

func validateNonEmpty(input string) error {
	if input == "" {
		return ErrEmptyInput
	}
	return nil
}

func validateBool(input string) error {
	if input == "" {
		return ErrEmptyInput
	}
	if input != "Y" && input != "N" {
		return ErrInvalidBoolInput
	}
	return nil
}

func validateEncryptionAlgorithm(input string) error {
	if input == "" {
		return ErrEmptyInput
	}
	if input != "ECDSA" && input != "RSA" {
		return ErrInvalidEncryptionAlgorithm
	}
	return nil
}

func validateEncryptionAlgorithmBitSize(input string) error {
	if input == "" {
		return ErrEmptyInput
	}
	if input != "256" && input != "384" && input != "521" {
		return ErrInvalidEncryptionAlgorithm
	}

	return nil
}
