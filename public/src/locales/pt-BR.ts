import app from './pt-BR/app';
import component from './pt-BR/component';
import globalHeader from './pt-BR/globalHeader';
import menu from './pt-BR/menu';
import pages from './pt-BR/pages';
import pwa from './pt-BR/pwa';
import roles from './pt-BR/roles';
import settingDrawer from './pt-BR/settingDrawer';
import settings from './pt-BR/settings';
import users from './pt-BR/users';

export default {
  'navBar.lang': 'Idiomas',
  'layout.user.link.help': 'ajuda',
  'layout.user.link.privacy': 'política de privacidade',
  'layout.user.link.terms': 'termos de serviços',
  'app.preview.down.block': 'Download this page to your local project',
  ...globalHeader,
  ...users,
  ...roles,
  ...app,
  ...menu,
  ...settingDrawer,
  ...settings,
  ...pwa,
  ...component,
  ...pages,
};
