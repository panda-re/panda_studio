import {EuiButton, EuiButtonEmpty, EuiButtonIcon, EuiConfirmModal, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSpacer, EuiText, formatDate} from '@elastic/eui';
import moment from 'moment';
import prettyBytes from 'pretty-bytes';
import { ReactElement, useState } from 'react';
import {useLocation, useNavigate} from 'react-router';
import { downloadRecordingFile, RecordingFile, useDownloadRecordingFile } from '../api';

function downloadHandler (event: React.MouseEvent, recordingId: string, file: RecordingFile){
  downloadRecordingFile(recordingId, file.id ?? "").then((data) => {
    const fileURL = window.URL.createObjectURL(data);
    let alink = document.createElement('a');
    alink.href = fileURL;
    alink.download = file.name ?? "defaultName";
    alink.click();
  });
}

function RecordingDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()
  
  const [isConfirmVisible, setIsConfirmVisible] = useState(false);

  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  function ConfirmModal(){
    return <EuiConfirmModal
      title="Are you sure you want to delete?"
      onCancel={() => setIsConfirmVisible(false)}
      onConfirm={() => navigate('/recordings', {state: {recordingId: location.state.item.id}})}
      cancelButtonText="Cancel"
      confirmButtonText="Delete Program"
      buttonColor="danger"
      defaultFocusedButton="confirm"
    ></EuiConfirmModal>;
  }

  function CreateRecordingFileRows(files: RecordingFile[]){
    var items: ReactElement[] = [];
    for(const file of files){
      items.push(<EuiFlexGroup justifyContent='center'>
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
              <EuiFlexItem>
                <EuiText textAlign='center'>
                  <EuiButtonIcon
                    iconType={"download"}
                    onClick={(value: React.MouseEvent) => {
                      downloadHandler(value, location.state.item.id, file)}}>
                  </EuiButtonIcon>
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
            <strong>Date Created:</strong> {formatDate(moment((location.state.item.date as String).slice(0, 19)), 'dateTime')}
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
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>Download:</strong>
                </EuiText>
              </EuiFlexItem>
            </EuiFlexGroup>
          {CreateRecordingFileRows(location.state.item.files)}
        </EuiPageTemplate.Section>
      </EuiFlexItem>

      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton style={buttonStyle} onClick={() => setIsConfirmVisible(true)}>Delete Recording</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
    {(isConfirmVisible) ? (ConfirmModal()) : null}
  </>)
}

export default RecordingDetailsPage;