import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Statistic, Typography } from 'antd';
import { WalletOutlined, ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons';
import { getSummary } from '../../api/transactions';

const { Title } = Typography;

const Dashboard: React.FC = () => {
  const [summary, setSummary] = useState<Record<string, any>>({});

  useEffect(() => {
    getSummary().then(res => setSummary(res.data)).catch(() => {});
  }, []);

  const summaryEntries = Object.entries(summary);
  const spent = summaryEntries.filter(([k]) => k.startsWith('spent-'));
  const earned = summaryEntries.filter(([k]) => k.startsWith('earned-'));

  return (
    <div>
      <Title level={4}>Dashboard</Title>
      <Row gutter={[16, 16]}>
        {spent.map(([key, val]: [string, any]) => (
          <Col span={8} key={key}>
            <Card><Statistic title={val.title || 'Spent'} value={val.monetary_value || '0'} prefix={<ArrowDownOutlined />} valueStyle={{ color: '#cf1322' }} suffix={val.currency_code} /></Card>
          </Col>
        ))}
        {earned.map(([key, val]: [string, any]) => (
          <Col span={8} key={key}>
            <Card><Statistic title={val.title || 'Earned'} value={val.monetary_value || '0'} prefix={<ArrowUpOutlined />} valueStyle={{ color: '#3f8600' }} suffix={val.currency_code} /></Card>
          </Col>
        ))}
        {summaryEntries.length === 0 && (
          <Col span={24}>
            <Card><Statistic title="Balance" value="0.00" prefix={<WalletOutlined />} /></Card>
          </Col>
        )}
      </Row>
    </div>
  );
};

export default Dashboard;
