from pandare import Panda
from enum import Enum
import queue
import pb.panda_agent_pb2 as pb
import socket

PANDA_IP = '127.0.0.1'

# Each enum represents the state panda was in that caused the exception
# Enum should match that in /cmd/panda_executor_test/panda_executor_test.go
class ErrorCode(Enum):
    RUNNING = 0
    NOT_RUNNING = 1
    RECORDING = 2
    NOT_RECORDING = 3
    REPLAYING = 4
    NOT_REPLAYING = 5

class PandaAgent:
    '''
    Agent that translates interaction requests into PANDA commands.
    
    Uses PyPANDA to produce desired results.
    '''
    # sentinel object to signal end of queue
    STOP_PANDA = object()
    FILE_PREFIX="data"

    def __init__(self, config: pb.PandaConfig):
        self.config = config
        self.panda = self.init_panda(config)
        self.current_recording = None
        self.serial_out = ""
    
    def init_panda(self, config: pb.PandaConfig):
        '''
        Construct PANDA instance based on the config

        Args:
            config: PANDA configuration parameters

        Returns:
            Panda: the created PANDA object
        '''
        return Panda(
            arch=config.arch,
            qcow=f"{self.FILE_PREFIX}/{config.qcowfilename}",
            mem=config.memory,
            os=config.os,
            expect_prompt=config.prompt,
            extra_args=config.extraargs)
    
    # This function is meant to run in a different thread
    def start(self):
        '''
        Starts the PANDA object and sets its snapshot.

        Supposed to run in a separate thread because QEMU is blocking.

        Returns:
            None

        Raises:
            RuntimeError if PANDA was already recording
        '''
        panda = self.panda
        if self.panda.started.is_set():
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")

        @panda.queue_blocking
        def panda_start():
            print("panda agent started with config:")
            # print config as struct
            print(self.config)
            
            # revert to the qcow's root snapshot
            message = panda.revert_sync(self.config.snapshot)
            if message != "":
                raise RuntimeError(ErrorCode.RUNNING, message)
        
        print("starting panda agent")
        panda.run()
        print("panda agent stopped")
    
    def stop(self):
        '''
        Stops the PANDA object.

        Returns:
            None

        Raises:
            RuntimeError if PANDA was not running
            RuntimeWarning if PANDA was not recording
        '''
        if self.panda.started.is_set() is False: 
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Cannot stop a PANDA instance when one is not running")
        if self.current_recording is not None:
            self.panda.end_analysis()
            self.current_recording = None
            raise RuntimeWarning(ErrorCode.RECORDING, "Request for PANDA stop before recording ended")    
        
        @self.panda.queue_blocking
        def panda_stop():
            self.panda.end_analysis()
    
    def _run_function(self, func, block=True, timeout=None):
        '''
        Wrapper function that utilizes PANDA queue to run commands and pass return value back to thread

        Args:
            func: function to run

        Returns:
            Output of the function
        
        Raises:
            RuntimeError if PANDA was not running
        '''
        # Since the queued function will be running in another thread, we need
        # a queue in order to pass the return value back to this thread
        if self.panda.running.is_set() is False: 
            raise RuntimeError(ErrorCode.NOT_RUNNING, "Can't run a function when PANDA isn't running")
        returnChannel = queue.Queue()

        @self.panda.queue_blocking
        def panda_queue_wrapper():
            returnChannel.put(func(self.panda))

        return returnChannel.get(block=block, timeout=timeout)
    
    def run_command(self, cmd):
        '''
        Runs a serial command in PANDA

        Args:
            cmd: serial command to be run
        
        Returns:
            String: all the output (stdout + stderr)

        Raises:
            RuntimeError if PANDA was not running
        '''
        def panda_run_command(panda: Panda):
            print(f'running command {cmd}')
            try:
                return panda.run_serial_cmd(cmd)
            except Exception as err:
                raise RuntimeError(ErrorCode.RUNNING, f"Unexpected {err=}, {type(err)=}")

        return self._run_function(panda_run_command)
    
    def revert_to_snapshot(self, snapshot):
        '''
        Reverts PANDA to a snapshot

        Args:
            snapshot: name of snapshot in the current qcow to load

        Returns:
            String: error message. Empty on success.
        '''
        def panda_revert_snapshot(panda: Panda):
            print(f'reverting to snapshot {snapshot}')
            return panda.revert_sync(snapshot)
        
        return self._run_function(panda_revert_snapshot)
    
    def start_recording(self, recording_name):
        '''
        Starts a PANDA recording

        Args:
            recording_name: name of recording to save

        Returns:
            gRPC response acknowledging the recording started
        
        Raises:
            RuntimeError if PANDA was not running
            RuntimeError if PANDA was already recording
        '''
        if self.current_recording is not None:
            raise RuntimeError(ErrorCode.RECORDING, "Cannot start new recording while recording in progress")

        self.current_recording = recording_name
        def panda_start_recording(panda: Panda):
            print(f'starting recording {recording_name}')
            try:
                return panda.record(f"{self.FILE_PREFIX}/{recording_name}")
            except Exception as err:
                raise RuntimeError(ErrorCode.RECORDING, f"Unexpected {err=}, {type(err)=}")
        
        return self._run_function(panda_start_recording)
    
    def stop_recording(self):
        '''
        Stops a PANDA recording

        Returns:
            str: name of the recording

        Raises:
            RuntimeError if PANDA was not running
            RuntimeWarning if PANDA was not recording
        '''
        if self.current_recording is None:
            raise RuntimeError(ErrorCode.NOT_RECORDING, "Must start a recording before stopping one")
        
        def panda_stop_recording(panda: Panda):
            print(f'stopping recording')
            try:
                panda.end_record()
            except Exception as err:
                raise RuntimeError(ErrorCode.RECORDING, f"Unexpected {err=}, {type(err)=}")
        
        recording_name = self.current_recording
        self._run_function(panda_stop_recording)
        self.current_recording = None
        return recording_name

    def start_replay(self, recording_name):
        '''
        Starts a PANDA replay with the specified configuration.

        Only one PANDA instance per agent can be running at one time.

        Args:
            recording_name: Replay name/path

        Returns:
            str: Serial output of the replay

        Raises:
            RuntimeError if PANDA was already running
            RuntimeError if recording files to replay to not exist
        '''
        panda = self.panda
        # Replay runs its own PANDA instance so PANDA should not be running beforehand
        if self.panda.started.is_set(): 
            raise RuntimeError(ErrorCode.RUNNING, "Cannot start another instance of PANDA while one is already running")
        if panda.recording_exists(f"{self.FILE_PREFIX}/{recording_name}") is not True:
            raise RuntimeError(ErrorCode.REPLAYING, f"Recording {recording_name} does not exist")
        # For user to see serial output. A way to remember what the recording did
        @panda.cb_replay_serial_write
        def serial_append(env, fifo_addr, part_addr, value):
            # Append serial to string as character
            self.serial_out += '%c' % value
        print(f'starting replay {recording_name}')
        try:
            panda.run_replay(f"{self.FILE_PREFIX}/{recording_name}")
        except Exception as err:
            raise RuntimeError(ErrorCode.REPLAYING, f"Unexpected {err=}, {type(err)=}")
        print("panda agent replay stopped")
        return self.serial_out

    def stop_replay(self):
        '''
        Stops a PANDA replay.

        PANDA replays stop naturally.

        Returns:
            str: Serial output of the replay

        Raises:
            RuntimeError if PANDA was not replaying
        '''
        if self.panda._in_replay is False:
            raise RuntimeError(ErrorCode.NOT_REPLAYING, "Must start a replay before stopping one")

        def panda_stop_replay(panda: Panda):
            print('stopping replay')
            try:
                panda.end_replay()
            except Exception as err:
                raise RuntimeError(ErrorCode.REPLAYING, f"Unexpected {err=}, {type(err)=}")
        
        self._run_function(panda_stop_replay)
        return self.serial_out

    def execute_network_command(self, request: pb.NetworkRequest):
        '''
        Handles execution of a network interaction. It will first construct the message based on the request
        sent to the backend. Then it will create a socket on the requested port and PANDA_IP using the specified
        sock_type. Finally it will send the message. 

        Note: The PANDA machine has to already be listening for a connection on the requested port, otherwise
        the connection will fail. If you want to start the network interaction first or have it try multiple
        times, look at the Non-Blocking-Network branch on Github which starts the connection in another process
        and retries the connection if it fails. This does bring with it the loss of the response however.

        Args:
            request: gRPC request containing the networking information such as application, command, etc.
        
        Returns:
            bytes: Response message of the PANDA Machine
        '''
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