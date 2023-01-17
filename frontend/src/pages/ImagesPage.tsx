import {EuiButton, EuiFieldText, EuiFlexGroup, EuiFlexItem, EuiPageTemplate, EuiSpacer, EuiText} from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';
import {EuiFlexGrid} from "@elastic/eui";

function ImagesPage() {

  return (<>
    <EuiPageTemplate.Header pageTitle="Image Dashboard" />

    <EuiPageTemplate.Section>

      <EuiFlexGrid columns={4}>
        <EuiFlexItem>
          <EuiText>Image ID</EuiText>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiText>Image Name</EuiText>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiText>Date</EuiText>
        </EuiFlexItem>
        <EuiFlexItem>
        </EuiFlexItem>
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
        <EuiFlexItem>
          <EuiButton iconType={'plusInCircle'}>Upload Base Image</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>

      <EuiSpacer size="xl" />

      <ImagesDataGrid />
    </EuiPageTemplate.Section>

  </>)
}

export default ImagesPage;