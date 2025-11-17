import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import '../App.css';

const normalizeBaseUrl = (value) => {
  if (!value) {
    return 'http://localhost:8080';
  }
  return value.endsWith('/') ? value.slice(0, -1) : value;
};

const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL);

export default function ArticleDetail() {
  const { id } = useParams();
  const [article, setArticle] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchArticle = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${API_BASE_URL}/article/${id}`);
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
        <div style={{ padding: '2rem', textAlign: 'center' }}>
          <p>Loading article...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app">
        <div style={{ padding: '2rem' }}>
          <h1 style={{ color: 'var(--accent-color)' }}>Error</h1>
          <p style={{ color: 'var(--muted-text)' }}>{error}</p>
          <Link to="/" className="secondary" style={{ display: 'inline-block', marginTop: '1rem' }}>
            ← Back to Articles
          </Link>
        </div>
      </div>
    );
  }

  if (!article) {
    return (
      <div className="app">
        <div style={{ padding: '2rem' }}>
          <h1>Article not found</h1>
          <Link to="/" className="secondary" style={{ display: 'inline-block', marginTop: '1rem' }}>
            ← Back to Articles
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <div style={{ maxWidth: '900px', margin: '0 auto' }}>
        <Link to="/" className="secondary" style={{ display: 'inline-block', marginBottom: '2rem' }}>
          ← Back to Articles
        </Link>

        <article className="article-card" style={{ marginTop: '1.5rem' }}>
          <header className="article-card__header">
            <span className="article-card__id">#{article.id ?? '—'}</span>
            <h1 style={{ margin: 0 }}>{article.title || 'Untitled article'}</h1>
            <span className="article-card__updated" title={`Last updated: ${article.updated ?? 'unknown'}`}>
              {article.updated ? new Date(article.updated).toLocaleString() : '—'}
            </span>
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
