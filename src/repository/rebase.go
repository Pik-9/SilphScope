package repository

import "os/exec"

func RebaseUnghostedBranchOnto(unghostedBranch string, baseBranch string, repositoryPath string) (string, error) {
	gitBin, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}

	abortCmd := exec.Command(gitBin, "-C", repositoryPath, "rebase", "--abort")
	rebaseCmd := exec.Command(gitBin, "-C", repositoryPath, "rebase", unghostedBranch, baseBranch)
	rebaseOut, rebaseErr := rebaseCmd.CombinedOutput()
	if rebaseErr != nil {
		abortCmd.Run()
		return string(rebaseOut), rebaseErr
	}

	return string(rebaseOut), nil
}
