import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, Popconfirm, Typography, message, Modal, Form, Input, Select, InputNumber, Switch } from 'antd';
import { PlusOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { fetchAccounts } from '../../store/slices/accountSlice';
import { createAccount, updateAccount, deleteAccount } from '../../api/accounts';
import { getCurrencies } from '../../api/general';

const { Title } = Typography;

const accountTypeMap: Record<string, string> = {
  asset: 'Asset account',
  expense: 'Expense account',
  revenue: 'Revenue account',
  cash: 'Cash account',
  liability: 'Debt',
};

const liabilityTypes = [
  { value: 'Debt', label: 'Debt' },
  { value: 'Loan', label: 'Loan' },
  { value: 'Mortgage', label: 'Mortgage' },
];

interface CurrencyOption { id: number; code: string; name: string; symbol: string }

const AccountList: React.FC = () => {
  const { t } = useTranslation();
  const { type } = useParams<{ type: string }>();
  const dispatch = useAppDispatch();
  const { items, pagination, loading } = useAppSelector(s => s.accounts);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [currencies, setCurrencies] = useState<CurrencyOption[]>([]);
  const [form] = Form.useForm();

  const accountType = type || 'asset';
  const isAsset = accountType === 'asset';
  const isLiability = accountType === 'liability';

  const accountRoles = [
    { value: 'defaultAsset', label: t('default_asset') },
    { value: 'sharedAsset', label: t('shared_asset') },
    { value: 'savingAsset', label: t('savings_account') },
    { value: 'ccAsset', label: t('credit_card') },
    { value: 'cashWalletAsset', label: t('cash_wallet') },
  ];

  const reload = () => {
    dispatch(fetchAccounts({ type: accountTypeMap[accountType] || `${accountType} account` }));
  };

  useEffect(() => { reload(); }, [dispatch, type]);

  useEffect(() => {
    getCurrencies().then(res => {
      const list = res.data?.data || [];
      setCurrencies(list.map((r: any) => ({
        id: Number(r.id),
        code: r.attributes?.code || r.code,
        name: r.attributes?.name || r.name,
        symbol: r.attributes?.symbol || r.symbol,
      })));
    }).catch(() => {});
  }, []);

  const openCreate = () => {
    setEditingId(null);
    form.resetFields();
    form.setFieldsValue({
      active: true,
      account_role: isAsset ? 'defaultAsset' : undefined,
      liability_type: isLiability ? 'Debt' : undefined,
    });
    setModalOpen(true);
  };

  const openEdit = (record: any) => {
    setEditingId(record.id);
    form.setFieldsValue({
      name: record.name,
      iban: record.iban,
      account_number: record.account_number,
      virtual_balance: record.virtual_balance,
      active: record.active !== false,
      account_role: record.account_role || undefined,
      currency_id: record.currency_id ? Number(record.currency_id) : undefined,
      notes: record.notes,
      order: record.order,
    });
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);

      if (editingId) {
        await updateAccount(editingId, {
          name: values.name,
          iban: values.iban || '',
          account_number: values.account_number || '',
          virtual_balance: values.virtual_balance || '0',
          active: values.active,
          account_role: values.account_role || '',
          notes: values.notes || '',
          order: values.order || 0,
        });
        message.success(t('success_updated'));
      } else {
        const resolvedType = isLiability ? (values.liability_type || 'Debt') : (accountTypeMap[accountType] || `${accountType} account`);
        await createAccount({
          name: values.name,
          type: resolvedType,
          iban: values.iban || '',
          account_number: values.account_number || '',
          virtual_balance: values.virtual_balance || '0',
          active: values.active,
          account_role: values.account_role || '',
          currency_id: values.currency_id ? String(values.currency_id) : '',
          currency_code: '',
          notes: values.notes || '',
        } as any);
        message.success(t('success_created'));
      }

      setModalOpen(false);
      reload();
    } catch (err: any) {
      if (err.response?.data?.message) {
        message.error(err.response.data.message);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: string) => {
    await deleteAccount(id);
    message.success(t('success_deleted'));
    reload();
  };

  const columns = [
    { title: t('name'), dataIndex: 'name', key: 'name' },
    { title: t('type'), dataIndex: 'type', key: 'type' },
    { title: t('balance'), dataIndex: 'current_balance', key: 'balance', render: (v: string, r: any) => `${r.currency_symbol || ''}${v || '0'}` },
    { title: t('iban'), dataIndex: 'iban', key: 'iban' },
    { title: t('active'), dataIndex: 'active', key: 'active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? t('yes') : t('no')}</Tag> },
    { title: t('action'), key: 'action', width: 120, render: (_: any, record: any) => (
      <Space>
        <Button icon={<EditOutlined />} size="small" onClick={() => openEdit(record)} />
        <Popconfirm title={t('delete_confirm')} onConfirm={() => handleDelete(record.id)}>
          <Button danger icon={<DeleteOutlined />} size="small" />
        </Popconfirm>
      </Space>
    )},
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t(accountType + '_accounts')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_account')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page,
          onChange: (page, pageSize) => dispatch(fetchAccounts({ page, limit: pageSize, type: accountTypeMap[accountType] || `${accountType} account` }))
        } : false}
      />

      <Modal
        title={editingId ? t('edit_account') : t('create_account')}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
        confirmLoading={submitting}
        destroyOnClose
        width={560}
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="name" label={t('name')} rules={[{ required: true, message: t('required') }]}>
            <Input />
          </Form.Item>

          {isLiability && !editingId && (
            <Form.Item name="liability_type" label={t('liability_type')} rules={[{ required: true, message: t('required') }]}>
              <Select options={liabilityTypes} />
            </Form.Item>
          )}

          {isAsset && (
            <Form.Item name="account_role" label={t('account_role')}>
              <Select options={accountRoles} allowClear />
            </Form.Item>
          )}

          {!editingId && (
            <Form.Item name="currency_id" label={t('currency')}>
              <Select
                showSearch
                allowClear
                placeholder={t('currency')}
                optionFilterProp="label"
                options={currencies.map(c => ({ value: c.id, label: `${c.code} - ${c.name} (${c.symbol})` }))}
              />
            </Form.Item>
          )}

          <Form.Item name="iban" label={t('iban')}>
            <Input />
          </Form.Item>

          <Form.Item name="account_number" label={t('account_number')}>
            <Input />
          </Form.Item>

          <Form.Item name="virtual_balance" label={t('virtual_balance')}>
            <Input placeholder="0.00" />
          </Form.Item>

          <Form.Item name="order" label={t('order')}>
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>

          <Form.Item name="active" label={t('active')} valuePropName="checked">
            <Switch />
          </Form.Item>

          <Form.Item name="notes" label={t('notes')}>
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AccountList;
