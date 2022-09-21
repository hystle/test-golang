package testGit

import (
	"fmt"
	gitv5 "github.com/go-git/go-git/v5"
	config "github.com/go-git/go-git/v5/config"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

func GitResetHead(r *gitv5.Repository) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	if err := w.Reset(&gitv5.ResetOptions{
		Mode:   gitv5.HardReset,
		Commit: head.Hash(),
	}); err != nil {
		return err
	}
	return nil
}

func GitPullRepo(r *gitv5.Repository) error {
	w, _ := r.Worktree()
	return w.Pull(&gitv5.PullOptions{RemoteName: "origin"})
}

func GitOpenRepo(path string) (*gitv5.Repository, error) {
	return gitv5.PlainOpen(path)
}

func GitGetRemoteUrl(r *gitv5.Repository, dir string) (string, error) {
	remote, err := r.Remote("origin")
	if err != nil {
		return "", err
	}
	// should only have 1 remote URL
	return remote.Config().URLs[0], nil
}

func GitCheckoutBranch(r *gitv5.Repository, branch string) error {
	w, err := r.Worktree()
	if err != nil {
		log.Errorln(err)
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := gitv5.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}
	if err := w.Checkout(&branchCoOpts); err != nil {
		log.Warnln("local checkout of branch '%s' failed, will attempt to fetch remote branch of same name.", branch)
		log.Warnln("like `git checkout <branch>` defaulting to `git checkout -b <branch> --track <remote>/<branch>`")

		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		if err = fetchOrigin(r, mirrorRemoteBranchRefSpec); err != nil {
			log.Errorln(err)
			return err
		}

		if err = w.Checkout(&branchCoOpts); err != nil {
			log.Errorln(err)
			return err
		}
	}
	return nil
}

func fetchOrigin(repo *gitv5.Repository, refSpecStr string) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		log.Errorln(err)
	}

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&gitv5.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil {
		if err == gitv5.NoErrAlreadyUpToDate {
			fmt.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

func GitCloneRepo(dir, gitUrl string) (*gitv5.Repository, error) {
	return gitv5.PlainClone(dir, false, &gitv5.CloneOptions{
		URL: gitUrl,
	})
}
