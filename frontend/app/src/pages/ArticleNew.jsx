import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { getAuthHeaders } from '../utils/api.js';
import './ArticleEdit.css';

function ArticleNew() {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!title.trim()) {
      setError('Title is required');
      return;
    }

    setSaving(true);
    setError(null);

    try {
      const response = await fetch('/api/articles', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          title: title.trim(),
          content: content.trim(),
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || errorData.message || 'Failed to create article');
      }

      const data = await response.json();
      
      // Navigate to the new article's detail page
      if (data.id) {
        navigate(`/article/${data.id}`);
      } else {
        // If no ID returned, go back to home
        navigate('/');
      }
    } catch (err) {
      setError(err.message);
      setSaving(false);
    }
  };

  return (
    <div className="article-edit">
      <div className="article-edit__container">
        <header className="article-edit__header">
          <h1>Create New Article</h1>
          <Link to="/" className="article-edit__cancel-link">
            ← Back to Articles
          </Link>
        </header>

        {error && (
          <div className="article-edit__error-banner" role="alert">
            <strong>Error:</strong> {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="article-edit__form">
          <div className="article-edit__field">
            <label htmlFor="title" className="article-edit__label">
              Title *
            </label>
            <input
              id="title"
              type="text"
              className="article-edit__input"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={saving}
              placeholder="Enter article title"
              autoFocus
              required
            />
          </div>

          <div className="article-edit__field">
            <label htmlFor="content" className="article-edit__label">
              Content (Markdown)
            </label>
            <textarea
              id="content"
              className="article-edit__textarea"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={saving}
              placeholder="Write your article content here using Markdown..."
            />
          </div>

          <div className="article-edit__actions">
            <button
              type="submit"
              className="article-edit__save-button"
              disabled={saving}
            >
              {saving ? 'Creating...' : 'Create Article'}
            </button>
            <Link to="/" className="article-edit__cancel-button">
              Cancel
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}

export default ArticleNew;
