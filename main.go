package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var verbose *bool

func main() {
	log.SetPrefix("⭕ Error ")
	log.SetFlags(log.Lshortfile)
	help := flag.Bool("h", false, "Shows commands usage")
	verbose = flag.Bool("verbose", false, "Show Verbose Output")
	workingDir := flag.String("dir", "/var/www/html", "Root Directory of Project")
	gitBranch := flag.String("branch", "master", "Git branch name")
	gitRemote := flag.String("git-remote", "origin", "Git remote name")
	gitReset := flag.Bool("git-reset", false, "Force Reset Git state")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *verbose {
		fmt.Println("Running with options:")
		fmt.Printf("  Working Directory: %s\n", *workingDir)
		fmt.Printf("  Git Branch: %s\n", *gitBranch)
		fmt.Printf("  Git Remote: %s\n", *gitRemote)
		fmt.Printf("  Git Reset: %v\n", *gitReset)
		fmt.Println("-----------------------------------")
	}

	folderName := strings.TrimRight(*workingDir, "/")[strings.LastIndex(strings.TrimRight(*workingDir, "/"), "/")+1:]

	fmt.Printf("~~~~~~ REDEPLOYING (%v) ~~~~~~\n", strings.ToUpper(regexp.MustCompile(`[-_]`).ReplaceAllString(folderName, " ")))
	if err := os.Chdir(*workingDir); err != nil {
		log.Fatalf("changing directory to %s: %v", *workingDir, err)
	}

	checkDir, _ := runCommandWithOutput("pwd")
	fmt.Printf("📁 Working dir : %v \n", checkDir)

	activeBranch, err := runCommandWithOutput("git", "branch", "--show-current")
	if err != nil {
		log.Fatalln("This Project is not a git repository!")
	}

	if *verbose {
		fmt.Printf("Current branch: %s\n", strings.TrimSpace(activeBranch))
		fmt.Printf("Target branch: %s\n", *gitBranch)
	}

	if *gitBranch != strings.TrimSpace(activeBranch) {
		fmt.Printf("⛓ Changing git branch to %v \n", *gitBranch)
		runCommand("git", "checkout", *gitBranch)
	}

	if *gitReset {
		fmt.Println("🧹 Forcing clean state with git reset...")
		runCommand("git", "reset", "--hard", "HEAD")
	}
	fetchOutput, err := runCommandWithOutput("git", "fetch")
	if err != nil {
		log.Fatalf("Failed to fetch: %v", err)
	}
	if *verbose && fetchOutput != "" {
		fmt.Printf("Fetch output:\n%s\n", fetchOutput)
	}

	fmt.Println("🖇 Checking Git Status")
	gitStatus, err := runCommandWithOutput("git", "rev-list", "--count", fmt.Sprintf("HEAD...%s/%s", *gitRemote, *gitBranch))
	if err != nil {
		log.Fatalf("Failed to check git status: %v ", err)
	}
	if *verbose {
		fmt.Printf("Commits behind remote: %s\n", strings.TrimSpace(gitStatus))
	}
	if strings.TrimSpace(gitStatus) == "0" {
		fmt.Println("✅ Already up to date. Stopping deployment.")
		os.Exit(0)

	}

	fmt.Println("🔄 Pulling latest changes...")
	gitOutput, err := runCommandWithOutput("git", "pull", *gitRemote, *gitBranch)
	if err != nil {
		log.Fatalf("Pulling git changes %v", err)
	}
	if strings.Contains(gitOutput, "Already up to date") {
		fmt.Println("✅ Already up to date. Stopping deployment.")
		os.Exit(0)
	}

	fmt.Println("✅ Deployment complete!")

}
func runCommand(name string, args ...string) {
	if *verbose {
		fmt.Printf("Executing: %s %s\n", name, strings.Join(args, " "))
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Command Failed: %v, %v\n", name, err)
	}
}

func runCommandWithOutput(name string, args ...string) (string, error) {
	if *verbose {
		fmt.Printf("Executing: %s %s\n", name, strings.Join(args, " "))
	}
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("%v: %v\n%s", name, err, stderr.String())
	}
	return stdout.String(), nil
}
