import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useCallback, useEffect, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import React, { Component } from 'react'
import Terminal from 'react-console-emulator'

function CreateRecordingPage() {
  const [name, setName] = useState('');
  const [volume, setVolume] = useState('');
  const [program, setProgram] = useState('');
  const [commands, setCommands] = useState('');
  const terminal = React.createRef()

  const sendAPICall = useCallback(function (name:string, volume:string, program:string) {
    setCommands(JSON.parse("[" + program + "]"))
    let recordingDetails = {
      volume: volume,
      commands: commands,
      name: name
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
  }, [name, volume, program])

  const displayResponse = useCallback(function (response: any) {
    const term:any = terminal.current
    for (let i =0; i < response['response'].length; i++) {
      term.pushToStdout('panda@panda:~$ ' + commands[i] + '\n')
      term.pushToStdout(response['response'][i])
    }
  }, [commands])

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
            value={name}
            onChange={(e) => setName(e.target.value)}
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
            value={volume}
            onChange={(e) => setVolume(e.target.value)}
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
            value={program}
            onChange={(e) => setProgram(e.target.value)}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>


    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={false}>
          <div>
            <EuiButton onClick={() => sendAPICall(name, volume, program)}>Create Recording</EuiButton>
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
            readOnly={true}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </>)
}
export default CreateRecordingPage;