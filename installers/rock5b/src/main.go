// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
)

func main() {
	adapter.Execute(&Rock5BInstaller{})
}

type Rock5BInstaller struct{}

type rock5BExtraOptions struct {
	Console    []string `json:"console"`
	ConfigFile string   `json:"configFile"`
}

func (i *Rock5BInstaller) GetOptions(extra rock5BExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=tty0",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	kernelArgs = append(kernelArgs, extra.Console...)

	return overlay.Options{
		Name:       "rock5b",
		KernelArgs: kernelArgs,
	}, nil
}

func (i *Rock5BInstaller) Install(options overlay.InstallOptions[rock5BExtraOptions]) error {
	// allows to copy a directory from the overlay to the target
	// err := copy.Dir(filepath.Join(options.ArtifactsPath, "arm64/firmware/boot"), filepath.Join(options.MountPrefix, "/boot/EFI"))
	// if err != nil {
	// 	return err
	// }

	// allows to copy a file from the overlay to the target
	err := copy.File(filepath.Join(options.ArtifactsPath, "arm64/u-boot/rock5b/u-boot.bin"), filepath.Join(options.MountPrefix, "/boot/EFI/u-boot.bin"))
	if err != nil {
		return err
	}

	if options.ExtraOptions.ConfigFile != "" {
		// do something with the config file
	}

	return nil
}
