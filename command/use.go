package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/shyiko/jabba/cfg"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

func Use(selector string) ([]string, error) {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return nil, err
	}
	return usePath(filepath.Join(cfg.Dir(), "jdk", ver))
}

func usePath(path string) ([]string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")

	goos := runtime.GOOS

	if goos == "darwin" {
		path = filepath.Join(path, "Contents", "Home")
	}

	systemJavaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		systemJavaHome, _ = os.LookupEnv("JAVA_HOME")
	}

	if goos == "windows" {
		log.Printf("JAVA_HOME has been set to %s", path)
		return []string{
			"[Environment]::SetEnvironmentVariable('JAVA_HOME', '" + path + "', 'User')",
			"$envPath = [Environment]::GetEnvironmentVariable('Path','User').replace('%JAVA_HOME%\\bin;', '')",
			"$newPath = '%JAVA_HOME%\\bin;' + $envPath",
			"[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')",
		}, nil
	}

	return []string{
		"export PATH=\"" + filepath.Join(path, "bin") + string(os.PathListSeparator) + pth + "\"",
		"export JAVA_HOME=\"" + path + "\"",
		"export JAVA_HOME_BEFORE_JABBA=\"" + systemJavaHome + "\"",
	}, nil
}
