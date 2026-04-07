import React, { useCallback, useEffect, useRef, useState } from 'react';
import { Table, Button, Space, Tag, Popconfirm, Typography, message, Modal, Form, Input, Select, DatePicker } from 'antd';
import { PlusOutlined, DeleteOutlined, RobotOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { fetchTransactions } from '../../store/slices/transactionSlice';
import { createTransaction, deleteTransaction } from '../../api/transactions';
import { aiSuggest, getCategories } from '../../api/general';
import dayjs from 'dayjs';

const { Title } = Typography;

const TransactionList: React.FC = () => {
  const { t } = useTranslation();
  const { type } = useParams<{ type: string }>();
  const dispatch = useAppDispatch();
  const { items, pagination, loading } = useAppSelector(s => s.transactions);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();
  const [aiSource, setAiSource] = useState<string | null>(null);
  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  const reload = () => dispatch(fetchTransactions({ type }));
  useEffect(() => { reload(); }, [dispatch, type]);

  useEffect(() => {
    getCategories(1, 1000).then(res => {
      const list = res.data?.data || [];
      setCategories(list.map((r: any) => ({ id: Number(r.id), name: r.attributes?.name || r.name })));
    }).catch(() => {});
  }, []);

  const handleDelete = async (id: string) => {
    await deleteTransaction(id);
    message.success(t('success_deleted'));
    reload();
  };

  const openCreate = () => {
    form.resetFields();
    form.setFieldsValue({ type: type || 'withdrawal', date: dayjs() });
    setAiSource(null);
    setModalOpen(true);
  };

  const handleDescriptionChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const desc = e.target.value;
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!desc || desc.length < 2) { setAiSource(null); return; }

    debounceRef.current = setTimeout(() => {
      aiSuggest(desc).then(res => {
        const data = res.data;
        if (data.source !== 'none' && data.category_id) {
          form.setFieldsValue({ category_id: data.category_id });
          setAiSource(data.source);
        }
      }).catch(() => {});
    }, 500);
  }, [form]);

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createTransaction({
        transactions: [{
          type: v.type,
          description: v.description,
          date: v.date.format('YYYY-MM-DD'),
          amount: v.amount,
          source_name: v.source_name || '',
          destination_name: v.destination_name || '',
          category_id: v.category_id || undefined,
          notes: v.notes || '',
        }],
      } as any);
      message.success(t('success_created'));
      setModalOpen(false);
      reload();
    } catch (err: any) {
      if (err.response?.data?.message) message.error(err.response.data.message);
    } finally {
      setSubmitting(false);
    }
  };

  const txTypes = [
    { value: 'withdrawal', label: t('withdrawal') },
    { value: 'deposit', label: t('deposit') },
    { value: 'transfer', label: t('transfer') },
  ];

  const columns = [
    { title: t('description'), key: 'desc', render: (_: any, r: any) => r.transactions?.[0]?.description || r.group_title || '-' },
    { title: t('amount'), key: 'amount', render: (_: any, r: any) => { const tx = r.transactions?.[0]; return tx ? `${tx.currency_symbol || ''}${tx.amount || '0'}` : '-'; }},
    { title: t('source_account'), key: 'source', render: (_: any, r: any) => r.transactions?.[0]?.source_name || '-' },
    { title: t('destination_account'), key: 'dest', render: (_: any, r: any) => r.transactions?.[0]?.destination_name || '-' },
    { title: t('date'), key: 'date', render: (_: any, r: any) => r.transactions?.[0]?.date?.substring(0, 10) || '-' },
    { title: t('type'), key: 'type', render: (_: any, r: any) => <Tag>{r.transactions?.[0]?.type || '-'}</Tag> },
    { title: t('action'), key: 'action', render: (_: any, record: any) => (
      <Popconfirm title={t('delete_confirm')} onConfirm={() => handleDelete(record.id)}>
        <Button danger icon={<DeleteOutlined />} size="small" />
      </Popconfirm>
    )},
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{type ? t(type === 'withdrawal' ? 'withdrawals' : type === 'deposit' ? 'deposits' : 'transfers') : t('transactions')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_transaction')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page,
          onChange: (page, pageSize) => dispatch(fetchTransactions({ page, limit: pageSize, type }))
        } : false}
      />
      <Modal title={t('create_transaction')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose width={520}>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="type" label={t('type')} rules={[{ required: true }]}>
            <Select options={txTypes} />
          </Form.Item>
          <Form.Item name="description" label={t('description')} rules={[{ required: true, message: t('required') }]}>
            <Input onChange={handleDescriptionChange} />
          </Form.Item>
          <Form.Item name="date" label={t('date')} rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="amount" label={t('amount')} rules={[{ required: true, message: t('required') }]}>
            <Input placeholder="0.00" />
          </Form.Item>
          <Form.Item name="source_name" label={t('source_account')}>
            <Input />
          </Form.Item>
          <Form.Item name="destination_name" label={t('destination_account')}>
            <Input />
          </Form.Item>
          <Form.Item
            name="category_id"
            label={
              <Space>
                {t('categories')}
                {aiSource && <Tag icon={<RobotOutlined />} color="blue">AI ({aiSource})</Tag>}
              </Space>
            }
          >
            <Select
              allowClear
              showSearch
              optionFilterProp="label"
              placeholder={t('categories')}
              options={categories.map(c => ({ value: c.id, label: c.name }))}
            />
          </Form.Item>
          <Form.Item name="notes" label={t('notes')}>
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default TransactionList;
