# Commands and Interactions Available for APIs

## Command Types and Meanings

### Serial Interactions

Run the given command in the serial console of the backend's PANDA instance.

cmd
    - type: "cmd"
    - cmd: serial command to be run ("touch FILE")

### Start and Stop Recording Interactions

Start or stop a recording. Starting a recording also requires a recording_name.

start_recording
    - type: "start_recording"
    - recording_name: name of the recording ("test_recording")

stop_recording
    - type: "stop_recording"

### Filesystem Interactions

Not yet Implemented

filesystem
    - type: "filesystem"
    - file: filename/disk image to add to system ("file.txt"/"image.io")

### Network Interactions

Specify a network packet for the backend to send to its PANDA instance. The packet requires a socket type, packet type, port number and packet data

network
    - type: "network"
    - sock_type: socket type, choices are "tcp" or "udp"
    - port: port number to send packet on (443)
    - packet_type: type of packet, supported applications are:
        - http: send an http packet, packet_data field is used to specify file  to GET, POST, PUT, etc. 
        - other: any other application is treated as a custom packet
    - packet_data: raw bytes/data to send, use may change depending on pakcet type
        - http usage: "/index"
        - example custom usage: bytes("GET /index  HTTP/1.1\r\n\r\n")
    - additional fields: Some application types will have extra fields specific to their operation
        - http fields: request_type
            - request_type: type of HTTP packet ("GET")
        
    

