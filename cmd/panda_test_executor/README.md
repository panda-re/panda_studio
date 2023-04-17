# PANDA Agent Test
The executor will test the PANDA agent to ensure proper operation of the backend. It both creates scenarios that will cause exceptions and verifies normal operation.

## Running Tests
Build the agent Docker container that will be tested using the makefile in the `panda_studio` directory:

```
make docker_agent
```

The test requires super user privileges and having golang installed, so there are two options for execution.
### In Docker
1.  Build the test executor with:

```
make docker_executor_test
```

2.  Run the test executor:

```
docker run -it --rm \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -v /root/.panda:/root/.panda \
     -v /tmp/panda-studio:/tmp/panda-studio \
     pandare/panda_test_executor
```

### In Terminal
1.  Install golang:
```
sudo apt install golang-go
```
2.  Run the test:
```
sudo go test -v -run TestMain ./cmd/panda_test_executor/
```
### Output
The terminal will print whenever a test starts and whether it passes or fails. Also, the number of tests run, the number of tests passed, and the success rate will be printed.
### Errors
The agent throws exceptions with a code and a message. The agent will test scenarios that will cause exceptions to be thrown and ensures the correct error code is returned. The codes are associated with the state PANDA is in and are described below:
| Name | Number | Description |
| ----- | ------ | ------- |
| RUNNING | 0 | PANDA is running, or an exception occurred |
| NOT_RUNNING | 1 | An attempt to use PANDA when it is not running |
| RECORDING | 2 | PANDA is recording, or an exception occurred while recording |
| NOT_RECORDING | 3 | PANDA is not recording |
| REPLAYING | 4 | PANDA is replaying, and an exception occurred while replaying |
| NOT_REPLAYING | 5 | PANDA is not replaying |

The enumeration in both the test executor and agent must match for the tests to work.

## Tests
The following section will explain various aspects of the tests so future developers can add tests to verify the functionality of new features and ensure previous features still function.

The test uses the following hierarchy:
* `Main` - starts agent tests
    * `Agent Test` – Focuses on a specific subset of functions of the agent
        * `Subtest` - Tests an edge case or function of the agent

Subtests exist outside of agent tests to reduce code duplication. However, using `go test` will run all tests, so use the `-run` flag to run all tests (`TestMain`) or a specific agent test (e.g. `TestAgent`, `TestRecording`, `TestReplay`).
### Agent Tests
Agent tests focus on a specific subset of agent functionality. For example, `TestAgent` tests for basic agent exceptions such as premature execution and I/O functionality. Each agent test should create, start, stop, and close its own agent. Stop and close are generally put in `t.Cleanup` because the cleanup function always runs last so if a test causes a fatal error the agent is stopped regardless.

### Subtests
There are two types of subtests: exception and functionality.

Exception tests will try to cause the agent to throw an exception. For example, `TestExtraStart` attempts to start the agent's PANDA instance again. Without an exception, this would cause errors, hence the test. The subtest ensures the error code matches what is expected, otherwise the test will fail.

Functionality tests will verify the agent is working. For example, `TestCommands` runs serial commands and checks whether the response matches what is expected. Whenever the agent returns something, the output should be verified. If the function fails or the output is incorrect the test will fail. 