import { Space } from 'antd';
import React from 'react';

import { AppstoreOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { useModel, history } from '@umijs/max';

import SelectLang from '../SelectLang';
import Avatar from './AvatarDropdown';
import styles from './index.less';

export const getActions = (role?: string) => {
  return [
    <div
      className={styles.action}
      hidden={!role}
      onClick={() => {
        history.push('/apps');
      }}
      title={'Dashboard'}
      key={'admin'}
    >
      <AppstoreOutlined />
    </div>,
    <div className={styles.action} title={'Help'} key={'help'}>
      <QuestionCircleOutlined />
    </div>,
    <Avatar menu={true} key={'avatar'} />,
    //{/* <NoticeIconView /> */}
    <SelectLang className={styles.action} key={'lang'} />,
  ];
};

const GlobalHeaderRight: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  if (!initialState || !initialState.settings) {
    return null;
  }
  const { role } = initialState.currentUser ?? {};
  return <Space className={styles.right}>{getActions(role)}</Space>;
};
export default GlobalHeaderRight;
