import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function AuthCallback() {
  const [searchParams] = useSearchParams();
  const nav = useNavigate();
  const { login } = useAuth();

  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      login(token);
      nav('/');
    } else {
      nav('/');
    }
  }, [searchParams, login, nav]);

  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <h2>Logging you in...</h2>
      <p>Please wait while we complete the authentication process.</p>
    </div>
  );
}
