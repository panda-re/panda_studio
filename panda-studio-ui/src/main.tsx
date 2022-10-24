import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import 'regenerator-runtime';
import '@elastic/eui/dist/eui_theme_dark.css'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
