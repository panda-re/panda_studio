import { EuiButton, EuiFieldSearch, EuiFlexGrid, EuiFlexItem, EuiPageTemplate, EuiSpacer } from '@elastic/eui';
import RecordingDataGrid from '../components/RecordingDataGrid';

function RecordingDashboardPage () {
  return (<>
    <EuiPageTemplate.Header pageTitle='Recording Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
    <EuiFlexGrid columns={4}>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Recording ID"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter File Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Image Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiButton iconType={'plusInCircle'}>Create New Recording</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>
      <EuiSpacer size="xl" />
      <RecordingDataGrid />
    </EuiPageTemplate.Section>
    </>
  );
};

export default RecordingDashboardPage;