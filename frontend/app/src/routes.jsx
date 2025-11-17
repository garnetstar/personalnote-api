import { Routes, Route } from 'react-router-dom';
import App from './App.jsx';
import ArticleDetail from './pages/ArticleDetail.jsx';

// 404 Not Found page
function NotFoundPage() {
  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <h1>404 - Page Not Found</h1>
      <p>The page you're looking for doesn't exist.</p>
    </div>
  );
}

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<App />} />
      <Route path="/article/:id" element={<ArticleDetail />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
