# Table of Contents
* [Introduction](#introduction)
* [Backend Testing](#back-end-testing)
  * [Testing Backend Functionality](#testing-functionality)
  * [Testing Backend Robustness](#testing-robustness)
* [API Testing](#api-testing)
* [Frontend Testing](#front-end-testing)
* [Full Integration Testing](#full-integration-testing)
  * [Test Cases](#test-cases)  
    * [Test Case: Recording Creation](#test-case-recording-creation-f1-f5)
    * [Test Case: Delete Recording](#test-case-delete-recording-f4)
    * [Test Case: Image Upload by File](#test-case-image-upload-by-file-f7)
    * [Test Case: Image Upload by Url](#test-case-image-upload-by-url-f7)
    * [Test Case: Derive Image](#test-case-derive-image-f2-f7)
    * [Test Case: Create Interaction Program](#test-case-create-interaction-program-f3-f11)
    * [Test Case: Download recording File](#test-case-download-recording-file-f9)
* [Future Automated API Unit Testing Ideas](#future-automated-api-unit-testing-ideas)


# Introduction:

Testing the PANDA Studio project will involve focused areas of testing in the front-end, API, and back-end. By separating our testing efforts into these three distinct areas, this allows for us to validate each layer separately and to different degrees of certainty. Currently, a majority of our testing capacity involves manually unit testing each piece of each layer with different strategies used to test each layer. In the future, an automated test suite would be ideal and increase overall test effectiveness.

Below are sections describing how to manually unit test each individual layer followed by full system integration test cases to manually test system wide functionality.

# Backend testing:
To manually test the backend you must execute interactions on a specified image and return a recording to the user to satisfy the desired functionality. The following process details how that is manually verified.

## Testing Functionality

1. Creating an interactions list

The test for this requirement is to specify a list of interactions. The [executor](https://github.com/panda-re/panda_studio/blob/backend-cleanup-documentation/cmd/panda_executor/panda_executor.go#L34-L39) already contains an interaction list that uses all implemented interactions. Currently, it runs the following serial commands:
- uname -a
- ls /
- touch NEW_FILE.txt
- ls /

The PANDA Executor acts as a frontend for the system so the backend of PANDA Studio can be tested without the need for the frontend/web UI.

2. Run the interactions list

The following command will start the backend of PANDA Studio and execute the interactions in the executor using an emulated x86_64 architecture. Two Docker containers should run, the executor, which is initiated in the command below, and an agent that will execute PANDA and return the output to the executor.

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
Both the files `test-rr-nondet.log` and `test-rr-snp` should be there. The files' exsitence proves a recording was returned. In order to test that this recording is valid, a user can modify the panda_executor script to call the replay gRPC request instead of using a set list of interactions. An example of this call can be found in the [TestRunReplay](https://github.com/panda-re/panda_studio/blob/backend-cleanup-documentation/cmd/panda_executor_test/panda_executor_test.go#L359) function in panda_executor_test. If all you want to test is if the system itself returns valid recordings, then following this [procedure](https://github.com/panda-re/panda_studio/blob/backend-cleanup-documentation/cmd/panda_executor_test/README.md) will indicate if the system is able to run a replay from a recording it created.

## Testing Robustness

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


# API testing:

API testing mainly revolves around unit testing the layer of application logic the endpoints call into. Using OpenAPI to generate the frontend endpoints allows us to unit test the application logic separately from the frontend access of the endpoints. Unlike the UI, unit testing the API is the best course of action as it provides us greater assurance that our code works as expected.

* Unit Testing
    * The API will primarily be tested through unit tests. Using OpenAPI to generate endpoints has sufficiently decoupled the application logic in such a way that enables unit testing of specific methods with mocked data. Integration tests may be written before unit tests to ensure that the system is working as expected on a larger scale.
    * Unit testing is currently completed manually by importing the panda_studio.yaml file into Postman or another program capable of sending network requests. After importing the yaml file into postman, each request's application logic can be tested separately. This mainly includes CRUD operations for all objects in the system.
* Integration Testing
    * Integration testing between the frontend and API layer logic involves using the frontend to perform tasks that accomplish various features and manually ensuring the proper functionality occurs. This testing strategy is discussed in further detail below.
* Acceptance Testing
    * The API will receive acceptance testing at the same time as the UI.


# Front-end testing:

Since this project is open-source, we are going to defer writing unit tests for UI pages and components. There are several reasons for this. First, UI tests generally only verify that components are rendered as expected and/or that certain actions happen when they are triggered in the UI. Both of these aspects can be verified manually, so there is no point in spending time on UI unit tests when there is more important work to be done. Additionally, with the project source code available to everyone in the PANDA community, others will be able to revisit this project and produce front-end unit tests as they see fit or point out bugs that may need to be resolved.

If there is time left at the end of the project, we may revisit front-end testing and produce unit tests for areas we feel are applicable, but it will not be a high priority. It is important to note here that the ECE team’s test plan and validation strategy is effectively all done by demonstration. This means that they will simply show the various features they plan to implement, running them from the web client, and manually verifying that the results are as expected. Our front-end testing strategy is effectively the same thing, as we plan to verify our results in the UI manually.

* Unit Testing
    * As mentioned above, the UI will not be unit tested at this time.
* Integration Testing
    * Integration testing with the API layer is covered with the feature based testing described below.
* Acceptance Testing
    * The UI will largely be tested through acceptance tests. This will involve performing the various actions defined in the use cases and verifying that the preconditions and postconditions are met. If the UI performs as expected in the use cases, the test will be considered a pass. Otherwise the test will be considered a failure.

# Full Integration Testing
## Test Cases:
### Test Case: Recording Creation
1. Navigate to the Recording Dashboard by selecting the Recordings tab on the far left
2. Select Create a New Recording
3. Specify the VM image to be used by doing one of the following:
    * Scroll through the options for the desired image and select
    * Search for a particular image and select
4. Specify the interaction program by doing one of the following:
    * Scroll through the options for the desired interaction program and select
    * Search for a particular interaction program and select
5. Specify a recording name
6. Click Create Recording
7. Verify that the new recording appears on the Recording dashboard

### Test Case: Delete Recording
1. Navigate to the Recording Dashboard
2. Select a recording
3. Click Delete Recording
4. Verify that the user has been navigated back to the Recording Dashboard and the recording has been deleted

### Test Case: Image Upload by File
5. Navigate to the Image Dashboard
6. Click the upload image button
7. Enter a name for the image
8. Click the x86_64 button to load all default config values
9. Delete url value that was entered by button logic
10. Select local image file to upload
11. Click create
12. Verify that the user has been navigated back to the Image Dashboard and the image has been created

### Test Case: Image Upload by Url
12. Navigate to the Image Dashboard
13. Click the upload image button
14. Enter a name for the image
15. Click the x86_64 button to load all default config values
17. Click create
18. Verify that the user has been navigated back to the Image Dashboard and the image has been created

### Test Case: Derive Image
1. Select the Image Details option on an image on the Image Dashboard
2. Select the Derive Image option
3. Fill out information for:
    1. Custom Image Name
    2. Docker Image name (reference for Docker hub)
    3. Amount of memory to allocate
    4. Extra args for PANDA (Optional)
4. Select the submit option
5. Verify a new image has been specified on the Image Dashboard

### Test Case: Create Interaction Program
1. Select the Program Dashboard
2. Select Create a New Interaction Program
3. Fill out the name field
4. Provide some interaction program, such as the following:
    ``` 
    START_RECORDING new_rec
    CMD ls /
    CMD touch /newFile123.txt
    CMD ls /
    STOP_RECORDING
    ```
5. Click the create button
6. Verify that the new interaction program shows on the dashboard with the correct information

### Test Case: Download recording File
1. Select a recording from the recording dashboard
2. Download the recording files
3. Verify the non-determinism log and recording output are as expected


# Future Automated API Unit Testing Ideas

As our API is written using the Go language, we have determined the following Go unit testing frameworks may be the best solutions for automated unit testing of the API layer:

* Go testing package
    * The advantage of the default Go testing package is that it is built into the Go runtime and does not involve installing any other software. However, it is relatively limited in its capabilities, and does not have the same level of assertion capabilities as do other more powerful frameworks.
* [httpexpect](https://github.com/gavv/httpexpect)
    * This is a REST API testing framework which would be advantageous for testing our API endpoints. It provides assertions for checking request and response payloads, which is something we will likely want to validate in our system. 
* [Testify](https://github.com/stretchr/testify)
    * Testify is another unit testing framework which provides more powerful assertions and mocking capabilities. We may determine the need for mocks when testing our system, so using this framework would be advantageous (as opposed to going with the default Go package). This framework could be applied to our lower-level application logic for verifying correctness. The main drawback is that it does not provide its own test coverage metrics, but these can be obtained elsewhere.
* [Mockery](https://github.com/vektra/mockery)
* [Gomock](https://github.com/golang/mock)

Using an automated test framework to build a unit testing suite would help with future project or open source work extending the capabilities of the Panda Studio application by ensuring functionality and increasing testing efficiency.
