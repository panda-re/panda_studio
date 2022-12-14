import grpc
import concurrent.futures as futures
import time
import threading
import queue
import multiprocessing

from pandare import Panda

import panda_interface_pb2_grpc as pb_grpc
import panda_interface_pb2 as pb

PORT = '[::]:50051'

executor = futures.ThreadPoolExecutor(max_workers=10)

class PandaAgent:
    def __init__(self, panda: Panda):
        self.panda = panda
    
    # This function is meant to run in a different thread
    def start(self):
        panda = self.panda

        agentPipe, pandaPipe = multiprocessing.Pipe(duplex=True)

        @panda.queue_blocking
        def process_commands():
            # revert to the qcow's root snapshot
            panda.revert_sync("root")
        
        panda.run()
    
    def run_command(cmd):
        # Send message to queue

        # Receive message from queue
        pass

class PandaExecutorServicer(pb_grpc.PandaExecutorServicer):
    def __init__(self, agent):
        self.agent = agent
        self.agent_thread = None
    
    def isAgentRunning(self):
        return self.agent_thread is not None
    
    def BootMachine(self, request, context):
        print("starting panda")
        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        print("started panda!")
        yield pb.BootMachineReply()
        return
    
    def RunCommand(self, request, context):
        print(request.command)
        return pb.RunCommandReply(statusCode=32)
        # return super().RunCommand(request, context)

def serve():
    panda = Panda(generic='x86_64')
    agent = PandaAgent(panda)
    server = grpc.server(executor)
    pb_grpc.add_PandaExecutorServicer_to_server(
        PandaExecutorServicer(agent), server)
    server.add_insecure_port(PORT)
    print(f'panda agent grpc server listening on port {PORT}')
    server.start()
    server.wait_for_termination()

def runPanda():
    panda = Panda(generic='i386')

    @panda.queue_blocking
    def run_cmd():
        # First revert to the qcow's root snapshot (synchronously)
        panda.revert_sync("root")
        # Then type a command via the serial port and print its results
        print(panda.run_serial_cmd("uname -a"))
        # When the command finishes, terminate the panda.run() call
        panda.end_analysis()
    
    panda.run()

if __name__ == "__main__":
    serve()
    # runPanda()