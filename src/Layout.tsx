import { EuiAvatar, EuiBadge, EuiButton, EuiCode, EuiCollapsibleNav, EuiCollapsibleNavGroup, EuiHeader, EuiHeaderLink, EuiHeaderLinks, EuiHeaderLogo, EuiHeaderSectionItemButton, EuiIcon, EuiPageTemplate, EuiPinnableListGroup, EuiPinnableListGroupItemProps, EuiProvider, EuiSideNav, EuiSpacer, EuiText, EuiTitle, useEuiTheme } from '@elastic/eui'

function Layout() {
  const { euiTheme } = useEuiTheme();

  const topNavLinks: EuiPinnableListGroupItemProps[] = [
    {
      label: 'Dashboard',
      pinned: true,
      iconType: 'home',
      isActive: true,
    },
    {
      label: 'Images',
      pinned: true,
    },
  ];

  return (<>
    <EuiHeader
      theme="dark"
      sections={[
        {
          items: [
            <EuiHeaderLogo>Elastic</EuiHeaderLogo>,
            <EuiHeaderLinks aria-label="App navigation dark theme example">
              <EuiHeaderLink isActive>Docs</EuiHeaderLink>
              <EuiHeaderLink>Code</EuiHeaderLink>
              <EuiHeaderLink iconType="help"> Help</EuiHeaderLink>
            </EuiHeaderLinks>,
          ],
          borders: 'right',
        },
        {
          items: [
            <EuiBadge
              color={euiTheme.colors.darkestShade}
              iconType="arrowDown"
              iconSide="right"
            >
              Production logs
            </EuiBadge>,
            <EuiHeaderSectionItemButton
              aria-label="2 Notifications"
              notification={'2'}
            >
              <EuiIcon type="cheer" size="m" />
            </EuiHeaderSectionItemButton>,
            <EuiHeaderSectionItemButton aria-label="Account menu">
              <EuiAvatar name="John Username" size="s" />
            </EuiHeaderSectionItemButton>,
          ],
          borders: 'none',
        },
      ]}
    />
    <EuiPageTemplate>
      <EuiCollapsibleNav
        isOpen={true}
        isDocked={true}
        size={240}
        button={
          <EuiButton onClick={() => { }}>Toggle nav</EuiButton>
        }
        onClose={() => { }}
      >
        <EuiCollapsibleNavGroup background='light'>
          <EuiPinnableListGroup
            listItems={topNavLinks}
            onPinClick={function (item: EuiPinnableListGroupItemProps): void {
              throw new Error('Function not implemented.');
            }} />
          <EuiIcon type="alert" />
        </EuiCollapsibleNavGroup>
      </EuiCollapsibleNav>

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
    </EuiPageTemplate>
  </>)
}

export default Layout;