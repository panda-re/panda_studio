from pandare import Panda

PANDA_RECORD = "Recording"
PANDA_DEST = "/tmp/panda-studio/"

panda = Panda(arch='x86_64', qcow='/root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2', mem='1024',
                 os='linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr', expect_prompt='root@ubuntu:.*# ',
                 extra_args='-display none')

# Start the guest
panda.run_replay(PANDA_DEST + PANDA_RECORD)
print("Finished Running Replay")