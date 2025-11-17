import { useState } from 'react';
import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import './ArticleList.css';

const formatDate = (value) => {
  if (!value) {
    return '—';
  }

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return value;
  }

  return parsed.toLocaleString('cs-CZ');
};

const getAnchorId = (article, index) => {
  if (article && (article.id || article.id === 0)) {
    return `article-${article.id}`;
  }

  return `article-${index}`;
};

export default function ArticleList({ articles, loading }) {
  const [expandedCards, setExpandedCards] = useState({});

  if (loading && (!articles || articles.length === 0)) {
    return (
      <div className="empty-state" role="alert" aria-live="polite">
        Fetching articles…
      </div>
    );
  }

  if (!articles || articles.length === 0) {
    return (
      <div className="empty-state" role="alert" aria-live="polite">
        No articles to display yet. Try refreshing or adjust your filters.
      </div>
    );
  }

  return (
    <ul className="article-list">
      {articles.map((article, index) => {
        const hasContent = Boolean(article.content && article.content.trim().length > 0);
        const key =
          article.id ??
          `${article.title ?? 'article'}-${article.updated ?? index}`;
        const shouldCollapse = hasContent && article.content.trim().length > 420;
        const isExpanded = Boolean(expandedCards[key]);
        const anchorId = getAnchorId(article, index);

        return (
          <li
            key={key}
            className="article-card"
            id={anchorId}
          >
            <header className="article-card__header">
              <Link 
                to={`/article/${article.id}`} 
                className="article-card__id"
                title={`View article #${article.id ?? '—'}`}
              >
                #{article.id ?? '—'}
              </Link>
              <h3>{article.title || 'Untitled article'}</h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem', alignItems: 'flex-end' }}>
                <span
                  className="article-card__updated"
                  title={`Last modified: ${article.updated ?? 'unknown'}`}
                >
                  Modified: {formatDate(article.updated)}
                </span>
                <span
                  className="article-card__updated"
                  title={`Created: ${article.created ?? 'unknown'}`}
                >
                  Created: {formatDate(article.created)}
                </span>
              </div>
            </header>

            {hasContent ? (
              <div
                className={`article-card__content-wrapper${
                  isExpanded ? ' article-card__content-wrapper--expanded' : ''
                }${shouldCollapse ? ' article-card__content-wrapper--collapsible' : ''}`}
              >
                <ReactMarkdown className="article-card__content" remarkPlugins={[remarkGfm]}>
                  {article.content}
                </ReactMarkdown>
              </div>
            ) : (
              <p className="article-card__no-content">
                <em>This article does not have any content yet.</em>
              </p>
            )}

            {shouldCollapse && (
              <button
                type="button"
                className="article-card__toggle secondary"
                onClick={() =>
                  setExpandedCards((prev) => ({
                    ...prev,
                    [key]: !isExpanded,
                  }))
                }
              >
                {isExpanded ? 'Show less' : 'Read more'}
              </button>
            )}
          </li>
        );
      })}
    </ul>
  );
}
