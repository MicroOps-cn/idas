import type { InputRef } from 'antd';
import { Checkbox, Row, Space, Switch, Input, Tabs, message } from 'antd';
import _ from 'lodash';
import { useEffect, useRef, useState } from 'react';

import { newId } from '@/utils/uuid';
import { StarFilled } from '@ant-design/icons';

const RoleViewTabTitle: React.FC<{
  role: API.AppRoleInfo;
  activeKey?: string;
  setActiveKey: (value: React.SetStateAction<string | undefined>) => void;
  setRole: (role: Partial<API.AppRole>) => Promise<void> | void;
}> = (props) => {
  const { role, setRole, activeKey, setActiveKey } = props;
  const activeTabNameEditRef = useRef<InputRef>(null);
  const [tabEditing, setTabEditing] = useState<boolean>(false);
  useEffect(() => {
    if (activeKey === role.id && !role.name) {
      setTabEditing(true);
    }
    if (activeKey === role.id && tabEditing) {
      activeTabNameEditRef.current?.focus({ cursor: 'all' });
    }
  }, [tabEditing, activeKey, role]);

  return (
    <>
      <Input
        hidden={!(tabEditing && role.id === activeKey)}
        style={{
          width: 80,
          padding: 0,
        }}
        ref={activeTabNameEditRef}
        onKeyDown={(e) => {
          if (e.key == 'Enter') {
            setRole({ id: role.id, name: e.currentTarget.value });
            setTabEditing(false);
          }
        }}
        onBlur={(e) => {
          setRole({ id: role.id, name: e.target.value });
          setTabEditing(false);
        }}
        autoFocus
        defaultValue={role.name ? role.name : 'New Role'}
      />

      <div
        hidden={tabEditing && role.id === activeKey}
        style={{
          width: 80,
          overflow: 'hidden',
          textShadow: '0 0 0 aliceblue',
        }}
        title={role.name}
        onDoubleClick={() => {
          setActiveKey(role.id);
          setTabEditing(true);
        }}
      >
        <StarFilled style={{ color: 'orange' }} hidden={!role.isDefault} />
        {role.name}
      </div>
    </>
  );
};

interface RoleViewProps {
  value?: API.AppRoleInfo[];
  onChange?: (roles: API.AppRoleInfo[]) => Promise<void> | void;
  urls: API.AppProxyUrl[];
}

export const RoleView: React.FC<RoleViewProps> = ({
  urls = [],
  value: roles = [],
  onChange: setRoles,
}) => {
  const [activeKey, setActiveKey] = useState<string>();
  const editor = {
    add: () => {
      const tabKey = `new-role-${newId()}`;
      setActiveKey(tabKey);
      setRoles?.([...roles, { name: '', id: tabKey }]);
    },
    remove: (targetKey: string | React.MouseEvent | React.KeyboardEvent) => {
      if (_.isString(targetKey)) {
        setRoles?.(roles.filter((role) => role.id != targetKey));
      } else {
        message.warning(`system error: ${targetKey} is not string`);
      }
    },
  };
  const setRole = async (role: Partial<API.AppRole>) => {
    await setRoles?.(
      roles.map((r) => {
        if (role.id == r.id) {
          return { ...r, ...role };
        }
        return r;
      }),
    );
  };
  return (
    <Tabs
      type="editable-card"
      tabPosition={'top'}
      style={{ height: 350 }}
      activeKey={activeKey}
      onChange={setActiveKey}
      onEdit={(targetKey, action) => {
        editor[action](targetKey);
      }}
    >
      {roles.map((role) => (
        <Tabs.TabPane
          tab={
            <RoleViewTabTitle
              role={role}
              setRole={setRole}
              activeKey={activeKey}
              setActiveKey={setActiveKey}
            />
          }
          key={role.id}
          tabKey={role.id}
        >
          <Space direction="vertical">
            <div style={{ display: urls.length === 0 ? 'none' : 'unset' }}>
              <div style={{ paddingBottom: 7 }}>Permission</div>
              <Checkbox.Group
                onChange={(vals) => {
                  setRole({ ...role, urls: vals as string[] });
                }}
                defaultValue={role.urls}
              >
                {urls.map((url) => {
                  return (
                    <Row key={url.id}>
                      <Checkbox value={url.id}>{url.name}</Checkbox>
                    </Row>
                  );
                })}
              </Checkbox.Group>
            </div>
            <div style={{ padding: 10 }}>
              默认:
              <Switch
                checked={role.isDefault}
                onChange={(checked) => {
                  setRoles?.(
                    roles.map((r) => {
                      if (role.id == r.id) {
                        return { ...r, isDefault: checked };
                      } else if (checked) {
                        return { ...r, isDefault: !checked };
                      }
                      return r;
                    }),
                  );
                }}
              />
            </div>
          </Space>
        </Tabs.TabPane>
      ))}
    </Tabs>
  );
};
