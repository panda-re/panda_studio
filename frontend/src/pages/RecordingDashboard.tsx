import { EuiButton, EuiFieldSearch, EuiFlexGrid, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSpacer } from '@elastic/eui';
import RecordingDataGrid from '../components/RecordingDataGrid';
import { useNavigate } from 'react-router';


function RecordingDashboardPage () {
  const navigate = useNavigate();
  return (<>
    <EuiPageTemplate.Header pageTitle='Recording Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
    <RecordingDataGrid />
    </EuiPageTemplate.Section>
    </>
  );
};

export default RecordingDashboardPage;