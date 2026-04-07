import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, Select, Switch, DatePicker } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getBills, createBill } from '../../api/general';
import dayjs from 'dayjs';
import { useTranslation } from 'react-i18next';

const { Title } = Typography;

const BillList: React.FC = () => {
  const { t } = useTranslation();

  const repeatFreqs = [
    { value: 'weekly', label: t('weekly') },
    { value: 'monthly', label: t('monthly') },
    { value: 'quarterly', label: t('quarterly') },
    { value: 'half-year', label: t('half_year') },
    { value: 'yearly', label: t('yearly') },
  ];

  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getBills(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ active: true, repeat_freq: 'monthly', date: dayjs() }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createBill({
        name: v.name,
        amount_min: v.amount_min,
        amount_max: v.amount_max,
        date: v.date.format('YYYY-MM-DD'),
        repeat_freq: v.repeat_freq,
        active: v.active,
        notes: v.notes || '',
      });
      message.success(t('success_created'));
      setModalOpen(false);
      fetchData();
    } catch (err: any) {
      if (err.response?.data?.message) message.error(err.response.data.message);
    } finally { setSubmitting(false); }
  };

  const columns = [
    { title: t('name'), dataIndex: 'name', key: 'name' },
    { title: t('amount_min'), dataIndex: 'amount_min', key: 'amount_min' },
    { title: t('amount_max'), dataIndex: 'amount_max', key: 'amount_max' },
    { title: t('repeat_freq'), dataIndex: 'repeat_freq', key: 'repeat_freq' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('bills')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_bill')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_bill')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label={t('name')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="amount_min" label={t('amount_min')} rules={[{ required: true, message: t('required') }]}><Input placeholder="0.00" /></Form.Item>
          <Form.Item name="amount_max" label={t('amount_max')} rules={[{ required: true, message: t('required') }]}><Input placeholder="0.00" /></Form.Item>
          <Form.Item name="date" label={t('date')} rules={[{ required: true }]}><DatePicker style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="repeat_freq" label={t('repeat_freq')} rules={[{ required: true }]}><Select options={repeatFreqs} /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="notes" label={t('notes')}><Input.TextArea rows={2} /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BillList;
