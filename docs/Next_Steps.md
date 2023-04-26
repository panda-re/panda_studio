# Next Steps

PANDA Studio has a lot of potential now that the foundation has been laid. This document hints at potential next steps for adding functionality to PANDA Studio.

## Partially Implemented

The project has some loose ends that can be connected for more functionality.

### Interactions

Interactions are the key way for controlling PANDA and making recordings. There is some implementation in the [agent](../panda_agent/agent.py) but is left unimplemented in the [program executor](../internal/panda_controller/program_executor.go).

With filesystem and network interactions, certain arguments must be passed to the backend during image specification in order for the interactions to be able to run as expected. The current way of handling this is to assume that the user specified these options when they created and uploaded an image. This was a solution that was accepted as the image specification was one of the last pieces of the system implemented, and in total, does not lead to users having the most flexible options when creating recordings. The currently proposed solution to this problem is to search the selected interaction list when the backend is being started and then add the required parameters to the image specification. This would let the user only have to worry about the general characteristics of an image, and not have to create a new image that is specific to an interaction list. 

There is a lot of potential for more interactions and adding more to the current interactions.

### Replay

The agent has the ability to replay recordings, that functionality is also not available in the frontend.

## Branches

One improvement that was abandoned was streaming. This can be found in the [23-PANDA-VM-logs branch](https://github.com/panda-re/panda_studio/tree/23-PANDA-VM-logs). The idea was to use gRPC streams to provide quicker updates for potentially long-running interactions such as a serial command or replays. Due to time constraints, conflicts, and because it was an extra feature it was left behind. The functionality works, but on an old version of the PandaAgent interface.

## Unimplemented

The agent uses PyPANDA to accomplish various interactions with PANDA and gRPC does not utilize all of that potential. More gRPC messages and ProtoBuf protocols could be added to decrease the gap between PANDA Studio and PyPANDA. Some current messages could also be expanded. For example, replays could add plugins and serial commands could add a timeout argument.