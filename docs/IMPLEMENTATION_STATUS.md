# Implementation Status of Planned Features

This document outlines the original planned features and their implementation status (implemented or not implemented).
Below is a list of the planned features that were fully implemented:

* Runs PANDA with a specified image and interactions in an automated way
* Recordings can be created or deleted from front end application
* Displays dashboard of recordings for users to view and select
* Displays dashboard of existing images to manage and use for recordings
* Displays dashboard of interaction programs for users to view and manage
* Specifies a list of interactions for a particular recording
* Images and recordings are downloadable
* Allows users to manually verify recordings


Below is a list of features with incomplete or no implementation:
* Recordings can be replayed from the frontend application
  * There is currently support for recording replay in the backend, but no UI work has been done
* Custom image derivation
  * Blocked by a bug in running Docker: following install, docker daemon will not start inside the qcow image.
  * Also currently only supports image derivation for Ubuntu images.
  * Otherwise, full stack implementation of this is complete. 
  * See the derive_image branch for implementation of this feature.
* Network and filesystem interactions
  * The framework for supporting these interaction types exists in the API, but a detailed solution does not yet exist
* To interface with the backend directly you can try out the [executor](/cmd/panda_executor/panda_executor.go).
  * Directions for use can be found in the root [README](README.md)
