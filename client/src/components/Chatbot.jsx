import { useState } from 'react';
import { MessageCircle, X, Send } from 'lucide-react';
import { chatAPI } from '../api/api';

const Chatbot = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);

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
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="fixed bottom-6 right-6 p-4 rounded-full shadow-lg hover:shadow-xl transition-all duration-200 z-50 hover:scale-110"
          style={{ backgroundColor: '#2563eb', color: 'white' }}
        >
          <MessageCircle className="w-6 h-6" />
        </button>
      )}

      {isOpen && (
        <div className="fixed bottom-6 right-6 w-96 h-[500px] rounded-xl shadow-2xl border flex flex-col z-50" style={{ backgroundColor: 'white', borderColor: '#e2e8f0' }}>
          <div className="text-white p-4 rounded-t-xl flex justify-between items-center" style={{ background: 'linear-gradient(to right, #2563eb, #1d4ed8)' }}>
            <div>
              <h3 className="font-semibold text-lg">HR Assistant</h3>
              <p className="text-xs text-primary-100">Always here to help</p>
            </div>
            <button 
              onClick={() => setIsOpen(false)} 
              className="hover:bg-white/20 p-2 rounded-lg transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3" style={{ backgroundColor: '#f8fafc' }}>
            {messages.length === 0 && (
              <div className="text-center mt-12" style={{ color: '#64748b' }}>
                <div className="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4" style={{ backgroundColor: '#dbeafe' }}>
                  <MessageCircle className="w-8 h-8" style={{ color: '#2563eb' }} />
                </div>
                <p className="text-sm font-medium mb-2">Welcome to HR Assistant!</p>
                <p className="text-xs">Ask me anything about HR policies, leave requests, or payroll</p>
              </div>
            )}
            {messages.map((msg, idx) => (
              <div
                key={idx}
                className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-[80%] p-3 rounded-lg shadow-sm ${
                    msg.role === 'user'
                      ? 'text-white rounded-br-none'
                      : 'rounded-bl-none border'
                  }`}
                  style={
                    msg.role === 'user'
                      ? { backgroundColor: '#2563eb', color: 'white' }
                      : { backgroundColor: 'white', color: '#0f172a', borderColor: '#e2e8f0' }
                  }
                >
                  <p className="text-sm">{msg.content}</p>
                </div>
              </div>
            ))}
            {loading && (
              <div className="flex justify-start">
                <div className="border p-3 rounded-lg rounded-bl-none shadow-sm" style={{ backgroundColor: 'white', borderColor: '#e2e8f0' }}>
                  <div className="flex space-x-1">
                    <div className="w-2 h-2 rounded-full animate-bounce" style={{ backgroundColor: '#2563eb' }}></div>
                    <div className="w-2 h-2 rounded-full animate-bounce" style={{ backgroundColor: '#2563eb', animationDelay: '0.2s' }}></div>
                    <div className="w-2 h-2 rounded-full animate-bounce" style={{ backgroundColor: '#2563eb', animationDelay: '0.4s' }}></div>
                  </div>
                </div>
              </div>
            )}
          </div>

          <div className="p-4 border-t rounded-b-xl" style={{ borderColor: '#e2e8f0', backgroundColor: 'white' }}>
            <div className="flex gap-2">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && !loading && handleSend()}
                placeholder="Type your message..."
                className="input-field text-sm"
                disabled={loading}
              />
              <button
                onClick={handleSend}
                disabled={loading || !input.trim()}
                className="btn btn-primary px-3 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Send className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default Chatbot;
