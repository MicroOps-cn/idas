import { Button, Card, message, Popconfirm, Space, Tabs, Tag } from 'antd';
import { isArray } from 'lodash';
import moment from 'moment';
import type { Tab } from 'rc-tabs/lib/interface';
import type { ReactNode } from 'react';
import { useEffect, useRef, useState } from 'react';

import { deleteAppKeys, getAppInfo, getAppKeys, updateApp } from '@/services/idas/apps';
import { GrantMode, GrantType } from '@/services/idas/enums';
import { enumToMap } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import { CloseOutlined } from '@ant-design/icons';
import type { ProDescriptionsItemProps } from '@ant-design/pro-descriptions';
import ProDescriptions from '@ant-design/pro-descriptions';
import type { ProColumns } from '@ant-design/pro-table';
import { ProTable } from '@ant-design/pro-table';
import type { ProCoreActionType } from '@ant-design/pro-utils';
import { useIntl, useParams, history } from '@umijs/max';

import AddKeyForm from './components/AddKeyForm';
import { GrantView } from './components/GrantView';

/**
 *  Delete node
 * @zh-CN 删除应用密钥
 *
 * @param id
 */
const handleRemoveKeyPair = async (id: string, appId: string) => {
  if (!id || !appId) {
    return true;
  }
  const hide = message.loading('Deleting ...');
  try {
    await deleteAppKeys({ appId: appId }, { appId, id });
    hide();
    message.success('Deleted successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again');
    return false;
  }
};

