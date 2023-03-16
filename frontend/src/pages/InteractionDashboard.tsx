import {
  EuiButton,
  EuiFieldSearch,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiPageTemplate,
  EuiSpacer,
  EuiText
} from '@elastic/eui';
import InteractionsDataGrid from '../components/InteractionsDataGrid';
import {EuiFlexGrid} from "@elastic/eui";
import RecordingDataGrid from "../components/RecordingDataGrid";
import {useNavigate} from "react-router";

function InteractionDashboard() {
  const navigate = useNavigate();

  return (<>
    <EuiPageTemplate.Header pageTitle='Interactions Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent='flexEnd'>
        <EuiFlexItem grow={false}>
          <EuiButton iconType={'plusInCircle'} onClick={() => navigate('/createInteractionProgram')}>Create New Interaction</EuiButton>
        </EuiFlexItem>
      </EuiFlexGroup>
      <InteractionsDataGrid />
    </EuiPageTemplate.Section>
  </>)
}

export default InteractionDashboard;