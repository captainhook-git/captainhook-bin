package git

import (
	"errors"
	"fmt"
	"github.com/captainhook-go/captainhook/git/config"
	"github.com/captainhook-go/captainhook/git/diff"
	"github.com/captainhook-go/captainhook/git/log"
	"github.com/captainhook-go/captainhook/git/revparse"
	"github.com/captainhook-go/captainhook/git/types"
	"github.com/captainhook-go/captainhook/io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	FilterDefault = "ACMR"
)

type Repository struct {
	root     string
	gitDir   string
	hooksDir string
}

func (r *Repository) Path() string {
	return r.root
}

// AbsPath returns the absolute path to the repository and the relative path as fallback
func (r *Repository) AbsPath() string {
	absPath, err := filepath.Abs(r.root)
	if err != nil {
		return r.root
	}
	return absPath
}

func (r *Repository) GitDir() string {
	return r.gitDir
}

func (r *Repository) HookExists(hook string) bool {
	return io.FileExists(r.HooksDir() + "/" + hook)
}

func (r *Repository) HooksDir() string {
	if r.hooksDir == "" {
		r.hooksDir = r.gitDir + "/hooks"
		var hooksPath = r.hooksPath()
		if hooksPath != "" {
			r.hooksDir = hooksPath
		}
	}
	return r.hooksDir
}

// hooksPath returns the core.hooksPath configured path
// If the path is relative it will be prefixed with the repository root
func (r *Repository) hooksPath() string {
	hooksPath := r.ConfigValue("core.hooksPath", "")
	if hooksPath == "" {
		return ""
	}
	if path.IsAbs(hooksPath) {
		return hooksPath
	}
	return r.root + "/" + hooksPath
}

func (r *Repository) CommitMessage(path string) (*types.CommitMessage, error) {
	if path == "" {
		return nil, errors.New("invalid file path")
	}
	commentChar := r.ConfigValue("core.commentchar", "#")
	return types.NewCommitMessageFromFile(path, commentChar)
}

func (r *Repository) PrepareCommitMessage(path string, msg *types.CommitMessage) error {
	return os.WriteFile(path, []byte(msg.Raw()), 0644)
}

func (r *Repository) ConfigValue(value string, defaultValue string) string {
	// git config --get VALUE
	out, err := Config(config.Get(value))
	if err != nil {
		return defaultValue
	}
	return out
}

func (r *Repository) IsMerging() bool {
	for _, file := range []string{"MERGE_MSG", "MERGE_HEAD", "MERGE_MODE"} {
		if _, err := os.Stat(r.gitDir + "/" + file); err == nil {
			return true
		}
	}
	return false
}

func (r *Repository) StagedFiles() ([]string, error) {
	// git diff-index --diff-algorithm=myers --no-ext-diff --cached --name-status HEAD
	out, err := DiffIndex(
		diff.Algorithm("myers"),
		diff.NoExtDiff,
		diff.Cached,
		log.NameOnly,
		diff.Filter(FilterDefault),
		diff.To("HEAD"),
	)
	if err != nil {
		return nil, err
	}
	return io.SplitLines(out), nil
}

func (r *Repository) ChangedFiles(from, to string) ([]string, error) {
	// git diff-tree --diff-algorithm=myers --no-ext-diff --no-commit-id --name-only -r --diff-filter=ADM FROM TO
	out, err := DiffTree(
		diff.NoExtDiff,
		log.NameOnly,
		diff.Recursive,
		diff.Filter(FilterDefault),
		diff.FromTo(from, to),
	)
	if err != nil {
		return []string{}, err
	}
	return io.SplitLines(out), nil
}

func (r *Repository) BranchName() string {
	// rev-parse --abbrev-ref HEAD
	out, err := RevParse(revparse.AbbrevRef, diff.To("HEAD"))
	if err != nil {
		return ""
	}
	return out
}

func (r *Repository) CommitsBetween(from string, to string) []*types.Commit {
	// git log --abbrev-commit --no-merges FROM TO
	out, err := Log(
		log.Format(log.XmlFormat),
		log.AbbrevCommit,
		log.NoMerges,
		log.FromTo(from, to),
	)
	if err != nil {
		return []*types.Commit{}
	}

	commits, _ := log.ParseXML("<log>" + out + "</log>")
	return commits
}

func NewRepository(gitDir string) (*Repository, error) {
	repoPath := path.Dir(gitDir)
	if !isPathARepository(gitDir) {
		err := fmt.Errorf("repository not found in: %s", gitDir)
		return nil, err
	}
	dotGitDir := gitDir
	if isGitSubmodule(dotGitDir) {
		// For submodules hooks are stored in the parents .git/modules directory
		dotGitContents, err := os.ReadFile(dotGitDir)
		if err != nil {
			return nil, err
		}
		r, _ := regexp.Compile("^gitdir:\\s*(.+)$")
		match := r.FindAllStringSubmatch(string(dotGitContents), -1)
		if len(match) > 0 {
			dotGitDir = fmt.Sprintf("%s/%s", repoPath, match[1])
		}
	}
	r := Repository{root: repoPath, gitDir: dotGitDir, hooksDir: ""}
	return &r, nil
}

// isPathARepository tries to answer the question if a directory is a git repository
// It checks at least the two most obvious files that should exist in a .git directory.
func isPathARepository(repoPath string) bool {
	pathSlice := []string{
		repoPath + "/config",
		repoPath + "/HEAD",
	}
	for _, p := range pathSlice {
		_, err := os.Stat(p)
		if err != nil {
			return false
		}
	}
	return true
}

func isGitSubmodule(dotGitDir string) bool {
	file, err := os.Open(dotGitDir)
	if err != nil {
		return false
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}
