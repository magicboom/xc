package config

import (
	"conductor"
	"github.com/chzyer/readline"
	"github.com/viert/properties"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type XcConfig struct {
	Readline  *readline.Config
	Conductor *conductor.ConductorConfig

	User         string
	SSHThreads   int
	PingCount    int
	RemoteTmpdir string
	Mode         string
	RaiseType    string
}

const (
	defaultConfigContents = `[main]
user = 
mode = parallel
history_file = ~/.xc_history
cache_dir = ~/.xc_cache
raise = none

[executer]
ssh_threads = 50
ping_count = 5
remote_tmpdir = /tmp

[inventoree]
url = http://c.inventoree.ru
work_groups = `
)

var (
	defaultConductorConfig = &conductor.ConductorConfig{
		CacheTTL:      time.Hour * 24,
		WorkGroupList: []string{},
		RemoteUrl:     "http://c.inventoree.ru",
	}

	defaultReadlineConfig = &readline.Config{
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}
	defaultHistoryFile = "~/.xc_history"
	defaultCacheDir    = "~/.xc_cache"
	defaultCacheTTL    = 24
	defaultUser        = os.Getenv("USER")
	defaultThreads     = 50
	defaultTmpDir      = "/tmp"
	defaultPingCount   = 5
	defaultMode        = "parallel"
	defaultRaiseType   = "none"
)

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = "$HOME/" + path[2:]
	}
	return os.ExpandEnv(path)
}

func ReadConfig(filename string) (*XcConfig, error) {
	return readConfig(filename, false)
}

func readConfig(filename string, secondPass bool) (*XcConfig, error) {
	var props *properties.Properties
	var err error

	props, err = properties.Load(filename)
	if err != nil {
		// infinite loop break
		if secondPass {
			return nil, err
		}

		if os.IsNotExist(err) {
			err = ioutil.WriteFile(filename, []byte(defaultConfigContents), 0644)
			if err != nil {
				return nil, err
			}
		}
		return readConfig(filename, true)
	}

	xc := new(XcConfig)
	xc.Readline = defaultReadlineConfig
	xc.Conductor = defaultConductorConfig

	hf, err := props.GetString("main.history_file")
	if err != nil {
		hf = defaultHistoryFile
	}
	xc.Readline.HistoryFile = expandPath(hf)

	cttl, err := props.GetInt("main.cache_ttl")
	if err != nil {
		cttl = defaultCacheTTL
	}
	xc.Conductor.CacheTTL = time.Hour * time.Duration(cttl)

	cd, err := props.GetString("main.cache_dir")
	if err != nil {
		cd = defaultCacheDir
	}
	xc.Conductor.CacheDir = expandPath(cd)

	user, err := props.GetString("main.user")
	if err != nil || user == "" {
		user = defaultUser
	}
	xc.User = user

	threads, err := props.GetInt("executer.ssh_threads")
	if err != nil {
		threads = defaultThreads
	}
	xc.SSHThreads = threads

	tmpdir, err := props.GetString("executer.remote_tmpdir")
	if err != nil {
		tmpdir = defaultTmpDir
	}
	xc.RemoteTmpdir = tmpdir

	pc, err := props.GetInt("executer.ping_count")
	if err != nil {
		pc = defaultPingCount
	}
	xc.PingCount = pc

	invURL, err := props.GetString("inventoree.url")
	if err == nil {
		xc.Conductor.RemoteUrl = invURL
	}

	wglist, err := props.GetString("inventoree.work_groups")
	if err == nil {
		xc.Conductor.WorkGroupList = strings.Split(wglist, ",")
	}

	rt, err := props.GetString("main.raise")
	if err != nil {
		rt = defaultRaiseType
	}
	xc.RaiseType = rt

	mode, err := props.GetString("main.mode")
	if err != nil {
		mode = defaultMode
	}
	xc.Mode = mode

	return xc, nil
}
