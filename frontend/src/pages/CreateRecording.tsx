import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import React, { Component } from 'react'
import Terminal from 'react-console-emulator'

const terminal = React.createRef()
let commands = ""

function CreateRecordingPage() {
  const nameRef = useRef(null)
  const volumeRef = useRef(null)
  const programRef = useRef(null)

  return (<>
    <EuiPageTemplate.Header pageTitle="Create Recording"/>

    <EuiPageTemplate.Section>
      <EuiFlexGroup>
        <EuiFlexItem grow={2}>
          <EuiText>Name: </EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EuiFieldText
            placeholder="eg, recording1"
            inputRef={nameRef}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup>
        <EuiFlexItem grow={2}>
          <EuiText>Image: </EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EuiFieldText
            placeholder="eg, guest.img"
            inputRef={volumeRef}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={2}>
          <EuiText>Specify commands to run:</EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EuiFieldText
            placeholder="eg, uname -a, ls..."
            inputRef={programRef}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>


    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={false}>
          <div>
            <EuiButton onClick={() => sendAPICall(nameRef, volumeRef, programRef)}>Create Recording</EuiButton>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem>
          <Terminal
            ref={terminal}
            commands={{}}
            promptLabel={'panda@panda:~$'}
            errorText={"\n"}
            readOnly={true}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </>)
}


let sendAPICall = function (nameInput: MutableRefObject<any>, volumeInput: MutableRefObject<any>, programInput: MutableRefObject<any>) {
  commands = JSON.parse("[" + programInput.current.value + "]")
  let recordingDetails = {
    volume: volumeInput.current.value,
    commands: commands,
    name: nameInput.current.value
  }
  console.log(recordingDetails)

  fetch('http://127.0.0.1:8080/panda', {
    method: 'POST',
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(recordingDetails)
  })
    .then(response => response.json())
    .then(response => displayResponse(response))
}

let displayResponse = function (response: any) {
  console.log(response)
  const term:any = terminal.current
  for (let i =0; i < response['response'].length; i++) {
    term.pushToStdout('panda@panda:~$ ' + commands[i] + '\n')
    term.pushToStdout(response['response'][i])
  }
}

export default CreateRecordingPage;