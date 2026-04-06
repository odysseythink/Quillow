import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Typography } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getTags } from '../../api/general';

const { Title } = Typography;

const TagList: React.FC = () => {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);

  const fetchData = (page = 1) => {
    setLoading(true);
    getTags(page).then(res => {
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
    { title: 'Tag', dataIndex: 'tag', key: 'tag' },
    { title: 'Description', dataIndex: 'description', key: 'description' },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>Tags</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Tag</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
    </div>
  );
};

export default TagList;
