# Test Plan

The following document details the test plan for this project. It is comprised of two halves, each of which were created by the two teams that were responsible for different aspects of this project.

Please refer to the [readme](../readme.md) for proper installation instructions on PANDA Studio. This should include Ubuntu, Docker, and more.

## Frontend

TODO

## Backend

The backend must execute interactions on a specified image and return a recording to the user to satisfy the desired functionality. The following process details how that is verified.

### Testing Functionality

1. Creating an interactions list

The test for this requirement is to specify a list of interactions. The [executor](../cmd/panda_executor/panda_executor.go) already contains an interaction list that uses all implemented interactions. Currently, it runs the following serial commands:
- uname -a
- ls /
- touch NEW_FILE.txt
- ls /

The PANDA Executor acts as a frontend for the system so the backend of PANDA Studio can be tested without the need for the frontend/web UI.

2. Run the interactions list

The following command will start the backend of PANDA Studio and execute the interactions in the executor using an x86_64 architecture. Two Docker containers should run, the executor, which is initiated in the command below, and an agent that will execute PANDA and return the output to the executor.

Run the command in the panda_studio directory in the Ubuntu terminal:
```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor
```

3. Wait for the program to finish executing and check for error output in the terminal from the script

After running the script, the execution will look like this in the console that it was run in:
```
Starting agent
Starting recording
> uname -a
Linux ubuntu 4.15.0-72-generic #81-Ubuntu SMP Tue Nov 26 12:20:02 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux
> ls /
bin   home            lib64       opt   sbin  tmp      vmlinuz.old
boot  initrd.img      lost+found  proc  snap  usr
dev   initrd.img.old  media       root  srv   var
etc   lib             mnt         run   sys   vmlinuz
> touch /NEW_FILE.txt
> ls /
NEW_FILE.txt  etc             lib         mnt   run   sys  vmlinuz
bin           home            lib64       opt   sbin  tmp  vmlinuz.old
boot          initrd.img      lost+found  proc  snap  usr
dev           initrd.img.old  media       root  srv   var
Stopping recording
Snapshot file: test-rr-snp
Nondet log file: test-rr-nondet.log

Starting replay
ls /
bin   home            lib64       opt   sbin  tmp      vmlinuz.old
boot  initrd.img      lost+found  proc  snap  usr
dev   initrd.img.old  media       root  srv   var
etc   lib             mnt         run   sys   vmlinuz
root@ubuntu:~# ls /
NEW_FILE.txt  etc             lib         mnt   run   sys  vmlinuz
bin           home            lib64       opt   sbin  tmp  vmlinuz.old
boot          initrd.img      lost+found  proc  snap  usr
dev           initrd.img.old  media       root  srv   var
root@ubuntu:~# 
loading snapshot
... done.
opening nondet log for read :   data/test-rr-nondet.log
data/test-rr-nondet.log:     6904677 instrs total.
Replay completed successfully
Exiting cpu_handle_execption loop
```

Each line with a `>` represents the serial command being run in PANDA with the following line being the output of the command. The output of the command should match what is expected and if the same command were run manually in a Docker container. For example, the above commands run `ls`, make a new file, then `ls` again and the new file appears which matches the expected outcome since the file `NEWFILE.txt` appears after running `ls` for the second time. Also, the executor will perform a replay using the recording that was just generated which will print the serial output that was recorded along with information about the replaying of the recording.

4. Check for the recording

The recording can be shown by running the following command:
```
ls /tmp/panda-studio
```
Both the files `test-rr-nondet.log` and `test-rr-snp` should be there. The files prove a recording was returned.

### Testing Robustness

There are potential ways to execute in an incorrect order such as starting a recording while one is in progress. Therefore, the program should prevent incorrect execution. This test checks those conditions on top of verifying normal execution. Without the two requirements being fulfilled, it will be hard to know if the system ran correctly or will run correctly from one successful test.

1. Setup

PANDA Studio has a test prepared to ensure the backend catches potential problems and propagates the errors up. More detailed information about the tests can be found in its corresponding [readme](../cmd/panda_executor_test/README.md).

Run the following command in the `panda_studio` directory:
```
make test
```

The command will build the necessary parts such as Docker containers to run the tests.

2. Run the tests

Run the following command:
```
docker run -it --rm \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -v /root/.panda:/root/.panda \
     -v /tmp/panda-studio:/tmp/panda-studio \
     pandare/panda_executor_test
```

3. Wait for it to finish executing

During execution, the test will print out which test is being run. If a test fails, it is accompanied by an error message. 

The output should look like:
```
=== RUN   TestMain
=== RUN   TestMain/Agent
=== RUN   TestMain/Agent/PreCommand
=== RUN   TestMain/Agent/PreStop
=== RUN   TestMain/Agent/ExtraStart
=== RUN   TestMain/Agent/Commands
=== RUN   TestMain/Agent/RunExtraReplay
=== RUN   TestMain/Recording
=== RUN   TestMain/Recording/PreStop
=== RUN   TestMain/Recording/StartRecording
=== RUN   TestMain/Recording/ExtraStart
=== RUN   TestMain/Recording/Commands
=== RUN   TestMain/Recording/StopRecording
=== RUN   TestMain/Recording/HangingStop
=== RUN   TestMain/Replay
=== RUN   TestMain/Replay/PreStop
=== RUN   TestMain/Replay/WrongReplay
=== RUN   TestMain/Replay/RunReplay
=== RUN   TestMain/Replay/RunExtraReplay
Number of tests: 15
Number passed: 15
Success rate: 100%
--- PASS: TestMain (69.67s)
    --- PASS: TestMain/Agent (18.60s)
        --- PASS: TestMain/Agent/PreCommand (0.01s)
        --- PASS: TestMain/Agent/PreStop (0.00s)
        --- PASS: TestMain/Agent/ExtraStart (0.00s)
        --- PASS: TestMain/Agent/Commands (1.88s)
        --- PASS: TestMain/Agent/RunExtraReplay (0.00s)
    --- PASS: TestMain/Recording (27.12s)
        --- PASS: TestMain/Recording/PreStop (0.00s)
        --- PASS: TestMain/Recording/StartRecording (1.20s)
        --- PASS: TestMain/Recording/ExtraStart (0.00s)
        --- PASS: TestMain/Recording/Commands (1.01s)
        --- PASS: TestMain/Recording/StopRecording (2.33s)
        --- PASS: TestMain/Recording/HangingStop (0.20s)
    --- PASS: TestMain/Replay (23.94s)
        --- PASS: TestMain/Replay/PreStop (0.00s)
        --- PASS: TestMain/Replay/WrongReplay (0.26s)
        --- PASS: TestMain/Replay/RunReplay (0.98s)
        --- PASS: TestMain/Replay/RunExtraReplay (0.00s)
PASS
```
Once all of the tests have run it will print out the number of tests, number of tests passed, the success rate, and each test along with whether it passed or failed. If the testing program errors out and does not complete, then the test fails. If a test case does not pass, then PANDA Studio does not function properly (or the test is not correct).
