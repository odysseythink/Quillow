import React, { useState, useRef, useEffect } from 'react';
import { Button, Input, Card, Tag, Space, message } from 'antd';
import { MessageOutlined, CloseOutlined, SendOutlined, RobotOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import client from '../api/client';

interface ChatMessage {
  id: number;
  role: 'user' | 'ai';
  content?: string;
  intent?: string;
  parsed?: any;
  confidence?: string;
  created?: boolean;
  transactionId?: number;
}

const ChatBubble: React.FC = () => {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [sending, setSending] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);
  const msgId = useRef(0);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const send = async () => {
    const text = input.trim();
    if (!text || sending) return;

    const userMsg: ChatMessage = { id: ++msgId.current, role: 'user', content: text };
    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setSending(true);

    try {
      const res = await client.post('/ai/chat', { message: text });
      const data = res.data;
      const aiMsg: ChatMessage = {
        id: ++msgId.current,
        role: 'ai',
        intent: data.intent,
        content: data.answer || '',
        parsed: data.parsed,
        confidence: data.confidence,
        created: data.created,
        transactionId: data.transaction_id,
      };

      if (data.intent === 'record' && data.confidence === 'high' && data.created) {
        aiMsg.content = `${t('chat_recorded')}: ${data.parsed?.category || ''} ¥${data.parsed?.amount || '0'} ${data.parsed?.description || ''} ${data.parsed?.date || ''}`;
      } else if (data.intent === 'record') {
        aiMsg.content = '';
      } else if (data.intent === 'query') {
        aiMsg.content = data.answer || t('chat_query_failed');
      }

      setMessages(prev => [...prev, aiMsg]);
    } catch {
      setMessages(prev => [...prev, { id: ++msgId.current, role: 'ai', content: t('error_occurred') }]);
    } finally {
      setSending(false);
    }
  };

  const confirmTransaction = async (msg: ChatMessage) => {
    if (!msg.parsed) return;
    try {
      const p = msg.parsed;
      await client.post('/transactions', {
        transactions: [{
          type: p.type || 'withdrawal',
          description: p.description || '',
          date: p.date || new Date().toISOString().slice(0, 10),
          amount: p.amount || '0',
          source_name: p.source_name || '',
          destination_name: p.destination_name || '',
          category_id: p.category_id || undefined,
        }],
      });
      message.success(t('chat_recorded'));
      setMessages(prev => prev.map(m => m.id === msg.id ? { ...m, created: true, content: `${t('chat_recorded')}: ¥${p.amount} ${p.description}` } : m));
    } catch {
      message.error(t('error_occurred'));
    }
  };

  const undoTransaction = async (msg: ChatMessage) => {
    if (!msg.transactionId) return;
    try {
      await client.delete(`/transactions/${msg.transactionId}`);
      message.success(t('chat_undo'));
      setMessages(prev => prev.map(m => m.id === msg.id ? { ...m, created: false, content: `[${t('chat_undo')}] ${m.content}` } : m));
    } catch {
      message.error(t('error_occurred'));
    }
  };

  const renderMessage = (msg: ChatMessage) => {
    if (msg.role === 'user') {
      return (
        <div key={msg.id} style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: 8 }}>
          <div style={{ background: '#1677ff', color: '#fff', padding: '6px 12px', borderRadius: 12, maxWidth: '80%' }}>{msg.content}</div>
        </div>
      );
    }

    // AI preview card (low confidence record)
    if (msg.intent === 'record' && !msg.created && msg.parsed) {
      const p = msg.parsed;
      return (
        <div key={msg.id} style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: 8 }}>
          <Card size="small" style={{ maxWidth: '85%' }} title={<Space><RobotOutlined />{t('chat_confirm')}</Space>}>
            <p>{t('type')}: <Tag>{p.type}</Tag></p>
            <p>{t('description')}: {p.description}</p>
            <p>{t('amount')}: ¥{p.amount}</p>
            <p>{t('date')}: {p.date}</p>
            {p.category && <p>{t('categories')}: {p.category}</p>}
            <Space style={{ marginTop: 8 }}>
              <Button type="primary" size="small" onClick={() => confirmTransaction(msg)}>{t('chat_confirm')}</Button>
              <Button size="small" onClick={() => setMessages(prev => prev.filter(m => m.id !== msg.id))}>{t('cancel')}</Button>
            </Space>
          </Card>
        </div>
      );
    }

    // AI text response
    return (
      <div key={msg.id} style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: 8 }}>
        <div style={{ background: '#f0f0f0', padding: '6px 12px', borderRadius: 12, maxWidth: '80%' }}>
          <RobotOutlined style={{ marginRight: 4 }} />
          {msg.content}
          {msg.created && msg.transactionId && (
            <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => undoTransaction(msg)}>{t('chat_undo')}</Button>
          )}
        </div>
      </div>
    );
  };

  return (
    <>
      {/* Floating button */}
      {!open && (
        <Button
          type="primary"
          shape="circle"
          size="large"
          icon={<MessageOutlined />}
          onClick={() => setOpen(true)}
          style={{ position: 'fixed', bottom: 24, right: 24, zIndex: 1000, width: 56, height: 56, fontSize: 24 }}
        />
      )}

      {/* Chat panel */}
      {open && (
        <div style={{ position: 'fixed', bottom: 24, right: 24, width: 360, height: 480, background: '#fff', borderRadius: 12, boxShadow: '0 4px 24px rgba(0,0,0,0.15)', zIndex: 1000, display: 'flex', flexDirection: 'column' }}>
          {/* Header */}
          <div style={{ padding: '12px 16px', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space><RobotOutlined /> Quillow AI</Space>
            <Button type="text" size="small" icon={<CloseOutlined />} onClick={() => setOpen(false)} />
          </div>

          {/* Messages */}
          <div ref={listRef} style={{ flex: 1, overflowY: 'auto', padding: 12 }}>
            {messages.length === 0 && (
              <div style={{ textAlign: 'center', color: '#999', marginTop: 40 }}>{t('chat_placeholder')}</div>
            )}
            {messages.map(renderMessage)}
          </div>

          {/* Input */}
          <div style={{ padding: '8px 12px', borderTop: '1px solid #f0f0f0' }}>
            <Input
              value={input}
              onChange={e => setInput(e.target.value)}
              onPressEnter={send}
              placeholder={t('chat_placeholder')}
              suffix={<Button type="text" icon={<SendOutlined />} onClick={send} loading={sending} />}
            />
          </div>
        </div>
      )}
    </>
  );
};

export default ChatBubble;
