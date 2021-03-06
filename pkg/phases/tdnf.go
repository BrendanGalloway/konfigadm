package phases

import (
	"fmt"
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

type TdnfPackageManager struct{}

func (p TdnfPackageManager) Install(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("tdnf install -y %s", strings.Join(pkg, " ")))
}

func (p TdnfPackageManager) Update() types.Commands {
	return types.Commands{}
}

func (p TdnfPackageManager) Uninstall(pkg ...string) types.Commands {
	return types.NewCommand(fmt.Sprintf("tdnf remove -y %s", strings.Join(pkg, " ")))
}

func (p TdnfPackageManager) Mark(pkg ...string) types.Commands {
	return types.Commands{}
}

func (p TdnfPackageManager) CleanupCaches() types.Commands {
	return types.Commands{}
}

func (p TdnfPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	_, ok := utils.SafeExec("tdnf info " + pkg)
	if !ok {
		log.Debugf("No matching package available in db for: %s", pkg)
		return ""
	}

	stdout, ok := utils.SafeExec("tdnf info installed " + pkg)
	if !ok {
		log.Debugf("%s package available in db but not installed", pkg)
		return ""
	}

	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "Version") {
			return strings.Split(line, ": ")[1]
		}
	}
	log.Debugf("Unable to find version info in " + stdout)
	return "Unknown Version"
}

func (p TdnfPackageManager) AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) types.Commands {
	repo := fmt.Sprintf(
		`[%s]
name=%s
baseurl=%s
enabled=1
`, name, name, url)

	if gpgKey != "" {
		repo += fmt.Sprintf(`gpgcheck=1
repo_gpgcheck=1
gpgkey=%s
`, gpgKey)
	} else {
		repo += `
gpgcheck=0
`
	}

	for k, v := range extraArgs {
		repo += fmt.Sprintf("%s = %s\n", k, v)
	}

	return types.NewCommand(fmt.Sprintf(`cat <<EOF >/etc/yum.repos.d/%s.repo
%s
EOF`, name, repo))
}
