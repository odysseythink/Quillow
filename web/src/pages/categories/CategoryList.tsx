import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Typography } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getCategories } from '../../api/general';

const { Title } = Typography;

const CategoryList: React.FC = () => {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);

  const fetchData = (page = 1) => {
    setLoading(true);
    getCategories(page).then(res => {
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
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>Categories</Title>
        <Button type="primary" icon={<PlusOutlined />}>Add Category</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
    </div>
  );
};

export default CategoryList;
