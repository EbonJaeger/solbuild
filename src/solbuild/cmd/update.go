//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"builder"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a solbuild profile",
	Long: `Update the base image of the specified solbuild profile, helping to
minimize the build times in future updates with this profile.`,
	Run: updateProfile,
}

func init() {
	updateCmd.Flags().StringVarP(&profile, "profile", "p", builder.DefaultProfile, "Build profile to use")
	RootCmd.AddCommand(updateCmd)
}

func updateProfile(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		profile = strings.TrimSpace(args[0])
	}

	if !builder.IsValidProfile(profile) {
		builder.EmitProfileError(profile)
		return
	}

	// Updating is handled all within the library itself
	bk := builder.NewBackingImage(profile)

	if !bk.IsInstalled() {
		fmt.Fprintf(os.Stderr, "Cannot find profile '%s'. Did you forget to run init?\n", profile)
		os.Exit(1)
	}

	if os.Geteuid() != 0 {
		fmt.Fprintf(os.Stderr, "You must be root to run init profiles\n")
		os.Exit(1)
	}

	if err := bk.Update(); err != nil {
		os.Exit(1)
	}
}