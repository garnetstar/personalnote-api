import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { getAuthHeaders } from '../utils/api.js';
import './ArticleEdit.css';

const API_BASE_URL = '/api';

export default function ArticleEdit() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [article, setArticle] = useState(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchArticle = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await fetch(`${API_BASE_URL}/article/${id}`, {
          headers: getAuthHeaders()
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        setArticle(data);
        setTitle(data.title || '');
        setContent(data.content || '');
      } catch (err) {
        console.error('Error fetching article:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchArticle();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!title.trim() || !content.trim()) {
      setError('Title and content are required');
      return;
    }

    try {
      setSaving(true);
      setError(null);

      const response = await fetch(`${API_BASE_URL}/article/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          title: title.trim(),
          content: content.trim(),
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const updatedArticle = await response.json();
      console.log('Article updated:', updatedArticle);
      
      // Navigate back to article detail page
      navigate(`/article/${id}`);
    } catch (err) {
      console.error('Error updating article:', err);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="article-edit">
        <div className="article-edit__container">
          <div className="article-edit__loading">Loading article...</div>
        </div>
      </div>
    );
  }

  if (error && !article) {
    return (
      <div className="article-edit">
        <div className="article-edit__container">
          <div className="article-edit__error">
            <h2>Error</h2>
            <p>{error}</p>
            <Link to="/" className="article-edit__back-button">
              ← Back to Home
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="article-edit">
      <div className="article-edit__container">
        <div className="article-edit__header">
          <h1>Edit Article #{id}</h1>
          <Link to={`/article/${id}`} className="article-edit__cancel-link">
            Cancel
          </Link>
        </div>

        {error && (
          <div className="article-edit__error-banner">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="article-edit__form">
          <div className="article-edit__field">
            <label htmlFor="title" className="article-edit__label">
              Title
            </label>
            <input
              id="title"
              type="text"
              className="article-edit__input"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Article title"
              required
              disabled={saving}
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
              placeholder="Article content in Markdown format..."
              rows={20}
              required
              disabled={saving}
            />
          </div>

          <div className="article-edit__actions">
            <button
              type="submit"
              className="article-edit__save-button"
              disabled={saving}
            >
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
            <Link
              to={`/article/${id}`}
              className="article-edit__cancel-button"
            >
              Cancel
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
