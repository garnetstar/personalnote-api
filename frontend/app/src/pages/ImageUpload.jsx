import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { getAuthHeaders } from '../utils/api.js';
import './ArticleEdit.css'; // Reusing styles

function ImageUpload() {
  const navigate = useNavigate();
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const handleFileChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
      setError(null);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!file) {
      setError('Please select a file');
      return;
    }

    setUploading(true);
    setError(null);
    setSuccess(null);

    const formData = new FormData();
    formData.append('file', file);

    const headers = getAuthHeaders();
    delete headers['Content-Type']; // Let browser set multipart/form-data with boundary

    try {
      // Use full URL or proxy. Assuming proxy is setup or base URL handling in App.jsx applies here too.
      // App.jsx uses VITE_API_BASE_URL.
      // I should probably import normalizeBaseUrl logic or just rely on relative path if proxy is set up.
      // App.jsx fetches from `${baseUrl}${path}`.
      // I'll try relative path '/api/upload' assuming the proxy redirects /api to backend.
      // Wait, App.jsx uses import.meta.env.VITE_API_BASE_URL.
      
      const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api';
      const cleanBaseUrl = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;

      const response = await fetch(`${cleanBaseUrl}/upload`, {
        method: 'POST',
        headers: headers,
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || errorData.message || 'Failed to upload file');
      }

      const data = await response.json();
      setSuccess(`File uploaded successfully! ID: ${data.fileId}`);
      setFile(null);
      // Reset file input
      document.getElementById('file-upload').value = '';
      
    } catch (err) {
      console.error(err);
      setError(err.message || 'An error occurred during upload');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="article-edit">
      <div className="article-edit__container">
        <header className="article-edit__header">
          <h1>Upload Image to Drive</h1>
          <Link to="/" className="article-edit__cancel-link">
            ← Back to Articles
          </Link>
        </header>

        {error && (
          <div className="article-edit__error-banner" role="alert">
            <strong>Error:</strong> {error}
          </div>
        )}

        {success && (
          <div className="article-edit__success-banner" role="status" style={{
            backgroundColor: '#d4edda', 
            color: '#155724', 
            padding: '1rem', 
            borderRadius: '4px', 
            marginBottom: '1rem'
          }}>
            <strong>Success:</strong> {success}
          </div>
        )}

        <form onSubmit={handleSubmit} className="article-edit__form">
          <div className="article-edit__field">
            <label htmlFor="file-upload" className="article-edit__label">
              Select Image *
            </label>
            <input
              id="file-upload"
              type="file"
              accept="image/*"
              className="article-edit__input"
              onChange={handleFileChange}
              disabled={uploading}
              required
              style={{ padding: '10px' }}
            />
            <p className="help-text" style={{ fontSize: '0.8rem', color: '#666', marginTop: '0.5rem' }}>
              Selected file: {file ? `${file.name} (${(file.size / 1024).toFixed(2)} KB)` : 'None'}
            </p>
          </div>

          <div className="article-edit__actions">
            <button
              type="submit"
              className="article-edit__save-button"
              disabled={uploading || !file}
            >
              {uploading ? 'Uploading...' : 'Upload to Google Drive'}
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

export default ImageUpload;
