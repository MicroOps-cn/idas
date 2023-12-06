import { Menu, Spin } from 'antd';
import { stringify } from 'querystring';
import type { MenuInfo } from 'rc-menu/lib/interface';
import React, { useCallback } from 'react';
import { history, useIntl, useModel } from 'umi';

import { loginPath } from '@/../config/env';
import Avatar from '@/components/Avatar';
import { userLogout as outLogin } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { LogoutOutlined, SecurityScanOutlined, SettingOutlined } from '@ant-design/icons';
import { useLocation, useSearchParams } from '@umijs/max';

import HeaderDropdown from '../HeaderDropdown';
import styles from './index.less';

export type GlobalHeaderRightProps = {
  menu?: boolean;
};

const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({ menu }) => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const intl = new IntlContext('', useIntl());
  const { pathname } = useLocation();
  const [searchParams] = useSearchParams();

  /**
   * 退出登录，并且将当前的 url 保存
   */
  const loginOut = useCallback(
    async (callback: () => void) => {
      await outLogin();
      callback();
      const { search } = history.location;
      const redirect = searchParams.get('redirect');
      // Note: There may be security issues, please note
      if (pathname !== loginPath && !redirect) {
        history.replace({
          pathname: loginPath,
          search: stringify({
            redirect: pathname + search,
          }),
        });
      }
    },
    [pathname, searchParams],
  );
  const onMenuClick = useCallback(
    (event: MenuInfo) => {
      const { key } = event;
      if (key === 'logout') {
        loginOut(() => {
          setInitialState((s) => ({ ...s, currentUser: undefined }));
        });
        return;
      }
      history.push(`/account/${key}`);
    },
    [setInitialState, loginOut],
  );

  const loading = (
    <span className={`${styles.action} ${styles.account}`}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  if (!initialState) {
    return loading;
  }

  const { currentUser } = initialState;

  if (!currentUser || !currentUser.username) {
    return loading;
  }

  const menuHeaderDropdown = (
    <Menu className={styles.menu} selectedKeys={[]} onClick={onMenuClick}>
      {menu && (
        <>
          <Menu.Item key="settings">
            <SettingOutlined />
            {intl.t('menu.account.setting', 'Account Settings')}
          </Menu.Item>
          <Menu.Item key="events">
            <SecurityScanOutlined />
            {intl.t('menu.account.events', 'Account Events')}
          </Menu.Item>
        </>
      )}
      {menu && <Menu.Divider />}

      <Menu.Item key="logout">
        <LogoutOutlined />
        {intl.t('menu.account.loginout', 'Logout')}
      </Menu.Item>
    </Menu>
  );
  return (
    <HeaderDropdown overlay={() => menuHeaderDropdown}>
      <span className={`${styles.action} ${styles.account}`}>
        <Avatar size="small" className={styles.avatar} src={`${currentUser.avatar}`} alt="avatar" />
        <span className={`${styles.name} anticon`}>
          {currentUser.fullName ?? currentUser.username}
        </span>
      </span>
    </HeaderDropdown>
  );
};

export default AvatarDropdown;
