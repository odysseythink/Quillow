import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Typography, message, Modal, Form, Input, InputNumber, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getCurrencies, createCurrency } from '../../api/general';

const { Title } = Typography;

const CurrencyList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getCurrencies(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); form.setFieldsValue({ decimal_places: 2, enabled: true }); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createCurrency({
        name: v.name,
        code: v.code,
        symbol: v.symbol,
        decimal_places: v.decimal_places,
        enabled: v.enabled,
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
    { title: t('code'), dataIndex: 'code', key: 'code' },
    { title: t('symbol'), dataIndex: 'symbol', key: 'symbol' },
    { title: t('decimal_places'), dataIndex: 'decimal_places', key: 'decimal_places' },
    { title: t('enabled'), dataIndex: 'enabled', key: 'enabled', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('currencies')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_currency')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_currency')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label={t('name')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="code" label={t('code')} rules={[{ required: true, message: t('required') }, { min: 3, max: 32 }]}><Input placeholder="e.g. USD" /></Form.Item>
          <Form.Item name="symbol" label={t('symbol')} rules={[{ required: true, message: t('required') }]}><Input placeholder="e.g. $" /></Form.Item>
          <Form.Item name="decimal_places" label={t('decimal_places')}><InputNumber style={{ width: '100%' }} min={0} max={10} /></Form.Item>
          <Form.Item name="enabled" label={t('enabled')} valuePropName="checked"><Switch /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CurrencyList;
