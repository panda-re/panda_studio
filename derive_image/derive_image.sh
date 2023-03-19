#!/bin/bash
# -- expand size of the main drive --
cp ${BASE_IMAGE_NAME} ${NEW_IMAGE_NAME}	
virt-resize --expand /dev/sda1 ${BASE_IMAGE_NAME} ${NEW_IMAGE_NAME}

# -- Open Guestfish --
sudo guestfish --network << EOF
add ${NEW_IMAGE_NAME}
set-network true
run
mount /dev/sda3 /
command "sudo apt update"
command "sudo apt install apt-transport-https curl gnupg-agent ca-certificates software-properties-common -y"
command "wget https://download.docker.com/linux/ubuntu/gpg" 
command "sudo apt-key add gpg"
command "sudo add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable'"
command "sudo apt install docker-ce docker-ce-cli containerd.io cgroupfs-mount -y"
command "sudo service docker start"
command "docker pull ${DOCKER_IMAGE_NAME}"
exit
EOF