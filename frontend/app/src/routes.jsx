import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext.jsx';
import App from './App.jsx';
import ArticleDetail from './pages/ArticleDetail.jsx';
import ArticleEdit from './pages/ArticleEdit.jsx';
import ArticleNew from './pages/ArticleNew.jsx';
import ImageUpload from './pages/ImageUpload.jsx';
import Login from './pages/Login.jsx';
import AuthCallback from './pages/AuthCallback.jsx';

// 404 Not Found page
function NotFoundPage() {
  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <h1>404 - Page Not Found</h1>
      <p>The page you're looking for doesn't exist.</p>
    </div>
  );
}

// Protected route wrapper
function ProtectedRoute({ children }) {
  const { isAuthenticated, loading } = useAuth();
  
  if (loading) {
    return <div style={{ padding: '2rem', textAlign: 'center' }}>Loading...</div>;
  }
  
  return isAuthenticated ? children : <Navigate to="/login" />;
}

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/auth/callback" element={<AuthCallback />} />
      <Route path="/" element={<ProtectedRoute><App /></ProtectedRoute>} />
      <Route path="/article/new" element={<ProtectedRoute><ArticleNew /></ProtectedRoute>} />
      <Route path="/upload" element={<ProtectedRoute><ImageUpload /></ProtectedRoute>} />
      <Route path="/article/:id" element={<ProtectedRoute><ArticleDetail /></ProtectedRoute>} />
      <Route path="/article/:id/edit" element={<ProtectedRoute><ArticleEdit /></ProtectedRoute>} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
