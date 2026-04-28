import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import './Audit.css';

function Audit() {
  const [logs, setLogs] = useState([]);
  const [stats, setStats] = useState({});
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState({ action: '', resource: '' });

  useEffect(() => {
    loadData();
  }, [filter]);

  const loadData = async () => {
    try {
      const params = {};
      if (filter.action) params.action = filter.action;
      if (filter.resource) params.resource = filter.resource;

      const [logsRes, statsRes] = await Promise.all([
        deployService.get('/audit/logs', params),
        deployService.get('/audit/stats'),
      ]);

      if (logsRes.logs) setLogs(logsRes.logs);
      if (statsRes) setStats(statsRes);
    } catch (error) {
      console.error('Failed to load audit data:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (time) => {
    if (!time) return '-';
    const d = new Date(time);
    return d.toLocaleString('zh-CN');
  };

  const getActionIcon = (action) => {
    switch (action) {
      case 'login': return '🔑';
      case 'logout': return '🔒';
      case 'create': return '➕';
      case 'update': return '✏️';
      case 'delete': return '🗑️';
      default: return '📝';
    }
  };

  const getActionName = (action) => {
    const names = {
      'login': '登录',
      'logout': '登出',
      'create': '创建',
      'update': '更新',
      'delete': '删除',
      'read': '查看',
    };
    return names[action] || action;
  };

  const getResourceName = (resource) => {
    const names = {
      'session': '会话',
      'deployment': '部署',
      'database': '数据库',
      'docker': '容器',
      'monitor': '监控',
      'user': '用户',
      'setting': '设置',
      'skill': 'Skill',
    };
    return names[resource] || resource;
  };

  return (
    <div className="audit-container">
      <div className="audit-header">
        <h1>📋 审计日志</h1>
      </div>

      {/* 统计卡片 */}
      <div className="audit-stats">
        <div className="stat-card">
          <div className="stat-value">{stats.total_actions || 0}</div>
          <div className="stat-label">总操作数</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.today_actions || 0}</div>
          <div className="stat-label">今日操作</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.week_actions || 0}</div>
          <div className="stat-label">本周操作</div>
        </div>
      </div>

      {/* 筛选 */}
      <div className="audit-filters">
        <select value={filter.action} onChange={(e) => setFilter({ ...filter, action: e.target.value })}>
          <option value="">全部操作</option>
          <option value="login">登录</option>
          <option value="create">创建</option>
          <option value="update">更新</option>
          <option value="delete">删除</option>
        </select>
        <select value={filter.resource} onChange={(e) => setFilter({ ...filter, resource: e.target.value })}>
          <option value="">全部资源</option>
          <option value="deployment">部署</option>
          <option value="database">数据库</option>
          <option value="docker">容器</option>
          <option value="user">用户</option>
          <option value="monitor">监控</option>
        </select>
      </div>

      {/* 日志列表 */}
      {loading ? (
        <div className="loading">加载中...</div>
      ) : (
        <div className="audit-table">
          <table>
            <thead>
              <tr>
                <th>时间</th>
                <th>用户</th>
                <th>操作</th>
                <th>资源</th>
                <th>详情</th>
                <th>IP</th>
                <th>状态</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id}>
                  <td>{formatTime(log.created_at)}</td>
                  <td>{log.username}</td>
                  <td>
                    <span className="action-badge">
                      {getActionIcon(log.action)} {getActionName(log.action)}
                    </span>
                  </td>
                  <td>{getResourceName(log.resource)}</td>
                  <td className="detail-cell">
                    {log.detail ? (
                      <span title={JSON.stringify(log.detail)}>
                        {typeof log.detail === 'string' ? JSON.parse(log.detail)?.name || log.resource_id || '-' : '-'}
                      </span>
                    ) : '-'}
                  </td>
                  <td>{log.ip || '-'}</td>
                  <td>
                    <span className={`status-badge ${log.status}`}>
                      {log.status === 'success' ? '✅ 成功' : '❌ 失败'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default Audit;
