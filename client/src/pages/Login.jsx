import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import axios from 'axios';

const { Title, Paragraph } = Typography;

const API_URL = window.location.hostname.includes('replit.dev')
  ? window.location.origin.replace('-5000', '-8080')
  : 'http://localhost:8080';

const Login = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const response = await axios.post(`${API_URL}/api/auth/login`, {
        username: values.username,
        password: values.password,
      });

      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      
      message.success('Login successful!');
      navigate('/');
    } catch (error) {
      console.error('Login error:', error);
      message.error(error.response?.data?.error || 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
      <Card className="w-full max-w-md shadow-lg">
        <div className="text-center mb-8">
          <Title level={2} className="mb-2">Welcome Back</Title>
          <Paragraph className="text-gray-600">
            Sign in to access the HCM system
          </Paragraph>
        </div>

        <Form
          name="login"
          onFinish={onFinish}
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: 'Please enter your username' }
            ]}
          >
            <Input
              prefix={<UserOutlined className="text-gray-400" />}
              placeholder="Username"
              autoComplete="username"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: 'Please enter your password' }
            ]}
          >
            <Input.Password
              prefix={<LockOutlined className="text-gray-400" />}
              placeholder="Password"
              autoComplete="current-password"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              className="w-full"
              loading={loading}
            >
              Sign In
            </Button>
          </Form.Item>
        </Form>

        <div className="text-center mt-4 text-gray-600">
          <Paragraph className="text-sm">
            Test accounts: alice, bob, carol, etc. | Password: password
          </Paragraph>
        </div>
      </Card>
    </div>
  );
};

export default Login;
