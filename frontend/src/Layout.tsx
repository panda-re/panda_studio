import { EuiAvatar, EuiBadge, EuiButton, EuiCode, EuiCollapsibleNav, EuiCollapsibleNavGroup, EuiEmptyPrompt, EuiHeader, EuiHeaderLink, EuiHeaderLinks, EuiHeaderLogo, EuiHeaderSectionItemButton, EuiIcon, EuiListGroup, EuiPage, EuiPageHeader, EuiPageSection, EuiPageSidebar, EuiPageTemplate, EuiPinnableListGroup, EuiPinnableListGroupItemProps, EuiProvider, EuiSideNav, EuiSpacer, EuiText, EuiTitle, useEuiTheme } from '@elastic/eui'
import { Navigate, Route, Routes, useNavigate } from 'react-router';
import DashboardPage from './pages/DashboardPage';
import NotFoundErrorPage from './pages/NotFoundErrorPage';

import pandaLogo from './assets/panda.svg';
import ImagesPage from './pages/ImagesPage';
import CreateRecordingPage from "./pages/CreateRecording";
import CreateImagePage from "./pages/CreateImage";
import RecordingDashboardPage from './pages/RecordingDashboard';
import RecordingDetailsPage from './pages/RecordingDetailsPage';
import ImageDetails from "./pages/ImageDetails";
import CreateInteractionProgramPage from './pages/CreateInteractionProgram';
import InteractionDashboard from "./pages/InteractionDashboard";
import InteractionDetails from "./pages/InteractionDetails";

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
      label:'Interactions',
      iconType: 'list',
      onClick: () => navigate('/interactions')
    },
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
        <Route path="/recordings" element={<RecordingDashboardPage />}/>
        <Route path="/recordingDetails" element={<RecordingDetailsPage />}/>
        <Route path="/createInteractionProgram" element={<CreateInteractionProgramPage />}/>
        <Route path="/" element={<Navigate to="/dashboard" replace={true} />} />
        <Route path="*" element={<NotFoundErrorPage />} />
        <Route path="/imageDetails" element={<ImageDetails />} />
        <Route path="/interactions" element={<InteractionDashboard />} />
        <Route path="/interactionDetails" element={<InteractionDetails />} />
      </Routes>

    </EuiPageTemplate>
  </>)
}

export default Layout;