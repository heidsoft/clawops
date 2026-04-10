import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 部署服务
export const deployService = {
  // 获取部署列表
  getDeployments: async (page = 1, pageSize = 20) => {
    const response = await apiClient.get('/deployments', {
      params: { page, page_size: pageSize },
    });
    return response;
  },

  // 获取部署详情
  getDeployment: async (id) => {
    const response = await apiClient.get(`/deployments/${id}`);
    return response.data;
  },

  // 创建部署
  createDeployment: async (data) => {
    const response = await apiClient.post('/deployments', data);
    return response.data;
  },

  // 更新部署
  updateDeployment: async (id, data) => {
    const response = await apiClient.put(`/deployments/${id}`, data);
    return response.data;
  },

  // 删除部署
  deleteDeployment: async (id) => {
    const response = await apiClient.delete(`/deployments/${id}`);
    return response.data;
  },

  // 启动部署
  startDeployment: async (id) => {
    const response = await apiClient.post(`/deployments/${id}/start`);
    return response.data;
  },

  // 停止部署
  stopDeployment: async (id) => {
    const response = await apiClient.post(`/deployments/${id}/stop`);
    return response.data;
  },

  // 获取监控数据
  getMetrics: async (id) => {
    const response = await apiClient.get(`/deployments/${id}/metrics`);
    return response.data;
  },

  // 获取日志
  getLogs: async (id, options = {}) => {
    const response = await apiClient.get(`/deployments/${id}/logs`, {
      params: options,
    });
    return response.data;
  },
};

// 用户服务
export const userService = {
  // 获取当前用户
  getCurrentUser: async () => {
    const response = await apiClient.get('/users/me');
    return response.data;
  },

  // 登录
  login: async (username, password) => {
    const response = await apiClient.post('/auth/login', { username, password });
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
    }
    return response.data;
  },

  // 登出
  logout: () => {
    localStorage.removeItem('token');
  },
};

// 监控服务
export const monitorService = {
  // 获取系统状态
  getSystemStatus: async () => {
    const response = await apiClient.get('/monitor/system');
    return response.data;
  },

  // 获取告警列表
  getAlerts: async () => {
    const response = await apiClient.get('/monitor/alerts');
    return response.data;
  },

  // 确认告警
  acknowledgeAlert: async (id) => {
    const response = await apiClient.post(`/monitor/alerts/${id}/ack`);
    return response.data;
  },
};

// 数据库部署服务
export const databaseService = {
  // 获取数据库部署列表
  getDatabases: async (page = 1, pageSize = 20) => {
    const response = await apiClient.get('/databases', {
      params: { page, page_size: pageSize },
    });
    return response;
  },

  // 获取数据库详情
  getDatabase: async (id) => {
    const response = await apiClient.get(`/databases/${id}`);
    return response.data;
  },

  // 创建数据库部署
  createDatabase: async (data) => {
    const response = await apiClient.post('/databases', data);
    return response.data;
  },

  // 删除数据库部署
  deleteDatabase: async (id) => {
    const response = await apiClient.delete(`/databases/${id}`);
    return response.data;
  },

  // 获取备份列表
  getDatabaseBackups: async (id) => {
    const response = await apiClient.get(`/databases/${id}/backups`);
    return response.data;
  },

  // 创建备份
  createDatabaseBackup: async (id) => {
    const response = await apiClient.post(`/databases/${id}/backups`);
    return response.data;
  },
};

// Docker 部署服务
export const dockerService = {
  // 获取 Docker 部署列表
  getDockerDeployments: async (page = 1, pageSize = 20) => {
    const response = await apiClient.get('/docker', {
      params: { page, page_size: pageSize },
    });
    return response;
  },

  // 获取容器详情
  getDockerDeployment: async (id) => {
    const response = await apiClient.get(`/docker/${id}`);
    return response.data;
  },

  // 创建 Docker 部署
  createDocker: async (data) => {
    const response = await apiClient.post('/docker', data);
    return response.data;
  },

  // 启动容器
  startDocker: async (id) => {
    const response = await apiClient.post(`/docker/${id}/start`);
    return response.data;
  },

  // 停止容器
  stopDocker: async (id) => {
    const response = await apiClient.post(`/docker/${id}/stop`);
    return response.data;
  },

  // 删除容器
  deleteDocker: async (id) => {
    const response = await apiClient.delete(`/docker/${id}`);
    return response.data;
  },

  // 获取容器日志
  getDockerLogs: async (id, lines = 100) => {
    const response = await apiClient.get(`/docker/${id}/logs`, {
      params: { lines },
    });
    return response.data;
  },

  // 获取容器状态
  getDockerStats: async (id) => {
    const response = await apiClient.get(`/docker/${id}/stats`);
    return response.data;
  },

  // 获取常用镜像列表
  getDockerImages: async () => {
    const response = await apiClient.get('/docker/images');
    return response;
  },
};

export default apiClient;
