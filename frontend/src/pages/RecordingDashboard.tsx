import { EuiButton, EuiFieldSearch, EuiFlexGrid, EuiFlexItem, EuiPageTemplate, EuiSpacer } from '@elastic/eui';
import RecordingDataGrid from '../components/RecordingDataGrid';
import { useNavigate } from 'react-router';


function RecordingDashboardPage () {
  const navigate = useNavigate();
  return (<>
    <EuiPageTemplate.Header pageTitle='Recording Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
    <EuiFlexGrid columns={4}>
        <EuiFlexItem>
          <EuiButton iconType={'plusInCircle'} onClick={() => navigate('/createRecording')}>Create New Recording</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>
      <EuiSpacer size="xl" />
      <RecordingDataGrid />
    </EuiPageTemplate.Section>
    </>
  );
};

export default RecordingDashboardPage;