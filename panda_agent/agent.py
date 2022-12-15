from pandare import Panda
import queue

class PandaAgent:
    # sentinel object to signal end of queue
    STOP_PANDA = object()

    def __init__(self, panda: Panda):
        self.panda = panda
        self.hasStarted = False
        self.isRunning = False
    
    # This function is meant to run in a different thread
    def start(self):
        self.hasStarted = False
        panda = self.panda

        @panda.queue_blocking
        def panda_start():
            print("panda agent started")
            self.isRunning = True
            # revert to the qcow's root snapshot
            panda.revert_sync("root")
        
        print("starting panda agent ")
        panda.run()
        print("panda agent stopped")
        self.isRunning = False
    
    def stop(self):
        @self.panda.queue_blocking
        def panda_stop():
            self.panda.end_analysis()
    
    def _run_function(self, func, block=True, timeout=None):
        # Since the queued function will be running in another thread, we need
        # a queue in order to pass the return value back to this thread
        returnChannel = queue.Queue()

        @self.panda.queue_blocking
        def panda_queue_wrapper():
            returnChannel.put(func(self.panda))

        return returnChannel.get(block=block, timeout=timeout)
    
    def run_command(self, cmd):
        def panda_run_command(panda: Panda):
            print(f'running command {cmd}')
            return panda.run_serial_cmd(cmd)

        return self._run_function(panda_run_command)
    
    def revert_to_snapshot(self, snapshot):
        def panda_revert_snapshot(panda: Panda):
            print(f'reverting to snapshot {snapshot}')
            return panda.revert_sync(snapshot)
        
        return self._run_function(panda_revert_snapshot)