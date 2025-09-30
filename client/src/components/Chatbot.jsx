import { useState, useRef, useEffect } from 'react';
import { FloatButton, Card, Input, Button, Avatar, Space, Typography, Spin } from 'antd';
import { 
  MessageOutlined, 
  SendOutlined, 
  CloseOutlined, 
  RobotOutlined, 
  UserOutlined 
} from '@ant-design/icons';
import { chatAPI } from '../api/api';

const { Text } = Typography;

const Chatbot = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;

    const userMessage = { role: 'user', content: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      const response = await chatAPI.sendMessage(input);
      const botMessage = { role: 'bot', content: response.data.response };
      setMessages((prev) => [...prev, botMessage]);
    } catch (error) {
      const errorMessage = { role: 'bot', content: 'Sorry, something went wrong. Please try again.' };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <FloatButton
        icon={<MessageOutlined />}
        type="primary"
        style={{ right: 24, bottom: 24 }}
        onClick={() => setIsOpen(!isOpen)}
      />

      {isOpen && (
        <Card
          title={
            <Space>
              <RobotOutlined style={{ fontSize: '20px' }} />
              <div>
                <div style={{ fontWeight: 600 }}>HR Assistant</div>
                <Text type="secondary" style={{ fontSize: '12px' }}>Always here to help</Text>
              </div>
            </Space>
          }
          extra={
            <Button 
              type="text" 
              icon={<CloseOutlined />} 
              onClick={() => setIsOpen(false)}
            />
          }
          style={{
            position: 'fixed',
            bottom: 80,
            right: 24,
            width: 400,
            height: 600,
            zIndex: 1000,
            display: 'flex',
            flexDirection: 'column'
          }}
          styles={{
            body: {
              padding: 0,
              display: 'flex',
              flexDirection: 'column',
              height: '100%',
              overflow: 'hidden'
            }
          }}
        >
          <div style={{ 
            flex: 1, 
            overflowY: 'auto', 
            padding: '16px',
            backgroundColor: '#f5f5f5'
          }}>
            {messages.length === 0 && (
              <div style={{ textAlign: 'center', marginTop: 60 }}>
                <Avatar 
                  size={64} 
                  icon={<RobotOutlined />} 
                  style={{ backgroundColor: '#1890ff', marginBottom: 16 }} 
                />
                <div>
                  <Text strong>Welcome to HR Assistant!</Text>
                </div>
                <Text type="secondary" style={{ fontSize: '12px', display: 'block', marginTop: 8 }}>
                  Ask me anything about HR policies, employees, leave requests, or payroll
                </Text>
              </div>
            )}
            
            {messages.map((msg, idx) => (
              <div
                key={idx}
                style={{
                  display: 'flex',
                  justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                  marginBottom: 12
                }}
              >
                <Space direction="horizontal" align="start">
                  {msg.role === 'bot' && (
                    <Avatar 
                      size="small" 
                      icon={<RobotOutlined />} 
                      style={{ backgroundColor: '#1890ff' }} 
                    />
                  )}
                  <div
                    style={{
                      maxWidth: '280px',
                      padding: '8px 12px',
                      borderRadius: '8px',
                      backgroundColor: msg.role === 'user' ? '#1890ff' : '#fff',
                      color: msg.role === 'user' ? '#fff' : '#000',
                      boxShadow: '0 1px 2px rgba(0,0,0,0.1)',
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-word'
                    }}
                  >
                    <Text style={{ color: msg.role === 'user' ? '#fff' : '#000', fontSize: '14px' }}>
                      {msg.content}
                    </Text>
                  </div>
                  {msg.role === 'user' && (
                    <Avatar 
                      size="small" 
                      icon={<UserOutlined />} 
                      style={{ backgroundColor: '#87d068' }} 
                    />
                  )}
                </Space>
              </div>
            ))}
            
            {loading && (
              <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: 12 }}>
                <Space>
                  <Avatar 
                    size="small" 
                    icon={<RobotOutlined />} 
                    style={{ backgroundColor: '#1890ff' }} 
                  />
                  <div style={{
                    padding: '8px 12px',
                    borderRadius: '8px',
                    backgroundColor: '#fff',
                    boxShadow: '0 1px 2px rgba(0,0,0,0.1)'
                  }}>
                    <Spin size="small" />
                  </div>
                </Space>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>

          <div style={{ 
            padding: '12px 16px', 
            borderTop: '1px solid #f0f0f0',
            backgroundColor: '#fff'
          }}>
            <Space.Compact style={{ width: '100%' }}>
              <Input
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onPressEnter={() => !loading && handleSend()}
                placeholder="Type your message..."
                disabled={loading}
                style={{ flex: 1 }}
              />
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleSend}
                disabled={loading || !input.trim()}
                loading={loading}
              />
            </Space.Compact>
          </div>
        </Card>
      )}
    </>
  );
};

export default Chatbot;
