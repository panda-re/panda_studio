import { EuiButton, EuiFieldSearch, EuiFieldText, EuiFlexGrid, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSearchBar, EuiSpacer, EuiText } from '@elastic/eui';
import RecordingDataGrid from '../components/RecordingDataGrid';

function RecordingDashboardPage () {
  return (<>
    <EuiPageTemplate.Header pageTitle='Recording Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
      <EuiSpacer size="xl" />
      <RecordingDataGrid />
    </EuiPageTemplate.Section>
    </>
  );
};

export default RecordingDashboardPage;