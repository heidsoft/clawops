import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import './Dashboard.css';

function Dashboard() {
  const [stats, setStats] = useState({
    deployments: 0,
    databases: 0,
    containers: 0,
    alerts: 0,
  });
  const [recentActivity, setRecentActivity] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [deployRes, dbRes, dockerRes, alertRes] = await Promise.all([
        deployService.get('/deployments').catch(() => ({ data: [] })),
        deployService.get('/databases').catch(() => ({ data: [] })),
        deployService.get('/docker').catch(() => ({ data: [] })),
        deployService.get('/monitor/alerts').catch(() => ({ alerts: [] })),
      ]);

      setStats({
        deployments: deployRes.data?.length || deployRes.total || 0,
        databases: dbRes.data?.length || 0,
        containers: dockerRes.data?.length || 0,
        alerts: (alertRes.alerts || []).filter(a => a.status === 'firing').length,
      });

      // 模拟最近活动
      setRecentActivity([
        { time: '10:30', action: '部署实例', target: 'prod-api', status: 'success' },
        { time: '09:45', action: '创建数据库', target: 'mysql-prod', status: 'success' },
        { time: '08:20', action: '启动容器', target: 'nginx-web', status: 'success' },
        { time: '昨天', action: '磁盘告警', target: 'prod-002', status: 'warning' },
      ]);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="dashboard-loading">加载中...</div>;
  }

  return (
    <div className="dashboard-container">
      <h1>📊 概览</h1>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-icon">🖥️</div>
          <div className="stat-info">
            <div className="stat-value">{stats.deployments}</div>
            <div className="stat-label">部署实例</div>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">🗄️</div>
          <div className="stat-info">
            <div className="stat-value">{stats.databases}</div>
            <div className="stat-label">数据库</div>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon">📦</div>
          <div className="stat-info">
            <div className="stat-value">{stats.containers}</div>
            <div className="stat-label">Docker容器</div>
          </div>
        </div>

        <div className="stat-card stat-alert">
          <div className="stat-icon">🚨</div>
          <div className="stat-info">
            <div className="stat-value">{stats.alerts}</div>
            <div className="stat-label">活跃告警</div>
          </div>
        </div>
      </div>

      <div className="dashboard-row">
        <div className="panel recent-activity">
          <h2>📋 最近活动</h2>
          <div className="activity-list">
            {recentActivity.map((item, index) => (
              <div key={index} className="activity-item">
                <span className="activity-time">{item.time}</span>
                <span className="activity-action">{item.action}</span>
                <span className="activity-target">{item.target}</span>
                <span className={`activity-status ${item.status}`}>
                  {item.status === 'success' ? '✅' : '⚠️'}
                </span>
              </div>
            ))}
          </div>
        </div>

        <div className="panel quick-actions">
          <h2>⚡ 快捷操作</h2>
          <div className="action-buttons">
            <button className="action-btn" onClick={() => window.location.href = '/ai'}>
              🤖 AI 对话
            </button>
            <button className="action-btn" onClick={() => window.location.href = '/deployments'}>
              ➕ 创建部署
            </button>
            <button className="action-btn" onClick={() => window.location.href = '/monitoring'}>
              📊 查看监控
            </button>
            <button className="action-btn" onClick={() => window.location.href = '/docker'}>
              🐳 创建容器
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
