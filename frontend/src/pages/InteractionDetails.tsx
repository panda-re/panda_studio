import {EuiButton, EuiConfirmModal, EuiPageTemplate, EuiSpacer, EuiText} from '@elastic/eui';
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import { useState } from 'react';
import {useLocation} from "react-router";
import {useNavigate} from "react-router-dom";

function InteractionDetails() {
  const location = useLocation();
  const navigate = useNavigate();

  const [isConfirmVisible, setIsConfirmVisible] = useState(false);

  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  function ConfirmModal(){
    return <EuiConfirmModal
      title="Are you sure you want to delete?"
      onCancel={() => setIsConfirmVisible(false)}
      onConfirm={() => navigate('/interactions', {state: {programId: location.state.item.id}})}
      cancelButtonText="Cancel"
      confirmButtonText="Delete Program"
      buttonColor="danger"
      defaultFocusedButton="confirm"
    ></EuiConfirmModal>;
  }

  return(<>
    <EuiPageTemplate.Header pageTitle="Interaction Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong> {location.state.item.id}
          </EuiText>
          <EuiSpacer size='xxl'></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Name:</strong> {location.state.item.name}
          </EuiText>
        </EuiPageTemplate.Section>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong><u>Interactions:</u></strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            <div style={{ whiteSpace: "pre-line" }}>{location.state.item.instructions}</div>
          </EuiText>
        </EuiPageTemplate.Section>
      </EuiFlexItem>
      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton style={buttonStyle} onClick={() => setIsConfirmVisible(true)}>Delete Interaction</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
    {(isConfirmVisible) ? (ConfirmModal()) : null}
  </>)
}


export default InteractionDetails;