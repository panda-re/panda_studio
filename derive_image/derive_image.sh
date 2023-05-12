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
    echo "KVM is not available. Execution will be slower"
    exit 5
fi

# -- expand size of the main drive --
cp ${BASE_IMAGE} ${NEW_IMAGE}	
# virt-resize --resize /dev/sda1=+${SIZE} ${BASE_IMAGE} ${NEW_IMAGE}
# expands into /dev/sda3

# -- Open Guestfish --
guestfish --network -a ${NEW_IMAGE} -i << EOF
copy-file-to-device /root/install-docker.sh /root/install-docker.sh
sh "/root/install-docker.sh"
exit
EOF