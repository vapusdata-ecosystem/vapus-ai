package pkg

import (
	"github.com/manifoldco/promptui"
)

var KubeConfigPath = promptui.Prompt{
	Label:   "Enter path to kubeconfig file, leave empty if you want to use default kubeconfig path of your home directory",
	Default: "",
}

var Namespace = promptui.Prompt{
	Label:   "Enter namespace to install vapusdata",
	Default: "vapusdata",
}

var VapusDevMode = promptui.Prompt{
	Label:   "Is vapusdata running in development mode? (Y/N)",
	Default: "Y",
}

var HelmCharurl = promptui.Prompt{
	Label:   "Enter helm chart url",
	Default: "",
}

var VapusInstallationName = promptui.Prompt{
	Label:   "Enter name for vapusdata installation",
	Default: "vapusdata",
}

var VapusPlatformVersion = promptui.Prompt{
	Label:   "Vapusdata platform version to install",
	Default: "",
}

func validateNonEmpty(input string) error {
	if input == "" {
		return ErrEmptyInput
	}
	return nil
}
