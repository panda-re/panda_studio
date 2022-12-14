from pandare import Panda
import queue

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
            print("panda agent started")
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
        
        print("starting panda agent ")
        panda.run()
        print("panda agent stopped")
        self.isRunning = False
    
    def stop(self):
        self._queue_command(PandaAgent.STOP_PANDA)
    
    def _queue_command(self, msg):
        if self.cmdQueue is None:
            raise RuntimeError("PANDA must be running to send commands")
        self.cmdQueue.put(msg)
    
    def run_command(self, cmd):
        retQueue = queue.Queue()

        def panda_run_command(panda: Panda):
            print(f'running command {cmd}')
            output = panda.run_serial_cmd(cmd)
            retQueue.put(output)
        
        # Send message to queue
        self._queue_command(panda_run_command)

        # Receive message from queue
        resp = retQueue.get(block=True, timeout=None)
        return resp