from pandare import Panda
import queue

class PandaAgent:
    # sentinel object to signal end of queue
    STOP_PANDA = object()

    def __init__(self, panda: Panda):
        self.panda = panda
        self.isRunning = False
        self.current_recording = None
        self.current_replay = None
    
    # This function is meant to run in a different thread
    def start(self):
        if self.isRunning: 
            raise RuntimeError("Cannot start another instance of PANDA while one is already running")
        panda = self.panda

        @panda.queue_blocking
        def panda_start():
            print("panda agent started")
            
            # revert to the qcow's root snapshot
            panda.revert_sync("root")
        
        print("starting panda agent ")
        self.isRunning = True
        panda.run()
        print("panda agent stopped")
        #self.isRunning = False
    
    def stop(self):
        if self.isRunning is False: 
            raise RuntimeError("Cannot stop a PANDA instance when one is not running")
        if self.current_recording is not None:
            # Stop recording for user, then stop PANDA
            print("Request for PANDA stop before recording ended. Stopping recording automatically")
            self.stop_recording()
        
        @self.panda.queue_blocking
        def panda_stop():
            self.panda.end_analysis()
        self.isRunning = False
    
    def _run_function(self, func, block=True, timeout=None):
        # Since the queued function will be running in another thread, we need
        # a queue in order to pass the return value back to this thread
        if self.isRunning is False: 
            raise RuntimeError("Can't run a function when PANDA isn't running")
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
    
    def start_recording(self, recording_name):
        if self.current_recording is not None:
            raise RuntimeError("Cannot start new recording while recording in progress")

        self.current_recording = recording_name
        def panda_start_recording(panda: Panda):
            print(f'starting recording {recording_name}')
            return panda.record(recording_name)
        
        return self._run_function(panda_start_recording)
    
    def stop_recording(self):
        if self.current_recording is None:
            raise RuntimeError("Must start a recording before stopping one")
        
        def panda_stop_recording(panda: Panda):
            print(f'stopping recording')
            panda.end_record()
        
        recording_name = self.current_recording
        self._run_function(panda_stop_recording)
        self.current_recording = None
        return recording_name

    def start_replay(self, recording_name):
        # TODO checks (recording?)
        if self.current_replay is not None:
            raise RuntimeError("Cannot start a new replay while replay is in progress")
        
        def panda_start_replay(panda: Panda):
            print(f'starting replay {recording_name}')
            panda.run_replay(recording_name)
        
        self.current_replay = recording_name
        return self._run_function(panda_start_replay)

    def stop_replay(self):
        # TODO checks
        if self.current_replay is None:
            raise RuntimeError("Must start a replay before stopping one")

        def panda_stop_replay(panda: Panda):
            print('stopping replay')
            panda.end_replay()
        
        self.current_replay = None
        return self._run_function(panda_stop_replay)
        

    def stop_replay(self):
        pass