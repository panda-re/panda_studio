import { EuiPageTemplate, EuiText } from '@elastic/eui';

function DashboardPage() {
  return (<>
    <EuiPageTemplate.Header pageTitle='header 1' rightSideItems={[<>hello</>]} />

    <EuiPageTemplate.Section>
      <EuiText textAlign='center'>
        <strong>Some strong text</strong>
      </EuiText>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Header pageTitle='Header 2' rightSideItems={[<>hello</>]} />

    <EuiPageTemplate.Section>
      <EuiText textAlign='center'>
        <strong>Some strong text</strong>
      </EuiText>
    </EuiPageTemplate.Section>
  </>);
}

export default DashboardPage;