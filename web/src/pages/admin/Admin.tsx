import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Typography, Descriptions, Table, Button, message, Popconfirm } from 'antd';
import { useTranslation } from 'react-i18next';
import client from '../../api/client';

const { Title } = Typography;

const Admin: React.FC = () => {
  const { t } = useTranslation();
  const [about, setAbout] = useState<Record<string, any>>({});
  const [users, setUsers] = useState<any[]>([]);
  const [config, setConfig] = useState<any[]>([]);

  const loadData = () => {
    client.get('/about').then(res => setAbout(res.data)).catch(() => {});
    client.get('/users').then(res => {
      const list = res.data?.data || [];
      setUsers(list.map((r: any) => ({ id: r.id, ...r.attributes })));
    }).catch(() => {});
    client.get('/configuration').then(res => {
      const list = res.data?.data || res.data || [];
      setConfig(Array.isArray(list) ? list : []);
    }).catch(() => {});
  };

  useEffect(() => { loadData(); }, []);

  const deleteUser = (id: string) => {
    client.delete(`/users/${id}`).then(() => {
      message.success(t('success_deleted'));
      loadData();
    }).catch(() => message.error(t('error_occurred')));
  };

  const userColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: t('email'), dataIndex: 'email', key: 'email' },
    { title: t('blocked'), dataIndex: 'blocked', key: 'blocked', render: (v: boolean) => v ? t('yes') : t('no') },
    { title: t('role'), dataIndex: 'role', key: 'role' },
    { title: t('created'), dataIndex: 'created_at', key: 'created_at' },
    {
      title: t('actions'), key: 'actions',
      render: (_: any, record: any) => (
        <Popconfirm title={t('delete_confirm')} onConfirm={() => deleteUser(record.id)}>
          <Button size="small" danger>{t('delete')}</Button>
        </Popconfirm>
      ),
    },
  ];

  return (
    <div>
      <Title level={4}>{t('administration')}</Title>
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card title={t('system_info')}>
            <Descriptions column={2}>
              <Descriptions.Item label={t('version')}>{about.version || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('api_version')}>{about.api_version || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('os')}>{about.os || '-'}</Descriptions.Item>
              <Descriptions.Item label={t('database_driver')}>{about.driver || '-'}</Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
        <Col span={24}>
          <Card title={t('users')}>
            <Table dataSource={users} columns={userColumns} rowKey="id" size="small" pagination={false} />
          </Card>
        </Col>
        {config.length > 0 && (
          <Col span={24}>
            <Card title={t('configuration')}>
              <Descriptions column={1}>
                {config.map((c: any, i: number) => (
                  <Descriptions.Item key={i} label={c.attributes?.name || c.name}>
                    {c.attributes?.value ?? c.data ?? '-'}
                  </Descriptions.Item>
                ))}
              </Descriptions>
            </Card>
          </Col>
        )}
      </Row>
    </div>
  );
};

export default Admin;
