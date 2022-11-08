from pandare import Panda
from pandare import autogen
import sys
import time
import os

base = "~/panda_class_materials"

#extra = [f"-kernel {base}/vmlinux-4.20-versatile", f"-initrd {base}/initramfz.eabi", "-append console=ttyAMA0 quiet", "-M versatilepb", f"-dtb {base}/versatile-pb-4.20.dtb", "-nographic"]
    #"-nographic", "-machine akita"]

extra = ["-nographic"]

architecture = 'arm'
#None
#"arm"
msize = "256"
#cow = "qcows/wheezy.qcow2"
cow = None

osv = None
kwargs = None
prompt = None

#  extra_args=f"-kernel {config['kernel']} \
#                      -initrd {config['initramfs']} \
#                      -append '{config['append']}' \
#                      -nographic {config['extra']}",


    #    'arm': { # We dynamically add DTB later if it exists for the version
    #      'mem': '256',
    #      'initramfs': f'{base}/initramfz.{arm_abi}',
    #      'kernel': f'{base}/vmlinux-{version}-versatile',
    #      'append': 'console=ttyAMA0 quiet',
    #      'extra': '-M versatilepb'},

#dtb = f'{base}/versatile-pb-{version}.dtb'

commands = []
if True:
    Image_File = open("/panda_class_materials/imageSpec.txt", "r")
    for line in Image_File :
        halves = line.split(",")
        print(halves)
        if(halves[0] == "Arch"):
            architecture = halves[1]
            print("arch = ", architecture)
        elif(halves[0] == "Msize"):
            msize = halves[1]
            print("Msize = ", msize)
        elif(halves[0] == "Qcow"):
            cow = halves[1]
            print("QCow = ", cow)
        elif(halves[0] == "prompt"):
            prompt = halves[1]
            print("Prompt = ", prompt)
        elif(halves[0] == "serial_kwargs"):
            kwargs = halves[1]
            print("Kwargs = ", kwargs)
        elif(halves[0] == "os_version"):
            osv = halves[1]
            print("OSV = ", osv)
        # Extra args expects multiple strings passed in as an array, therefore inthe file they are split at the
        elif(halves[0] == "Extra"):
            extra = halves[1]
            # for i in range(len(halves) - 2):
            #     print("Extra = ", extra)
            #     extra.append(halves[i+1])
            # print("Extra = ", extra)
        elif(halves[0] == "Other Descriptions Passed in"):
           architecture = halves[1]    

    Interaction_File = open("/panda_class_materials/interactions.txt", "r")

    for line in Interaction_File :
        halves = line.split(",")
        print("Interaction Line: ", halves)
        if halves[0] == "serial":
            print("Serial Interaction = ", halves[1])
            commands.append(halves[1])
    print("Interactions = ", commands)

# f = open("Image Spec", r)

# The qcow is the only real thing required, we will change the number of parameters based on what the CSSE team wants

# panda = Panda(qcow="qcows/wheezy.qcow2", extra_args=["-nographic"]) # Create an instance of panda
#panda = Panda(arch = architecture, mem=msize, qcow=cow, extra_args=extra, os="linux") # Create an instance of panda
# panda = Panda(arch = architecture, qcow=cow, extra_args=extra, mem= msize, expect_prompt=b'/ #', serial_kwargs={"unansi": False},
#          os_version='linux-32-os:0')
# f'{base}/vmlinux-{version}-versatile'
# f"-kernel /panda_class_materials/arm//vmlinux-4.20-versatile"
# panda = Panda(arch = 'arm', qcow=None, extra_args=f"-kernel /panda_class_materials/arm//vmlinux-4.20-versatile \
#         -initrd /panda_class_materials/arm//initramfz.eabi\
#         -append 'console=ttyAMA0 quiet' \
#         -nographic \
#         -M versatilepb \
#         -dtb /panda_class_materials/arm//versatile-pb-4.20.dtb", 
#          mem= '256', expect_prompt=b'/ #', serial_kwargs={"unansi": False},
#          os_version='linux-32-os:0')


panda = Panda(arch = architecture, qcow=cow, extra_args=extra, 
         mem= '256', expect_prompt=bytes(prompt,'utf-8'), serial_kwargs={"unansi": False},
         os_version=osv)
# Counter of the number of basic blocks
# blocks = 0

# Register a callback to run before_block_exec and increment blocks
# @panda.cb_before_block_exec
# def before_block_execute(cpustate, transblock):
#     global blocks
#     blocks += 1

# This 'blocking' function is queued to run in a seperate thread from the main CPU loop
# which allows for it to wait for the guest to complete commands
@panda.queue_blocking
def run_cmd():
    #start = time.time()
    # Wait until the first prompt
    panda.serial_read_until(bytes(prompt,'utf_8'))
    #panda.serial_console.use_unansi = True
    #print(f"System started after {time.time()-start:.2f}s")
    # First revert to the qcow's root snapshot (synchronously)
    #panda.revert_sync("root")
    # Then type a command via the serial port and print its results
    print(panda.run_serial_cmd("uname -a"))

    for command in commands:
        print(panda.run_serial_cmd(command))

    # When the command finishes, terminate the panda.run() call
    panda.end_analysis()

# Start the guest
panda.run()
print("Finished Running")
# print("Finished. Saw a total of {} basic blocks during execution".format(blocks))