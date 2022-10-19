import { EuiPageTemplate, EuiText } from '@elastic/eui';

function NotFoundErrorPage() {
  return (<>
    <EuiPageTemplate.Header pageTitle="Page not Found" />
    <EuiPageTemplate.Section>
      <EuiText>You seem to be lost.</EuiText>
    </EuiPageTemplate.Section>
  </>)
}

export default NotFoundErrorPage