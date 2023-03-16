import { EuiButton, EuiFieldSearch, EuiFlexGrid, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSpacer } from '@elastic/eui';
import RecordingDataGrid from '../components/RecordingDataGrid';
import { useNavigate } from 'react-router';


function RecordingDashboardPage () {
  const navigate = useNavigate();
  return (<>
    <EuiPageTemplate.Header pageTitle='Recording Dashboard' rightSideItems={[]} />
    <EuiPageTemplate.Section>
    <EuiFlexGroup justifyContent='flexEnd'>
        <EuiFlexItem grow={false}>
          <EuiButton iconType={'plusInCircle'} onClick={() => navigate('/createRecording')}>Create New Recording</EuiButton>
        </EuiFlexItem>
      </EuiFlexGroup>
      <RecordingDataGrid />
    </EuiPageTemplate.Section>
    </>
  );
};

export default RecordingDashboardPage;