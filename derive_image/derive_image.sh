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

# -- expand size of the main drive --
cp ${BASE_IMAGE} ${NEW_IMAGE}	
virt-resize --resize /dev/sda1=+${SIZE} ${BASE_IMAGE} ${NEW_IMAGE}
# expands into /dev/sda3

# -- Open Guestfish --
guestfish --network << EOF
add ${NEW_IMAGE}
set-network true
run
mount /dev/sda3 /
command "sudo apt update"
command "sudo apt install apt-transport-https curl gnupg-agent ca-certificates software-properties-common -y"
command "wget https://download.docker.com/linux/ubuntu/gpg" 
command "sudo apt-key add gpg"
command "sudo add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable'"
command "sudo apt install docker-ce docker-ce-cli containerd.io cgroupfs-mount -y"
command "sudo systemctl enable docker"
command "sudo systemctl start docker"
command "docker pull ${DOCKER_IMAGE}"
exit
EOF