import grpc
import concurrent.futures as futures
import time
import threading
import queue
import multiprocessing

# important for pickling functions
import dill

from pandare import Panda

import panda_interface_pb2_grpc as pb_grpc
import panda_interface_pb2 as pb

PORT = '[::]:50051'

executor = futures.ThreadPoolExecutor(max_workers=10)

class PandaAgent:
    # sentinel object to signal end of queue
    STOP_PANDA = object()

    def __init__(self, panda: Panda):
        self.panda = panda
        # Queue your thunks you want executed in the PANDA loop here!
        self.cmdQueue = None
        self.hasStarted = False
        self.isRunning = False
    
    # This function is meant to run in a different thread
    def start(self):
        self.hasStarted = False
        panda = self.panda

        self.cmdQueue = queue.Queue()

        @panda.queue_blocking
        def process_commands():
            print("started panda!")
            self.isRunning = True
            # revert to the qcow's root snapshot
            panda.revert_sync("root")

            # receive and process command
            while True:
                cmd = self.cmdQueue.get(block=True)
                if cmd is PandaAgent.STOP_PANDA:
                    break
                
                cmd(panda)
            
            panda.end_analysis()
        
        print("starting panda")
        panda.run()
    
    def _queue_command(self, msg):
        if self.cmdQueue is None:
            raise RuntimeError("PANDA must be running to send commands")
        self.cmdQueue.put(msg)
    
    def run_command(self, cmd):
        retQueue = queue.Queue()

        def panda_run_command(panda: Panda):
            print(f'running command {cmd}')
            output = panda.run_serial_cmd(cmd)
            print(f'output: {output}')
            retQueue.put(output)
        
        # Send message to queue
        self._queue_command(panda_run_command)

        # Receive message from queue
        resp = retQueue.get(block=True, timeout=None)
        print(f'got response {resp}')
        return resp

class PandaExecutorServicer(pb_grpc.PandaExecutorServicer):
    def __init__(self, agent: PandaAgent):
        self.agent = agent
        self.agent_thread = None
    
    def isAgentRunning(self):
        return self.agent_thread is not None
    
    def BootMachine(self, request, context):
        if self.agent.hasStarted:
            return

        # start panda in a new thread, because qemu blocks this thread otherwise
        executor.submit(self.agent.start)
        yield pb.BootMachineReply()
        return
    
    def RunCommand(self, request, context):
        print(request.command)
        output = self.agent.run_command(request.command)
        return pb.RunCommandReply(statusCode=0, output=output)
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