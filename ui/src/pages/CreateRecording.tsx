import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';


function CreateRecordingPage() {
  const nameRef = useRef(null)
  const memoryRef = useRef(null)
  const imageRef = useRef(null)
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
          <EuiText>Memory size: </EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EuiFieldText
            placeholder="eg, 2G"
            inputRef={memoryRef}
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
            inputRef={imageRef}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={2}>
          <EuiText>Specify program to run:</EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          <EuiFieldText
            placeholder="eg, python script.py arg1 arg2..."
            inputRef={programRef}
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>


    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={false}>
          <div>
            <EuiButton onClick={() => sendAPICall(nameRef, memoryRef, imageRef, programRef)}>Create Recording</EuiButton>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </>)
}

let sendAPICall = function (nameInput: MutableRefObject<any>, memoryInput: MutableRefObject<any>, imageInput: MutableRefObject<any>, programInput: MutableRefObject<any>) {
  let recordingDetails = {
    name: nameInput.current.value,
    image: imageInput.current.value,
    memory: memoryInput.current.value,
    commands: programInput.current.value
  }
  console.log(recordingDetails)
  console.log("This is where we run the recording in PANDA")

  fetch('http://127.0.0.1:5000/runPanda', {
    method: 'POST',
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(recordingDetails)
  })
    .then(response => response.json())
    .then(response => console.log(JSON.stringify(response)))
}

export default CreateRecordingPage;