import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Typography, Descriptions, Form, Input, Button, message, Tabs } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import client from '../../api/client';

const { Title } = Typography;

const Profile: React.FC = () => {
  const { t } = useTranslation();
  const [user, setUser] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(true);
  const [passwordForm] = Form.useForm();
  const [emailForm] = Form.useForm();
  const [changingPassword, setChangingPassword] = useState(false);
  const [changingEmail, setChangingEmail] = useState(false);

  const loadUser = () => {
    setLoading(true);
    client.get('/profile').then(res => {
      const data = res.data?.data;
      setUser(data ? { id: data.id, ...data.attributes } : {});
    }).catch(() => {
      message.error(t('error_occurred'));
    }).finally(() => setLoading(false));
  };

  useEffect(() => { loadUser(); }, []);

  const onChangePassword = async (values: any) => {
    setChangingPassword(true);
    try {
      await client.post('/profile/change-password', {
        current_password: values.current_password,
        new_password: values.new_password,
      });
      message.success(t('password_changed'));
      passwordForm.resetFields();
    } catch (err: any) {
      message.error(err.response?.data?.message || t('error_occurred'));
    } finally {
      setChangingPassword(false);
    }
  };

  const onChangeEmail = async (values: any) => {
    setChangingEmail(true);
    try {
      await client.post('/profile/change-email', {
        password: values.password,
        new_email: values.new_email,
      });
      message.success(t('email_changed'));
      emailForm.resetFields();
      loadUser();
    } catch (err: any) {
      message.error(err.response?.data?.message || t('error_occurred'));
    } finally {
      setChangingEmail(false);
    }
  };

  const tabItems = [
    {
      key: 'info',
      label: t('profile'),
      icon: <UserOutlined />,
      children: (
        <Card loading={loading}>
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label={t('user_id')}>{user.id || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('email')}>{user.email || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('role')}>{user.role || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('blocked')}>{user.blocked ? t('yes') : t('no')}</Descriptions.Item>
            <Descriptions.Item label={t('created')}>{user.created_at || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('updated')}>{user.updated_at || '-'}</Descriptions.Item>
          </Descriptions>
        </Card>
      ),
    },
    {
      key: 'password',
      label: t('change_password'),
      icon: <LockOutlined />,
      children: (
        <Card>
          <Form form={passwordForm} layout="vertical" onFinish={onChangePassword} style={{ maxWidth: 400 }}>
            <Form.Item name="current_password" label={t('current_password')} rules={[{ required: true, message: t('required') }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item name="new_password" label={t('new_password')} rules={[{ required: true, message: t('required') }, { min: 8, message: t('min_password') }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item
              name="confirm_password"
              label={t('confirm_password')}
              dependencies={['new_password']}
              rules={[
                { required: true, message: t('required') },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue('new_password') === value) return Promise.resolve();
                    return Promise.reject(new Error(t('passwords_not_match')));
                  },
                }),
              ]}
            >
              <Input.Password />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={changingPassword}>{t('change_password')}</Button>
            </Form.Item>
          </Form>
        </Card>
      ),
    },
    {
      key: 'email',
      label: t('change_email'),
      icon: <MailOutlined />,
      children: (
        <Card>
          <Form form={emailForm} layout="vertical" onFinish={onChangeEmail} style={{ maxWidth: 400 }}>
            <Form.Item name="new_email" label={t('new_email')} rules={[{ required: true, message: t('required') }, { type: 'email', message: t('invalid_email') }]}>
              <Input />
            </Form.Item>
            <Form.Item name="password" label={t('current_password')} rules={[{ required: true, message: t('required') }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={changingEmail}>{t('change_email')}</Button>
            </Form.Item>
          </Form>
        </Card>
      ),
    },
  ];

  return (
    <div>
      <Title level={4}>{t('profile')}</Title>
      <Row>
        <Col span={24}>
          <Tabs items={tabItems} />
        </Col>
      </Row>
    </div>
  );
};

export default Profile;
