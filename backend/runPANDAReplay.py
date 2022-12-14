from pandare import Panda
import parsePANDA

PANDA_RECORD = "exampleRecording"
PANDA_DEST = "/panda_studio/recordings/"

inputs = parsePANDA.pandaInputs

if True:
    Image_File = open("/panda_studio/backend/imageSpec.txt", "r")
    parsePANDA.parseImage(Image_File, inputs)

panda = Panda(arch=inputs.architecture, qcow=inputs.cow, extra_args=inputs.extra,
              mem=inputs.msize, expect_prompt=bytes(inputs.prompt, 'utf-8'), serial_kwargs={"unansi": False},
              os_version=inputs.osv)

# Start the guest
panda.run_replay(PANDA_DEST + PANDA_RECORD)
print("Finished Running Replay")
