from pandare import Panda
from enum import Enum
import queue
import pb.panda_agent_pb2 as pb
import socket

PANDA_IP = '127.0.0.1'

# Each enum represents the state panda was in that caused the exception
# Enum should match that in /cmd/panda_test_executor/panda_executor_test.go
class ErrorCode(Enum):
    RUNNING = 0
    NOT_RUNNING = 1
    RECORDING = 2
    NOT_RECORDING = 3
    REPLAYING = 4
    NOT_REPLAYING = 5

class PandaAgent:
    # sentinel object to signal end of queue
    STOP_PANDA = object()

    def __init__(self, panda: Panda):
        self.panda = panda
        self.current_recording = None
        self.serial_out = ""
    
    # This function is meant to run in a different thread
    def start(self):
        panda = self.panda
        if self.panda.started.is_set():
            raise RuntimeError(ErrorCode.RUNNING.value, "Cannot start another instance of PANDA while one is already running")

        @panda.queue_blocking
        def panda_start():
            print("panda agent started")
            
            # revert to the qcow's root snapshot
            message = panda.revert_sync("root")
            if message != "":
                raise RuntimeError(ErrorCode.RUNNING, message)
        
        print("starting panda agent ")
        panda.run()
        print("panda agent stopped")
    
    def stop(self):
        if self.panda.started.is_set() is False: 
            raise RuntimeError(ErrorCode.NOT_RUNNING.value, "Cannot stop a PANDA instance when one is not running")
        if self.current_recording is not None:
            self.panda.end_analysis()
            self.current_recording = None
            raise RuntimeWarning(ErrorCode.RECORDING.value, "Request for PANDA stop before recording ended")    
        
        @self.panda.queue_blocking
        def panda_stop():
            self.panda.end_analysis()
    
    def _run_function(self, func, block=True, timeout=None):
        # Since the queued function will be running in another thread, we need
        # a queue in order to pass the return value back to this thread
        if self.panda.running.is_set() is False: 
            raise RuntimeError(ErrorCode.NOT_RUNNING.value, "Can't run a function when PANDA isn't running")
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
            raise RuntimeError(ErrorCode.RECORDING.value, "Cannot start new recording while recording in progress")

        self.current_recording = recording_name
        def panda_start_recording(panda: Panda):
            print(f'starting recording {recording_name}')
            try:
                return panda.record(recording_name)
            except Exception as err:
                raise RuntimeError(ErrorCode.RECORDING.value, f"Unexpected {err=}, {type(err)=}")
        
        return self._run_function(panda_start_recording)
    
    def stop_recording(self):
        if self.current_recording is None:
            raise RuntimeError(ErrorCode.NOT_RECORDING.value, "Must start a recording before stopping one")
        
        def panda_stop_recording(panda: Panda):
            print(f'stopping recording')
            try:
                panda.end_record()
            except Exception as err:
                raise RuntimeError(ErrorCode.RECORDING.value, f"Unexpected {err=}, {type(err)=}")
        
        recording_name = self.current_recording
        self._run_function(panda_stop_recording)
        self.current_recording = None
        return recording_name

    def start_replay(self, recording_name):
        panda = self.panda
        # Replay runs its own PANDA instance so PANDA should not be running beforehand
        if self.panda.started.is_set(): 
            raise RuntimeError(ErrorCode.RUNNING.value, "Cannot start another instance of PANDA while one is already running")
        if panda.recording_exists(recording_name) is False:
            raise RuntimeError(ErrorCode.REPLAYING.value, f"Recording {recording_name} does not exist")
        # For user to see serial output. A way to remember what the recording did
        @panda.cb_replay_serial_write
        def serial_append(env, fifo_addr, part_addr, value):
            # Append serial to string as character
            self.serial_out += '%c' % value
        print("starting panda replay agent ")
        print(f'starting replay {recording_name}')
        panda.run_replay(recording_name)
        print("panda agent replay stopped")
        return self.serial_out

    def stop_replay(self):
        if self.panda._in_replay is False:
            raise RuntimeError(ErrorCode.NOT_REPLAYING.value, "Must start a replay before stopping one")

        def panda_stop_replay(panda: Panda):
            print('stopping replay')
            try:
                panda.end_replay()
            except Exception as err:
                raise RuntimeError(ErrorCode.REPLAYING.value, f"Unexpected {err=}, {type(err)=}")
        
        self._run_function(panda_stop_replay)
        return self.serial_out

    def execute_network_command(self, request: pb.NetworkRequest):
        
        response = b''
        message = b''

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
            message += bytes(' ', encoding='utf-8')
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