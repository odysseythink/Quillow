import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import { useAppSelector } from './store/hooks';
import AppLayout from './components/layout/AppLayout';
import Login from './pages/auth/Login';
import Dashboard from './pages/dashboard/Dashboard';
import AccountList from './pages/accounts/AccountList';
import TransactionList from './pages/transactions/TransactionList';
import BudgetList from './pages/budgets/BudgetList';
import BillList from './pages/bills/BillList';
import CategoryList from './pages/categories/CategoryList';
import TagList from './pages/tags/TagList';
import PiggyBankList from './pages/piggyBanks/PiggyBankList';
import RuleList from './pages/rules/RuleList';
import RecurrenceList from './pages/recurring/RecurrenceList';
import CurrencyList from './pages/currencies/CurrencyList';
import WebhookList from './pages/webhooks/WebhookList';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

const App: React.FC = () => (
  <ConfigProvider theme={{ token: { colorPrimary: '#1677ff' } }}>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<ProtectedRoute><AppLayout /></ProtectedRoute>}>
          <Route index element={<Dashboard />} />
          <Route path="accounts/:type" element={<AccountList />} />
          <Route path="transactions/:type" element={<TransactionList />} />
          <Route path="budgets" element={<BudgetList />} />
          <Route path="bills" element={<BillList />} />
          <Route path="categories" element={<CategoryList />} />
          <Route path="tags" element={<TagList />} />
          <Route path="piggy-banks" element={<PiggyBankList />} />
          <Route path="rules" element={<RuleList />} />
          <Route path="recurring" element={<RecurrenceList />} />
          <Route path="currencies" element={<CurrencyList />} />
          <Route path="webhooks" element={<WebhookList />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </ConfigProvider>
);

export default App;
