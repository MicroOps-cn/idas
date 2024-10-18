import { Menu } from 'antd';
import React, { useEffect, useRef, useState } from 'react';

import defaultSettings from '@/../config/defaultSettings';
import { getActions } from '@/components/RightContent';
import { currentUser as fetchCurrentUser } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { getPublicPath } from '@/utils/request';
import { GridContent, ProLayout } from '@ant-design/pro-components';
import { useModel, useIntl, history } from '@umijs/max';

import BaseView from './components/BaseView';
import BindingView from './components/BindingView';
import SecurityView from './components/SecurityView';
import SessionView from './components/SessionView';
import styles from './style.less';

const { Item } = Menu;
const UserSetting: React.FC = ({}) => {
  const intl = new IntlContext('user.settings', useIntl());
  const mainRef = useRef<HTMLDivElement | null>(null);
  const menuMap = {
    base: intl.t('menuMap.basic', 'Basic Settings'),
    security: intl.t('menuMap.security', 'Security Settings'),
    binding: intl.t('menuMap.binding', 'Account Binding'),
    sessions: intl.t('menuMap.sessions', 'Sessions'),
  };

  const {
    initialState = {
      settings: { navTheme: 'light' },
      globalConfig: undefined,
      currentUser: undefined,
    },
    setInitialState,
  } = useModel('@@initialState');
  // @ts-expect-error
  const { navTheme = 'light' } = initialState.settings;
  const [currentUser, setCurrentUser] = useState<
    | (API.UserInfo & {
        fetchTime: number;
      })
    | undefined
  >(initialState?.currentUser);

  const reload = () => {
    const { fetchTime } = currentUser ?? {};
    if (fetchTime && fetchTime + 1 > Number(new Date()) / 1000) {
      return;
    }
    fetchCurrentUser().then((resp) => {
      if (resp.data) {
        setInitialState({
          ...initialState,
          currentUser: { ...resp.data, fetchTime: Number(new Date()) / 1000 },
        });
        setCurrentUser({ ...resp.data, fetchTime: Number(new Date()) / 1000 });
      }
    });
  };

  useEffect(() => {
    reload();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  const [selectKey, setSelectKey] = useState<keyof typeof menuMap>('base');
  const renderChildren = () => {
    switch (selectKey) {
      case 'base':
        return <BaseView parentIntl={intl} currentUser={currentUser} reload={reload} />;
      case 'security':
        return <SecurityView parentIntl={intl} currentUser={currentUser} reload={reload} />;
      case 'binding':
        return <BindingView />;
      case 'sessions':
        return <SessionView parentIntl={intl} currentUser={currentUser} />;
      default:
        break;
    }

    return null;
  };

  const getMenu = () => {
    return Object.keys(menuMap).map((item) => <Item key={item}>{menuMap[item]}</Item>);
  };
  const getRightTitle = () => {
    return menuMap[selectKey];
  };

  const globalConfig = initialState?.globalConfig ?? null;
  return (
    <ProLayout
      logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
      title={globalConfig?.title ?? defaultSettings.title}
      onMenuHeaderClick={() => {
        history.push('/');
      }}
      layout="top"
      navTheme={navTheme}
      actionsRender={() => getActions(initialState?.currentUser?.role)}
    >
      <GridContent>
        <div className={styles.main} ref={mainRef}>
          <div className={styles.leftMenu}>
            <Menu
              mode={'inline'}
              selectedKeys={[selectKey]}
              onClick={({ key }) => setSelectKey(key as typeof selectKey)}
            >
              {getMenu()}
            </Menu>
          </div>
          <div className={styles.right}>
            <div className={styles.title}>{getRightTitle()}</div>
            {renderChildren()}
          </div>
        </div>
      </GridContent>
    </ProLayout>
  );
};

export default UserSetting;
