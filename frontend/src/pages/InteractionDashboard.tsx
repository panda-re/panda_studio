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

function InteractionDashboard() {

  return (<>
    <EuiPageTemplate.Header pageTitle='Interactions Dashboard' rightSideItems={[]} />

    <EuiPageTemplate.Section>
      <EuiFlexGrid columns={4}>
        {/* <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Interaction ID"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Interaction Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Date"
          />
        </EuiFlexItem> */}
        <EuiFlexItem>
          <EuiButton iconType={'plusInCircle'}>Create New Interaction</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>
      <EuiSpacer size="xl" />

      <InteractionsDataGrid />
    </EuiPageTemplate.Section>
  </>)
}

export default InteractionDashboard;