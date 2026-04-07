import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, InputNumber, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getWebhooks, createWebhook } from '../../api/general';

const { Title } = Typography;

const WebhookList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getWebhooks(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ active: true, trigger: 1, response: 1, delivery: 1 }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createWebhook({
        title: v.title,
        url: v.url,
        active: v.active,
        trigger: v.trigger,
        response: v.response,
        delivery: v.delivery,
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
    { title: t('url'), dataIndex: 'url', key: 'url' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('webhooks')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_webhook')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_webhook')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="title" label={t('title')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="url" label={t('url')} rules={[{ required: true, message: t('required') }, { type: 'url', message: 'Enter a valid URL' }]}><Input placeholder="https://example.com/webhook" /></Form.Item>
          <Form.Item name="trigger" label={t('trigger')} rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item>
          <Form.Item name="response" label={t('response')} rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item>
          <Form.Item name="delivery" label={t('delivery')} rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default WebhookList;
