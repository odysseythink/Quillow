import client from './client';

export const getBills = (page = 1, limit = 50) => client.get('/bills', { params: { page, limit } });
export const getBill = (id: string) => client.get(`/bills/${id}`);
export const createBill = (data: any) => client.post('/bills', data);
export const updateBill = (id: string, data: any) => client.put(`/bills/${id}`, data);
export const deleteBill = (id: string) => client.delete(`/bills/${id}`);

export const getCategories = (page = 1, limit = 50) => client.get('/categories', { params: { page, limit } });
export const getCategory = (id: string) => client.get(`/categories/${id}`);
export const createCategory = (data: any) => client.post('/categories', data);
export const updateCategory = (id: string, data: any) => client.put(`/categories/${id}`, data);
export const deleteCategory = (id: string) => client.delete(`/categories/${id}`);

export const getTags = (page = 1, limit = 50) => client.get('/tags', { params: { page, limit } });
export const getTag = (id: string) => client.get(`/tags/${id}`);
export const createTag = (data: any) => client.post('/tags', data);
export const deleteTag = (id: string) => client.delete(`/tags/${id}`);

export const getPiggyBanks = (page = 1, limit = 50) => client.get('/piggy-banks', { params: { page, limit } });
export const getPiggyBank = (id: string) => client.get(`/piggy-banks/${id}`);
export const createPiggyBank = (data: any) => client.post('/piggy-banks', data);
export const deletePiggyBank = (id: string) => client.delete(`/piggy-banks/${id}`);

export const getCurrencies = (page = 1, limit = 50) => client.get('/currencies', { params: { page, limit } });
export const getRules = (page = 1, limit = 50) => client.get('/rules', { params: { page, limit } });
export const getRuleGroups = (page = 1, limit = 50) => client.get('/rule-groups', { params: { page, limit } });
export const getWebhooks = (page = 1, limit = 50) => client.get('/webhooks', { params: { page, limit } });
export const getRecurrences = (page = 1, limit = 50) => client.get('/recurrences', { params: { page, limit } });

export const getAbout = () => client.get('/about');
export const getPreferences = () => client.get('/preferences');
