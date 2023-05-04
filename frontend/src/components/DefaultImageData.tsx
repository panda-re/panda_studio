export const archOptions = [
  { 
    value: 'x86_64', 
    text: 'x86_64',
    defaultConfig: {
      os: "linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr",
      prompt: "root@ubuntu:.*#",
      cdrom: "ide1-cd0",
      snapshot: "root",
      url: "https://panda-re.mit.edu/qcows/linux/ubuntu/1804/x86_64/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2",
      default_mem: "1024",
      extra_args: "-display none",
    }
  },
  { 
    value: 'i386', 
    text: 'i386',
    defaultConfig: {
      os: "linux-32-ubuntu:4.4.200-170-generic",
      prompt: "root@instance-1:.*#",
      cdrom: "ide1-cd0",
      snapshot: "root",
      url: "https://panda-re.mit.edu/qcows/linux/ubuntu/1604/x86/ubuntu_1604_x86.qcow",
      default_mem: "1024",
      extra_args: "-display none"
    }
  },
  { 
    value: 'arm', 
    text: 'arm',
    defaultConfig: {
      os: "linux-32-debian:3.2.0-4-versatile-arm",
      prompt: "root@debian-armel:.*# ",
      cdrom: "scsi0-cd2",
      snapshot: "root",
      url: "https://panda-re.mit.edu/qcows/linux/debian/7.3/arm/debian_7.3_arm.qcow",
      default_mem: "128M",
      extra_args: "-display none"
    }
  },
  { 
    value: 'aarch64', 
    text: 'aarch64',
    defaultConfig: {
      os: "linux-64-ubuntu:5.4.0-58-generic-arm64",
      prompt: "root@ubuntu-panda:.*# ",
      cdrom: "",
      snapshot: "root",
      url: "https://panda-re.mit.edu/qcows/linux/ubuntu/2004/aarch64/ubuntu20_04-aarch64.qcow",
      default_mem: "1G",
      extra_args: "-nographic -machine virt -cpu cortex-a57 -drive file=~/.panda/ubuntu20_04-aarch64-flash0.qcow,if=pflash,readonly=on",
    }
  },
];