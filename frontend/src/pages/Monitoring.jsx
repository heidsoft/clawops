import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import { monitorService } from '../services/deploy';
import './Monitoring.css';

function Monitoring() {
  const [overview, setOverview] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [stats, setStats] = useState({});
  const [loading, setLoading] = useState(true);
  const [selectedInstance, setSelectedInstance] = useState(null);
  const [instanceMetrics, setInstanceMetrics] = useState(null);

  useEffect(() => {
    loadData();
    // 定时刷新
    const interval = setInterval(loadData, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      // 并行加载所有数据
      const [overviewRes, alertsRes, statsRes] = await Promise.all([
        deployService.get('/monitor/overview'),
        deployService.get('/monitor/alerts'),
        deployService.get('/monitor/stats'),
      ]);

      if (overviewRes.overview) setOverview(overviewRes.overview);
      if (overviewRes.data) setOverview(overviewRes.data);
      if (alertsRes.alerts) setAlerts(alertsRes.alerts);
      if (alertsRes.data) setAlerts(alertsRes.data);
      if (statsRes.data) setStats(statsRes.data);
      if (statsRes) setStats(statsRes);
    } catch (error) {
      console.error('Failed to load monitoring data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleInstanceClick = async (instance) => {
    setSelectedInstance(instance);
    try {
      const res = await deployService.get(`/deployments/${instance.instance_id}/metrics`);
      setInstanceMetrics(res);
    } catch (error) {
      console.error('Failed to load metrics:', error);
    }
  };

  const getSeverityClass = (severity) => {
    switch (severity) {
      case 'critical': return 'severity-critical';
      case 'warning': return 'severity-warning';
      default: return 'severity-info';
    }
  };

  const getStatusClass = (status) => {
    switch (status) {
      case 'firing': return 'status-firing';
      case 'acknowledged': return 'status-acknowledged';
      case 'resolved': return 'status-resolved';
      default: return '';
    }
  };

  const formatTime = (time) => {
    if (!time) return '-';
    const d = new Date(time);
    return d.toLocaleString('zh-CN');
  };

  if (loading) {
    return <div className="monitoring-loading">加载中...</div>;
  }

  return (
    <div className="monitoring-container">
      <div className="monitoring-header">
        <h1>📊 监控中心</h1>
        <button onClick={loadData} className="btn-refresh">🔄 刷新</button>
      </div>

      {/* 统计卡片 */}
      <div className="stats-cards">
        <div className="stat-card">
          <div className="stat-value">{stats.total || 0}</div>
          <div className="stat-label">总告警数</div>
        </div>
        <div className="stat-card stat-firing">
          <div className="stat-value">{stats.firing || 0}</div>
          <div className="stat-label">🔥 进行中</div>
        </div>
        <div className="stat-card stat-ack">
          <div className="stat-value">{stats.acknowledged || 0}</div>
          <div className="stat-label">👀 已确认</div>
        </div>
        <div className="stat-card stat-resolved">
          <div className="stat-value">{stats.resolved || 0}</div>
          <div className="stat-label">✅ 已解决</div>
        </div>
      </div>

      <div className="monitoring-content">
        {/* 左侧：实例概览 */}
        <div className="instances-panel">
          <h2>�instance 实例状态</h2>
          <div className="instances-list">
            {overview.map((inst) => (
              <div
                key={inst.instance_id}
                className={`instance-card ${selectedInstance?.instance_id === inst.instance_id ? 'selected' : ''}`}
                onClick={() => handleInstanceClick(inst)}
              >
                <div className="instance-header">
                  <span className="instance-name">{inst.instance_name}</span>
                  <span className={`instance-status ${inst.status}`}>{inst.status}</span>
                </div>
                <div className="instance-metrics">
                  <div className="metric">
                    <span className="metric-label">CPU</span>
                    <div className="metric-bar">
                      <div
                        className="metric-fill cpu"
                        style={{ width: `${inst.cpu}%` }}
                      />
                    </div>
                    <span className="metric-value">{inst.cpu?.toFixed(1)}%</span>
                  </div>
                  <div className="metric">
                    <span className="metric-label">内存</span>
                    <div className="metric-bar">
                      <div
                        className="metric-fill memory"
                        style={{ width: `${inst.memory}%` }}
                      />
                    </div>
                    <span className="metric-value">{inst.memory?.toFixed(1)}%</span>
                  </div>
                  <div className="metric">
                    <span className="metric-label">磁盘</span>
                    <div className="metric-bar">
                      <div
                        className="metric-fill disk"
                        style={{ width: `${inst.disk}%` }}
                      />
                    </div>
                    <span className="metric-value">{inst.disk?.toFixed(1)}%</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 右侧：告警列表 */}
        <div className="alerts-panel">
          <h2>🚨 活跃告警</h2>
          <div className="alerts-list">
            {alerts.length === 0 ? (
              <div className="no-alerts">🎉 暂无活跃告警</div>
            ) : (
              alerts.map((alert) => (
                <div key={alert.id} className={`alert-card ${getSeverityClass(alert.severity)}`}>
                  <div className="alert-header">
                    <span className={`severity-badge ${alert.severity}`}>
                      {alert.severity === 'critical' ? '🔴 严重' : '🟡 警告'}
                    </span>
                    <span className={`status-badge ${getStatusClass(alert.status)}`}>
                      {alert.status === 'firing' ? '进行中' : alert.status === 'acknowledged' ? '已确认' : '已解决'}
                    </span>
                  </div>
                  <div className="alert-body">
                    <div className="alert-message">{alert.message}</div>
                    <div className="alert-meta">
                      <span>📍 {alert.instance_name}</span>
                      <span>⏰ {formatTime(alert.triggered_at)}</span>
                    </div>
                  </div>
                  <div className="alert-actions">
                    {alert.status === 'firing' && (
                      <>
                        <button className="btn-ack">确认</button>
                        <button className="btn-resolve">解决</button>
                      </>
                    )}
                    {alert.status === 'acknowledged' && (
                      <button className="btn-resolve">解决</button>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* 底部：告警趋势 */}
      <div className="trend-panel">
        <h2>📈 近7天告警趋势</h2>
        <div className="trend-chart">
          {stats.daily_stats && Object.entries(stats.daily_stats).map(([day, count]) => (
            <div key={day} className="trend-bar">
              <div className="trend-value">{count}</div>
              <div className="trend-bar-fill" style={{ height: `${(count / 25) * 100}%` }} />
              <div className="trend-label">{day}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Monitoring;
