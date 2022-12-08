from pandare import Panda
import shutil

base = "~/panda_studio"
extra = ["-nographic"]
architecture = 'arm'
msize = "256"
cow = None
osv = None
kwargs = None
prompt = None

PANDA_RECORD = "exampleRecording"
PANDA_DEST = "/panda_studio/recordings/"

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
#   dtb = f'{base}/versatile-pb-{version}.dtb'

commands = []
if True:
    Image_File = open("/panda_studio/backend/imageSpec.txt", "r")
    for line in Image_File:
        halves = line.split(",")
        print(halves)
        if (halves[0] == "Arch"):
            architecture = halves[1]
            print("arch = ", architecture)
        elif (halves[0] == "Msize"):
            msize = halves[1]
            print("Msize = ", msize)
        elif (halves[0] == "Qcow"):
            cow = halves[1]
            print("QCow = ", cow)
        elif (halves[0] == "prompt"):
            prompt = halves[1]
            print("Prompt = ", prompt)
        elif (halves[0] == "serial_kwargs"):
            kwargs = halves[1]
            print("Kwargs = ", kwargs)
        elif (halves[0] == "os_version"):
            osv = halves[1]
            print("OSV = ", osv)
        elif (halves[0] == "Extra"):
            extra = halves[1]
        elif (halves[0] == "Other Descriptions Passed in"):
            architecture = halves[1]
        elif (halves[0] == "Name"):
            PANDA_RECORD = halves[1]

    Interaction_File = open("/panda_studio/backend/interactions.txt", "r")

    for line in Interaction_File:
        halves = line.split(",")
        print("Interaction Line: ", halves)
        if halves[0] == "serial":
            print("Serial Interaction = ", halves[1])
            commands.append(halves[1])
    print("Interactions = ", commands)

# The qcow is the only real thing required, we will change the number of parameters based on what the CSSE team wants
# panda = Panda(qcow="qcows/wheezy.qcow2", extra_args=["-nographic"]) # Create an instance of panda

# panda = Panda(arch = 'arm', qcow=None, extra_args=f"-kernel /panda_class_materials/arm//vmlinux-4.20-versatile \
#         -initrd /panda_class_materials/arm//initramfz.eabi\
#         -append 'console=ttyAMA0 quiet' \
#         -nographic \
#         -M versatilepb \
#         -dtb /panda_class_materials/arm//versatile-pb-4.20.dtb",
#          mem= '256', expect_prompt=b'/ #', serial_kwargs={"unansi": False},
#          os_version='linux-32-os:0')


panda = Panda(arch=architecture, qcow=cow, extra_args=extra,
              mem=msize, expect_prompt=bytes(prompt, 'utf-8'), serial_kwargs={"unansi": False},
              os_version=osv)


# This 'blocking' function is queued to run in a seperate thread from the main CPU loop
# which allows for it to wait for the guest to complete commands
@panda.queue_blocking
def run_cmd():

    # Wait until the first prompt
    panda.serial_read_until(bytes(prompt, 'utf_8'))

    print(panda.run_serial_cmd("uname -a"))

    panda.record(PANDA_RECORD, snapshot_name=None)
    for command in commands:
        print(panda.run_serial_cmd(command))
    panda.end_record()
    # When the command finishes, terminate the panda.run() call
    panda.end_analysis()


# Start the guest
panda.run()
print("Finished Running")

# Move the now complete recording to the backend folder (shared volume)
shutil.move(PANDA_RECORD + "-rr-nondet.log",
            PANDA_DEST + PANDA_RECORD + "-rr-nondet.log")
shutil.move(PANDA_RECORD + "-rr-snp", 
            PANDA_DEST + PANDA_RECORD + "-rr-snp")
