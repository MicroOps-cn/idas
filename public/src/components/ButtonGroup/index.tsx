import { Dropdown, Button as AntdButton, DropdownProps, Space } from 'antd';
import type { ButtonProps as AntdButtonProps } from 'antd';
import { message } from 'antd';
import { ButtonType } from 'antd/lib/button';
import { MenuItemType } from 'antd/lib/menu/interface';
import { isBoolean } from 'lodash';
import { MenuInfo } from 'rc-menu/lib/interface';
import { ReactNode } from 'react';
import { useState } from 'react';

import { DownOutlined } from '@ant-design/icons';

interface ButtonProps extends AntdButtonProps {
  success?: React.ReactElement;
  failed?: React.ReactElement;
  onClick: (event: React.MouseEvent<HTMLElement, MouseEvent>) => Promise<void | boolean>;
}

export const Button: React.FC<ButtonProps> = ({ success, failed, onClick, ...props }) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [status, setStatus] = useState<0 | 1 | 2>(0);
  return (
    <AntdButton
      {...props}
      loading={loading}
      onClick={async (e) => {
        if (onClick) {
          try {
            setLoading(true);
            const ok = await onClick(e);
            if (isBoolean(ok) && !ok) {
              setStatus(2);
            } else {
              setStatus(1);
            }
          } catch (error) {
            message.error(`${error}`, 3);
            setStatus(2);
          } finally {
            setLoading(false);
          }
        }
      }}
    >
      {status === 1 && success ? success : status === 2 && failed ? failed : props.children}
    </AntdButton>
  );
};

interface ButtonItemType extends Omit<MenuItemType, 'type'> {
  hidden?: boolean;
  failed?: React.ReactElement;
  success?: React.ReactElement;
  type?: ButtonType;
  onClick?: (info: Omit<MenuInfo, 'item' | 'keyPath'>) => void;
}

interface ButtonGroupProps extends DropdownProps {
  maxItems?: number;
  items: ButtonItemType[];
  moreLabel?: ReactNode;
}

const GroupButtonItem: React.FC<ButtonItemType> = ({
  success,
  failed,
  key,
  style,
  type = 'link',
  onClick,
  label,
}) => {
  return (
    <Button
      success={success}
      failed={failed}
      key={key}
      style={{
        padding: '4px 0px',
        ...style,
      }}
      type={type}
      onClick={async (e) => {
        return onClick?.({
          key: key.toString(),
          domEvent: e,
        });
      }}
    >
      {label}
    </Button>
  );
};

export const ButtonGroup: React.FC<ButtonGroupProps> = ({
  maxItems,
  items,
  moreLabel = 'More',
}) => {
  const visibleItems = items.filter((item) => item && !item.hidden);
  if (maxItems === undefined || maxItems >= visibleItems.length) {
    const buttons = visibleItems.map(GroupButtonItem);
    if (maxItems === undefined) {
      return <AntdButton.Group>{buttons}</AntdButton.Group>;
    }
    return <Space>{buttons}</Space>;
  }
  const dropdownItems = visibleItems.slice(maxItems - 1).map((item) => {
    return {
      ...item,
      label: <div style={{ padding: '0px 15px' }}>{item.label}</div>,
    };
  });
  return (
    <Space>
      {visibleItems.slice(0, maxItems - 1).map(GroupButtonItem)}
      <Dropdown
        menu={{ items: dropdownItems.map((item) => ({ ...item, type: 'item' })) }}
        trigger={['click']}
      >
        <a onClick={(e) => e.preventDefault()} style={{ gap: 3, display: 'inline-flex' }}>
          {moreLabel}
          <DownOutlined />
        </a>
      </Dropdown>
    </Space>
  );
};

export default ButtonGroup;
