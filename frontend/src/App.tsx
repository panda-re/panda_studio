import { EuiProvider } from '@elastic/eui';
import { createBrowserRouter, Link, RouterProvider } from 'react-router-dom';
import Layout from './Layout';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const router = createBrowserRouter([
  {
    path: '*',
    element: <Layout />,
  },
])

const queryClient = new QueryClient()

function App() {

  return (
    <QueryClientProvider  client={queryClient}>
      <EuiProvider colorMode='dark'>
        <RouterProvider router={router} />
      </EuiProvider>
    </QueryClientProvider>
  )
}

export default App;
