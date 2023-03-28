import {EuiButton, EuiPageTemplate, EuiSelectableOption, EuiText} from '@elastic/eui';
import {useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import React from 'react'
import EntitySearchBar from '../components/EntitySearchBar';

import prettyBytes from 'pretty-bytes';
import { ExecuteProgramRequest, useExecuteProgramById, useFindAllImages, useFindAllPrograms } from '../api';
import { useNavigate } from 'react-router';

function CreateRecordingPage() {
  const navigate = useNavigate();
  const [name, setName] = useState('');

  var programEntities: EuiSelectableOption[] = [];
  var imageEntities: EuiSelectableOption[] = [];

  const {isLoading: imagesLoading, data: images} = useFindAllImages();
  const {isLoading: programsLoading, data: programs} = useFindAllPrograms();
  const executeFn = useExecuteProgramById({mutation: {onSuccess: () => navigate('/recordings')}});

  if(images != null){
    images.map((r) =>{
      var size = 0;
      if(r.files != null){
        for(var f of r.files){
          size+= (f.size != null) ? +f.size: 0;
        }
      }
      imageEntities.push({label: `Image Name: ${r.name}  ----   Image Id: ${r.id}  ----   Image Size: ${prettyBytes(size, { maximumFractionDigits: 2 })}`,
      data: r});
    })
  }

  // Generate selectable options for Interaction Program search component
  if(programs != null){
    programs.map((r) =>
      programEntities.push({label: `Program Name: ${r.name} ------  Id: ${r.id}`, data: r})
    );
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
          {imagesLoading && <div>Loading...</div> ||
          <EntitySearchBar name="Image" entities={imageEntities} returnSelectedOption={(returnSelectedImage)}></EntitySearchBar>}
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={2}>
          <EuiText>Interaction Program:</EuiText>
        </EuiFlexItem>
        <EuiFlexItem grow={8}>
          {programsLoading && <div>Loading...</div> ||
          <EntitySearchBar name="Interaction Program" entities={programEntities} returnSelectedOption={(returnSelectedProgram)}></EntitySearchBar>}
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceAround"}>
        <EuiFlexItem grow={false}>
          <div>
            {/* <EuiButton onClick={sendAPICall}>Create Recording</EuiButton> */}
            <EuiButton onClick={() => {
              if(selectedImage==null || selectedProgram==null){
                alert("Please select an image and program");
                return;
              }
              const req: ExecuteProgramRequest = {
                imageId: selectedImage.data!.id,
              }
              alert(`Creating recording with name: ${name}, image: ${selectedImage.data!.id}, program: ${selectedProgram.data?.id}`);
              executeFn.mutate({programId: selectedProgram.data!.id, data: req})
            }}>Create Recording</EuiButton>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>

  </>)
}
export default CreateRecordingPage;