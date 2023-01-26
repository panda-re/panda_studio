import {
  EuiButton,
  EuiFieldSearch,
  EuiFlexItem,
  EuiPageTemplate,
  EuiSpacer,
  EuiText
} from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';
import {EuiFlexGrid} from "@elastic/eui";

function ImagesPage() {

  return (<>
    <EuiPageTemplate.Header pageTitle='Image Dashboard' rightSideItems={[]} />

    <EuiPageTemplate.Section>
      <EuiFlexGrid columns={4}>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Image ID"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Image Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Date"
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