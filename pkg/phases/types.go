package phases

import (
	"fmt"
	"strings"

	cloudinit "github.com/moshloop/configadm/pkg/cloud-init"
	. "github.com/moshloop/configadm/pkg/systemd"
	yaml "gopkg.in/yaml.v3"
)

type Port struct {
	Port   int `json:"port,omitempty"  validate:"min=1,max=65536"`
	Target int `json:"target,omitempty"  validate:"min=1,max=65536"`
}

//Container represents a container to be run using systemd
type Container struct {
	//The name of the service (e.g systemd unit name or deployment name)
	Service string `json:"service,omitempty"`
	Image   string `json:"image"`
	//A map of environment variables to pass through
	Env map[string]string `json:"env,omitempty"`
	//A map of labels to add to the container
	Labels map[string]string `json:"labels,omitempty"`
	//Additional arguments to the docker run command e.g. -p 8080:8080
	DockerOpts string `json:"docker_opts,omitempty"`
	//Additional options to the docker client e.g. -H unix:///tmp/var/run/docker.sock
	DockerClientArgs string `json:"docker_client_args,omitempty"`
	//Additional arguments to the container
	Args     string   `json:"args,omitempty"`
	Ports    []Port   `json:"ports,omitempty"`
	Commands []string `json:"commands,omitempty"`
	//Map of files to mount into the container
	Files map[string]string `json:"files,omitempty"`
	//Map of templates to mount into the container
	Templates map[string]string `json:"templates,omitempty"`
	Volumes   []string          `json:"volumes,omitempty"`
	//CPU limit in cores (Defaults to 1 )
	CPU int `json:"cpu,omitempty" validate:"min=0,max=32"`
	//	Memory Limit in MB. (Defaults to 1024)
	Mem int `json:"mem,omitempty" validate:"min=0,max=1048576"`
	//default:	user-bridge	 only
	Network string `json:"network,omitempty"`
	// default: 1
	Replicas int `json:"replicas,omitempty"`
}

//ContainerRuntime installs a container runtime such as docker or CRI-O
type ContainerRuntime struct {
	Type    string `json:"type,omitempty"`
	Arg     string `json:"arg,omitempty"`
	Options string `json:"options,omitempty"`
	Version string `json:"version,omitempty"`
}

//Kubernetes installs the packages and configures the system for kubernetes, it does not actually bootstrap and configure kuberntes itself
//Use kubeadm in a `command` to actually configure and start kubernetes
type Kubernetes struct {
	Version      string `json:"version,omitempty"`
	DownloadPath string
	ImagePrefix  string
}

