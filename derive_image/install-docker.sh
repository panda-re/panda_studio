APT_PACKAGES="apt-transport-https curl gnupg-agent ca-certificates software-properties-common"
DOCKER_APT_PACKAGES="docker-ce docker-ce-cli containerd.io cgroupfs-mount"
export DEBIAN_FRONTEND=noninteractive

apt update
apt install -y ${APT_PACKAGES}
wget https://download.docker.com/linux/ubuntu/gpg
apt-key add gpg
add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable'
apt install -y ${DOCKER_APT_PACKAGES}
service docker start && docker pull ${DOCKER_IMAGE}t