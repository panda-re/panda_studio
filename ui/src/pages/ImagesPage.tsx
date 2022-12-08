import { EuiPageTemplate } from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';

function ImagesPage() {

  return (<>
    <EuiPageTemplate.Header pageTitle="Images" />

    <EuiPageTemplate.Section>
      <ImagesDataGrid />
    </EuiPageTemplate.Section>
  </>)
}

export default ImagesPage;