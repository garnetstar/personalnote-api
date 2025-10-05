import React, { useState, useEffect } from 'react';
import './App.css';

// Types for our API data
interface Article {
  id: number;
  title: string;
  content: string;
  category: string;
  author: string;
}

function App() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [filterCategory, setFilterCategory] = useState('');
  const [newUser, setNewUser] = useState({ name: '', id: 0 });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const API_BASE = 'http://localhost:8080';

  // Fetch all articles
  const fetchArticles = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/articles`);
      if (response.ok) {
        const data = await response.json();
        setArticles(data.articles || []); // Access the articles array from the response
      } else {
        setMessage('Failed to fetch articles');
      }
    } catch (error) {
      setMessage('Error connecting to API');
    }
    setLoading(false);
  };

  // Fetch article by ID
  const fetchArticleById = async (id: number) => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/article/${id}`);
      if (response.ok) {
        const data = await response.json();
        setSelectedArticle(data);
        setMessage(`Fetched article: ${data.title}`);
      } else {
        setMessage('Article not found');
      }
    } catch (error) {
      setMessage('Error fetching article');
    }
    setLoading(false);
  };

  // Filter articles by category
  const filterArticles = async () => {
    if (!filterCategory.trim()) {
      setMessage('Please enter a category to filter');
      return;
    }
    
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/article/filter/all/${filterCategory}`);
      if (response.ok) {
        const data = await response.json();
        setArticles(data.articles || []); // Access the articles array from the response
        setMessage(`Filtered by category: ${filterCategory}`);
      } else {
        setMessage('No articles found for this category');
      }
    } catch (error) {
      setMessage('Error filtering articles');
    }
    setLoading(false);
  };

  // Create user
  const createUser = async () => {
    if (!newUser.name.trim() || newUser.id <= 0) {
      setMessage('Please enter valid user name and ID');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/user`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newUser),
      });
      
      if (response.ok) {
        const data = await response.json();
        setMessage(`User created: ${data.name} (ID: ${data.id})`);
        setNewUser({ name: '', id: 0 });
      } else {
        setMessage('Failed to create user');
      }
    } catch (error) {
      setMessage('Error creating user');
    }
    setLoading(false);
  };

  // Test API health
  const testHealth = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/`);
      if (response.ok) {
        const text = await response.text();
        setMessage(`API Health: ${text}`);
      } else {
        setMessage('API health check failed');
      }
    } catch (error) {
      setMessage('API is not responding');
    }
    setLoading(false);
  };

  // Load articles on component mount
  useEffect(() => {
    fetchArticles();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Simple Go API Frontend</h1>
        <p>React UI for testing Go API endpoints</p>
      </header>

      <main className="App-main">
        {/* Status Message */}
        {message && (
          <div className="message">
            {message}
          </div>
        )}

        {loading && <div className="loading">Loading...</div>}

        {/* API Health Check */}
        <section className="section">
          <h2>API Health</h2>
          <button onClick={testHealth} disabled={loading}>
            Test API Health
          </button>
        </section>

        {/* Articles Display */}
        <section className="section">
          <h2>Articles ({articles.length})</h2>
          <div className="actions">
            <button onClick={fetchArticles} disabled={loading}>
              Refresh Articles
            </button>
            <div className="filter-group">
              <input
                type="text"
                placeholder="Filter by category"
                value={filterCategory}
                onChange={(e) => setFilterCategory(e.target.value)}
              />
              <button onClick={filterArticles} disabled={loading}>
                Filter
              </button>
            </div>
          </div>
          
          <div className="articles-grid">
            {articles.map((article) => (
              <div key={article.id} className="article-card">
                <h3>{article.title}</h3>
                <p><strong>Category:</strong> {article.category}</p>
                <p><strong>Author:</strong> {article.author}</p>
                <p className="content">{article.content}</p>
                <button 
                  onClick={() => fetchArticleById(article.id)}
                  disabled={loading}
                >
                  View Details
                </button>
              </div>
            ))}
          </div>
        </section>

        {/* Selected Article Detail */}
        {selectedArticle && (
          <section className="section">
            <h2>Article Details</h2>
            <div className="article-detail">
              <h3>{selectedArticle.title}</h3>
              <p><strong>ID:</strong> {selectedArticle.id}</p>
              <p><strong>Category:</strong> {selectedArticle.category}</p>
              <p><strong>Author:</strong> {selectedArticle.author}</p>
              <p className="content">{selectedArticle.content}</p>
              <button onClick={() => setSelectedArticle(null)}>
                Close Details
              </button>
            </div>
          </section>
        )}

        {/* Create User */}
        <section className="section">
          <h2>Create User</h2>
          <div className="form-group">
            <input
              type="text"
              placeholder="User name"
              value={newUser.name}
              onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
            />
            <input
              type="number"
              placeholder="User ID"
              value={newUser.id || ''}
              onChange={(e) => setNewUser({ ...newUser, id: parseInt(e.target.value) || 0 })}
            />
            <button onClick={createUser} disabled={loading}>
              Create User
            </button>
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;
