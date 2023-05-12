#!/bin/bash

if [ -z ${BASE_IMAGE} ];
then 
    echo "base image not set"
    exit 1
fi

if [ -z ${NEW_IMAGE} ];
then 
    echo "new image name not set"
    exit 2
fi

if [ -z ${DOCKER_IMAGE} ];
then 
    echo "docker image name not set"
    exit 3
fi

if [ -z ${SIZE} ];
then 
    echo "size name not set"
    exit 4
fi

if kvm-ok; then
    echo "KVM is installed and operational"
    export LIBGUESTFS_HV=$(which qemu-system-x86_64)
else
    echo "KVM is not installed... Defaulting to slower emulation mode"
    exit 5
fi

# -- expand size of the main drive --
cp ${BASE_IMAGE} ${NEW_IMAGE}	
# virt-resize --resize /dev/sda1=+${SIZE} ${BASE_IMAGE} ${NEW_IMAGE}
# expands into /dev/sda3

APT_PACKAGES="apt-transport-https curl gnupg-agent ca-certificates software-properties-common"
DOCKER_APT_PACKAGES="docker-ce docker-ce-cli containerd.io cgroupfs-mount"

# -- Open Guestfish --
guestfish --network -a ${NEW_IMAGE} -i << EOF
setenv DEBIAN_FRONTEND noninteractive
command "sudo apt update"
command "apt install -y ${APT_PACKAGES}"
command "wget https://download.docker.com/linux/ubuntu/gpg" 
command "apt-key add gpg"
command "add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable'"
command "apt install -y ${DOCKER_APT_PACKAGES}"
command "service docker start"
command "docker pull ${DOCKER_IMAGE}"
exit
EOF