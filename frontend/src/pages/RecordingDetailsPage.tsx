import {EuiButton, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSpacer, EuiText} from '@elastic/eui';
import prettyBytes from 'pretty-bytes';
import { ReactElement } from 'react';
import {useLocation, useNavigate} from 'react-router';
import { RecordingFile } from '../api';


function RecordingDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  const deleteCurrentRecording = () => {
    navigate('/recordings', {state: {recordingId: location.state.item.id}});
  }

  function CreateRecordingFileRows(files: RecordingFile[]){
    var items: ReactElement[] = [];
    for(var file of files){
      items.push(<EuiFlexGroup>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  {file.id}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  {file.name}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  {file.file_type}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  {(file.size != null) ? prettyBytes(file.size, { maximumFractionDigits: 2 }) : "0"}
                </EuiText>
              </EuiFlexItem>
            </EuiFlexGroup>)
    }
    return items;
  }

  return (<>
    <EuiPageTemplate.Header pageTitle="Recording Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong>  {location.state.item.id}
          </EuiText>
          <EuiSpacer size='l'></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Recording Name:</strong> {location.state.item.name}
          </EuiText>
          <EuiSpacer size='l'></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Image Id:</strong> {location.state.item.image_id}
          </EuiText>
          <EuiSpacer size='l'></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Interactions Id:</strong> {location.state.item.program_id}
          </EuiText>
          <EuiSpacer size='l'></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Date Created:</strong> {(location.state.item.date as String).substring(0, 19)}
          </EuiText>
        <EuiSpacer size='xxl'></EuiSpacer>
        <EuiText textAlign='center'><strong><u>Recording Files</u></strong></EuiText>
        <EuiFlexGroup>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>ID:</strong>
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>Name:</strong>
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>Type:</strong>
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>File Size:</strong>
                </EuiText>
              </EuiFlexItem>
            </EuiFlexGroup>
          {CreateRecordingFileRows(location.state.item.files)}
        </EuiPageTemplate.Section>
      </EuiFlexItem>

      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton style={buttonStyle} onClick={deleteCurrentRecording}>Delete Recording</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
  </>)
}

export default RecordingDetailsPage;