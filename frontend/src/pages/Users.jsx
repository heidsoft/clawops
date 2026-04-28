import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import './Users.css';

function Users() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState(null);

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const res = await deployService.get('/users');
      if (res.users) {
        setUsers(res.users);
      }
    } catch (error) {
      console.error('Failed to load users:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (userData) => {
    try {
      await deployService.post('/users', userData);
      alert('用户创建成功');
      setShowModal(false);
      loadUsers();
    } catch (error) {
      alert('创建失败：' + error.message);
    }
  };

  const handleDelete = async (userId) => {
    if (!confirm('确定要删除这个用户吗？')) return;
    try {
      await deployService.delete(`/users/${userId}`);
      alert('用户删除成功');
      loadUsers();
    } catch (error) {
      alert('删除失败：' + error.message);
    }
  };

  const getRoleBadgeClass = (role) => {
    switch (role) {
      case 'super_admin': return 'role-super';
      case 'admin': return 'role-admin';
      case 'manager': return 'role-manager';
      case 'user': return 'role-user';
      case 'viewer': return 'role-viewer';
      default: return '';
    }
  };

  const getRoleName = (role) => {
    const names = {
      'super_admin': '超级管理员',
      'admin': '管理员',
      'manager': '运维经理',
      'user': '普通用户',
      'viewer': '访客',
    };
    return names[role] || role;
  };

  return (
    <div className="users-container">
      <div className="users-header">
        <h1>👥 用户管理</h1>
        <button className="btn-create" onClick={() => setShowModal(true)}>
          ➕ 创建用户
        </button>
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : (
        <div className="users-table">
          <table>
            <thead>
              <tr>
                <th>用户名</th>
                <th>昵称</th>
                <th>邮箱</th>
                <th>角色</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.username}</td>
                  <td>{user.nickname}</td>
                  <td>{user.email}</td>
                  <td>
                    <span className={`role-badge ${getRoleBadgeClass(user.role)}`}>
                      {getRoleName(user.role)}
                    </span>
                  </td>
                  <td>
                    <span className={`status-badge ${user.status}`}>
                      {user.status === 'active' ? '✅ 正常' : '❌ 停用'}
                    </span>
                  </td>
                  <td>
                    <button className="btn-edit">编辑</button>
                    <button className="btn-delete" onClick={() => handleDelete(user.id)}>删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>创建用户</h2>
            <form onSubmit={(e) => {
              e.preventDefault();
              handleCreate({
                username: e.target.username.value,
                email: e.target.email.value,
                password: e.target.password.value,
                nickname: e.target.nickname.value,
                role: e.target.role.value,
              });
            }}>
              <div className="form-group">
                <label>用户名</label>
                <input type="text" name="username" required />
              </div>
              <div className="form-group">
                <label>昵称</label>
                <input type="text" name="nickname" />
              </div>
              <div className="form-group">
                <label>邮箱</label>
                <input type="email" name="email" required />
              </div>
              <div className="form-group">
                <label>密码</label>
                <input type="password" name="password" required />
              </div>
              <div className="form-group">
                <label>角色</label>
                <select name="role">
                  <option value="user">普通用户</option>
                  <option value="manager">运维经理</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setShowModal(false)}>取消</button>
                <button type="submit" className="btn-submit">创建</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Users;
