from pandare import Panda
import shutil
import parsePANDA

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

inputs = parsePANDA.pandaInputs

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
    parsePANDA.parseImage(Image_File, inputs)
   
    Interaction_File = open("/panda_studio/backend/interactions.txt", "r")
    parsePANDA.parseInteractions(Interaction_File, commands)
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
panda = Panda(arch=inputs.architecture, qcow=inputs.cow, extra_args=inputs.extra,
              mem=inputs.msize, expect_prompt=bytes(inputs.prompt, 'utf-8'), serial_kwargs={"unansi": False},
              os_version=inputs.osv)



# This 'blocking' function is queued to run in a seperate thread from the main CPU loop
# which allows for it to wait for the guest to complete commands
@panda.queue_blocking
def run_cmd():

    # Wait until the first prompt
    panda.serial_read_until(bytes(inputs.prompt, 'utf_8'))

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
