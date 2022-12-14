import grpc
import concurrent.futures as futures

from pandare import Panda

# protobuf compiler is stupid
import sys
sys.path.append(".")
sys.path.append("pb")

import pb.panda_agent_pb2_grpc as pb_grpc
import pb.panda_agent_pb2 as pb
from agent import PandaAgent

PORT = '[::]:50051'

executor = futures.ThreadPoolExecutor(max_workers=10)

class PandaAgentServicer(pb_grpc.PandaAgentServicer):
    def __init__(self, server, agent: PandaAgent):
        self.server = server
        self.agent = agent
        self.agent_thread = None
    
    def isAgentRunning(self):
        return self.agent_thread is not None
    
    def StartAgent(self, request, context):
        if self.agent.hasStarted:
            return

        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        return pb.StartAgentResponse()
    
    def StopAgent(self, request, context):
        self.agent.stop()
        self.server.stop(grace=5)
        return pb.StopAgentResponse()
    
    def RunCommand(self, request, context):
        output = self.agent.run_command(request.command)
        return pb.RunCommandReply(statusCode=0, output=output)

def serve():
    panda = Panda(generic='x86_64')
    agent = PandaAgent(panda)
    server = grpc.server(executor)
    pb_grpc.add_PandaAgentServicer_to_server(
        PandaAgentServicer(server, agent), server)
    server.add_insecure_port(PORT)
    print(f'panda agent grpc server listening on port {PORT}')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()