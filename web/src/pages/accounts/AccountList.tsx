import React, { useEffect } from 'react';
import { Table, Button, Space, Tag, Popconfirm, Typography, message } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { fetchAccounts } from '../../store/slices/accountSlice';
import { deleteAccount } from '../../api/accounts';

const { Title } = Typography;

const AccountList: React.FC = () => {
  const { type } = useParams<{ type: string }>();
  const dispatch = useAppDispatch();
  const { items, pagination, loading } = useAppSelector(s => s.accounts);

  useEffect(() => {
    dispatch(fetchAccounts({ type: type ? `${type} account` : undefined }));
  }, [dispatch, type]);

  const handleDelete = async (id: string) => {
    await deleteAccount(id);
    message.success('Account deleted');
    dispatch(fetchAccounts({ type: type ? `${type} account` : undefined }));
  };

  const columns = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Type', dataIndex: 'type', key: 'type' },
    { title: 'Balance', dataIndex: 'current_balance', key: 'balance', render: (v: string, r: any) => `${r.currency_symbol || ''}${v || '0'}` },
    { title: 'IBAN', dataIndex: 'iban', key: 'iban' },
    { title: 'Active', dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? 'Yes' : 'No'}</Tag> },
    { title: 'Action', key: 'action', render: (_: any, record: any) => (
      <Popconfirm title="Delete?" onConfirm={() => handleDelete(record.id)}>
        <Button danger icon={<DeleteOutlined />} size="small" />
      </Popconfirm>
    )},
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{type ? `${type.charAt(0).toUpperCase() + type.slice(1)} Accounts` : 'Accounts'}</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Account</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page,
          onChange: (page, pageSize) => dispatch(fetchAccounts({ page, limit: pageSize, type: type ? `${type} account` : undefined }))
        } : false}
      />
    </div>
  );
};

export default AccountList;
