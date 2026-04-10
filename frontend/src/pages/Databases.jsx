import React, { useState, useEffect } from 'react';
import { databaseService } from '../services/deploy';
import './Databases.css';

function Databases() {
  const [databases, setDatabases] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedDb, setSelectedDb] = useState(null);

  useEffect(() => {
    loadDatabases();
  }, []);

  const loadDatabases = async () => {
    try {
      setLoading(true);
      const data = await databaseService.getDatabases();
      setDatabases(data.data || []);
    } catch (error) {
      console.error('Failed to load databases:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
      await databaseService.createDatabase({
        user_id: 'u001',
        name: formData.get('name'),
        database_type: formData.get('database_type'),
        version: formData.get('version'),
        plan: formData.get('plan'),
      });
      setShowCreateModal(false);
      loadDatabases();
      alert('创建成功');
    } catch (error) {
      alert('创建失败：' + error.message);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('确定要删除这个数据库部署吗？此操作不可恢复！')) {
      return;
    }
    try {
      await databaseService.deleteDatabase(id);
      loadDatabases();
      alert('删除成功');
    } catch (error) {
      alert('删除失败：' + error.message);
    }
  };

  const getStatusBadge = (status) => {
    const statusMap = {
      running: { class: 'status-running', text: '运行中', icon: 'fa-circle' },
      deploying: { class: 'status-deploying', text: '部署中', icon: 'fa-spinner fa-spin' },
      stopped: { class: 'status-stopped', text: '已停止', icon: 'fa-circle' },
      error: { class: 'status-error', text: '错误', icon: 'fa-exclamation-circle' },
    };
    const statusInfo = statusMap[status] || statusMap.stopped;
    return (
      <span className={`status-badge ${statusInfo.class}`}>
        <i className={`fas ${statusInfo.icon}`}></i> {statusInfo.text}
      </span>
    );
  };

  if (loading) {
    return <div className="loading-container"><div className="loading">加载中...</div></div>;
  }

  return (
    <div className="databases-page">
      <div className="page-header">
        <h1><i className="fas fa-database"></i> 数据库部署</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowCreateModal(true)}
        >
          <i className="fas fa-plus"></i> 新建数据库
        </button>
      </div>

      <div className="databases-list">
        {databases.map(db => (
          <div key={db.id} className="database-card">
            <div className="database-header">
              <div className="database-info">
                <h3>
                  <i className="fas fa-database"></i>
                  {db.name}
                </h3>
                <span className="db-type">{db.database_type} {db.version}</span>
              </div>
              {getStatusBadge(db.status)}
            </div>

            <div className="database-details">
              <div className="detail-item">
                <i className="fas fa-server"></i>
                <span>{db.host}:{db.port}</span>
              </div>
              <div className="detail-item">
                <i className="fas fa-memory"></i>
                <span>{db.memory_size}MB</span>
              </div>
              <div className="detail-item">
                <i className="fas fa-hdd"></i>
                <span>{db.disk_size}GB</span>
              </div>
            </div>

            <div className="database-actions">
              <button className="btn btn-sm btn-outline-primary" title="备份">
                <i className="fas fa-backup"></i> 备份
              </button>
              <button className="btn btn-sm btn-outline-secondary" title="配置">
                <i className="fas fa-cog"></i>
              </button>
              <button className="btn btn-sm btn-outline-danger" title="删除" onClick={() => handleDelete(db.id)}>
                <i className="fas fa-trash"></i>
              </button>
            </div>
          </div>
        ))}
      </div>

      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>新建数据库部署</h2>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>数据库名称</label>
                <input type="text" name="name" placeholder="如: mysql-prod" required />
              </div>
              <div className="form-group">
                <label>数据库类型</label>
                <select name="database_type" required>
                  <option value="">请选择</option>
                  <option value="mysql">MySQL</option>
                  <option value="postgresql">PostgreSQL</option>
                </select>
              </div>
              <div className="form-group">
                <label>版本</label>
                <select name="version" required>
                  <option value="8.0">MySQL 8.0</option>
                  <option value="5.7">MySQL 5.7</option>
                  <option value="14">PostgreSQL 14</option>
                  <option value="13">PostgreSQL 13</option>
                </select>
              </div>
              <div className="form-group">
                <label>套餐</label>
                <select name="plan" required>
                  <option value="small">小型 (1核2G) - 40GB</option>
                  <option value="medium">中型 (2核4G) - 100GB</option>
                  <option value="large">大型 (4核8G) - 200GB</option>
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn" onClick={() => setShowCreateModal(false)}>取消</button>
                <button type="submit" className="btn btn-primary">创建</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Databases;
