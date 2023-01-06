import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {EuiInnerText} from "@elastic/eui";

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
        <EuiFlexItem grow={false}>
          <div>
            <EuiInnerText>
              {(ref, innerText) => (
                <span id="inner" title={innerText}>
                </span>
              )}
            </EuiInnerText>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </>)
}

let sendAPICall = function (nameInput: MutableRefObject<any>, volumeInput: MutableRefObject<any>, programInput: MutableRefObject<any>) {
  let commands = JSON.parse("[" + programInput.current.value + "]")
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
  const element = document.getElementById("inner")
  if (element != null) {
    element.innerText = response['response']
  }
}

export default CreateRecordingPage;