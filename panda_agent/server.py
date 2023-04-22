import grpc
import concurrent.futures as futures
import os

from pandare import Panda

# protobuf compiler is stupid
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
    def __init__(self, server):
        self.server = server
        self.agent = None
    
    def StartAgent(self, request: pb.StartAgentRequest, context):
        if self.agent is not None:
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")

        print("Starting agent")
        self.agent = PandaAgent(request.config)

        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        sleep(0.5) # ensures internal flags get set
        return pb.StartAgentResponse()
    
    def StopAgent(self, request: pb.StopAgentRequest, context):
        if self.agent is not None:
            self.agent.stop()
            self.server.stop(grace=5)
            return pb.StopAgentResponse()
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot stop a PANDA instance when one is not running")
    
    def RunCommand(self, request: pb.RunCommandRequest, context):
        if self.agent is not None:
            output = self.agent.run_command(request.command)
            return pb.RunCommandResponse(output=output)
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot run a function when PANDA is not running")
    
    def StartRecording(self, request: pb.StartRecordingRequest, context):
        if self.agent is not None:
            self.agent.start_recording(recording_name=request.recording_name)
            return pb.StartRecordingResponse()
        else:
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot start a recording when PANDA is not running")
    
    def StopRecording(self, request: pb.StopRecordingRequest, context):
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
        if self.agent is not None and self.agent.panda.started.is_set():
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")
        
        print("Starting replay")
        self.agent = PandaAgent(request.config)

        serial = self.agent.start_replay(request.recording_name)
        log = read_execution_log()
        
        return pb.StartReplayResponse(serial=serial, replay=log)

    def StopReplay(self, request: pb.StopReplayRequest, context):
        if self.agent is not None:    
            serial = self.agent.stop_replay()
            log = read_execution_log()
            return pb.StopReplayResponse(serial=serial, replay=log)
        else:
            raise RuntimeError(ErrorCode.NOT_REPLAYING, "Must start a replay before stopping one")

    def SendNetworkCommand(self, request: pb.NetworkRequest, context):
        
        response = self.agent.execute_network_command(request)

        return pb.NetworkResponse(0, response)

# Reads the execution log that is teed from the container startup
# Log contains PyPANDA and Agent output
def read_execution_log(self):
    try:
        os.stat(f"{FILE_PREFIX}/{EXECUTION_LOG}")
        with (open(f"{FILE_PREFIX}/{EXECUTION_LOG}")) as file:
            return file.read()
    except (OSError, ValueError):
        return "Execution log disabled. See Dockerfile.panda-agent"

def serve():
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