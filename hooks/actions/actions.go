package actions

import (
	"errors"
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/git"
	"github.com/captainhook-go/captainhook/hooks"
	"github.com/captainhook-go/captainhook/hooks/actions/branch"
	"github.com/captainhook-go/captainhook/hooks/actions/debug"
	"github.com/captainhook-go/captainhook/hooks/actions/file"
	"github.com/captainhook-go/captainhook/hooks/actions/message"
	"github.com/captainhook-go/captainhook/hooks/actions/notify"
	"github.com/captainhook-go/captainhook/io"
	"strings"
)

var (
	actionCreationConfig = map[string]map[string]func(appIO io.IO, conf *configuration.Configuration, repo git.Repo) hooks.Action{
		"branch": {
			"ensurenaming":                       branch.NewEnsureNaming,
			"preventpushoffixupandsquashcommits": branch.NewPreventPushOfFixupAndSquashCommits,
		},
		"debug": {
			"fail":    debug.NewFail,
			"success": debug.NewSuccess,
		},
		"file": {
			"blocksecrets":        file.NewBlockSecrets,
			"doesnotcontainregex": file.NewDoesNotContainRegex,
			"isnotempty":          file.NewIsNotEmpty,
			"maxsize":             file.NewMaxSize,
		},
		"message": {
			"injectissuekeyfrombranch": message.NewInjectIssueKeyFromBranch,
			"cacheonfail":              message.NewCacheOnFail,
			"mustfollowbeamsrules":     message.NewBeamsRules,
			"mustcontainsregex":        message.NewContainsRegex,
			"preparefromfile":          message.NewPrepareFromFile,
			"prepare":                  message.NewPrepare,
		},
		"notify": {
			"gitnotify": notify.NewGitNotify,
		},
	}
)

// ActionCreationFunc is returning a function to create the configured action
func ActionCreationFunc(path []string) (func(appIO io.IO, conf *configuration.Configuration, repo git.Repo) hooks.Action, error) {
	if len(path) != 2 {
		return nil, errors.New("invalid action functionality")
	}

	for index, value := range path {
		path[index] = strings.ToLower(value)
	}
	group, ok := actionCreationConfig[path[0]]
	if !ok {
		return nil, errors.New("invalid action functionality group")
	}

	intFunc, ok := group[path[1]]
	if !ok {
		return nil, errors.New("invalid action functionality")
	}

	return intFunc, nil

}
