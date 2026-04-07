import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, Switch, DatePicker, InputNumber } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getPiggyBanks, createPiggyBank } from '../../api/general';

const { Title } = Typography;

const PiggyBankList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getPiggyBanks(page).then(res => {
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
      await createPiggyBank({
        name: v.name,
        account_id: v.account_id,
        target_amount: v.target_amount || '0',
        target_date: v.target_date?.format('YYYY-MM-DD') || '',
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
    { title: t('target_amount'), dataIndex: 'target_amount', key: 'target_amount' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('piggy_banks')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_piggy_bank')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_piggy_bank')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label={t('name')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="account_id" label={t('account_id')} rules={[{ required: true, message: t('required') }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item>
          <Form.Item name="target_amount" label={t('target_amount')}><Input placeholder="0.00" /></Form.Item>
          <Form.Item name="target_date" label={t('target_date')}><DatePicker style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="active" label={t('active')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="notes" label={t('notes')}><Input.TextArea rows={2} /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PiggyBankList;
