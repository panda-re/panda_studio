import {EuiPageTemplate, EuiText} from '@elastic/eui';
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation} from "react-router";

function InteractionDetails() {
  const location = useLocation()
  const interactions = [
    "uname -a\n",
    "ls /\n",
    "touch NEWFILE.txt\n",
    "ls /\n",
    "cd /\n",
    "sudo rm -rf bin\n"
  ]

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
            <strong>Date Created:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.date.toString()}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Interactions:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[0]}
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[1]}
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[2]}
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[3]}
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[4]}
          </EuiText>
          <EuiText textAlign={"center"}>
            {interactions[5]}
          </EuiText>
        </EuiPageTemplate.Section>

      </EuiFlexItem>


    </EuiFlexGroup>

  </>)
}


export default InteractionDetails;