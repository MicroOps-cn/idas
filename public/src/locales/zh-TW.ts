import app from './zh-TW/app';
import component from './zh-TW/component';
import globalHeader from './zh-TW/globalHeader';
import menu from './zh-TW/menu';
import pages from './zh-TW/pages';
import pwa from './zh-TW/pwa';
import roles from './zh-TW/roles';
import settingDrawer from './zh-TW/settingDrawer';
import settings from './zh-TW/settings';
import users from './zh-TW/users';

export default {
  'navBar.lang': '語言',
  'layout.user.link.help': '幫助',
  'layout.user.link.privacy': '隱私',
  'layout.user.link.terms': '條款',
  'app.preview.down.block': '下載此頁面到本地項目',
  ...globalHeader,
  ...users,
  ...roles,
  ...pages,
  ...app,
  ...menu,
  ...settingDrawer,
  ...settings,
  ...pwa,
  ...component,
};
