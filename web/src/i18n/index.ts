import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import en_US from './locales/en_US.json';
import zh_CN from './locales/zh_CN.json';
import de_DE from './locales/de_DE.json';

export const languages = [
  { code: 'en_US', label: 'English (US)' },
  { code: 'zh_CN', label: '简体中文' },
  { code: 'de_DE', label: 'Deutsch' },
];

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en_US: { translation: en_US },
      zh_CN: { translation: zh_CN },
      de_DE: { translation: de_DE },
    },
    fallbackLng: 'en_US',
    interpolation: { escapeValue: false },
    detection: {
      order: ['localStorage', 'navigator'],
      lookupLocalStorage: 'firefly_language',
      caches: ['localStorage'],
    },
  });

export default i18n;
