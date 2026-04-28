import React, { useState, useEffect } from 'react';
import { deployService } from '../services/deploy';
import './Skills.css';

function Skills() {
  const [skills, setSkills] = useState([]);
  const [categories, setCategories] = useState([]);
  const [userSkills, setUserSkills] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('market'); // market / my-skills
  const [selectedSkill, setSelectedSkill] = useState(null);

  useEffect(() => {
    loadSkills();
    loadUserSkills();
  }, [selectedCategory, searchQuery]);

  const loadSkills = async () => {
    try {
      setLoading(true);
      const params = {};
      if (selectedCategory) params.category = selectedCategory;
      if (searchQuery) params.search = searchQuery;
      
      const res = await deployService.get('/skills', params);
      if (res.skills) {
        setSkills(res.skills);
        setCategories(res.categories || []);
      }
    } catch (error) {
      console.error('Failed to load skills:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadUserSkills = async () => {
    try {
      const res = await deployService.get('/my-skills');
      if (res.user_skills) {
        setUserSkills(res.user_skills);
      }
    } catch (error) {
      console.error('Failed to load user skills:', error);
    }
  };

  const handleInstall = async (skillId) => {
    try {
      await deployService.post('/skills/install', { skill_id: skillId });
      alert('安装成功！');
      loadUserSkills();
    } catch (error) {
      alert('安装失败：' + error.message);
    }
  };

  const handleUninstall = async (skillId) => {
    if (!confirm('确定要卸载这个 Skill 吗？')) return;
    try {
      await deployService.post(`/skills/${skillId}/uninstall`);
      alert('卸载成功！');
      loadUserSkills();
    } catch (error) {
      alert('卸载失败：' + error.message);
    }
  };

  const handleStar = async (skillId) => {
    try {
      await deployService.post(`/skills/${skillId}/star`);
      loadSkills(); // 刷新显示最新点赞数
    } catch (error) {
      console.error('Failed to star:', error);
    }
  };

  const formatNumber = (num) => {
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k';
    }
    return num;
  };

  return (
    <div className="skills-container">
      <div className="skills-header">
        <h1>🔌 Skill 市场</h1>
        <div className="skills-tabs">
          <button 
            className={`tab ${activeTab === 'market' ? 'active' : ''}`}
            onClick={() => setActiveTab('market')}
          >
            🏪 市场
          </button>
          <button 
            className={`tab ${activeTab === 'my-skills' ? 'active' : ''}`}
            onClick={() => setActiveTab('my-skills')}
          >
            📦 我的 Skills ({userSkills.length})
          </button>
        </div>
      </div>

      {activeTab === 'market' && (
        <>
          {/* 搜索和筛选 */}
          <div className="skills-filters">
            <div className="search-box">
              <input
                type="text"
                placeholder="🔍 搜索 Skills..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            <div className="category-filters">
              <button
                className={`category-btn ${selectedCategory === '' ? 'active' : ''}`}
                onClick={() => setSelectedCategory('')}
              >
                全部
              </button>
              {categories.map((cat) => (
                <button
                  key={cat.id}
                  className={`category-btn ${selectedCategory === cat.name ? 'active' : ''}`}
                  onClick={() => setSelectedCategory(cat.name)}
                >
                  {cat.icon} {cat.name} ({cat.count})
                </button>
              ))}
            </div>
          </div>

          {/* Skills 列表 */}
          {loading ? (
            <div className="skills-loading">加载中...</div>
          ) : (
            <div className="skills-grid">
              {skills.map((skill) => (
                <div key={skill.id} className="skill-card">
                  <div className="skill-header">
                    <span className="skill-icon">{skill.icon}</span>
                    <div className="skill-title">
                      <h3>{skill.name}</h3>
                      <span className="skill-version">v{skill.version}</span>
                    </div>
                    {skill.is_official && <span className="official-badge">官方</span>}
                  </div>
                  <p className="skill-description">{skill.description}</p>
                  <div className="skill-meta">
                    <span className="skill-author">👤 {skill.author}</span>
                    <span className="skill-stars">⭐ {formatNumber(skill.stars)}</span>
                    <span className="skill-installs">📥 {formatNumber(skill.installs)}</span>
                  </div>
                  <div className="skill-actions">
                    <button className="btn-detail" onClick={() => setSelectedSkill(skill)}>
                      查看详情
                    </button>
                    <button className="btn-star" onClick={() => handleStar(skill.id)}>
                      ⭐
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {skills.length === 0 && !loading && (
            <div className="skills-empty">
              <p>没有找到相关的 Skills</p>
              <button onClick={() => { setSearchQuery(''); setSelectedCategory(''); }}>
                清除筛选
              </button>
            </div>
          )}
        </>
      )}

      {activeTab === 'my-skills' && (
        <div className="my-skills-section">
          {userSkills.length === 0 ? (
            <div className="skills-empty">
              <p>你还没有安装任何 Skill</p>
              <button onClick={() => setActiveTab('market')}>去市场看看</button>
            </div>
          ) : (
            <div className="skills-grid">
              {userSkills.map((us) => (
                <div key={us.id} className="skill-card installed">
                  <div className="skill-header">
                    <span className="skill-icon">🚀</span>
                    <div className="skill-title">
                      <h3>{us.skill_name}</h3>
                      <span className="skill-version">v{us.version}</span>
                    </div>
                    <span className={`status-badge ${us.status}`}>{us.status}</span>
                  </div>
                  <div className="skill-info">
                    <span>安装时间：{new Date(us.installed_at).toLocaleDateString()}</span>
                  </div>
                  <div className="skill-actions">
                    <button className="btn-config">⚙️ 配置</button>
                    <button className="btn-uninstall" onClick={() => handleUninstall(us.skill_id)}>
                      🗑️ 卸载
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Skill 详情弹窗 */}
      {selectedSkill && (
        <div className="modal-overlay" onClick={() => setSelectedSkill(null)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <div className="modal-title">
                <span className="skill-icon">{selectedSkill.icon}</span>
                <h2>{selectedSkill.name}</h2>
                <span className="skill-version">v{selectedSkill.version}</span>
                {selectedSkill.is_official && <span className="official-badge">官方</span>}
              </div>
              <button className="modal-close" onClick={() => setSelectedSkill(null)}>×</button>
            </div>
            <div className="modal-body">
              <p className="skill-description">{selectedSkill.description}</p>
              
              <div className="skill-stats">
                <span>⭐ {selectedSkill.stars} 点赞</span>
                <span>📥 {selectedSkill.installs} 安装</span>
                <span>👤 {selectedSkill.author}</span>
              </div>

              <div className="skill-readme">
                <h3>📖 详细说明</h3>
                <pre>{selectedSkill.readme || `# ${selectedSkill.name}\n\n${selectedSkill.description}`}</pre>
              </div>

              <div className="skill-tags">
                {selectedSkill.tags && selectedSkill.tags.split(',').map((tag) => (
                  <span key={tag} className="tag">{tag.trim()}</span>
                ))}
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn-cancel" onClick={() => setSelectedSkill(null)}>
                关闭
              </button>
              <button className="btn-install" onClick={() => {
                handleInstall(selectedSkill.id);
                setSelectedSkill(null);
              }}>
                📥 安装
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Skills;
