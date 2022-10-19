import { EuiProvider } from '@elastic/eui';
import { createBrowserRouter, Link, RouterProvider } from 'react-router-dom';
import Layout from './Layout';

const router = createBrowserRouter([
  {
    path: '/nut',
    element: <div><Link to="/" >NUT</Link></div>
  },
  {
    path: '/',
    element: <Layout />,
  },
])

function App() {

  return (
    <EuiProvider colorMode='dark'>
      <RouterProvider router={router} />
    </EuiProvider>
  )
}

export default App;
