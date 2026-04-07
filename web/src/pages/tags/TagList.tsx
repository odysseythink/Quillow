import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Typography, message, Modal, Form, Input } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getTags, createTag } from '../../api/general';

const { Title } = Typography;

const TagList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState<any>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchData = (page = 1) => {
    setLoading(true);
    getTags(page).then(res => {
      const data = res.data;
      if (data.data) {
        setItems(data.data.map((r: any) => ({ id: r.id, ...r.attributes })));
        setPagination(data.meta?.pagination);
      }
    }).catch(() => {}).finally(() => setLoading(false));
  };

  useEffect(() => { fetchData(); }, []);

  const openCreate = () => { form.resetFields(); setModalOpen(true); };

  const handleSubmit = async () => {
    try {
      const v = await form.validateFields();
      setSubmitting(true);
      await createTag({ tag: v.tag, description: v.description || '' });
      message.success(t('success_created'));
      setModalOpen(false);
      fetchData();
    } catch (err: any) {
      if (err.response?.data?.message) message.error(err.response.data.message);
    } finally { setSubmitting(false); }
  };

  const columns = [
    { title: t('tag'), dataIndex: 'tag', key: 'tag' },
    { title: t('description'), dataIndex: 'description', key: 'description' },
  ];

  return (
    <div>
      <Space style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Title level={4} style={{ margin: 0 }}>{t('tags')}</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>{t('add_tag')}</Button>
      </Space>
      <Table dataSource={items} columns={columns} rowKey="id" loading={loading}
        pagination={pagination ? { total: pagination.total, pageSize: pagination.per_page, current: pagination.current_page, onChange: fetchData } : false}
      />
      <Modal title={t('create_tag')} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="tag" label={t('tag')} rules={[{ required: true, message: t('required') }]}><Input /></Form.Item>
          <Form.Item name="description" label={t('description')}><Input.TextArea rows={2} /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default TagList;
