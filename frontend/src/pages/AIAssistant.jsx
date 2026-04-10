import React, { useState, useEffect, useRef } from 'react';
import { aiService } from '../services/ai';
import './AIAssistant.css';

function AIAssistant() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [sessionId, setSessionId] = useState('');
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);

  useEffect(() => {
    // 初始化欢迎消息
    setMessages([
      {
        id: 'welcome',
        role: 'assistant',
        content: '你好！我是 ClawOps 数字员工 🤖\n\n我可以帮你：\n• 查询和管理部署实例\n• 创建数据库（MySQL/PostgreSQL）\n• 管理 Docker 容器\n• 查看系统状态和监控\n\n有什么可以帮你的吗？',
        timestamp: new Date(),
      },
    ]);
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSend = async () => {
    if (!input.trim() || loading) return;

    const userMessage = {
      id: Date.now().toString(),
      role: 'user',
      content: input.trim(),
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      const response = await aiService.chat({
        session_id: sessionId,
        message: userMessage.content,
      });

      const assistantMessage = {
        id: response.id || Date.now().toString() + '-a',
        role: 'assistant',
        content: response.content,
        intent: response.intent,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);

      if (response.session_id && !sessionId) {
        setSessionId(response.session_id);
      }
    } catch (error) {
      console.error('Chat error:', error);
      const errorMessage = {
        id: Date.now().toString() + '-e',
        role: 'assistant',
        content: '抱歉，遇到了点问题。请稍后再试。',
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const formatTime = (date) => {
    return new Date(date).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderContent = (content) => {
    // 支持简单的 markdown 格式
    return content
      .split('\n')
      .map((line, i) => {
        // 处理粗体
        line = line.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
        // 处理代码
        line = line.replace(/`(.*?)`/g, '<code>$1</code>');
        return (
          <div key={i} dangerouslySetInnerHTML={{ __html: line || '&nbsp;' }} />
        );
      });
  };

  return (
    <div className="ai-assistant">
      <div className="ai-header">
        <div className="ai-header-info">
          <div className="ai-avatar">🤖</div>
          <div>
            <h2>ClawOps 数字员工</h2>
            <span className="ai-status">
              <span className="status-dot"></span>
              在线
            </span>
          </div>
        </div>
        <div className="ai-header-actions">
          <button className="btn btn-sm btn-outline-secondary" onClick={() => setMessages([messages[0]])}>
            <i className="fas fa-redo"></i> 新对话
          </button>
        </div>
      </div>

      <div className="ai-messages">
        {messages.map((msg) => (
          <div key={msg.id} className={`message ${msg.role}`}>
            <div className="message-avatar">
              {msg.role === 'assistant' ? '🤖' : '👤'}
            </div>
            <div className="message-content">
              <div className="message-text">{renderContent(msg.content)}</div>
              <div className="message-time">{formatTime(msg.timestamp)}</div>
            </div>
          </div>
        ))}

        {loading && (
          <div className="message assistant">
            <div className="message-avatar">🤖</div>
            <div className="message-content">
              <div className="message-text typing">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      <div className="ai-input-area">
        <div className="ai-input-container">
          <textarea
            ref={inputRef}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="输入你的问题或命令..."
            rows="1"
          />
          <button
            className="btn btn-primary send-btn"
            onClick={handleSend}
            disabled={!input.trim() || loading}
          >
            {loading ? (
              <i className="fas fa-spinner fa-spin"></i>
            ) : (
              <i className="fas fa-paper-plane"></i>
            )}
          </button>
        </div>
        <div className="ai-input-hint">
          按 Enter 发送，Shift + Enter 换行
        </div>
      </div>

      <div className="ai-quick-actions">
        <button onClick={() => setInput('查看我的部署实例')}>
          📋 查看部署
        </button>
        <button onClick={() => setInput('创建一个 MySQL 数据库')}>
          🗄️ 创建数据库
        </button>
        <button onClick={() => setInput('部署一个 Nginx 容器')}>
          🐳 创建容器
        </button>
        <button onClick={() => setInput('系统状态怎么样')}>
          📊 系统状态
        </button>
      </div>
    </div>
  );
}

export default AIAssistant;
