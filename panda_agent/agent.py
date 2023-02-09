from pandare import Panda
import queue
import pb.panda_agent_pb2 as pb
import socket

PANDA_IP = '127.0.0.1'

class PandaAgent:
    # sentinel object to signal end of queue
    STOP_PANDA = object()

    def __init__(self, panda: Panda):
        self.panda = panda
        # TODO PANDA has running and replay variables already, use those
        self.isRunning = False
        self.current_recording = None
        self.current_replay = None
        self.serial_out = ""
    
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
        # Replay runs its own PANDA instance so PANDA should not be running beforehand
        if self.isRunning: 
            raise RuntimeError("Cannot start another instance of PANDA while one is already running")
        panda = self.panda
        if panda.recording_exists(recording_name) is False:
            raise RuntimeError(f"Recording {recording_name} does not exist")
        # For user to see serial output. A way to remember what the recording did
        @panda.cb_replay_serial_write
        def serial_append(env, fifo_addr, part_addr, value):
            # Append serial to string as character
            self.serial_out += '%c' % value
        print("starting panda replay agent ")
        self.isRunning = True
        self.current_replay = recording_name
        print(f'starting replay {recording_name}')
        panda.run_replay(recording_name)
        self.current_replay = None
        print("panda agent replay stopped")
        # self.isRunning = False
        return self.serial_out

    def stop_replay(self):
        if self.current_replay is None:
            raise RuntimeError("Must start a replay before stopping one")

        def panda_stop_replay(panda: Panda):
            print('stopping replay')
            panda.end_replay()
        
        self.current_replay = None
        self._run_function(panda_stop_replay)
        return self.serial_out

    def execute_network_command(self, request: pb.NetworkRequest):
        
        response = []
        message = bytearray[b'']

        if request.application == "HTTP" or request.application == "http":
            if request.command == "GET":
                message = b'GET'
            elif request.command == "HEAD":
                message = b'HEAD'
            elif request.command == "POST":
                message = b'POST'
            elif request.command == "PUT":
                message = b'PUT'
            elif request.command == "DELETE":
                message = b'DELETE'
            elif request.command == "CONNECT":
                message = b'CONNECT'
            elif request.command == "OPTIONS":
                message = b'OPTIONS'
            elif request.command == "TRACE":
                message = b'TRACE'
            elif request.command == "PATCH":
                message = b'PATCH'
            else:
                return "Invalid HTTP Command, make sure its somehting like GET, POST, or PUT"
            
            # Add a space then add the custome message (may be '' which is fine as long as its something)
            message += bytes(' ')
            message += bytes(request.customPacket, encoding='utf-8')
            message += b' HTTP/1.1\r\n\r\n'

        else:
            message = bytes(request.customPacket, encoding='utf-8')





        with socket.socket(socket.AF_INET, request.socketType) as con:
            con.connect((PANDA_IP, request.port))
            con.send(message)
            response = con.recv(4096)
            con.close()
        
        return response