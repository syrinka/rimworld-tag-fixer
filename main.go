package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
)

var NVW_SteamID = "2599504692" // No Version Warning

func detectWorkshopPath() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\Valve\\Steam", registry.ALL_ACCESS)
	if err != nil {
		return "", err
	}
	path, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "steamapps", "workshop", "content", "294100"), nil
}

func verifyRunnable(workshop string) {
	info, err := os.Stat(workshop)
	if err != nil || info.Name() != "294100" || !info.IsDir() {
		color.Red.Println("[FATAL] invalid workshop path")
		os.Exit(1)
	}

	var nvw = filepath.Join(workshop, NVW_SteamID)
	_, err = os.Stat(nvw)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		color.Yellow.Println("seems like you haven't subscribed the \"No Version Warning\" mod. Without that, this program will essentially do nothing.")
	}
}

func init() {
	command.Flags().StringVarP(&FlagVersion, "version", "v", "", "the target version you're trying to fix")
	command.Flags().StringArrayVarP(&FlagFiles, "file", "f", nil, "external ModIdsToFix files")
	command.Flags().BoolVarP(&FlagYes, "yes", "y", false, "don't re-confirm")
	command.Flags().BoolVarP(&FlagQuiet, "quiet", "q", false, "hide skip info")
}

func main() {
	command.Execute()
}

var (
	FlagVersion string
	FlagFiles   []string
	FlagYes     bool
	FlagQuiet   bool
)

var command = &cobra.Command{
	Use:   "tag-fixer",
	Short: "fix those missing tags!",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var workshop string
		if len(args) == 1 {
			workshop = args[0]
		} else if runtime.GOOS == "windows" {
			var err error
			workshop, err = detectWorkshopPath()
			if err != nil {
				color.Red.Println("[FATAL] fail to detect steam workshop path, please pass it as argument")
				os.Exit(1)
			}
		} else {
			color.Red.Println("[FATAL] steam workshop path is required as first argument")
			os.Exit(1)
		}
		verifyRunnable(workshop)

		if FlagVersion == "" {
			color.Red.Println("[FATAL] use -v to specify a version to fix")
			os.Exit(1)
		}

		basefile := filepath.Join(workshop, NVW_SteamID, FlagVersion, "ModIdsToFix.xml")

		if !FlagYes {
			fmt.Printf("target version tag: %s\n", FlagVersion)
			fmt.Println("ModIdsToFix files:")
			fmt.Printf("- %s\n", basefile)
			for _, file := range FlagFiles {
				fmt.Printf("- %s\n", file)
			}
			fmt.Printf("Enter to continue, or Ctrl+C to abort")
			os.Stdin.Read(make([]byte, 1))
		}

		var ids = collectFixable(basefile)
		for _, file := range FlagFiles {
			ids = append(ids, collectFixable(file)...)
		}

		entries, _ := os.ReadDir(workshop)
		for _, e := range entries {
			var meta ModMetaData
			err := meta.Init(filepath.Join(workshop, e.Name(), "About", "About.xml"))
			if err != nil {
				color.LightRed.Printf("[WARN] fail to operating on %s\n", e.Name())
				continue
			}

			if meta.ContainVersionTag(FlagVersion) {
				if !FlagQuiet {
					fmt.Printf("(%s) %s ", e.Name(), meta.Name())
					color.Gray.Println("[tag existed, skip]")
				}
			} else if slices.Contains(ids, meta.Id()) {
				fmt.Printf("(%s) %s ", e.Name(), meta.Name())
				color.LightGreen.Println("[no tag, fixable, fix!]")
				meta.AddVersionTag(FlagVersion)
				meta.Update()
			} else {
				fmt.Printf("(%s) %s ", e.Name(), meta.Name())
				color.Red.Println("[no tag, not fixable, skip]")
			}
		}
	},
}
