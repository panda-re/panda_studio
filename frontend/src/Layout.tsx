import { EuiAvatar, EuiBadge, EuiButton, EuiCode, EuiCollapsibleNav, EuiCollapsibleNavGroup, EuiEmptyPrompt, EuiHeader, EuiHeaderLink, EuiHeaderLinks, EuiHeaderLogo, EuiHeaderSectionItemButton, EuiIcon, EuiListGroup, EuiPage, EuiPageHeader, EuiPageSection, EuiPageSidebar, EuiPageTemplate, EuiPinnableListGroup, EuiPinnableListGroupItemProps, EuiProvider, EuiSideNav, EuiSpacer, EuiText, EuiTitle, useEuiTheme } from '@elastic/eui'
import { Navigate, Route, Routes, useNavigate } from 'react-router';
import DashboardPage from './pages/DashboardPage';
import NotFoundErrorPage from './pages/NotFoundErrorPage';

import pandaLogo from './assets/panda.svg';
import ImagesPage from './pages/ImagesPage';
import CreateRecordingPage from "./pages/CreateRecording";
import CreateImagePage from "./pages/CreateImage";

function Layout() {
  const { euiTheme } = useEuiTheme();
  const navigate = useNavigate();

  const topNavLinks: EuiPinnableListGroupItemProps[] = [
    {
      label: 'Dashboard',
      iconType: 'home',
      isActive: true,
      onClick: () => navigate('/dashboard'),
    },
    {
      label: 'Images',
      iconType: 'storage',
      onClick: () => navigate('/images'),
    },
    {
      label: 'Recordings',
      iconType: 'layers',
      onClick: () => navigate('/recordings'),
    },
    {
      label:'Create Recording',
      iconType: 'play',
      onClick: () => navigate('/createRecording')
    },
    {
      label:'Create Image',
      iconType: 'documentEdit',
      onClick: () => navigate('/createImage')
    }
  ];

  return (<>
    <EuiPageTemplate>
      <EuiHeader
        sections={[
          {
            items: [
              <EuiHeaderLogo iconType={pandaLogo}>PANDA Studio</EuiHeaderLogo>,
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
      {/* Breadcrumb example
      <EuiHeader
        sections={[{
          breadcrumbs: [
            {
              text: 'Management',
              href: '#',
              onClick: (e) => {
                e.preventDefault();
              },
            },
            {
              text: 'Users',
            },
          ]
        }]}
      />
      */}
      <EuiPageTemplate.Sidebar>
        <EuiCollapsibleNavGroup background='light'>
          <EuiListGroup
            listItems={topNavLinks}
          />
        </EuiCollapsibleNavGroup>
      </EuiPageTemplate.Sidebar>

      <Routes>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/images" element={<ImagesPage />} />
        <Route path="/createRecording" element={<CreateRecordingPage />} />
        <Route path='/createImage' element={<CreateImagePage />} />
        <Route path="/" element={<Navigate to="/dashboard" replace={true} />} />
        <Route path="*" element={<NotFoundErrorPage />} />
      </Routes>

    </EuiPageTemplate>
  </>)
}

export default Layout;