import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getRules } from '../../api/general';

const { Title } = Typography;

const RuleList: React.FC = () => {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);

  const fetchData = (page = 1) => {
    setLoading(true);
    getRules(page).then(res => {
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
    { title: 'Title', dataIndex: 'title', key: 'title' },
    { title: 'Active', dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? 'Yes' : 'No'}</Tag> },
    { title: 'Order', dataIndex: 'order', key: 'order' },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>Rules</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Rule</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
    </div>
  );
};

export default RuleList;
