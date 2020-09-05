package util

import (
	"fmt"
	"github.com/ruckstack/ruckstack/internal"
	"gopkg.in/yaml.v2"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	packageConfig *internal.PackageConfig
	systemConfig  *internal.SystemConfig
	localConfig   *internal.LocalConfig
	installDir    string
)

func InstallDir() string {
	if installDir == "" {

		installDir = os.Getenv("RUCKSTACK_HOME")

		if installDir == "" {
			ex, exErr := os.Executable()
			if exErr != nil {
				panic(exErr)
			}
			exPath := filepath.Dir(ex)
			installDir = filepath.Dir(exPath)
		}
	}
	return installDir
}

func SetInstallDir(newInstallDir string) {
	installDir = newInstallDir
}

func SetPackageConfig(passedPackageConfig *internal.PackageConfig) {
	packageConfig = passedPackageConfig
}

func GetPackageConfig() (*internal.PackageConfig, error) {
	if packageConfig != nil {
		return packageConfig, nil
	}

	file, err := os.OpenFile(InstallDir()+"/.package.config", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	packageConfig = new(internal.PackageConfig)
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(packageConfig)
	if err != nil {
		return nil, err
	}

	return packageConfig, nil
}

func GetSystemConfig() (*internal.SystemConfig, error) {
	if systemConfig != nil {
		return systemConfig, nil
	}

	file, err := os.Open(InstallDir() + "/config/system.config")
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	systemConfig = new(internal.SystemConfig)
	err = decoder.Decode(systemConfig)
	if err != nil {
		return nil, err
	}

	return systemConfig, nil
}

func SetSystemConfig(passedSystemConfig *internal.SystemConfig) {
	systemConfig = passedSystemConfig
}

func GetLocalConfig() (*internal.LocalConfig, error) {
	if localConfig != nil {
		return localConfig, nil
	}

	file, err := os.Open(InstallDir() + "/config/local.config")
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	localConfig = new(internal.LocalConfig)
	if err := decoder.Decode(localConfig); err != nil {
		return nil, err
	}

	return localConfig, nil
}

func SetLocalConfig(passedLocalConfig *internal.LocalConfig) {
	localConfig = passedLocalConfig
}

func ExpectNoError(err error) {
	if err != nil {
		fmt.Printf("Unexpected error %s", err)
		//panic(err)
		os.Exit(15)
	}
}

func ExecBash(bashCommand string) {
	command := exec.Command("bash", "-c", bashCommand)
	command.Dir = InstallDir()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		panic(err)
	}
}

func GetAbsoluteName(object meta.Object) string {
	return object.GetNamespace() + "/" + object.GetName()
}
