#!/usr/bin/env/python3
import time
import os
from sys import argv
from pandare import Panda

# CONFIG - read from argv, or hardcode here
base =    argv[1] if len(argv) > 1 else './generic/'
arch =    argv[2] if len(argv) > 2    else 'mipsel'
version = argv[3] if len(argv) > 3 else '2.6.39'

PANDA_RECORD = "1337H@cker"

print(f"Configuration: {base}, {arch} {version}")
# END CONFIG

base += f'/{arch}/'
versions = ['2.6.39' , '3.0' , '3.10' , '3.19' , '4.0' , '4.10' , '4.20' , '5.0' , '5.10' , '5.16']
arm_abi  = 'eabi'
if versions.index(version) < 2:
    arm_abi  = 'oabi'
assert(version in version)

configs = {'x86_64': {
         'mem': '1g',
         'kernel': f'{base}/linux{version}',
         'initramfs': f'{base}/initramfz',
         'append': 'console=ttyS0 quiet',
         'extra': ''},
       'arm': { # We dynamically add DTB later if it exists for the version
         'mem': '256',
         'initramfs': f'{base}/initramfz.{arm_abi}',
         'kernel': f'{base}/vmlinux-{version}-versatile',
         'append': 'console=ttyAMA0 quiet',
         'extra': '-M versatilepb'},
       'mips': {
         'mem': '512',
         'initramfs': f'{base}/initramfz',
         'kernel': f'{base}/vmlinux-{version}-malta',
         'append': 'console=ttyS0 quiet',
         'extra': '-M malta'},
       'mipsel': {
         'mem': '512',
         'initramfs': f'{base}/initramfz',
         'kernel': f'{base}/vmlinux-{version}-malta',
         'append': 'console=ttyS0 quiet',
         'extra': '-M malta'},
      }

config = configs[arch]

if arch == 'arm':
    # If there's a DTB available, it's because you need it
    dtb = f'{base}/versatile-pb-{version}.dtb'
    if os.path.isfile(dtb):
        config['extra'] += f' -dtb {dtb}'

panda = Panda(arch=arch, qcow=None, mem=config['mem'],
        extra_args=f"-kernel {config['kernel']} \
                     -initrd {config['initramfs']} \
                     -append '{config['append']}' \
                     -nographic {config['extra']}",
         expect_prompt=b'/ #', serial_kwargs={"unansi": False},
         os_version='linux-32-os:0')

@panda.queue_blocking
def driver():
    start = time.time()
    # Wait until the first prompt
    panda.serial_read_until(b"/ #")
    panda.serial_console.use_unansi = True
    print(f"System started after {time.time()-start:.2f}s")

    panda.record(PANDA_RECORD, snapshot_name=None)
    os = panda.run_serial_cmd("uname -a")
    username = panda.run_serial_cmd("whoami")
    print(f"Running as {username} on {os}")
    print( panda.run_serial_cmd("echo 42069, 1337 C0d3r H@ck3r Skl11s"))
    panda.end_record()

    panda.end_analysis()

panda.run()