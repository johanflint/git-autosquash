package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

const fixupPrefix = "fixup!"

func main() {
	repository, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		printf("No git repository found")
		return
	} else if err != nil {
		printf("Unable to open git repository: %s", err)
		return
	}

	var upstreamCommit *object.Commit
	autosquashCommitMessages := make([]string, 0)

	log, err := repository.Log(&git.LogOptions{})
	log.ForEach(func(commit *object.Commit) error {
		if commit.NumParents() > 1 {
			return storer.ErrStop
		}

		if strings.HasPrefix(commit.Message, fixupPrefix) {
			messageWithoutPrefix := strings.TrimLeft(strings.TrimPrefix(commit.Message, fixupPrefix), " ")
			autosquashCommitMessages = append(autosquashCommitMessages, messageWithoutPrefix)
		}

		if matchesAutosquashCommit(commit, autosquashCommitMessages) {
			upstreamCommit = commit
		}

		return nil
	})

	if len(autosquashCommitMessages) == 0 {
		printf("Nothing to autosquash: no commits found that start with %s", fixupPrefix)
		return
	}

	upstreamParentCommit, err := upstreamCommit.Parent(0)
	if err == object.ErrParentNotFound {
		printf("Unable to find the parent commit of %s", upstreamCommit.Hash)
		return
	}

	gitArguments := []string{"rebase", "--autosquash", "-i", fmt.Sprintf("%s", upstreamParentCommit.Hash)}
	printf("Executing git %s", strings.Join(gitArguments, " "))
	command := exec.Command("git", gitArguments...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Run()
}

func printf(format string, args ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}

func matchesAutosquashCommit(commit *object.Commit, autosquashCommitMessages []string) bool {
	for _, commitMessage := range autosquashCommitMessages {
		if strings.Contains(commit.Message, commitMessage) {
			return true
		}
	}

	return false
}
