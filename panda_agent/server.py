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

PORTS = [
    "[::]:50051",
    "unix:///panda/shared/panda-agent.sock"
]

executor = futures.ThreadPoolExecutor(max_workers=10)

class PandaAgentServicer(pb_grpc.PandaAgentServicer):
    def __init__(self, server, agent: PandaAgent):
        self.server = server
        self.agent = agent
        self.agent_thread = None
    
    def isAgentRunning(self):
        return self.agent_thread is not None
    
    def StartAgent(self, request: pb.StartAgentRequest, context):
        if self.agent.hasStarted:
            return

        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        return pb.StartAgentResponse()
    
    def StopAgent(self, request: pb.StopAgentRequest, context):
        self.agent.stop()
        self.server.stop(grace=5)
        return pb.StopAgentResponse()
    
    def RunCommand(self, request: pb.RunCommandRequest, context):
        output = self.agent.run_command(request.command)
        return pb.RunCommandReply(statusCode=0, output=output)
    
    def StartRecording(self, request: pb.StartRecordingRequest, context):
        self.agent.start_recording(recording_name=request.recording_name)
        return pb.StartRecordingResponse()
    
    def StopRecording(self, request: pb.StopRecordingRequest, context):
        recording_name = self.agent.stop_recording()
        return pb.StopRecordingResponse(
            recording_name=recording_name,
            snapshot_filename=f"{recording_name}-rr-snp",
            ndlog_filename=f"{recording_name}-rr-nondet.log"
        )

def serve():
    panda = Panda(arch='x86_64', qcow='/panda/shared/system_image.qcow2', mem='1024',
                 os='linux-64-ubuntu:4.15.0-72-generic-noaslr-nokaslr', expect_prompt='root@ubuntu:.*# ',
                 extra_args='-display none')
    agent = PandaAgent(panda)
    server = grpc.server(executor)
    pb_grpc.add_PandaAgentServicer_to_server(
        PandaAgentServicer(server, agent), server)

    for port in PORTS:
        server.add_insecure_port(port)
        print(f'panda agent grpc server listening on port {port}')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()