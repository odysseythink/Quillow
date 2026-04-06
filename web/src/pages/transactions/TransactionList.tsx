import React, { useEffect } from 'react';
import { Table, Button, Space, Tag, Popconfirm, Typography, message } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { fetchTransactions } from '../../store/slices/transactionSlice';
import { deleteTransaction } from '../../api/transactions';

const { Title } = Typography;

const TransactionList: React.FC = () => {
  const { type } = useParams<{ type: string }>();
  const dispatch = useAppDispatch();
  const { items, pagination, loading } = useAppSelector(s => s.transactions);

  useEffect(() => {
    dispatch(fetchTransactions({ type }));
  }, [dispatch, type]);

  const handleDelete = async (id: string) => {
    await deleteTransaction(id);
    message.success('Transaction deleted');
    dispatch(fetchTransactions({ type }));
  };

  const columns = [
    { title: 'Description', key: 'desc', render: (_: any, r: any) => r.transactions?.[0]?.description || r.group_title || '-' },
    { title: 'Amount', key: 'amount', render: (_: any, r: any) => { const t = r.transactions?.[0]; return t ? `${t.currency_symbol || ''}${t.amount || '0'}` : '-'; }},
    { title: 'Source', key: 'source', render: (_: any, r: any) => r.transactions?.[0]?.source_name || '-' },
    { title: 'Destination', key: 'dest', render: (_: any, r: any) => r.transactions?.[0]?.destination_name || '-' },
    { title: 'Date', key: 'date', render: (_: any, r: any) => r.transactions?.[0]?.date?.substring(0, 10) || '-' },
    { title: 'Type', key: 'type', render: (_: any, r: any) => <Tag>{r.transactions?.[0]?.type || '-'}</Tag> },
    { title: 'Action', key: 'action', render: (_: any, record: any) => (
      <Popconfirm title="Delete?" onConfirm={() => handleDelete(record.id)}>
        <Button danger icon={<DeleteOutlined />} size="small" />
      </Popconfirm>
    )},
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{type ? `${type.charAt(0).toUpperCase() + type.slice(1)}s` : 'Transactions'}</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Transaction</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page,
          onChange: (page, pageSize) => dispatch(fetchTransactions({ page, limit: pageSize, type }))
        } : false}
      />
    </div>
  );
};

export default TransactionList;
