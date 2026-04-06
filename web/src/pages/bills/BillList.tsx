import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getBills } from '../../api/general';

const { Title } = Typography;

const BillList: React.FC = () => {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);

  const fetchData = (page = 1) => {
    setLoading(true);
    getBills(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
      setLoading(false);
    }).catch(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const columns = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Amount Min', dataIndex: 'amount_min', key: 'amount_min' },
    { title: 'Amount Max', dataIndex: 'amount_max', key: 'amount_max' },
    { title: 'Repeat Freq', dataIndex: 'repeat_freq', key: 'repeat_freq' },
    { title: 'Active', dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? 'Yes' : 'No'}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>Bills</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Bill</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
    </div>
  );
};

export default BillList;
