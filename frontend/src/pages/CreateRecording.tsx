import {EuiButton, EuiPageTemplate, EuiSelectableOption, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useCallback, useEffect, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import React, { Component } from 'react'
import EntitySearchBar from '../components/EntitySearchBar';
import { Image, InteractionProgram } from '../components/Interfaces';

import prettyBytes from 'pretty-bytes';
import { useNavigate } from 'react-router';

function CreateRecordingPage() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [volume, setVolume] = useState('');
  const [program, setProgram] = useState('');
  const [commands, setCommands] = useState('');
  const terminal = useRef(null)

  const sendAPICall = useCallback(function () {
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

  const data: Image[] = [
    {
      id: 'record_1',
      name: 'test_recording',
      operatingSystem: 'Ubuntu',
      date: new Date(),
      size: 150*1024*1024,
    },
    {
      id: 'record_2',
      name: 'test_recording2',
      operatingSystem: 'Ubuntu',
      date: new Date(),
      size: 150*1024*1024,
    }
  ];

  // Generate selectable options for Image search component
  let imageEntities: EuiSelectableOption[] = [];
  data.map((r) =>
    imageEntities.push({label: `${r.id} - ${r.name} - ${r.operatingSystem} - ${r.date.toLocaleDateString()} - ${prettyBytes(r.size, { maximumFractionDigits: 2 })}`,
  data: r})
  );

  const [selectedImage, setSelectedImage] = React.useState<EuiSelectableOption | undefined>(undefined);
  function returnSelectedImage(message: EuiSelectableOption){
   setSelectedImage(message);
  }

  const programs: InteractionProgram[] = [
    {
      id: 'program1',
      name: 'Test Program 1',
      date: new Date(),
    },
    {
      id: 'program2',
      name: 'Test Program 2',
      date: new Date(),
    },
    {
      id: 'program3',
      name: 'Test Program 3',
      date: new Date(),
    },
    {
      id: 'program4',
      name: 'Test Program 4',
      date: new Date(),
    },
    {
      id: 'program5',
      name: 'Test Program 5',
      date: new Date(),
    },
    {
      id: 'program6',
      name: 'Test Program 6',
      date: new Date(),
    },
    {
      id: 'program7',
      name: 'Test Program 7',
      date: new Date(),
    },
    {
      id: 'program8',
      name: 'Test Program 8',
      date: new Date(),
    },
  ]

  // Generate selectable options for Interaction Program search component
  let interactionProgramEntities: EuiSelectableOption[] = [];
  programs.map((r) =>
    interactionProgramEntities.push({label: `${r.id} - ${r.name} - ${r.date.toLocaleDateString()}`})
  );

  const [selectedProgram, setSelectedProgram] = React.useState<EuiSelectableOption | undefined>(undefined);
  function returnSelectedProgram(message: EuiSelectableOption){
   setSelectedProgram(message);
  }

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
          <EntitySearchBar name="Image" entities={imageEntities} returnSelectedOption={(returnSelectedImage)}></EntitySearchBar>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={2}>
          <EuiText>Specify commands to run:</EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EntitySearchBar name="Interaction Program" entities={interactionProgramEntities} returnSelectedOption={(returnSelectedProgram)}></EntitySearchBar>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={false}>
          <div>
            {/* <EuiButton onClick={sendAPICall}>Create Recording</EuiButton> */}
            <EuiButton onClick={() => {
              alert(`Creating recording with name: ${name}, image: ${selectedImage?.data?.name}, program: ${selectedProgram}`); 
              // navigate('/recordings');
            }}>Create Recording</EuiButton>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

  </>)
}
export default CreateRecordingPage;