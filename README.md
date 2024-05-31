# talos-sbc-rk3588

Talos overlay and custom kernel for Rockchip RK3588 ARM64 single-board computer (SBC) devices such as the Radxa Rock 5 series.

# Why does this exist?

Currently, mainline support for the RK3588 chipset is still evolving and requires a newer kernel than the 6.6 LTS available in upstream Talos Linux.
Additionally, not all patches have been merged.

This project relies on the forks of U-Boot and Linux kernel for RK3588 maintained by Collabora.

# Device Support

* [Rock 5B](https://wiki.radxa.com/Rock5/5B)
* [Rock 5A](https://wiki.radxa.com/Rock5/5a)

# Install

Flashable images are available from the [releases](https://github.com/milas/talos-sbc-rk3588/releases/latest).

You can write this to your eMMC/SD card using `dd`, Balena Etcher, etc.
(I prefer [Caligula](https://github.com/ifd3f/caligula).)

# Machine Configuration

Use the `ghcr.io/milas/talos-rk3588` images instead of the upstream Talos Linux images.
These include a device appropriate U-Boot and kernel with RK3588 hardware support.

```yaml
machine:
  install:
    # for eMMC, use /dev/mmcblk0
    # for SD card, use /dev/mmcblk1
    disk: /dev/mmcblk0
    image: ghcr.io/milas/talos-rk3588:v1.7.4-rk3588.alpha.4-rock-5b
    wipe: false
```

# Resources

* [Collabora RK3588 upstreaming](https://gitlab.collabora.com/hardware-enablement/rockchip-3588)
* [siderolabs/talos](https://github.com/siderolabs/talos/)
* [milas/rock5-talos](https://github.com/milas/rock5-talos) - forked version of Talos v1.4 using BSP kernel

# Disclaimer

This is NOT supported or endorsed by Rockchip, Radxa, Sidero Labs, or Collabora - please do not go to them with support requests!
