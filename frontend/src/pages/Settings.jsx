import React from 'react';
import './Settings.css';

function Settings() {
  return (
    <div className="settings-container">
      <h1>⚙️ 系统设置</h1>
      
      <div className="settings-section">
        <h2>👤 账户信息</h2>
        <div className="setting-item">
          <label>用户名</label>
          <input type="text" defaultValue="admin" disabled />
        </div>
        <div className="setting-item">
          <label>邮箱</label>
          <input type="email" defaultValue="admin@openclaw.cn" disabled />
        </div>
      </div>

      <div className="settings-section">
        <h2>🔔 通知设置</h2>
        <div className="setting-item">
          <label>
            <input type="checkbox" defaultChecked /> 启用钉钉通知
          </label>
        </div>
        <div className="setting-item">
          <label>
            <input type="checkbox" defaultChecked /> 启用邮件通知
          </label>
        </div>
      </div>

      <div className="settings-section">
        <h2>🤖 AI 设置</h2>
        <div className="setting-item">
          <label>AI Provider</label>
          <select defaultValue="openai">
            <option value="openai">OpenAI</option>
            <option value="ollama">Ollama (本地)</option>
          </select>
        </div>
        <div className="setting-item">
          <label>模型</label>
          <input type="text" defaultValue="gpt-4o-mini" placeholder="输入模型名称" />
        </div>
      </div>

      <div className="settings-actions">
        <button className="btn-save">保存设置</button>
      </div>
    </div>
  );
}

export default Settings;
