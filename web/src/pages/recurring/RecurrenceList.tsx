import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, Switch, DatePicker } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getRecurrences, createRecurrence } from '../../api/general';
import dayjs from 'dayjs';

const { Title } = Typography;

const RecurrenceList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getRecurrences(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ active: true, first_date: dayjs() }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createRecurrence({
        title: v.title,
        description: v.description || '',
        first_date: v.first_date.format('YYYY-MM-DD'),
        repeat_until: v.repeat_until?.format('YYYY-MM-DD') || '',
        active: v.active,
        apply_rules: v.apply_rules ?? true,
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
    { title: t('first_date'), dataIndex: 'first_date', key: 'first_date' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('recurring_transactions')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_recurrence')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_recurrence')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="title" label={t('title')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="description" label={t('description')}><Input.TextArea rows={2} /></Form.Item>
          <Form.Item name="first_date" label={t('first_date')} rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="repeat_until" label={t('repeat_until')}><DatePicker style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="apply_rules" label={t('apply_rules')} valuePropName="checked"><Switch /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default RecurrenceList;
