// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"fmt"
	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

const (
	off int64 = 512 * 64
	// https://github.com/u-boot/u-boot/blob/4de720e98d552dfda9278516bf788c4a73b3e56f/configs/rock-pi-4c-rk3399_defconfig#L7=
	dtb = "rockchip/rk3588-rock-5b.dtb"
)

func main() {
	adapter.Execute[rock5BExtraOptions](&Rock5BInstaller{})
}

type Rock5BInstaller struct{}

type rock5BExtraOptions struct {
	Console    []string `json:"console"`
	ConfigFile string   `json:"configFile"`
}

func (i *Rock5BInstaller) GetOptions(extra rock5BExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
		"slab_nomerge",
		"earlycon=uart8250,mmio32,0xfeb50000",
		"console=ttyFIQ0,1500000n8",
		"consoleblank=0",
		"console=ttyS2,1500000n8",
		"console=tty1",
		"loglevel=7",
		"cgroup_enable=cpuset",
		"swapaccount=1",
		"irqchip.gicv3_pseudo_nmi=0",
		"coherent_pool=2M",
	}

	kernelArgs = append(kernelArgs, extra.Console...)

	return overlay.Options{
		Name:       "rock5b",
		KernelArgs: kernelArgs,
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *Rock5BInstaller) Install(options overlay.InstallOptions[rock5BExtraOptions]) error {
	var f *os.File

	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", options.InstallDisk, err)
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(filepath.Join(options.ArtifactsPath, "arm64/u-boot/rock5b/u-boot-rockchip.bin"))
	if err != nil {
		return err
	}

	if _, err = f.WriteAt(uboot, off); err != nil {
		return err
	}

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "/boot/EFI/dtb", dtb)

	err = os.MkdirAll(filepath.Dir(dst), 0o600)
	if err != nil {
		return err
	}

	return copy.File(src, dst)
}
