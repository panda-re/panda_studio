import {
  EuiTextArea,
  EuiButton,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiPageTemplate,
  EuiSpacer,
  EuiLink} from "@elastic/eui";
import React from "react";
import { useState } from "react";
import {CreateProgramRequest, useCreateProgram} from "../api";
import {useNavigate} from "react-router-dom";
import {useQueryClient} from "@tanstack/react-query";

function CreateInteractionProgramPage (){
  // const makeId = htmlIdGenerator();
  const [instructionsValue, setInstructions] = useState('');
  const [nameValue, setName] = useState('');
  const navigate = useNavigate();

  const onInstructionsChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInstructions(e.target.value);
  }

  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  }

  const queryClient = useQueryClient();
  const createInteractionProgram = useCreateProgram({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries();
        navigate('/interactions')
      },
      onError: (response) => alert(response),
    }
  });

  return (<>
    <EuiPageTemplate.Header pageTitle='Create Interaction Program' />
    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceEvenly"}>
        <EuiFlexItem grow={2}>
          <EuiFieldText
            placeholder={"Program Name"}
            value={nameValue}
            onChange={e => onNameChange(e)}>
          </EuiFieldText>
        </EuiFlexItem>
        <EuiFlexItem grow={4}>
          <div>
            <EuiButton onClick={() => {
              if (instructionsValue == "") {
                alert('No instructions?')
                return
              }
              const createProgramRequest: CreateProgramRequest = {
                name: nameValue,
                instructions: instructionsValue
              }
              createInteractionProgram.mutate({data: createProgramRequest})
            }}>Create Interaction Program</EuiButton>
          </div>
        </EuiFlexItem>
        <EuiFlexItem grow>
          <EuiLink href="https://github.com/panda-re/panda_studio/blob/main/cmd/panda_api/commands.md" target="_blank">Program Format Help</EuiLink>
        </EuiFlexItem>
      </EuiFlexGroup>
      <EuiSpacer size="m"></EuiSpacer>
      <EuiTextArea
        fullWidth={true}
        value={instructionsValue}
        onChange={e => onInstructionsChange(e)}>
      </EuiTextArea>
    </EuiPageTemplate.Section>
  </>);
}

export default CreateInteractionProgramPage;