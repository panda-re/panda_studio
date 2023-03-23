import {EuiButton, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiText} from '@elastic/eui';
import {useLocation, useNavigate} from 'react-router';

import prettyBytes from 'pretty-bytes';
import {useDeleteRecordingById} from "../api";
import {useQueryClient} from "@tanstack/react-query";
import {render} from "react-dom";

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

  return (<>
    <EuiPageTemplate.Header pageTitle="Recording Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.id}
          </EuiText>
        </EuiPageTemplate.Section>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>File Name:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.name}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Image Name:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.imageName}
          </EuiText>
        </EuiPageTemplate.Section>
        {/*
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Size:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {prettyBytes(location.state.item.size, { maximumFractionDigits: 2 })}
          </EuiText>
        </EuiPageTemplate.Section>
        */}
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Date Created:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.date}
          </EuiText>
        </EuiPageTemplate.Section>
      </EuiFlexItem>

      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton
              style={buttonStyle}
              onClick={() => {
                navigate('/recordings')
              }}
            >
              Recording Dashboard
            </EuiButton>
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
            <EuiButton style={buttonStyle} onClick={deleteCurrentRecording}>Delete Recording</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
  </>)
}

export default RecordingDetailsPage;