import { useCallback, useEffect, useMemo, useState } from 'react';
import ArticleList from './components/ArticleList.jsx';
import './App.css';

const normalizeBaseUrl = (value) => {
  if (!value) {
    return 'http://localhost:8080';
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

export default function App() {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [status, setStatus] = useState('Loading articles…');
  const [filterMode, setFilterMode] = useState('all');
  const [keyword, setKeyword] = useState('');
  const [articleId, setArticleId] = useState('');
  const [lastUpdatedAt, setLastUpdatedAt] = useState(null);

  const baseUrl = useMemo(() => normalizeBaseUrl(API_BASE_URL), []);

  const fetchArticles = useCallback(
    async (path, { emptyMessage } = {}) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${baseUrl}${path}`);
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

  const handleIdSubmit = async (event) => {
    event.preventDefault();
    const trimmedId = articleId.trim();

    if (!trimmedId) {
      setError('Please provide an article ID to look up.');
      setStatus('Awaiting a valid article ID.');
      return;
    }

    if (!/^\d+$/.test(trimmedId)) {
      setError('Article IDs must be numeric.');
      setStatus('Awaiting a valid article ID.');
      return;
    }

    await fetchArticles(`/article/${trimmedId}`, {
      emptyMessage: `No article found with ID ${trimmedId}.`,
    });
  };

  const handleResetFilters = () => {
    setKeyword('');
    setArticleId('');
    setFilterMode('all');
    setError(null);
    loadAllArticles();
  };

  return (
    <div className="app">
      <header className="app__header">
        <button
          type="button"
          className="secondary"
          onClick={loadAllArticles}
          disabled={loading}
        >
          Refresh
        </button>
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
            {lastUpdatedAt ? lastUpdatedAt.toLocaleString() : '—'}
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

        <form className="toolbar" onSubmit={handleIdSubmit}>
          <label className="toolbar__label" htmlFor="article-id">
            Lookup by ID
          </label>
          <input
            id="article-id"
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            placeholder="e.g. 42"
            value={articleId}
            onChange={(event) => setArticleId(event.target.value)}
            disabled={loading}
          />
          <button type="submit" disabled={loading}>
            Fetch
          </button>
        </form>

        <div className="toolbar toolbar--actions">
          <button type="button" className="secondary" onClick={handleResetFilters} disabled={loading}>
            Clear filters
          </button>
        </div>
      </section>

      <div className={`status-bar${error ? ' status-bar--error' : ''}`} role="status" aria-live="polite">
        <span className="status-chip">{loading ? 'Loading' : error ? 'Error' : 'Ready'}</span>
        {loading && <span className="loader" aria-hidden="true" />}
        <span className="message">{error || status}</span>
      </div>

      <ArticleList articles={articles} loading={loading} />
    </div>
  );
}
