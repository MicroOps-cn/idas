import { Menu } from 'antd';
import React, { useRef, useState } from 'react';

import { IntlContext } from '@/utils/intl';
import { GridContent } from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';

import SecurityView from './components/SecurityView';
import styles from './style.less';

const { Item } = Menu;
const UserSetting: React.FC = ({}) => {
  const intl = new IntlContext('settings', useIntl());
  const mainRef = useRef<HTMLDivElement | null>(null);
  const menuMap = {
    base: intl.t('menuMap.basic', 'Basic Settings'),
    security: intl.t('menuMap.security', 'Security Settings'),
  };

  const [selectKey, setSelectKey] = useState<keyof typeof menuMap>('security');
  const renderChildren = () => {
    switch (selectKey) {
      case 'security':
        return <SecurityView parentIntl={intl} />;
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
  return (
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
  );
};

export default UserSetting;