const AppDetail: React.FC = ({}) => {
  const { aid } = useParams();
  const [currentTabKey, setCurrentTabKey] = useState<'baisc' | 'users'>('baisc');
  const [currentExtra, setCurrentExtra] = useState<React.ReactNode>();
  const [appInfo, setAppInfo] = useState<API.AppInfo>();
  const intl = new IntlContext('pages.apps', useIntl());
  const descripterRef = useRef<ProCoreActionType<any>>();
  const keyPairRef = useRef<ProCoreActionType<any>>();
  const [granting, setGranting] = useState<boolean>(false);
  const [appRoles, setAppRoles] = useState<API.AppRoleInfo[]>([]);
  const [grantedUserList, setGrantedUserList] = useState<API.UserInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [updating, setUpdating] = useState<boolean>(false);
  const [addKeyFormVisible, setAddKeyFormVisible] = useState<boolean>(false);
  const getAppRoles = () => {
    return appRoles.map(({ id, name, isDefault, urls }) => {
      return { id, name, isDefault, urls };
    });
  };
  const detailColumns: ProDescriptionsItemProps<API.AppInfo>[] = [
    {
      title: intl.t('name.label', 'Name'),
      dataIndex: 'name',
    },
    {
      title: intl.t('description.label', 'Description'),
      dataIndex: 'description',
      span: 2,
      ellipsis: true,
      contentStyle: { display: 'grid' },
    },
    {
      title: intl.t('displayName.label', 'Display Name'),
      dataIndex: 'displayName',
    },
    {
      title: intl.t('url.label', 'URL'),
      dataIndex: 'url',
      span: 2,
      copyable: true,
      render: (dom, entry) => {
        return (
          <a target="_blank" href={entry.url} rel="noreferrer">
            {dom}
          </a>
        );
      },
      ellipsis: true,
      contentStyle: { display: 'contents' },
    },
    {
      title: intl.t('grantType.label', 'Grant Type'),
      dataIndex: 'grantType',
      valueEnum: enumToMap(GrantType, intl, 'grantType.value'),
      render(_, entity) {
        if (isArray(entity.grantType)) {
          if (entity.grantType.length > 0) {
            return entity.grantType.map((t) => (
              <Tag key={t}>
                {intl.formatMessage({
                  id: `grantType.value.${GrantType[t]}`,
                  defaultMessage: GrantType[t] as string,
                })}
              </Tag>
            ));
          } else {
            return (
              <Tag>
                {intl.formatMessage({
                  id: `grantType.value.${GrantType[0]}`,
                  defaultMessage: GrantType[0] as string,
                })}
              </Tag>
            );
          }
        }
        return <Tag>{entity.grantType}</Tag>;
      },
    },
    {
      title: intl.t('title.grantMode', 'Grant Mode'),
      dataIndex: 'grantMode',
      valueEnum: enumToMap(GrantMode, intl, 'grantMode.value'),
    },
  ];

  const fetchAppInfo = async () => {
    try {
      if (!aid) {
        return { success: false };
      }
      setLoading(true);
      const resp = await getAppInfo({ id: aid });
      if (resp.data) {
        const { data: app } = resp;
        setAppInfo(app);
        setAppRoles(app.roles ? app.roles : []);
        // @ts-ignore
        setGrantedUserList(app.users ? app.users : []);
      }
      keyPairRef.current.reload();
      return resp;
    } catch (error) {
      if (!(error as RequestError).handled) {
        console.error(`failed to get app info: ${error}`);
      }
      return { success: false };
    } finally {
      setLoading(false);
    }
  };
  const appProxyColumns: ProColumns<API.AppProxyInfo>[] = [
    {
      title: intl.t('proxy.title.domain', 'Domain'),
      dataIndex: 'domain',
    },
    {
      title: intl.t('proxy.title.insecureSkipVerify', 'Skip TLS Verify'),
      dataIndex: 'insecureSkipVerify',
      valueType: 'switch',
    },
    {
      title: intl.t('proxy.title.transparentServerName', 'Transparent Server Name'),
      dataIndex: 'transparentServerName',
      valueType: 'switch',
    },
    {
      title: intl.t('proxy.title.hstsOffload', 'HSTS Offload'),
      dataIndex: 'hstsOffload',
      valueType: 'switch',
    },
    {
      title: intl.t('proxy.title.jwtProvider', 'JWT provider'),
      dataIndex: 'jwtProvider',
      render: (_, entity) => {
        if (entity.jwtProvider && entity.jwtCookieName) {
          return `[Cookie]${entity.jwtCookieName}`;
        }
        return '';
      },
    },
  ];
  const proxyURLColumns: ProColumns<API.AppProxyUrl>[] = [
    {
      title: intl.t('proxy.urls.title.name', 'Name'),
      dataIndex: 'name',
      ellipsis: true,
    },
    {
      title: intl.t('proxy.urls.title.method', 'Method'),
      dataIndex: 'method',
      ellipsis: true,
    },
    {
      title: intl.t('proxy.urls.title.url', 'URL'),
      dataIndex: 'url',
      ellipsis: true,
    },
    {
      title: intl.t('proxy.urls.title.upstream', 'Upstream'),
      dataIndex: 'upstream',
      ellipsis: true,
      render: (_, entry) => {
        return entry.upstream ?? appInfo?.proxy?.upstream;
      },
    },
    {
      title: intl.t('proxy.urls.title.role', 'Role'),
      ellipsis: true,
      render: (_, entry) => {
        const roles = appInfo?.roles ?? [];
        const doms: ReactNode[] = [];
        for (const role of roles) {
          for (const url of role.urls ?? []) {
            if (url == entry.id) {
              doms.push(<Tag key={role.id}>{role.name}</Tag>);
            }
          }
        }
        return doms;
      },
    },
  ];
  const appKeyPairColumns: ProColumns<API.SimpleAppKeyInfo>[] = [
    {
      title: intl.t('keypair.title.name', 'Name'),
      dataIndex: 'name',
      ellipsis: true,
    },
    {
      title: intl.t('keypair.title.key', 'Key'),
      dataIndex: 'key',
      ellipsis: true,
    },
    {
      title: intl.t('keypair.title.createTime', 'Create Time'),
      dataIndex: 'createTime',
      width: 200,
      render: (_, item) => moment(item.createTime).locale(intl.locale).format('LLL'),
    },
    {
      width: 40,
      render: (_, record) => [
        <Popconfirm
          key="delete"
          title={intl.t(
            'keypair.delete.popconfirm',
            `Are you sure you want to delete the key named {name}? After deletion, the service cannot be accessed using this key pair.`,
            undefined,
            { name: record.name },
          )}
          onConfirm={async () => {
            if (await handleRemoveKeyPair(record.id, record.appId)) {
              keyPairRef.current?.reload();
            }
          }}
        >
          <a>
            <CloseOutlined />
          </a>
        </Popconfirm>,
      ],
    },
  ];
  const tabsItems: Tab[] = [
    {
      label: intl.t('title.basic', 'Basic Info'),
      key: 'basic',
      children: (
        <>
          <ProDescriptions<API.AppInfo>
            column={{ xs: 1, sm: 2, md: 3 }}
            actionRef={descripterRef}
            request={fetchAppInfo}
            columns={detailColumns}
          />
          <Card
            bodyStyle={{ padding: 0 }}
            headStyle={{ padding: 0 }}
            title={<b>{intl.t('title.keys', 'Key-pairs')}</b>}
            bordered={false}
            extra={
              <Button
                type="primary"
                onClick={() => {
                  setAddKeyFormVisible(true);
                }}
              >
                {intl.t('button.create-key', 'Create Key-pair')}
              </Button>
            }
          >
            <ProTable
              actionRef={keyPairRef}
              toolBarRender={false}
              request={async (params) => {
                if (!aid) {
                  return { success: false };
                }
                return getAppKeys({ ...params, appId: aid });
              }}
              columns={appKeyPairColumns}
              rowKey="id"
              search={false}
            />
          </Card>
        </>
      ),
    },
    {
      label: intl.t('users', 'User'),
      key: 'users',
      children: (
        <div
          style={{
            height: 'calc( 100vh - 300px )',
          }}
        >
          <GrantView
            users={grantedUserList}
            roles={appRoles}
            onChange={setGrantedUserList}
            loading={loading}
            granting={granting}
            parentIntl={intl}
            type="horizontal"
          />
        </div>
      ),
    },
    {
      label: intl.t('proxy', 'Proxy'),
      key: 'proxy',
      disabled: !((appInfo?.grantType ?? []) as any).includes(GrantType.proxy),
      children: (
        <>
          <ProDescriptions<API.AppProxyInfo>
            column={{ xs: 1, sm: 2, md: 3 }}
            columns={appProxyColumns}
            dataSource={appInfo?.proxy}
          />
          <Card
            bodyStyle={{ padding: 0 }}
            headStyle={{ padding: 0 }}
            title={<b>{intl.t('title.urls', 'URLs')}</b>}
            bordered={false}
          >
            <ProTable
              actionRef={keyPairRef}
              toolBarRender={false}
              dataSource={appInfo?.proxy?.urls ?? []}
              columns={proxyURLColumns}
              rowKey="id"
              search={false}
            />
          </Card>
        </>
      ),
    },
  ];

  const handleUpdateApp = async (newInfo: Partial<API.UpdateAppRequest>) => {
    if (appInfo) {
      try {
        setUpdating(true);
        const { isDelete, createTime, updateTime, ...info } = appInfo;
        const resp = await updateApp(
          { id: appInfo.id },
          {
            ...info,
            roles: getAppRoles(),
            ...newInfo,
          },
        );
        if (resp.errorMessage) {
          message.error(`failed to update: ${resp.errorMessage}`);
        }
        message.info('Successfully updated app information.');
        setUpdating(false);
        return resp;
      } catch (error) {
        if (!(error as RequestError).handled) {
          console.error(`failed to get user list: ${error}`);
        }
      } finally {
        setUpdating(false);
      }
    }
    return { success: false };
  };

  useEffect(() => {
    let buttons = [
      <Button
        hidden={!appInfo?.id}
        key={'edit'}
        type="primary"
        onClick={() => history.push(`/apps/${appInfo?.id}/edit`)}
      >
        {intl.t('button.edit', 'Edit')}
      </Button>,
      <Button key={'refresh'} type="primary" onClick={() => descripterRef.current.reload()}>
        {intl.t('button.refresh', 'Refresh')}
      </Button>,
    ];
    if (currentTabKey === 'users') {
      if (granting) {
        buttons = [
          <Button
            type="primary"
            onClick={async () => {
              handleUpdateApp({
                users: grantedUserList.map((user) => ({
                  id: user.id,
                  roleId: user.roleId,
                })),
              }).then((val) => {
                if (val.success) {
                  setGranting(false);
                  descripterRef.current.reload();
                }
              });
            }}
            loading={updating}
            key="save"
          >
            {intl.t('button.save', 'Save')}
          </Button>,
          <Button type="link" onClick={() => setGranting(false)} key="cancel">
            {intl.t('button.cancel', 'Cancel')}
          </Button>,
        ];
      } else {
        buttons.push(
          <Button type="primary" onClick={() => setGranting(true)}>
            {intl.t('button.grant', 'Grant')}
          </Button>,
        );
      }
    }
    setCurrentExtra(buttons);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentTabKey, granting, grantedUserList, updating]);

  return (
    <>
      <Card
        style={{
          height: '100%',
        }}
        bodyStyle={{ height: '100%' }}
      >
        <Tabs
          items={tabsItems}
          onChange={(key) => {
            setCurrentTabKey(key as 'baisc' | 'users');
          }}
          tabBarExtraContent={<Space>{currentExtra}</Space>}
        />
      </Card>
      <AddKeyForm
        visible={addKeyFormVisible}
        app={appInfo}
        onClose={() => {
          setAddKeyFormVisible(false);
          keyPairRef.current.reload();
        }}
        parentIntl={intl}
      />
    </>
  );
};
export default AppDetail;
