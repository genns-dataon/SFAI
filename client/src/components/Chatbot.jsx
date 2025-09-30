import { useState, useRef, useEffect } from 'react';
import { FloatButton, Card, Input, Button, Avatar, Space, Typography, Spin, Modal, message as antMessage, Tooltip } from 'antd';
import { 
  MessageOutlined, 
  SendOutlined, 
  CloseOutlined, 
  RobotOutlined, 
  UserOutlined,
  LikeOutlined,
  DislikeOutlined,
  LikeFilled,
  DislikeFilled
} from '@ant-design/icons';
import { chatAPI, feedbackAPI } from '../api/api';

const { Text } = Typography;
const { TextArea } = Input;

const Chatbot = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [feedbackModal, setFeedbackModal] = useState({ visible: false, messageIdx: null, rating: null });
  const [feedbackComment, setFeedbackComment] = useState('');
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;

    const userMessage = { role: 'user', content: input, question: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      const response = await chatAPI.sendMessage(input);
      const botMessage = { 
        role: 'bot', 
        content: response.data.response,
        question: input,
        feedback: null 
      };
      setMessages((prev) => [...prev, botMessage]);
    } catch (error) {
      const errorMessage = { role: 'bot', content: 'Sorry, something went wrong. Please try again.', question: input, feedback: null };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  const handleFeedback = async (messageIdx, rating) => {
    const message = messages[messageIdx];
    
    if (rating === 'positive') {
      // Thumbs up - immediately save feedback
      try {
        await feedbackAPI.create({
          question: message.question,
          response: message.content,
          rating: 'positive',
          comment: ''
        });
        
        // Update message feedback state
        setMessages(prev => prev.map((msg, idx) => 
          idx === messageIdx ? { ...msg, feedback: 'positive' } : msg
        ));
        
        antMessage.success('Thanks for your feedback!');
      } catch (error) {
        antMessage.error('Failed to save feedback');
      }
    } else {
      // Thumbs down - show modal for comment
      setFeedbackModal({ visible: true, messageIdx, rating: 'negative' });
    }
  };

  const submitNegativeFeedback = async () => {
    const { messageIdx } = feedbackModal;
    const message = messages[messageIdx];

    try {
      await feedbackAPI.create({
        question: message.question,
        response: message.content,
        rating: 'negative',
        comment: feedbackComment
      });
      
      // Update message feedback state
      setMessages(prev => prev.map((msg, idx) => 
        idx === messageIdx ? { ...msg, feedback: 'negative' } : msg
      ));
      
      antMessage.success('Thanks for your feedback! We\'ll improve based on your input.');
      setFeedbackModal({ visible: false, messageIdx: null, rating: null });
      setFeedbackComment('');
    } catch (error) {
      antMessage.error('Failed to save feedback');
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
                  flexDirection: 'column',
                  alignItems: msg.role === 'user' ? 'flex-end' : 'flex-start',
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
                
                {msg.role === 'bot' && (
                  <div style={{ marginTop: 4, marginLeft: 32 }}>
                    <Space size="small">
                      <Tooltip title="Helpful">
                        <Button
                          type="text"
                          size="small"
                          icon={msg.feedback === 'positive' ? <LikeFilled /> : <LikeOutlined />}
                          onClick={() => !msg.feedback && handleFeedback(idx, 'positive')}
                          disabled={msg.feedback !== null}
                          style={{ 
                            color: msg.feedback === 'positive' ? '#52c41a' : undefined,
                            fontSize: '12px'
                          }}
                        />
                      </Tooltip>
                      <Tooltip title="Not helpful">
                        <Button
                          type="text"
                          size="small"
                          icon={msg.feedback === 'negative' ? <DislikeFilled /> : <DislikeOutlined />}
                          onClick={() => !msg.feedback && handleFeedback(idx, 'negative')}
                          disabled={msg.feedback !== null}
                          style={{ 
                            color: msg.feedback === 'negative' ? '#ff4d4f' : undefined,
                            fontSize: '12px'
                          }}
                        />
                      </Tooltip>
                    </Space>
                  </div>
                )}
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

      <Modal
        title="Help us improve"
        open={feedbackModal.visible}
        onOk={submitNegativeFeedback}
        onCancel={() => {
          setFeedbackModal({ visible: false, messageIdx: null, rating: null });
          setFeedbackComment('');
        }}
        okText="Submit"
        cancelText="Cancel"
      >
        <div style={{ marginTop: 16 }}>
          <Text>What went wrong? Your feedback helps us improve.</Text>
          <TextArea
            rows={4}
            value={feedbackComment}
            onChange={(e) => setFeedbackComment(e.target.value)}
            placeholder="Tell us what went wrong or what you expected..."
            style={{ marginTop: 8 }}
          />
        </div>
      </Modal>
    </>
  );
};

export default Chatbot;
