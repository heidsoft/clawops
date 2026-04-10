import React, { useState, useEffect } from 'react';
import { dockerService } from '../services/deploy';
import './Docker.css';

function Docker() {
  const [deployments, setDeployments] = useState([]);
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedDeployment, setSelectedDeployment] = useState(null);

  useEffect(() => {
    loadDeployments();
    loadImages();
  }, []);

  const loadDeployments = async () => {
    try {
      setLoading(true);
      const data = await dockerService.getDockerDeployments();
      setDeployments(data.data || []);
    } catch (error) {
      console.error('Failed to load docker deployments:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadImages = async () => {
    try {
      const data = await dockerService.getDockerImages();
      setImages(data.data || []);
    } catch (error) {
      console.error('Failed to load images:', error);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    try {
      await dockerService.createDocker({
        user_id: 'u001',
        name: formData.get('name'),
        image: formData.get('image'),
        plan: formData.get('plan'),
        container_port: parseInt(formData.get('container_port')),
        command: formData.get('command'),
      });
      setShowCreateModal(false);
      loadDeployments();
      alert('创建成功');
    } catch (error) {
      alert('创建失败：' + error.message);
    }
  };

  const handleStart = async (id) => {
    try {
      await dockerService.startDocker(id);
      loadDeployments();
      alert('启动成功');
    } catch (error) {
      alert('启动失败：' + error.message);
    }
  };

  const handleStop = async (id) => {
    try {
      await dockerService.stopDocker(id);
      loadDeployments();
      alert('停止成功');
    } catch (error) {
      alert('停止失败：' + error.message);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('确定要删除这个容器吗？此操作不可恢复！')) {
      return;
    }
    try {
      await dockerService.deleteDocker(id);
      loadDeployments();
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
    <div className="docker-page">
      <div className="page-header">
        <h1><i className="fab fa-docker"></i> Docker 容器</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowCreateModal(true)}
        >
          <i className="fas fa-plus"></i> 新建容器
        </button>
      </div>

      <div className="docker-list">
        {deployments.map(docker => (
          <div key={docker.id} className="docker-card">
            <div className="docker-header">
              <div className="docker-info">
                <h3>
                  <i className="fab fa-docker"></i>
                  {docker.name}
                </h3>
                <span className="docker-image">{docker.image}</span>
              </div>
              {getStatusBadge(docker.status)}
            </div>

            <div className="docker-details">
              <div className="detail-item">
                <i className="fas fa-server"></i>
                <span>{docker.host}:{docker.external_port}</span>
              </div>
              <div className="detail-item">
                <i className="fas fa-microchip"></i>
                <span>{docker.cpu} 核</span>
              </div>
              <div className="detail-item">
                <i className="fas fa-memory"></i>
                <span>{docker.memory}MB</span>
              </div>
            </div>

            <div className="docker-actions">
              {docker.status === 'running' ? (
                <button className="btn btn-sm btn-outline-warning" onClick={() => handleStop(docker.id)}>
                  <i className="fas fa-stop"></i> 停止
                </button>
              ) : (
                <button className="btn btn-sm btn-outline-success" onClick={() => handleStart(docker.id)}>
                  <i className="fas fa-play"></i> 启动
                </button>
              )}
              <button className="btn btn-sm btn-outline-info" title="日志">
                <i className="fas fa-file-alt"></i> 日志
              </button>
              <button className="btn btn-sm btn-outline-primary" title="监控">
                <i className="fas fa-chart-line"></i>
              </button>
              <button className="btn btn-sm btn-outline-danger" onClick={() => handleDelete(docker.id)} title="删除">
                <i className="fas fa-trash"></i>
              </button>
            </div>
          </div>
        ))}
      </div>

      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h2>新建 Docker 容器</h2>
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label>容器名称</label>
                <input type="text" name="name" placeholder="如: nginx-web" required />
              </div>
              <div className="form-group">
                <label>镜像</label>
                <select name="image" required>
                  <option value="">请选择镜像</option>
                  {images.map(img => (
                    <option key={img.name} value={`${img.name}:latest`}>
                      {img.name} - {img.description}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label>套餐</label>
                <select name="plan" required>
                  <option value="small">小型 (1核1G)</option>
                  <option value="medium">中型 (2核4G)</option>
                  <option value="large">大型 (4核8G)</option>
                </select>
              </div>
              <div className="form-group">
                <label>容器端口</label>
                <input type="number" name="container_port" placeholder="如: 80" required />
              </div>
              <div className="form-group">
                <label>启动命令（可选）</label>
                <input type="text" name="command" placeholder="如: nginx -g 'daemon off;'" />
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

export default Docker;
