import express from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';

const app = express();

// Serve static files from dist/ directory
app.use(express.static('dist'));

const BACKEND_HOST = process.env.BACKEND_URL || 'http://localhost:8080';

// Proxy all requests to the backend at /api
app.use('/api', createProxyMiddleware({
    target: BACKEND_HOST,
    changeOrigin: true,
}));

// Serve index.html for all other requests
app.get('*', (req, res) => {
    res.sendFile('dist/index.html', { root: '.' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Frontend server listening on port ${PORT}`);
});
