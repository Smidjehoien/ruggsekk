package environment

import (
	"github.com/ruckstack/ruckstack/common/ui"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	RuckstackHome string
	resourceRoot  string
	cacheRoot     string
	tempDir       string

	OutDir string

	isRunningTests = false
)

func init() {
	executable, err := os.Executable()
	if err == nil {
		ui.VPrintln("Cannot determine if we are running tests: ", err)
	}
	isRunningTests = strings.HasPrefix(executable, "/tmp/")

	//find RuckstackHome
	if isRunningTests {
		//work from the working directory
		RuckstackHome, err = os.Getwd()

		if err != nil {
			ui.Fatalf("Cannot determine working directory: %s", err)
		}
	} else {
		RuckstackHome = filepath.Dir(executable)
	}

	//search back until we fine the root containing the LICENSE file
	for RuckstackHome != "/" {
		if _, err := os.Stat(filepath.Join(RuckstackHome, "LICENSE")); os.IsNotExist(err) {
			RuckstackHome = filepath.Dir(RuckstackHome)
			continue
		}
		break
	}

	if RuckstackHome == "/" {
		ui.Fatal("Cannot determine Ruckstack home")
	}
	ui.VPrintf("Ruckstack home: %s\n", RuckstackHome)

	//find resourceRoot
	if isRunningTests {
		resourceRoot = RuckstackHome + "/builder/cli/install_root/resources"
	} else {
		resourceRoot = RuckstackHome + "/resources"
	}
	ui.VPrintf("Ruckstack resource root: %s", resourceRoot)

	//find cacheRoot
	cacheRoot = os.Getenv("RUCKSTACK_CACHE_DIR")
	if cacheRoot == "" {
		cacheRoot = RuckstackHome + "/cache"
	}
	ui.VPrintf("Ruckstack cache root: %s", cacheRoot)

	//find tempDir
	if isRunningTests {
		tempDir = RuckstackHome + "/tmp"
	} else {
		tempDir, err = ioutil.TempDir("", "ruckstack")
		if err != nil {
			ui.Fatalf("Cannot determine temp directory: %s", err)
		}
	}

	ui.VPrintf("Ruckstack temp dir: %s", tempDir)

	if !isRunningTests {
		//when running all tests, the init method is called too often and clearing the temp dir interferes with other tests
		err = os.RemoveAll(tempDir)
		if err != nil {
			ui.VPrintf("Cannot clear temp dir: %s", err)
		}
		err = os.MkdirAll(tempDir, 0755)
		if err != nil {
			ui.VPrintf("Cannot create temp dir: %s", err)
		}
	}
}

/**
Returns true if ruckstack is running via the launcher
*/
func IsRunningLauncher() bool {
	return os.Getenv("RUCKSTACK_DOCKERIZED") == "true"
}

/**
Returns the full path to the given subpath of "resources" in RuckstackHome.
Returns an error if the file does not exist
*/
func ResourcePath(path string) (string, error) {
	resourcePath := filepath.Join(resourceRoot, path)

	if _, err := os.Stat(resourcePath); err != nil {
		return "", err
	}

	return resourcePath, nil
}

/**
Returns the given path as a sub-path of the Ruckstack "temporary" directory.
*/
func TempPath(pathInTmp string) string {
	return filepath.Join(tempDir, pathInTmp)
}

/**
Returns the given path as a sub-path of the Ruckstack "cache" dir.
The cache directory is preserved from one run to the next
*/
func CachePath(pathInCache string) string {
	return filepath.Join(cacheRoot, pathInCache)
}

/**
Returns the given path as a sub-path of the Ruckstack "out" dir.
*/
func OutPath(path string) string {
	if OutDir == "" {
		ui.Fatal("out directory not specified")
	}
	return filepath.Join(OutDir, path)
}
