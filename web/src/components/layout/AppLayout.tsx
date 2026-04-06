import React, { useState } from 'react';
import { Layout, Menu, theme, Avatar, Dropdown, Space } from 'antd';
import {
  DashboardOutlined, BankOutlined, SwapOutlined, FundOutlined,
  FileTextOutlined, TagsOutlined, TagOutlined, SaveOutlined,
  SettingOutlined, LogoutOutlined, DollarOutlined, ThunderboltOutlined,
  SyncOutlined, ApiOutlined, AppstoreOutlined, UserOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch } from '../../store/hooks';
import { logout } from '../../store/slices/authSlice';

const { Header, Sider, Content } = Layout;

const AppLayout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken();

  const menuItems = [
    { key: '/', icon: <DashboardOutlined />, label: 'Dashboard' },
    { key: '/accounts/asset', icon: <BankOutlined />, label: 'Accounts' },
    { key: '/transactions/withdrawal', icon: <SwapOutlined />, label: 'Transactions' },
    { key: '/budgets', icon: <FundOutlined />, label: 'Budgets' },
    { key: '/bills', icon: <FileTextOutlined />, label: 'Bills' },
    { key: '/categories', icon: <TagsOutlined />, label: 'Categories' },
    { key: '/tags', icon: <TagOutlined />, label: 'Tags' },
    { key: '/piggy-banks', icon: <SaveOutlined />, label: 'Piggy Banks' },
    { key: '/rules', icon: <ThunderboltOutlined />, label: 'Rules' },
    { key: '/recurring', icon: <SyncOutlined />, label: 'Recurring' },
    { key: '/currencies', icon: <DollarOutlined />, label: 'Currencies' },
    { key: '/webhooks', icon: <ApiOutlined />, label: 'Webhooks' },
    { key: '/object-groups', icon: <AppstoreOutlined />, label: 'Groups' },
  ];

  const userMenu = [
    { key: 'profile', icon: <UserOutlined />, label: 'Profile', onClick: () => navigate('/profile') },
    { key: 'admin', icon: <SettingOutlined />, label: 'Admin', onClick: () => navigate('/admin') },
    { type: 'divider' as const },
    { key: 'logout', icon: <LogoutOutlined />, label: 'Logout', onClick: () => { dispatch(logout()); navigate('/login'); } },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
        <div style={{ height: 32, margin: 16, background: 'rgba(255,255,255,.2)', borderRadius: 6, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 'bold' }}>
          {collapsed ? 'FF' : 'Firefly III'}
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
        <Header style={{ padding: '0 24px', background: colorBgContainer, display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
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
    </Layout>
  );
};

export default AppLayout;
