import { EuiAvatar, EuiBadge, EuiCollapsibleNavGroup, EuiHeader, EuiHeaderLink, EuiHeaderLinks, EuiHeaderLogo, EuiHeaderSectionItemButton, EuiIcon, EuiListGroup, EuiPageTemplate, EuiPinnableListGroupItemProps, useEuiTheme } from '@elastic/eui'
import { Navigate, Route, Routes, useNavigate } from 'react-router';
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
      label: 'Recordings',
      iconType: 'layers',
      onClick: () => navigate('/recordings'),
    },
    {
      label: 'Images',
      iconType: 'storage',
      onClick: () => navigate('/images'),
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
                <EuiHeaderLink isActive href='https://github.com/panda-re/panda_studio'>Code</EuiHeaderLink>
              </EuiHeaderLinks>,
            ],
            borders: 'right',
          },
        ]}
      />
      <EuiPageTemplate.Sidebar>
        <EuiCollapsibleNavGroup background='light'>
          <EuiListGroup
            listItems={topNavLinks}
          />
        </EuiCollapsibleNavGroup>
      </EuiPageTemplate.Sidebar>

      <Routes>
        <Route path="/images" element={<ImagesPage />} />
        <Route path="/createRecording" element={<CreateRecordingPage />} />
        <Route path='/createImage' element={<CreateImagePage />} />
        <Route path="/recordings" element={<RecordingDashboardPage />}/>
        <Route path="/recordingDetails" element={<RecordingDetailsPage />}/>
        <Route path="/createInteractionProgram" element={<CreateInteractionProgramPage />}/>
        <Route path="/" element={<Navigate to="/recordings" replace={true} />} />
        <Route path="*" element={<NotFoundErrorPage />} />
        <Route path="/imageDetails" element={<ImageDetails />} />
        <Route path="/interactions" element={<InteractionDashboard />} />
        <Route path="/interactionDetails" element={<InteractionDetails />} />
      </Routes>

    </EuiPageTemplate>
  </>)
}

export default Layout;