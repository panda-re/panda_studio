import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation} from "react-router";
import {useNavigate} from "react-router-dom";

function InteractionDetails() {
  const location = useLocation();
  const navigate = useNavigate();

  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  const deleteCurrentInteractionProgram = () => {
    navigate('/interactions', {state: {programId: location.state.item.id}});
  }

  return(<>
    <EuiPageTemplate.Header pageTitle="Interaction Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong> {location.state.item.id}
          </EuiText>
        </EuiPageTemplate.Section>
        <EuiPageTemplate.Section>
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
            <EuiButton style={buttonStyle} onClick={deleteCurrentInteractionProgram}>Delete Interaction</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>

  </>)
}


export default InteractionDetails;