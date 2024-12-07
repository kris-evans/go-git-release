package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

type config struct {
	filePath         string
	releaseNotesPath string
	releaseTag       string
	releaseProject   string
	dryRun           bool
}

func main() {
	config := &config{}

	app := &cli.App{
		Name:   "go-git-release",
		Usage:  "Simple opinionated release tooling for monorepos.",
		Action: run(config),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				EnvVars:     []string{"RELEASE_PATH"},
				Usage:       "Specifies the name of the path to create release notes.",
				Destination: &config.filePath,
				Value:       ".",
			},
			&cli.StringFlag{
				Name:        "notes",
				EnvVars:     []string{"RELEASE_NOTES_PATH"},
				Usage:       "Specifies the name of the file to export release notes.",
				Destination: &config.releaseNotesPath,
				Value:       "RELEASE_NOTES.md",
			},
			&cli.StringFlag{
				Name:        "project",
				EnvVars:     []string{"RELEASE_PROJECT"},
				Usage:       "Specifies the name of the project for release notes and release commits. (e.g. project-name)",
				Destination: &config.releaseProject,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "tag",
				EnvVars:     []string{"RELEASE_TAG"},
				Usage:       "Specifies the name of the tag for release notes and release commits. (e.g. v2024.12.01)",
				Destination: &config.releaseTag,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				EnvVars:     []string{"RELEASE_DRY_RUN"},
				Usage:       "Prints the commands that would be executed without running them.",
				Destination: &config.dryRun,
				Value:       false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log := slog.Default()
		log.Error(err.Error())
	}
}

func run(config *config) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log := slog.Default()

		projectReleaseTag := fmt.Sprintf("%s-%s", strings.ToLower(config.releaseProject), strings.ToLower(config.releaseTag))

		releaseBranch := filepath.Join("release", projectReleaseTag)

		log.Info("creating release branch", "branch", releaseBranch)

		if !config.dryRun {
			cmd := exec.Command(
				"git",
				"checkout",
				"-b",
				releaseBranch,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error creating release branch: %w", err)
			}

			log.Info(string(out))
		}

		releaseNotesFilePath := filepath.Join(config.filePath, config.releaseNotesPath)

		if !config.dryRun {
			_, err := os.Stat(releaseNotesFilePath)
			if os.IsNotExist(err) {
				log.Info("creating release notes file", "path", releaseNotesFilePath)
				file, err := os.Create(releaseNotesFilePath)
				if err != nil {
					return fmt.Errorf("error creating file: %w", err)
				}
				defer file.Close()
			}
		}

		wildcardPath := filepath.Join(config.filePath, "**", "*")

		log.Info("creating release notes using git cliff", "command",
			fmt.Sprintf("git cliff --include-path %s --strip all --tag %s --output %s", wildcardPath, projectReleaseTag, releaseNotesFilePath))

		if !config.dryRun {
			cmd := exec.Command(
				"git",
				"cliff",
				"--include-path", filepath.Join(config.filePath, "**", "*"),
				"--strip", "all",
				"--tag",
				projectReleaseTag,
				"--output", releaseNotesFilePath,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error running git cliff: %w", err)
			}

			log.Info(string(out))
		} else {
			cmd := exec.Command(
				"git",
				"cliff",
				"--include-path", filepath.Join(config.filePath, "**", "*"),
				"--strip", "all",
				"--tag",
				projectReleaseTag,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error running git cliff: %w", err)
			}

			log.Info(string(out))
		}

		log.Info("adding files to be committed", "files", releaseNotesFilePath)

		if !config.dryRun {
			cmd := exec.Command(
				"git",
				"add",
				releaseNotesFilePath,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error adding release notes: %w", err)
			}

			log.Info(string(out))
		}

		commitMessage := fmt.Sprintf("Release %s %s", config.releaseProject, config.releaseTag)

		log.Info("commiting the release notes", "message", commitMessage)

		if !config.dryRun {
			cmd := exec.Command(
				"git",
				"commit",
				"-m",
				commitMessage,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error adding release notes: %w", err)
			}

			log.Info(string(out))
		}

		log.Info("creating git tag", "tag", projectReleaseTag)

		if !config.dryRun {
			cmd := exec.Command(
				"git",
				"tag",
				projectReleaseTag,
			)

			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("error adding release notes: %w", err)
			}

			log.Info(string(out))
		}

		return nil
	}
}
