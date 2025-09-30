import { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import axios from 'axios';

const API_URL = window.location.hostname.includes('replit.dev')
  ? window.location.origin.replace('-5000', '-8080')
  : 'http://localhost:8080';

const ProtectedRoute = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const token = localStorage.getItem('token');

  useEffect(() => {
    const verifyToken = async () => {
      if (!token) {
        setIsAuthenticated(false);
        return;
      }

      try {
        await axios.get(`${API_URL}/api/me`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        setIsAuthenticated(true);
      } catch (error) {
        console.error('Token verification failed:', error);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setIsAuthenticated(false);
      }
    };

    verifyToken();
  }, [token]);

  if (isAuthenticated === null) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Spin size="large" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default ProtectedRoute;
