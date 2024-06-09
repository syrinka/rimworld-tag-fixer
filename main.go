package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var NVW_SteamID = "2599504692" // No Version Warning

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
	command.Flags().StringP("version", "v", "", "the target version you're trying to fix")
	command.Flags().StringArrayP("file", "f", []string{}, "external ModIdsToFix files")
	command.Flags().BoolP("yes", "y", false, "don't re-confirm")
}

func main() {
	command.Execute()
}

var command = &cobra.Command{
	Use:   "tag-fixer",
	Short: "fix those missing tags!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var workshop = args[0]
		verifyRunnable(workshop)

		tag, _ := cmd.Flags().GetString("version")
		if tag == "" {
			color.Red.Println("[FATAL] no version specified?")
			os.Exit(1)
		}

		basefile := filepath.Join(workshop, NVW_SteamID, tag, "ModIdsToFix.xml")
		extfiles, _ := cmd.Flags().GetStringArray("file")

		confirm, _ := cmd.Flags().GetBool("yes")
		if !confirm {
			fmt.Printf("target version tag: %s\n", tag)
			fmt.Println("ModIdsToFix files:")
			fmt.Printf("- %s\n", basefile)
			for _, e := range extfiles {
				fmt.Printf("- %s\n", e)
			}
			fmt.Printf("Enter to continue, or Ctrl+C to abort")
			os.Stdin.Read(make([]byte, 1))
		}

		var ids = collectFixable(basefile)
		for _, e := range extfiles {
			ids = append(ids, collectFixable(e)...)
		}

		entries, _ := os.ReadDir(workshop)
		for _, e := range entries {
			var meta ModMetaData
			err := meta.Init(filepath.Join(workshop, e.Name(), "About", "About.xml"))
			if err != nil {
				fmt.Printf("[WARN] fail to operating on %s\n", e.Name())
				continue
			}

			fmt.Printf("(%s) %s ", e.Name(), meta.Name())
			if meta.ContainVersionTag(tag) {
				color.Gray.Println("[tag existed, skip]")
			} else if slices.Contains(ids, meta.Id()) {
				color.LightGreen.Println("[no tag, fixable, fix!]")
				meta.AddVersionTag(tag)
				meta.Update()
			} else {
				color.Red.Println("[no tag, not fixable, skip]")
			}
		}
	},
}
