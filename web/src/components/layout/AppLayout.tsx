import React, { useState } from 'react';
import { Layout, Menu, theme, Avatar, Dropdown, Space, Select } from 'antd';
import {
  DashboardOutlined, BankOutlined, SwapOutlined, FundOutlined,
  FileTextOutlined, TagsOutlined, TagOutlined, SaveOutlined,
  SettingOutlined, LogoutOutlined, DollarOutlined, ThunderboltOutlined,
  SyncOutlined, ApiOutlined, AppstoreOutlined, UserOutlined, GlobalOutlined, UploadOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch } from '../../store/hooks';
import { logout } from '../../store/slices/authSlice';
import { useTranslation } from 'react-i18next';
import { languages } from '../../i18n';
import ChatBubble from '../ChatBubble';

const { Header, Sider, Content } = Layout;

const AppLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken();
  const { t, i18n } = useTranslation();

  const menuItems = [
    { key: '/', icon: <DashboardOutlined />, label: t('dashboard') },
    { key: '/accounts/asset', icon: <BankOutlined />, label: t('accounts') },
    { key: '/transactions/withdrawal', icon: <SwapOutlined />, label: t('transactions') },
    { key: '/budgets', icon: <FundOutlined />, label: t('budgets') },
    { key: '/bills', icon: <FileTextOutlined />, label: t('bills') },
    { key: '/categories', icon: <TagsOutlined />, label: t('categories') },
    { key: '/tags', icon: <TagOutlined />, label: t('tags') },
    { key: '/piggy-banks', icon: <SaveOutlined />, label: t('piggy_banks') },
    { key: '/rules', icon: <ThunderboltOutlined />, label: t('rules') },
    { key: '/recurring', icon: <SyncOutlined />, label: t('recurring') },
    { key: '/currencies', icon: <DollarOutlined />, label: t('currencies') },
    { key: '/webhooks', icon: <ApiOutlined />, label: t('webhooks') },
    { key: '/object-groups', icon: <AppstoreOutlined />, label: t('groups') },
    { key: '/import', icon: <UploadOutlined />, label: t('import') },
  ];

  const userMenu = [
    { key: 'profile', icon: <UserOutlined />, label: t('profile'), onClick: () => navigate('/profile') },
    { key: 'admin', icon: <SettingOutlined />, label: t('admin'), onClick: () => navigate('/admin') },
    { type: 'divider' as const },
    { key: 'logout', icon: <LogoutOutlined />, label: t('logout'), onClick: () => { dispatch(logout()); navigate('/login'); } },
  ];

  const handleLanguageChange = (lng: string) => {
    i18n.changeLanguage(lng);
    localStorage.setItem('firefly_language', lng);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
        <div style={{ height: 32, margin: 16, background: 'rgba(255,255,255,.2)', borderRadius: 6, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 'bold' }}>
          {collapsed ? 'Q' : 'Quillow'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: '0 24px', background: colorBgContainer, display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 16 }}>
          <Space>
            <GlobalOutlined />
            <Select
              value={i18n.language}
              onChange={handleLanguageChange}
              style={{ width: 140 }}
              size="small"
              options={languages.map(l => ({ value: l.code, label: l.label }))}
            />
          </Space>
          <Dropdown menu={{ items: userMenu }} placement="bottomRight">
            <Space style={{ cursor: 'pointer' }}>
              <Avatar icon={<UserOutlined />} />
            </Space>
          </Dropdown>
        </Header>
        <Content style={{ margin: 24 }}>
          <div style={{ padding: 24, background: colorBgContainer, borderRadius: borderRadiusLG, minHeight: 360 }}>
            <Outlet />
          </div>
        </Content>
      </Layout>
      <ChatBubble />
    </Layout>
  );
};

export default AppLayout;
