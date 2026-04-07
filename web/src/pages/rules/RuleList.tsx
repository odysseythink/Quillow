import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, Switch, InputNumber } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getRules, createRule } from '../../api/general';

const { Title } = Typography;

const RuleList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getRules(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ active: true, strict: false }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createRule({
        title: v.title,
        description: v.description || '',
        rule_group_id: v.rule_group_id,
        active: v.active,
        strict: v.strict,
        triggers: v.trigger_type ? [{ type: v.trigger_type, value: v.trigger_value || '' }] : [],
        actions: v.action_type ? [{ type: v.action_type, value: v.action_value || '' }] : [],
      });
      message.success(t('success_created'));
      setModalOpen(false);
      fetchData();
    } catch (err: any) {
      if (err.response?.data?.message) message.error(err.response.data.message);
    } finally { setSubmitting(false); }
  };

  const columns = [
    { title: t('title'), dataIndex: 'title', key: 'title' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
    { title: t('order'), dataIndex: 'order', key: 'order' },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('rules')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_rule')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_rule')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose width={520}>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="title" label={t('title')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="description" label={t('description')}><Input.TextArea rows={2} /></Form.Item>
          <Form.Item name="rule_group_id" label={t('rule_group_id')} rules={[{ required: true, message: t('required') }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="strict" label={t('strict')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="trigger_type" label={t('trigger_type')}><Input placeholder="e.g. description_contains" /></Form.Item>
          <Form.Item name="trigger_value" label={t('trigger_value')}><Input /></Form.Item>
          <Form.Item name="action_type" label={t('action_type')}><Input placeholder="e.g. set_category" /></Form.Item>
          <Form.Item name="action_value" label={t('action_value')}><Input /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default RuleList;
