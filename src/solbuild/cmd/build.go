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
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var buildCmd = &cobra.Command{
	Use:   "build [package.yml|pspec.xml]",
	Short: "build a package",
	Long: `Build the given package in a chroot environment, and upon success,
store those packages in the current directory`,
	RunE: buildPackage,
}

func init() {
	buildCmd.Flags().StringVarP(&profile, "profile", "p", builder.DefaultProfile, "Build profile to use")
	RootCmd.AddCommand(buildCmd)
}

func buildPackage(cmd *cobra.Command, args []string) error {
	pkgPath := ""

	if len(args) == 1 {
		pkgPath = args[0]
	} else {
		// Try to find the logical path..
		pkgPath = FindLikelyArg()
	}

	if !builder.IsValidProfile(profile) {
		builder.EmitProfileError(profile)
		return nil
	}

	pkgPath = strings.TrimSpace(pkgPath)

	if pkgPath == "" {
		return errors.New("Require a filename to build")
	}

	// Complain about missing profile
	bk := builder.NewBackingImage(profile)
	if !bk.IsInstalled() {
		fmt.Fprintf(os.Stderr, "Cannot find profile '%s'. Did you forget to run init?\n", profile)
		return nil
	}

	pkg, err := builder.NewPackage(pkgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load package: %v\n", err)
		return nil
	}

	log.WithFields(log.Fields{
		"profile": profile,
		"version": pkg.Version,
		"package": pkg.Name,
		"type":    pkg.Type,
		"release": pkg.Release,
	}).Info("Building package")
	if pkg.Type != builder.PackageTypeYpkg {
		log.Warning("Full sandboxing is not possible with legacy format")
	}
	return nil
}