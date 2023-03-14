import {EuiPageTemplate, EuiText} from '@elastic/eui';
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation} from "react-router";

function InteractionDetails() {
  const location = useLocation()

  return(<>
    <EuiPageTemplate.Header pageTitle="Interaction Details" />

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
            <strong>Name:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.name}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Interactions:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            <div style={{ whiteSpace: "pre-line" }}>{location.state.item.instructions}</div>
          </EuiText>
        </EuiPageTemplate.Section>

      </EuiFlexItem>


    </EuiFlexGroup>

  </>)
}


export default InteractionDetails;