import {EuiButton, EuiFieldText, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiText} from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';

function ImagesPage() {

  return (<>
    <EuiPageTemplate.Header pageTitle="Image Dashboard" />

    <EuiPageTemplate.Section>
      <EuiFlexGroup>
        <EuiFlexItem>
          <EuiText>Image ID</EuiText>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiText>Image Name</EuiText>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiText>Date</EuiText>
        </EuiFlexItem>
      </EuiFlexGroup>

      <EuiFlexGroup>
        <EuiFlexItem>
          <EuiFieldText
            placeholder="Enter Image ID"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldText
            placeholder="Enter Image Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldText
            placeholder="Date"
          />
        </EuiFlexItem>
        <EuiFlexItem >
          <EuiButton>Upload Base Image</EuiButton>
        </EuiFlexItem>
      </EuiFlexGroup>

      <ImagesDataGrid />
    </EuiPageTemplate.Section>
  </>)
}

export default ImagesPage;