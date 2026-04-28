import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import './Notifications.css';

function Notifications() {
  const [settings, setSettings] = useState({
    dingtalk: {
      enabled: true,
      webhook: '',
      secret: '',
      atMobiles: [],
      isAtAll: false,
    },
    email: {
      enabled: false,
      smtpHost: '',
      smtpPort: 587,
      from: '',
      to: [],
    },
    webhook: {
      enabled: false,
      url: '',
      method: 'POST',
    },
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [testStatus, setTestStatus] = useState('');

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    try {
      const res = await deployService.get('/notifications/settings');
      if (res.settings) {
        setSettings(res.settings);
      }
    } catch (error) {
      console.error('Failed to load settings:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await deployService.put('/notifications/settings', settings);
      alert('设置已保存');
    } catch (error) {
      alert('保存失败：' + error.message);
    } finally {
      setSaving(false);
    }
  };

  const handleTest = async (type) => {
    setTestStatus('sending');
    try {
      await deployService.post('/notifications/test', {
        type,
        webhook: type === 'dingtalk' ? settings.dingtalk.webhook : undefined,
      });
      setTestStatus('success');
      setTimeout(() => setTestStatus(''), 3000);
    } catch (error) {
      setTestStatus('error');
      alert('测试失败：' + error.message);
    }
  };

  const handleAddMobile = (e) => {
    if (e.key === 'Enter' && e.target.value) {
      setSettings({
        ...settings,
        dingtalk: {
          ...settings.dingtalk,
          atMobiles: [...settings.dingtalk.atMobiles, e.target.value],
        },
      });
      e.target.value = '';
    }
  };

  const handleRemoveMobile = (mobile) => {
    setSettings({
      ...settings,
      dingtalk: {
        ...settings.dingtalk,
        atMobiles: settings.dingtalk.atMobiles.filter((m) => m !== mobile),
      },
    });
  };

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  return (
    <div className="notifications-container">
      <h1>🔔 通知设置</h1>

      {/* 钉钉通知 */}
      <div className="settings-section">
        <div className="section-header">
          <div className="section-title">
            <span className="icon">💬</span>
            <h2>钉钉通知</h2>
          </div>
          <label className="toggle">
            <input
              type="checkbox"
              checked={settings.dingtalk.enabled}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  dingtalk: { ...settings.dingtalk, enabled: e.target.checked },
                })
              }
            />
            <span className="toggle-slider"></span>
          </label>
        </div>

        {settings.dingtalk.enabled && (
          <div className="section-content">
            <div className="form-group">
              <label>Webhook 地址</label>
              <input
                type="text"
                placeholder="https://oapi.dingtalk.com/robot/send?access_token=xxx"
                value={settings.dingtalk.webhook}
                onChange={(e) =>
                  setSettings({
                    ...settings,
                    dingtalk: { ...settings.dingtalk, webhook: e.target.value },
                  })
                }
              />
              <p className="hint">
                在钉钉群中添加自定义机器人，复制 Webhook 地址
              </p>
            </div>

            <div className="form-group">
              <label>加签密钥（可选）</label>
              <input
                type="password"
                placeholder="SEC开头的密钥"
                value={settings.dingtalk.secret}
                onChange={(e) =>
                  setSettings({
                    ...settings,
                    dingtalk: { ...settings.dingtalk, secret: e.target.value },
                  })
                }
              />
            </div>

            <div className="form-group">
              <label>@ 指定人手机号（回车添加）</label>
              <div className="tags-input">
                {settings.dingtalk.atMobiles.map((mobile) => (
                  <span key={mobile} className="tag">
                    {mobile}
                    <button onClick={() => handleRemoveMobile(mobile)}>×</button>
                  </span>
                ))}
                <input
                  type="text"
                  placeholder="输入手机号按回车"
                  onKeyPress={handleAddMobile}
                />
              </div>
            </div>

            <div className="form-group">
              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={settings.dingtalk.isAtAll}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      dingtalk: { ...settings.dingtalk, isAtAll: e.target.checked },
                    })
                  }
                />
                @ 所有人
              </label>
            </div>

            <button
              className={`btn-test ${testStatus}`}
              onClick={() => handleTest('dingtalk')}
              disabled={!settings.dingtalk.webhook}
            >
              {testStatus === 'sending' ? '发送中...' : testStatus === 'success' ? '✅ 发送成功' : '发送测试消息'}
            </button>
          </div>
        )}
      </div>

      {/* 邮件通知 */}
      <div className="settings-section">
        <div className="section-header">
          <div className="section-title">
            <span className="icon">📧</span>
            <h2>邮件通知</h2>
          </div>
          <label className="toggle">
            <input
              type="checkbox"
              checked={settings.email.enabled}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  email: { ...settings.email, enabled: e.target.checked },
                })
              }
            />
            <span className="toggle-slider"></span>
          </label>
        </div>

        {settings.email.enabled && (
          <div className="section-content">
            <div className="form-row">
              <div className="form-group">
                <label>SMTP 服务器</label>
                <input
                  type="text"
                  placeholder="smtp.example.com"
                  value={settings.email.smtpHost}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      email: { ...settings.email, smtpHost: e.target.value },
                    })
                  }
                />
              </div>
              <div className="form-group" style={{ width: '120px' }}>
                <label>端口</label>
                <input
                  type="number"
                  placeholder="587"
                  value={settings.email.smtpPort}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      email: { ...settings.email, smtpPort: parseInt(e.target.value) },
                    })
                  }
                />
              </div>
            </div>

            <div className="form-group">
              <label>发件人</label>
              <input
                type="email"
                placeholder="alert@example.com"
                value={settings.email.from}
                onChange={(e) =>
                  setSettings({
                    ...settings,
                    email: { ...settings.email, from: e.target.value },
                  })
                }
              />
            </div>
          </div>
        )}
      </div>

      {/* Webhook 通知 */}
      <div className="settings-section">
        <div className="section-header">
          <div className="section-title">
            <span className="icon">🔗</span>
            <h2>Webhook 通知</h2>
          </div>
          <label className="toggle">
            <input
              type="checkbox"
              checked={settings.webhook.enabled}
              onChange={(e) =>
                setSettings({
                  ...settings,
                  webhook: { ...settings.webhook, enabled: e.target.checked },
                })
              }
            />
            <span className="toggle-slider"></span>
          </label>
        </div>

        {settings.webhook.enabled && (
          <div className="section-content">
            <div className="form-group">
              <label>Webhook URL</label>
              <input
                type="text"
                placeholder="https://your-webhook-endpoint.com/alert"
                value={settings.webhook.url}
                onChange={(e) =>
                  setSettings({
                    ...settings,
                    webhook: { ...settings.webhook, url: e.target.value },
                  })
                }
              />
            </div>

            <div className="form-group">
              <label>请求方法</label>
              <select
                value={settings.webhook.method}
                onChange={(e) =>
                  setSettings({
                    ...settings,
                    webhook: { ...settings.webhook, method: e.target.value },
                  })
                }
              >
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
              </select>
            </div>
          </div>
        )}
      </div>

      <div className="actions">
        <button className="btn-save" onClick={handleSave} disabled={saving}>
          {saving ? '保存中...' : '💾 保存设置'}
        </button>
      </div>
    </div>
  );
}

export default Notifications;
