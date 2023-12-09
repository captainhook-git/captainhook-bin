package log

import (
	"github.com/captainhook-go/captainhook/git/types"
)

const (
	XmlFormat = "<commit>%n" +
		"<hash>%h</hash>%n" +
		"<names><![CDATA[%d]]></names>%n" +
		"<date>%ci</date>%n" +
		"<author><![CDATA[%an]]></author>%n" +
		"<subject><![CDATA[%s]]></subject>%n" +
		"<body><![CDATA[%n%b%n]]></body>%n" +
		"</commit>"
)

// AbbrevCommit is used to display only shortened commit hashes
func AbbrevCommit(g *types.Cmd) {
	g.AddOption("--abbrev-commit")
}

func AuthoredBy(name string) func(g *types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOption("--author=" + name)
	}
}

// Format is used to set the --format option
func Format(format string) func(g *types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOption("--format='" + format + "'")
	}
}

func ParseXML(out string) ([]*types.Commit, error) {
	var log []*types.Commit
	xmlLog, err := types.ParseLogXml(out)
	if err != nil {
		return log, err
	}
	for _, c := range xmlLog.Commits {
		log = append(log, types.CreateCommitFromXML(c))
	}
	return log, nil
}

// FromTo defines the range of the git log
func FromTo(from, to string) func(g *types.Cmd) {
	if to == "" {
		to = "HEAD"
	}
	return func(g *types.Cmd) {
		g.AddOption(from + ".." + to)
	}
}

// InTimeFrame adds --after and --before options
func InTimeFrame(after, before string) func(g *types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOption("--after=" + after)
		g.AddOption("--before=" + before)
	}
}

// NameOnly is used to output only file names
func NameOnly(g *types.Cmd) {
	g.AddOption("--name-only")
}

func NameStatus(g *types.Cmd) {
	g.AddOption("--name-status")
}

func NoCommitID(g *types.Cmd) {
	g.AddOption("--no-commit-id")
}

// NoMerges is used to exclude merges from the log
func NoMerges(g *types.Cmd) {
	g.AddOption("--no-merges")
}
