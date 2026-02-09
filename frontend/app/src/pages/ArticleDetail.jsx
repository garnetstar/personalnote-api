import { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { getAuthHeaders } from '../utils/api.js';
import { useAuth } from '../context/AuthContext.jsx';
import '../App.css';

const normalizeBaseUrl = (value) => {
  if (!value) {
    return '/api';
  }
  return value.endsWith('/') ? value.slice(0, -1) : value;
};

const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL);

export default function ArticleDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [article, setArticle] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this article? This action cannot be undone.')) {
      return;
    }

    setDeleting(true);
    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/article/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || errorData.message || 'Failed to delete article');
      }

      navigate('/');
    } catch (err) {
      setError(err.message);
      setDeleting(false);
    }
  };

  useEffect(() => {
    const fetchArticle = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${API_BASE_URL}/article/${id}`, {
          headers: getAuthHeaders()
        });
        const contentType = response.headers.get('content-type') ?? '';
        const data = contentType.includes('application/json')
          ? await response.json()
          : await response.text();

        if (!response.ok) {
          throw new Error(
            typeof data === 'string' ? data : data.message || `Failed to load article ${id}`
          );
        }

        setArticle(data);
      } catch (err) {
        setError(err.message || 'Failed to load article');
      } finally {
        setLoading(false);
      }
    };

    fetchArticle();
  }, [id]);

  if (loading) {
    return (
      <div className="app">
        <div style={{ padding: '1.25rem 0.75rem', textAlign: 'center' }}>
          <p>Loading article...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app">
        <div style={{ padding: '1.25rem 0.75rem' }}>
          <h1 style={{ color: 'var(--accent-color)' }}>Error</h1>
          <p style={{ color: 'var(--muted-text)' }}>{error}</p>
          <Link to="/" className="secondary" style={{ display: 'inline-block', marginTop: '0.75rem' }}>
            ← Back to Articles
          </Link>
        </div>
      </div>
    );
  }

  if (!article) {
    return (
      <div className="app">
        <div style={{ padding: '1.25rem 0.75rem' }}>
          <h1>Article not found</h1>
          <Link to="/" className="secondary" style={{ display: 'inline-block', marginTop: '0.75rem' }}>
            ← Back to Articles
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <div style={{ maxWidth: '100%', margin: '0 auto', padding: '0 0.75rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem', flexWrap: 'wrap', gap: '0.75rem' }}>
          <Link to="/" className="secondary" style={{ display: 'inline-block' }}>
            ← Back to Articles
          </Link>
          {user && article && user.id === article.user_id && (
            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
              <Link 
                to={`/article/${id}/edit`} 
                style={{ 
                  padding: '0.4rem 0.8rem', 
                  background: 'var(--accent-color)', 
                  color: 'white', 
                  textDecoration: 'none', 
                  borderRadius: '4px',
                  transition: 'all 0.2s',
                  fontSize: '0.9rem'
                }}
              >
                Edit Article
              </Link>
              <button
                onClick={handleDelete}
                disabled={deleting}
                style={{
                  padding: '0.4rem 0.8rem',
                  background: '#dc2626',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: deleting ? 'not-allowed' : 'pointer',
                  opacity: deleting ? 0.6 : 1,
                  transition: 'all 0.2s',
                  fontSize: '0.9rem'
                }}
              >
                {deleting ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          )}
        </div>

        <article className="article-card" style={{ marginTop: '0.5rem' }}>
          <header className="article-card__header">
            <span className="article-card__id">#{article.id ?? '—'}</span>
            <h1 style={{ margin: 0 }}>{article.title || 'Untitled article'}</h1>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem', alignItems: 'flex-end' }}>
              <span className="article-card__updated" title={`Last modified: ${article.updated ?? 'unknown'}`}>
                Modified: {article.updated ? new Date(article.updated).toLocaleString('cs-CZ') : '—'}
              </span>
              <span className="article-card__updated" title={`Created: ${article.created ?? 'unknown'}`}>
                Created: {article.created ? new Date(article.created).toLocaleString('cs-CZ') : '—'}
              </span>
            </div>
          </header>

          {article.content && article.content.trim() ? (
            <div className="article-card__content-wrapper article-card__content-wrapper--expanded">
              <ReactMarkdown className="article-card__content" remarkPlugins={[remarkGfm]}>
                {article.content}
              </ReactMarkdown>
            </div>
          ) : (
            <p className="article-card__no-content">
              <em>This article does not have any content yet.</em>
            </p>
          )}
        </article>
      </div>
    </div>
  );
}
