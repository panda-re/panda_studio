import grpc
import concurrent.futures as futures
import os

import sys
sys.path.append(".")
sys.path.append("pb")

import pb.panda_agent_pb2_grpc as pb_grpc
import pb.panda_agent_pb2 as pb
from agent import PandaAgent
from agent import ErrorCode
from time import sleep

FILE_PREFIX="data"
EXECUTION_LOG = "execution.log"

PORTS = [
    "[::]:50051",
    "unix:///panda/shared/panda-agent.sock"
]

executor = futures.ThreadPoolExecutor(max_workers=10)

class PandaAgentServicer(pb_grpc.PandaAgentServicer):
    '''
    Servicer that accepts gRPC messages and calls the corresponding agent functions.
    '''
    def __init__(self, server):
        self.server = server
        self.agent = None
    
    def StartAgent(self, request: pb.StartAgentRequest, context):
        '''
        Starts a PANDA instance with the specified configuration. It also loads the snapshot specified.

        Only one PANDA instance per agent can be running at one time.

        Args:
            request: gRPC request that contains the PANDA configuration parameters

        Returns:
            gRPC response acknowledging PANDA starting
        
        Raises:
            RuntimeError if PANDA was already running
        '''
        if self.agent is not None:
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")

        print("Starting agent")
        self.agent = PandaAgent(request.config)

        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        sleep(0.5) # ensures internal flags get set
        return pb.StartAgentResponse()
    
    def StopAgent(self, request: pb.StopAgentRequest, context):
        '''
        Stops the server and the agent

        Args:
            request: gRPC request

        Returns:
            gRPC response acknowledging PANDA stopping

        Raises:
            RuntimeError if PANDA was already stopped
            RuntimeWarning if a recording was in progress
        '''
        if self.agent is not None:
            self.agent.stop()
            self.server.stop(grace=5)
            return pb.StopAgentResponse()
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot stop a PANDA instance when one is not running")
    
    def RunCommand(self, request: pb.RunCommandRequest, context):
        '''
        Runs a serial command in PANDA

        Args:
            request: gRPC request containing the serial command
        
        Returns:
            gRPC response containing the serial output

        Raises:
            RuntimeError if PANDA was not running
        '''
        if self.agent is not None:
            output = self.agent.run_command(request.command)
            return pb.RunCommandResponse(output=output)
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot run a function when PANDA is not running")
    
    def StartRecording(self, request: pb.StartRecordingRequest, context):
        '''
        Starts a PANDA recording

        Args:
            request: gRPC request containing the recording name

        Returns:
            gRPC response acknowledging the recording started
        
        Raises:
            RuntimeError if PANDA was not running
            RuntimeError if PANDA was already recording
        '''
        if self.agent is not None:
            self.agent.start_recording(recording_name=request.recording_name)
            return pb.StartRecordingResponse()
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot start a recording when PANDA is not running")
    
    def StopRecording(self, request: pb.StopRecordingRequest, context):
        '''
        Stops a PANDA recording

        Args:
            request: gRPC request

        Returns:
            gRPC response containing the recording name and the filename of the recording files

        Raises:
            RuntimeError if PANDA was not running
            RuntimeWarning if PANDA was not recording
        '''
        if self.agent is not None:
            recording_name = self.agent.stop_recording()
            return pb.StopRecordingResponse(
                recording_name=recording_name,
                snapshot_filename=f"{recording_name}-rr-snp",
                ndlog_filename=f"{recording_name}-rr-nondet.log"
            )
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot stop a recording when PANDA is not running")

    def StartReplay(self, request: pb.StartReplayRequest, context):
        '''
        Starts a PANDA replay with the specified configuration.

        Only one PANDA instance per agent can be running at one time.

        Args:
            request: gRPC request containing the replay name

        Returns:
            gRPC response containing the serial and PANDA output of the replay

        Raises:
            RuntimeError if PANDA was already running
            RuntimeError if recording files to replay to not exist
        '''
        if self.agent is not None and self.agent.panda.started.is_set():
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")
        
        print("Starting replay")
        self.agent = PandaAgent(request.config)

        serial = self.agent.start_replay(request.recording_name)
        log = read_execution_log()
        
        return pb.StartReplayResponse(serial=serial, replay=log)

    def StopReplay(self, request: pb.StopReplayRequest, context):
        '''
        Stops a PANDA replay.

        PANDA replays stop naturally.

        Args:
            request: gRPC request

        Returns:
            gRPC response containing the serial and PANDA output of the replay

        Raises:
            RuntimeError if PANDA was not replaying
        '''
        if self.agent is not None:    
            serial = self.agent.stop_replay()
            log = read_execution_log()
            return pb.StopReplayResponse(serial=serial, replay=log)
        else:
            raise RuntimeError(ErrorCode.NOT_REPLAYING, "Must start a replay before stopping one")

    def SendNetworkCommand(self, request: pb.NetworkRequest, context):
        '''
        Executes a network command

        Args:
            request: gRPC request containing the networking information such as application, command, etc.
        
        Returns:
            gRPC response containing the status and response
        '''
        response = self.agent.execute_network_command(request)

        return pb.NetworkResponse(0, response)

def read_execution_log():
    '''
    Reads the log that is teed from the container startup.

    Log contains the PyPANDA and Agent output.

    Returns:
        Data in the file `data/execution.log`
    '''
    try:
        os.stat(f"{FILE_PREFIX}/{EXECUTION_LOG}")
        with (open(f"{FILE_PREFIX}/{EXECUTION_LOG}")) as file:
            return file.read()
    except (OSError, ValueError):
        return "Execution log disabled. See Dockerfile.panda-agent"

def serve():
    '''
    Starts the gRPC server and the servicer
    '''
    print("Starting server...")
    server = grpc.server(executor)
    pb_grpc.add_PandaAgentServicer_to_server(
        PandaAgentServicer(server), server)

    for port in PORTS:
        server.add_insecure_port(port)
        print(f'panda agent grpc server listening on port {port}')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()