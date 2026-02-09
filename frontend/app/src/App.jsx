import { useCallback, useEffect, useMemo, useState, memo } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from './context/AuthContext.jsx';
import { getAuthHeaders } from './utils/api.js';
import ArticleList from './components/ArticleList.jsx';
import './App.css';

const normalizeBaseUrl = (value) => {
  if (!value) {
    return '/api';
  }
  return value.endsWith('/') ? value.slice(0, -1) : value;
};

const API_BASE_URL = normalizeBaseUrl(import.meta.env.VITE_API_BASE_URL);

const normaliseArticles = (payload) => {
  if (!payload) {
    return [];
  }

  if (Array.isArray(payload)) {
    return payload;
  }

  if (Array.isArray(payload.articles)) {
    return payload.articles;
  }

  if (Array.isArray(payload.data)) {
    return payload.data;
  }

  if (typeof payload === 'object' && payload !== null) {
    return [payload];
  }

  return [];
};

const extractMessage = (payload) => {
  if (!payload || typeof payload !== 'object') {
    return undefined;
  }

  return payload.message || payload.detail || payload.error;
};

const getAnchorId = (article, index) => {
  if (article && (article.id || article.id === 0)) {
    return `article-${article.id}`;
  }

  return `article-${index}`;
};

const Sidebar = memo(({ loading, articles, sidebarEntries, isCompactLayout }) => {
  return (
    <aside className="app__sidebar" aria-label="Article titles">
      <header className="app__sidebar-header">
        <h2>Articles</h2>
        <span className="app__sidebar-count" aria-live="polite">
          {loading ? '…' : articles.length}
        </span>
      </header>

      {isCompactLayout && (
        <p className="app__sidebar-tip" role="status">
          This panel is stacked above the articles on small screens. Enlarge your window to keep it docked to the left.
        </p>
      )}

      {loading && articles.length === 0 ? (
        <p className="app__sidebar-placeholder">Loading titles…</p>
      ) : sidebarEntries.length === 0 ? (
        <p className="app__sidebar-placeholder">No articles to display yet. (articles: {articles.length}, entries: {sidebarEntries.length})</p>
      ) : (
        <ul className="app__sidebar-list">
          {sidebarEntries.map((entry) => (
            <li key={entry.anchorId} className="app__sidebar-item">
              <Link to={`/article/${entry.id}`} className="app__sidebar-link">
                <span className="app__sidebar-title">{entry.label}</span>
                <span className="app__sidebar-date">
                  {entry.updated ? new Date(entry.updated).toLocaleDateString('cs-CZ') : '—'}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </aside>
  );
});

Sidebar.displayName = 'Sidebar';

export default function App() {
  const { logout, user } = useAuth();
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [status, setStatus] = useState('Loading articles…');
  const [filterMode, setFilterMode] = useState('all');
  const [keyword, setKeyword] = useState('');
  const [lastUpdatedAt, setLastUpdatedAt] = useState(null);
  const [isCompactLayout, setIsCompactLayout] = useState(() => {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
      return false;
    }

    return window.matchMedia('(max-width: 768px)').matches;
  });

  const baseUrl = useMemo(() => normalizeBaseUrl(API_BASE_URL), []);

  const fetchArticles = useCallback(
    async (path, { emptyMessage } = {}) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${baseUrl}${path}`, {
          headers: getAuthHeaders()
        });
        const contentType = response.headers.get('content-type') ?? '';
        const payload = contentType.includes('application/json')
          ? await response.json()
          : await response.text();

        if (!response.ok) {
          const reason =
            typeof payload === 'string'
              ? payload
              : extractMessage(payload) || `Request failed with status ${response.status}`;
          throw new Error(reason);
        }

        const articlesFromResponse = normaliseArticles(payload);
        setArticles(articlesFromResponse);
        setLastUpdatedAt(new Date());

        const messageFromPayload = extractMessage(payload);

        if (articlesFromResponse.length === 0) {
          setStatus(messageFromPayload || emptyMessage || 'No articles to display yet.');
        } else {
          setStatus(
            messageFromPayload ||
              `Showing ${articlesFromResponse.length} article${
                articlesFromResponse.length === 1 ? '' : 's'
              }.`
          );
        }
      } catch (err) {
        setError(err.message || 'Unexpected error while contacting the API.');
        setStatus('Unable to fetch articles.');
        setArticles([]);
      } finally {
        setLoading(false);
      }
    },
    [baseUrl]
  );

  const loadAllArticles = useCallback(() => {
    fetchArticles('/articles', { emptyMessage: 'The database is empty. Add some articles to begin.' });
  }, [fetchArticles]);

  useEffect(() => {
    loadAllArticles();
  }, [loadAllArticles]);

  useEffect(() => {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
      return undefined;
    }

    const mediaQuery = window.matchMedia('(max-width: 768px)');

    const updateCompactLayout = (event) => {
      setIsCompactLayout(event.matches);
    };

    // Ensure the initial value stays in sync with the current viewport.
    setIsCompactLayout(mediaQuery.matches);

    if (typeof mediaQuery.addEventListener === 'function') {
      mediaQuery.addEventListener('change', updateCompactLayout);
    } else {
      mediaQuery.addListener(updateCompactLayout);
    }

    return () => {
      if (typeof mediaQuery.removeEventListener === 'function') {
        mediaQuery.removeEventListener('change', updateCompactLayout);
      } else {
        mediaQuery.removeListener(updateCompactLayout);
      }
    };
  }, []);

  const handleKeywordSubmit = async (event) => {
    event.preventDefault();
    const trimmedKeyword = keyword.trim();

    if (!trimmedKeyword) {
      setError('Please enter a keyword to search.');
      setStatus('Awaiting a keyword to search.');
      return;
    }

    await fetchArticles(`/article/filter/${filterMode}/${encodeURIComponent(trimmedKeyword)}`, {
      emptyMessage: 'No articles matched your search. Try a different term.',
    });
  };

  const sidebarEntries = useMemo(
    () =>
      articles.map((article, index) => ({
        anchorId: getAnchorId(article, index),
        label: (article.title && article.title.trim()) || `Untitled article ${index + 1}`,
        id: article.id ?? index,
        updated: article.updated,
      })),
    [articles]
  );

  return (
    <div className="app">
      <div className={`app__layout${isCompactLayout ? ' app__layout--stacked' : ''}`}>
        <Sidebar 
          loading={loading}
          articles={articles}
          sidebarEntries={sidebarEntries}
          isCompactLayout={isCompactLayout}
        />

        <main className="app__main">
          <header className="app__header">
            <div
              className="app__header-inner"
              style={{ display: 'flex', gap: '0.75rem', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}
            >
              <div className="app__header-actions" style={{ display: 'flex', gap: '0.75rem' }}>
                <button
                  type="button"
                  className="secondary"
                  onClick={loadAllArticles}
                  disabled={loading}
                >
                  Refresh
                </button>
                <Link to="/article/new" style={{ textDecoration: 'none' }}>
                  <button type="button">
                    New Article
                  </button>
                </Link>
              </div>
              <div className="app__header-user" style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                {user && (
                  <div style={{ textAlign: 'right', fontSize: '0.85rem', lineHeight: '1.3' }}>
                    <div style={{ fontWeight: '500', color: '#333' }}>
                      {user.name || user.email}
                    </div>
                    <div style={{ color: '#666', fontSize: '0.8rem' }}>
                      {user.email}
                    </div>
                  </div>
                )}
                <button
                  type="button"
                  className="secondary"
                  onClick={logout}
                  style={{ padding: '0.5rem 1rem' }}
                >
                  Logout
                </button>
              </div>
            </div>
          </header>

          <section className="stat-card" aria-live="polite">
            <div className="stat-card__primary">
              <h2>{articles.length}</h2>
              <p>{articles.length === 1 ? 'article loaded' : 'articles loaded'}</p>
            </div>
            <div className="stat-card__meta">
              <span>
                API base URL: <code>{baseUrl}</code>
              </span>
              <span>
                Last updated:{' '}
                {lastUpdatedAt ? lastUpdatedAt.toLocaleString('cs-CZ') : '—'}
              </span>
            </div>
          </section>

          <section className="controls" aria-label="Article filters and search">
            <form className="toolbar" onSubmit={handleKeywordSubmit}>
              <label className="toolbar__label" htmlFor="keyword-input">
                Keyword search
              </label>
              <select
                id="filter-mode"
                value={filterMode}
                onChange={(event) => setFilterMode(event.target.value)}
                disabled={loading}
              >
                <option value="all">Title &amp; content</option>
                <option value="title">Title only</option>
              </select>
              <input
                id="keyword-input"
                type="text"
                placeholder="e.g. technology"
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
                disabled={loading}
                autoComplete="off"
              />
              <button type="submit" disabled={loading}>
                Search
              </button>
            </form>
          </section>

          <ArticleList articles={articles} loading={loading} />
        </main>
      </div>
    </div>
  );
}
