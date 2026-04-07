import React, { useState } from 'react';
import { Upload, Table, Button, Typography, message, Tag, Space, Checkbox } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import client from '../../api/client';

const { Title } = Typography;
const { Dragger } = Upload;

const Import: React.FC = () => {
  const { t } = useTranslation();
  const [previewing, setPreviewing] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [source, setSource] = useState('');
  const [transactions, setTransactions] = useState<any[]>([]);
  const [imported, setImported] = useState<number | null>(null);

  const handleUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    setPreviewing(true);
    setImported(null);
    try {
      const res = await client.post('/import/preview', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      setSource(res.data.source);
      setTransactions(res.data.transactions.map((tx: any) => ({ ...tx, selected: tx.selected !== false })));
    } catch (err: any) {
      message.error(err.response?.data?.message || t('error_occurred'));
    } finally {
      setPreviewing(false);
    }
    return false; // prevent default upload
  };

  const toggleSelect = (index: number) => {
    setTransactions(prev => prev.map(tx => tx.index === index ? { ...tx, selected: !tx.selected } : tx));
  };

  const toggleAll = (selected: boolean) => {
    setTransactions(prev => prev.map(tx => ({ ...tx, selected })));
  };

  const handleConfirm = async () => {
    const selected = transactions.filter(tx => tx.selected);
    if (selected.length === 0) { message.warning(t('required')); return; }
    setConfirming(true);
    try {
      const res = await client.post('/import/confirm', { transactions: selected });
      setImported(res.data.imported);
      message.success(`${t('import_success')}: ${res.data.imported}`);
    } catch (err: any) {
      message.error(err.response?.data?.message || t('error_occurred'));
    } finally {
      setConfirming(false);
    }
  };

  const columns = [
    {
      title: <Checkbox onChange={e => toggleAll(e.target.checked)} />,
      key: 'select', width: 40,
      render: (_: any, record: any) => <Checkbox checked={record.selected} onChange={() => toggleSelect(record.index)} />,
    },
    { title: t('date'), dataIndex: 'date', key: 'date', width: 110 },
    { title: t('type'), dataIndex: 'type', key: 'type', width: 80, render: (v: string) => <Tag>{v}</Tag> },
    { title: t('description'), dataIndex: 'description', key: 'description' },
    { title: t('import_counterparty'), dataIndex: 'counterparty', key: 'counterparty' },
    { title: t('amount'), dataIndex: 'amount', key: 'amount', width: 100 },
    { title: t('import_status'), dataIndex: 'status', key: 'status', width: 100 },
  ];

  return (
    <div>
      <Title level={4}>{t('import_bills')}</Title>

      {transactions.length === 0 && imported === null && (
        <Dragger
          accept=".csv"
          showUploadList={false}
          beforeUpload={handleUpload}
          disabled={previewing}
        >
          <p className="ant-upload-drag-icon"><InboxOutlined /></p>
          <p className="ant-upload-text">{t('upload_csv')}</p>
          <p className="ant-upload-hint">{t('upload_csv_hint')}</p>
        </Dragger>
      )}

      {transactions.length > 0 && (
        <>
          <Space style={{ marginBottom: 16 }}>
            <Tag color="blue">{t('import_source')}: {source}</Tag>
            <span>{transactions.filter(tx => tx.selected).length} / {transactions.length} {t('import_preview')}</span>
          </Space>
          <Table dataSource={transactions} columns={columns} rowKey="index" size="small" pagination={false} scroll={{ y: 400 }} />
          <Space style={{ marginTop: 16 }}>
            <Button type="primary" onClick={handleConfirm} loading={confirming}>{t('import_confirm')}</Button>
            <Button onClick={() => { setTransactions([]); setSource(''); }}>{t('cancel')}</Button>
          </Space>
        </>
      )}

      {imported !== null && transactions.length === 0 && (
        <div style={{ textAlign: 'center', marginTop: 40 }}>
          <Title level={5}>{t('import_success')}: {imported}</Title>
          <Button onClick={() => setImported(null)}>{t('import_bills')}</Button>
        </div>
      )}
    </div>
  );
};

export default Import;
