import React, { useState } from 'react';
import { userService } from '../services/deploy';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    
    try {
      const user = await userService.login(username, password);
      onLogin(user);
    } catch (err) {
      setError('用户名或密码错误');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1><i className="fas fa-rocket"></i> ClawOps</h1>
        <p className="subtitle">AI DevOps Platform</p>
        
        {error && <div className="alert alert-error">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>用户名</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="admin"
              required
            />
          </div>
          <div className="form-group">
            <label>密码</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>
          <button type="submit" className="btn-login" disabled={loading}>
            {loading ? '登录中...' : '登录'}
          </button>
        </form>
        
        <div className="login-hint">
          <small>默认账号: admin / admin123</small>
        </div>
      </div>
      
      <style>{`
        .login-container {
          min-height: 100vh;
          display: flex;
          align-items: center;
          justify-content: center;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .login-card {
          background: white;
          padding: 3rem;
          border-radius: 16px;
          box-shadow: 0 20px 60px rgba(0,0,0,0.3);
          width: 100%;
          max-width: 400px;
          text-align: center;
        }
        .login-card h1 {
          color: #667eea;
          margin-bottom: 0.5rem;
        }
        .subtitle {
          color: #666;
          margin-bottom: 2rem;
        }
        .form-group {
          margin-bottom: 1.5rem;
          text-align: left;
        }
        .form-group label {
          display: block;
          margin-bottom: 0.5rem;
          color: #333;
          font-weight: 500;
        }
        .form-group input {
          width: 100%;
          padding: 0.75rem 1rem;
          border: 2px solid #e1e1e1;
          border-radius: 8px;
          font-size: 1rem;
          transition: border-color 0.3s;
        }
        .form-group input:focus {
          outline: none;
          border-color: #667eea;
        }
        .btn-login {
          width: 100%;
          padding: 1rem;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
          border: none;
          border-radius: 8px;
          font-size: 1rem;
          font-weight: 600;
          cursor: pointer;
          transition: transform 0.2s;
        }
        .btn-login:hover {
          transform: translateY(-2px);
        }
        .btn-login:disabled {
          opacity: 0.7;
          cursor: not-allowed;
        }
        .alert-error {
          background: #fee;
          color: #c00;
          padding: 0.75rem;
          border-radius: 8px;
          margin-bottom: 1rem;
        }
        .login-hint {
          margin-top: 1.5rem;
          color: #999;
        }
      `}</style>
    </div>
  );
}

export default Login;
