package upgrade

import (
	"fmt"
	"github.com/op/go-logging"
	"golang.org/x/mod/semver"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var log = logging.MustGetLogger("util.upgrade")

type Updater interface {
	IsNewerVersionAvailable() (bool, string)
	Upgrade() error
}

type DefaultUpdater struct {
	CurrentVersion string
	BaseURL        *url.URL
	ReleasePath    *url.URL
	VersionPath    *url.URL
}

func SprintReleaseKey(os string, arch string) string {
	if os != "windows" {
		return fmt.Sprintf("polarctl-%s-%s", os, arch)
	} else {
		return fmt.Sprintf("polarctl-%s-%s.exe", os, arch)
	}
}

func NewUpdater(baseUrl string, os string, arch string, versionPath string, currentVersion string) (Updater, error) {
	releasePath := SprintReleaseKey(os, arch)

	if !semver.IsValid(currentVersion) {
		return nil, fmt.Errorf("invalid currentVersion: %s", currentVersion)
	}

	parsedBaseURL, err := url.Parse(baseUrl)
	if err != nil || !parsedBaseURL.IsAbs() {
		return nil, fmt.Errorf("invalid baseUrl: %s", baseUrl)
	}

	parsedReleasePath, err := url.Parse(releasePath)
	if err != nil {
		return nil, fmt.Errorf("invalid releasePath: %s", releasePath)
	}

	parsedVersionPath, err := url.Parse(versionPath)
	if err != nil {
		return nil, fmt.Errorf("invalid versionPath: %s", versionPath)
	}

	return &DefaultUpdater{
		CurrentVersion: currentVersion,
		BaseURL:        parsedBaseURL,
		ReleasePath:    parsedReleasePath,
		VersionPath:    parsedVersionPath,
	}, nil
}

func (updater *DefaultUpdater) IsNewerVersionAvailable() (bool, string) {
	resp, err := http.Get(updater.BaseURL.ResolveReference(updater.VersionPath).String())
	if err != nil {
		log.Errorf("Error searching for new version of polarctl", err)
		return false, ""
	}
	defer resp.Body.Close()

	versionBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading new version of polarctl", err)
		return false, ""
	}

	remoteVersion := strings.TrimSpace(string(versionBytes))
	if err != nil || !semver.IsValid(remoteVersion) {
		log.Errorf("New version of polarctl was invalid, '%s'", remoteVersion)
		return false, ""
	}
	return semver.Compare(updater.CurrentVersion, remoteVersion) < 0, remoteVersion
}

func (updater *DefaultUpdater) Upgrade() error {
	available, remoteVersion := updater.IsNewerVersionAvailable()
	if available {
		s := updater.BaseURL.ResolveReference(updater.ReleasePath).String()
		log.Infof("Downloading new version from %s", s)
		resp, err := http.Get(s)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error fetching new version, status was %d", resp.StatusCode)
		}

		dest, err := os.Executable()
		if err != nil {
			return err
		}

		// Move the old version to a backup path that we can recover from
		// in case the upgrade fails
		destBackup := dest + ".bak"
		if _, err := os.Stat(dest); err == nil {
			if err := os.Rename(dest, destBackup); os.IsPermission(err) {
				return fmt.Errorf("permission denied, see https://git.smith.care/smith/uc-phep/polar/polar-control-2#permission-denied")
			} else if err != nil {
				return err
			}
		}

		// Use the same flags that ioutil.WriteFile uses
		f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			os.Rename(destBackup, dest)
			return err
		}
		defer f.Close()

		log.Infof("Downloading new version to %s", dest)
		if _, err := io.Copy(f, resp.Body); err != nil {
			os.Rename(destBackup, dest)
			return err
		}
		// The file must be closed already so we can execute it in the next step
		f.Close()

		// Removing backup
		os.Remove(destBackup)

		log.Infof("Updated to version %s", remoteVersion)
	} else {
		log.Infof("No new version available")
	}
	return nil
}