//Service is a systemd service to be installed and started
type Service struct {
	Name        string            `json:"name,omitempty"`
	ExecStart   string            `json:"exec_start,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Extra       SystemD           `json:"extra,omitempty"`
}

//User mirrors the CloudInit User struct.
type User struct {
	// The user's login name
	Name string `yaml:"name,omitempty"`
	//The user name's real name, i.e. "Bob B. Smith"
	Gecos string `yaml:"gecos,omitempty"`
	//	Optional. The SELinux user for the user's login, such as
	//          "staff_u". When this is omitted the system will select the default
	//           SELinux user.
	SeLinuxUser string `yaml:"selinux_user,omitempty"`
	ExpireDate  string `yaml:"expiredate,omitempty"`
	//	Defaults to none. Accepts a sudo rule string, a list of sudo rule
	//         strings or False to explicitly deny sudo usage. Examples:
	//
	//         Allow a user unrestricted sudo access.
	//             sudo:  ALL=(ALL) NOPASSWD:ALL
	//
	//         Adding multiple sudo rule strings.
	//             sudo:
	//               - ALL=(ALL) NOPASSWD:/bin/mysql
	//               - ALL=(ALL) ALL
	//
	//         Prevent sudo access for a user.
	//             sudo: False
	//
	//         Note: Please double check your syntax and make sure it is valid.
	//               cloud-init does not parse/check the syntax of the sudo
	//               directive.
	Sudo string `yaml:"sudo,omitempty"`
	//	The hash -- not the password itself -- of the password you want
	//           to use for this user. You can generate a safe hash via:
	//               mkpasswd --method=SHA-512 --rounds=4096
	//           (the above command would create from stdin an SHA-512 password hash
	//           with 4096 salt rounds)
	//
	//           Please note: while the use of a hashed password is better than
	//               plain text, the use of this feature is not ideal. Also,
	//               using a high number of salting rounds will help, but it should
	//               not be relied upon.
	//
	//               To highlight this risk, running John the Ripper against the
	//               example hash above, with a readily available wordlist, revealed
	//               the true password in 12 seconds on a i7-2620QM.
	//
	//               In other words, this feature is a potential security risk and is
	//               provided for your convenience only. If you do not fully trust the
	//               medium over which your cloud-config will be transmitted, then you
	//               should use SSH authentication only.
	//
	//               You have thus been warned.
	Passwd string `yaml:"passwd,omitempty"`
	// define the primary group. Defaults to a new group created named after the user.
	PrimaryGroup string `yaml:"primary_group,omitempty"`
	Groups       string `yaml:"groups,omitempty"`
	// Optional. Import SSH ids
	SSHImportID string `yaml:"ssh_import_id,omitempty"`
	//Defaults to true. Lock the password to disable password login
	LockPasswd bool `yaml:"lock_passwd,omitempty"`
	//When set to true, do not create home directory
	NoCreateHome bool `yaml:"no_create_home,omitempty"`
	//When set to true, do not create a group named after the user.
	NoUserGroup bool `yaml:"no_user_group,omitempty"`
	//When set to true, do not initialize lastlog and faillog database.
	NoLogInit bool `yaml:"no_log_init,omitempty"`
	//Add keys to user's authorized keys file
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	//Create the user as inactive
	Inactive bool `yaml:"inactive,omitempty"`
	// Create the user as a system user. This means no home directory.
	System bool `yaml:"system,omitempty"`
	//Create a Snappy (Ubuntu-Core) user via the snap create-user
	//             command available on Ubuntu systems.  If the user has an account
	//             on the Ubuntu SSO, specifying the email will allow snap to
	//             request a username and any public ssh keys and will import
	//             these into the system with username specified by SSO account./
	//             If 'username' is not set in SSO, then username will be the
	//             shortname before the email domain.
	Snapuser string `yaml:"snapuser,omitempty"`
	//	Set true to block ssh logins for cloud
	//      ssh public keys and emit a message redirecting logins to
	//      use <default_username> instead. This option only disables cloud
	//      provided public-keys. An error will be raised if ssh_authorized_keys
	//      or ssh_import_id is provided for the same user.
	SSHRedirectUser bool `yaml:"ssh_redirect_user,omitempty"`
}

type File struct {
	Content        string
	ContentFromURL string
	Permissions    string
	Owner          string
	Flags          []string
}

type Filesystem map[string]File

type Command struct {
	Cmd   string
	Flags []Flag
}

func (c Command) String() string {
	return c.Cmd
}

func (c *Command) UnmarshalYAML(node *yaml.Node) error {
	c.Cmd = node.Value
	comment := node.LineComment
	if !strings.Contains(comment, "#") {
		return nil
	}
	comment = comment[1:]
	for _, flag := range strings.Split(comment, " ") {
		if FLAG, ok := FLAG_MAP[flag]; ok {
			c.Flags = append(c.Flags, FLAG)
		} else {
			return fmt.Errorf("Unknown flag: %s", flag)
		}
	}
	return nil
}

type Package struct {
	Name      string
	Mark      bool
	Uninstall bool
	Flags     []Flag
}

func (p Package) String() string {
	return p.Name
}

func (p *Package) UnmarshalYAML(node *yaml.Node) error {
	p.Name = node.Value
	if strings.HasPrefix(node.Value, "!") {
		p.Name = node.Value[1:]
		p.Uninstall = true
	}
	if strings.HasPrefix(node.Value, "=") {
		p.Name = node.Value[1:]
		p.Mark = true
	}
	comment := node.LineComment
	if !strings.Contains(comment, "#") {
		return nil
	}
	comment = comment[1:]
	for _, flag := range strings.Split(comment, " ") {
		if FLAG, ok := FLAG_MAP[flag]; ok {
			p.Flags = append(p.Flags, FLAG)
		} else {
			return fmt.Errorf("Unknown flag: %s", flag)
		}
	}
	return nil
}

type PackageRepo struct {
	URL    string
	GPGKey string
	Flags  []string
}

//Config is the logical model after runtime tags have been applied
type Config struct {
	Commands         []Command            `yaml:"commands,omitempty"`
	PreCommands      []Command            `yaml:"pre_commands,omitempty"`
	PostCommands     []Command            `yaml:"post_commands,omitempty"`
	Filesystem       Filesystem           `yaml:"filesystem,omitempty"`
	Files            map[string]string    `yaml:"files,omitempty"`
	Templates        map[string]string    `yaml:"templates,omitempty"`
	Sysctls          map[string]string    `yaml:"sysctls,omitempty"`
	Packages         []Package            `yaml:"packages,omitempty"`
	PackageRepos     []PackageRepo        `yaml:"package_repos,omitempty"`
	Images           []string             `yaml:"images,omitempty"`
	Containers       []Container          `yaml:"containers,omitempty"`
	ContainerRuntime *ContainerRuntime    `yaml:"container_runtime,omitempty"`
	Kubernetes       *Kubernetes          `yaml:"kubernetes,omitempty"`
	Environment      map[string]string    `yaml:"environment,omitempty"`
	Timezone         string               `yaml:"timezone,omitempty"`
	Extra            *cloudinit.CloudInit `yaml:"extra,omitempty"`
	Services         map[string]Service   `yaml:"services,omitempty"`
	Users            []User               `yaml:"users,omitempty"`
	Context          *SystemContext       `yaml:"context,omitempty"`
}

type Applier interface {
	Apply(ctx SystemContext)
}
type SystemContext struct {
	Vars  map[string]interface{}
	Flags []Flag
	Name  string
}

type Transformer func(cfg *Config, ctx *SystemContext) (commands []Command, files Filesystem, err error)

type FlagProcessor func(cfg *Config, flags ...Flag)

type AllPhases interface {
	Phase
	ProcessFlagsPhase
}

type Phase interface {
	ApplyPhase(cfg *Config, ctx *SystemContext) (commands []Command, files Filesystem, err error)
}

type ProcessFlagsPhase interface {
	ProcessFlags(cfg *Config, flags ...Flag)
}