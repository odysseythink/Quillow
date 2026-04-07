import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getBudgets, createBudget } from '../../api/budgets';
import { useTranslation } from 'react-i18next';

const { Title } = Typography;

const BudgetList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getBudgets(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ active: true }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createBudget({ name: v.name, active: v.active, notes: v.notes || '' } as any);
      message.success(t('success_created'));
      setModalOpen(false);
      fetchData();
    } catch (err: any) {
      if (err.response?.data?.message) message.error(err.response.data.message);
    } finally { setSubmitting(false); }
  };

  const columns = [
    { title: t('name'), dataIndex: 'name', key: 'name' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
    { title: t('order'), dataIndex: 'order', key: 'order' },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('budgets')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_budget')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_budget')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label={t('name')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="notes" label={t('notes')}><Input.TextArea rows={2} /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BudgetList;
