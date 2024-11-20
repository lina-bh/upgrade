FROM docker.io/library/archlinux:base-devel

RUN pacman -Syu --noconfirm go buildah git fuse-overlayfs just
ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot
