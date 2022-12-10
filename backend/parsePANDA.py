
class pandaInputs:
    base = "~/panda_studio"
    extra = ["-nographic"]
    architecture = 'arm'
    msize = "256"
    cow = None
    osv = None
    kwargs = None
    prompt = None

def parseImage(Image_File, inputs):
    #Image_File = open("/panda_studio/backend/imageSpec.txt", "r")
    for line in Image_File:
        halves = line.split(",")
        print(halves)
        if (halves[0] == "Arch"):
            inputs.architecture = halves[1]
            print("arch = ", inputs.architecture)
        elif (halves[0] == "Msize"):
            inputs.msize = halves[1]
            print("Msize = ", inputs.msize)
        elif (halves[0] == "Qcow"):
            inputs.cow = halves[1]
            print("QCow = ", inputs.cow)
        elif (halves[0] == "prompt"):
            inputs.prompt = halves[1]
            print("Prompt = ", inputs.prompt)
        elif (halves[0] == "serial_kwargs"):
            inputs.kwargs = halves[1]
            print("Kwargs = ", inputs.kwargs)
        elif (halves[0] == "os_version"):
            inputs.osv = halves[1]
            print("OSV = ", inputs.osv)
        elif (halves[0] == "Extra"):
            inputs.extra = halves[1]
        elif (halves[0] == "Other Descriptions Passed in"):
            inputs.architecture = halves[1]



 # for line in Image_File:
    #     halves = line.split(",")
    #     print(halves)
    #     if (halves[0] == "Arch"):
    #         architecture = halves[1]
    #         print("arch = ", architecture)
    #     elif (halves[0] == "Msize"):
    #         msize = halves[1]
    #         print("Msize = ", msize)
    #     elif (halves[0] == "Qcow"):
    #         cow = halves[1]
    #         print("QCow = ", cow)
    #     elif (halves[0] == "prompt"):
    #         prompt = halves[1]
    #         print("Prompt = ", prompt)
    #     elif (halves[0] == "serial_kwargs"):
    #         kwargs = halves[1]
    #         print("Kwargs = ", kwargs)
    #     elif (halves[0] == "os_version"):
    #         osv = halves[1]
    #         print("OSV = ", osv)
    #     elif (halves[0] == "Extra"):
    #         extra = halves[1]
    #     elif (halves[0] == "Other Descriptions Passed in"):
    #         architecture = halves[1]
    #     elif (halves[0] == "Name"):
    #         PANDA_RECORD = halves[1]

def parseInteractions(Interaction_File, commands):
#Interaction_File = open("/panda_studio/backend/interactions.txt", "r")

    for line in Interaction_File:
        halves = line.split(",")
        print("Interaction Line: ", halves)
        if halves[0] == "serial":
            print("Serial Interaction = ", halves[1])
            commands.append(halves[1])
    print("Interactions = ", commands)


        # for line in Interaction_File:
    #     halves = line.split(",")
    #     print("Interaction Line: ", halves)
    #     if halves[0] == "serial":
    #         print("Serial Interaction = ", halves[1])
    #         commands.append(halves[1])