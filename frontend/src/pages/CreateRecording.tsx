import {EuiButton, EuiPageTemplate, EuiSelectableOption, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useCallback, useEffect, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import React, { Component } from 'react'
import EntitySearchBar from '../components/EntitySearchBar';
import { InteractionProgram } from '../components/Interfaces';

import prettyBytes from 'pretty-bytes';
import { useNavigate } from 'react-router';
import { AxiosRequestConfig } from 'axios';
import { findAllImages, Image, useFindAllImages } from '../api';

// const sendAPICall = useCallback(function () {
//   setCommands(JSON.parse("[" + program + "]"))
//   let recordingDetails = {
//     volume: volume,
//     commands: commands,
//     name: name
//   }
//   console.log(recordingDetails)

//   fetch('http://127.0.0.1:8080/panda', {
//     method: 'POST',
//     headers: {
//       'Access-Control-Allow-Origin': '*',
//       'Accept': 'application/json',
//       'Content-Type': 'application/json'
//     },
//     body: JSON.stringify(recordingDetails)
//   })
//     .then(response => response.json())
//     .then(response => displayResponse(response))
// }, [name, volume, program])

// const displayResponse = useCallback(function (response: any) {
//   const term:any = terminal.current
//   for (let i =0; i < response['response'].length; i++) {
//     term.pushToStdout('panda@panda:~$ ' + commands[i] + '\n')
//     term.pushToStdout(response['response'][i])
//   }
// }, [commands])

function CreateRecordingPage() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [volume, setVolume] = useState('');
  const [program, setProgram] = useState('');
  const [commands, setCommands] = useState('');
  const terminal = useRef(null);

  const reqConfig: AxiosRequestConfig = {
    baseURL: "http://localhost:8080/api"
  }

  var imageEntities: EuiSelectableOption[] = [];

  const {isLoading, error, data} = useFindAllImages({ axios: reqConfig});

  if(data?.data != null){
    data.data.map((r) =>{
      var size = 0;
      if(r.files != null){
        for(var f of r.files){
          size+= (f.size != null) ? +f.size: 0;
        }
      }
      imageEntities.push({label: `Image Name: ${r.name}  ----   Image Id: ${r.id}  ----   Image Size: ${prettyBytes(size, { maximumFractionDigits: 2 })}`,
      // entities.push({label: `id:${r.id} - name:${r.name} - size:${prettyBytes(size, { maximumFractionDigits: 2 })}`,
      data: r});
    })
  }

  function getProgramEntities(){
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
    return interactionProgramEntities;
  }

  const [selectedImage, setSelectedImage] = React.useState<EuiSelectableOption | undefined>(undefined);
  function returnSelectedImage(message: EuiSelectableOption){
    setSelectedImage(message);
  }

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
          {isLoading && <div>Loading...</div> ||
          <EntitySearchBar name="Image" entities={imageEntities} returnSelectedOption={(returnSelectedImage)}></EntitySearchBar>}
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={2}>
          <EuiText>Specify commands to run:</EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EntitySearchBar name="Interaction Program" entities={getProgramEntities()} returnSelectedOption={(returnSelectedProgram)}></EntitySearchBar>
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