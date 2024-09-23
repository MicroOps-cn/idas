import { Button, message, Drawer, Modal, List, Tabs } from 'antd';
import moment from 'moment';
import React, { useState, useRef, useEffect } from 'react';
import { useIntl } from 'umi';

import ButtonGroup from '@/components/ButtonGroup';
import { UserStatus } from '@/services/idas/enums';
import { getSessions as getUserSessions, deleteSession } from '@/services/idas/sessions';
import {
  getUsers as getUsersInfo,
  createUser as addUser,
  updateUser,
  deleteUsers,
  patchUsers,
  sendActivateMail,
  getUserInfo,
} from '@/services/idas/users';
import { enumToStatusEnum } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import {
  CheckOutlined,
  ExclamationCircleOutlined,
  LogoutOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { LightFilter, ProFormSelect } from '@ant-design/pro-components';
import { FooterToolbar, PageContainer } from '@ant-design/pro-components';
import type { ProDescriptionsItemProps } from '@ant-design/pro-descriptions';
import ProDescriptions from '@ant-design/pro-descriptions';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';

import type { FormValueType } from './components/CreateOrUpdateForm';
import CreateOrUpdateForm from './components/CreateOrUpdateForm';
import GrantView from './components/GrantView';
import styles from './index.less';

/**
 * @en-US Add node
 * @zh-CN 添加用户
 * @param fields
 */
const handleAdd = async (fields: FormValueType) => {
  const hide = message.loading('Adding ...');
  try {
    delete fields.id;
    await addUser(fields);
    hide();
    message.success('Added successfully');
    return true;
  } catch (error) {
    hide();
    if (!(error as RequestError).handled) {
      message.error('Adding failed, please try again!');
    }
    return false;
  }
};

/**
 * @en-US Update node
 * @zh-CN 更新用户
 *
 * @param fields
 */
const handleUpdate = async (fields: FormValueType) => {
  const hide = message.loading('Configuring');
  try {
    if (fields.id) {
      //@ts-ignore
      await updateUser({ id: fields.id }, fields);
      hide();
      message.success('update is successful');
      return true;
    } else {
      message.success('update failed, system error');
      return false;
    }
  } catch (error) {
    hide();
    if (!(error as RequestError).handled) {
      message.error('Update failed, please try again!');
    }
    return false;
  }
};
/**
 * @en-US Enable user
 * @zh-CN 启用用户
 *
 * @param fields
 */
const handleEnable = async (selectedRows: API.UserInfo[]) => {
  const hide = message.loading('enabling ...');
  if (!selectedRows) {
    return true;
  }
  try {
    await patchUsers(
      selectedRows.map((row): API.PatchUserRequest => {
        return { id: row.id, status: UserStatus.normal };
      }),
    );
    hide();
    message.success('Enabled successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error('Enabled failed, please try again');
    return false;
  }
};

/**
 * @en-US Disable user
 * @zh-CN 禁用用户
 *
 * @param fields
 */
const handleDisable = async (selectedRows: API.UserInfo[]) => {
  const hide = message.loading('Disabling ...');
  if (!selectedRows) {
    return true;
  }
  try {
    await patchUsers(
      selectedRows.map((row): API.PatchUserRequest => {
        return { id: row.id, status: UserStatus.disabled };
      }),
    );
    hide();
    message.success('Disabled successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error('Disable failed, please try again');
    return false;
  }
};

/**
 *  Delete node
 * @zh-CN 删除用户
 *
 * @param selectedRows
 */
const handleRemove = async (selectedRows: API.UserInfo[]) => {
  const hide = message.loading('Deleting ...');
  if (!selectedRows) {
    return true;
  }
  try {
    await deleteUsers(
      selectedRows.map(
        (row): API.DeleteUserRequest => ({
          id: row.id,
        }),
      ),
    );
    hide();
    message.success('Deleted successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error(`Delete failed, please try again: ${error}`);
    return false;
  }
};

/**
 *  Delete node
 * @zh-CN 删除用户会话
 *
 * @param id
 */
const handleRemoveSession = async (id: string) => {
  if (!id) {
    return true;
  }
  const hide = message.loading('Deleting ...');
  try {
    await deleteSession({
      id: id,
    });
    hide();
    message.success('Deleted successfully and will refresh soon');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again');
    return false;
  }
};

const UserList: React.FC = () => {
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建/修改窗口的弹窗
   *  */

  const [modalVisible, handleModalVisible] = useState<boolean>(false);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [granting, setGranting] = useState<boolean>(false);
  const [detailActiveTabKey, setDetailActiveTabKey] = useState<string>('sessions');
  const actionRef = useRef<ActionType>();
  const sessionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.UserInfo>();
  const [selectedRowsState, setSelectedRows] = useState<API.UserInfo[]>([]);
  const [currentStatus, setCurrentStatus] = useState<UserStatus | 'all'>('all');
  const [keywords, setKeywords] = useState<string>();

  const [currentUserApps, setCurrentUserApps] = useState<API.UserApp[]>([]);

  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */
  const intl = new IntlContext('pages.users', useIntl());
  const getUsers = async (
    params: Omit<API.getUsersParams, 'storage' | 'status'>,
  ): Promise<API.GetUsersResponse> => {
    let queryParams: API.getUsersParams = {
      ...params,
      status: currentStatus !== 'all' ? currentStatus : undefined,
    };
    if (keywords) {
      queryParams = { keywords, ...queryParams };
    }
    return getUsersInfo(queryParams);
  };

  useEffect(() => {
    if (showDetail) {
      setDetailActiveTabKey('sessions');
      setGranting(false);
      if (currentRow?.id) {
        getUserInfo({ id: currentRow.id }).then((ret) => {
          if (ret.success && ret.data?.apps) {
            setCurrentUserApps(ret.data.apps);
          }
        });
      }
    }
  }, [currentRow, showDetail]);
  const sesionColumns: ProColumns<API.SessionInfo>[] = [
    {
      title: intl.t('session.title.lastSeen', 'Last Seen'),
      dataIndex: 'lastSeen',
      render: (_, item) => {
        return moment(item.lastSeen).locale(intl.locale).fromNow();
      },
    },
    {
      title: intl.t('session.title.expiry', 'Expiry'),
      dataIndex: 'expiry',
      render: (_, item) => {
        return moment(item.expiry).locale(intl.locale).fromNow();
      },
    },
    {
      title: intl.t('session.title.loggedOn', 'Logged on'),
      dataIndex: 'createTime',
      render: (_, item) => {
        return moment(item.createTime).locale(intl.locale).format('LLL');
      },
    },
    {
      render: (_, record) => [
        <a
          key="delete"
          onClick={async () => {
            if (await handleRemoveSession(record.id)) {
              sessionRef.current?.reload();
            }
          }}
        >
          <LogoutOutlined />
        </a>,
      ],
    },
  ];

  const userStatusEnum = enumToStatusEnum(UserStatus, intl, 'status.value', {
    normal: 'Success',
    disable: 'Error',
    inactive: 'Warning',
  });

  const columns: ProColumns<API.UserInfo>[] = [
    {
      title: intl.t('updateForm.userName.nameLabel', 'User name'),
      hideInSearch: true,
      dataIndex: 'username',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDetail(true);
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: intl.t('title.fullName', 'FullName'),
      dataIndex: 'fullName',
      hideInSearch: true,
    },
    {
      title: intl.t('title.phoneNumber', 'Telephone number'),
      dataIndex: 'phoneNumber',
      hideInSearch: true,
    },
    {
      title: intl.t('title.email', 'Email'),
      dataIndex: 'email',
      hideInSearch: true,
    },
    {
      title: intl.t('title.status', 'Status'),
      dataIndex: 'status',
      hideInForm: true,
      valueEnum: userStatusEnum,
    },
    {
      title: intl.t('title.updatedTime', 'Last update time'),
      dataIndex: 'updateTime',
      valueType: 'dateTime',
      hideInTable: true,
      hideInSearch: true,
      hideInForm: true,
    },
    {
      title: intl.t('title.loginTime', 'Last login time'),
      dataIndex: 'loginTime',
      valueType: 'dateTime',
      hideInSearch: true,
      hideInForm: true,
    },
    {
      title: intl.t('title.createTime', 'Create time'),
      dataIndex: 'createTime',
      valueType: 'dateTime',
      hideInSearch: true,
      hideInTable: true,
      hideInForm: true,
    },
    {
      title: intl.t('title.option', 'Operate'),
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => {
        return [
          <ButtonGroup
            maxItems={2}
            moreLabel={intl.t('button.more', 'More')}
            items={[
              {
                key: 'edit',
                label: intl.t('button.edit', 'Edit'),
                style: { flex: 'unset' },
                onClick: async () => {
                  handleModalVisible(true);
                  setCurrentRow(record);
                },
              },
              {
                key: 'activate',
                label: intl.t('button.activate', 'Activate'),
                hidden: record.status !== UserStatus.user_inactive,
                success: <CheckOutlined color="green" />,
                style: { flex: 'unset' },
                onClick: async () => {
                  if (!record.email) {
                    throw new Error(intl.t('activate.no-email', ' The user has no email.'));
                  }
                  const resp = await sendActivateMail({
                    userId: record.id,
                  });
                  if (resp.success) {
                    message.success(intl.t('activate.succcess', 'Email sent successfully.'));
                  } else {
                    throw new Error(intl.t('activate.failed', 'Email sent failed.'));
                  }
                },
              },
              {
                key: 'delete',
                label: intl.t('button.delete', 'Delete'),
                style: { flex: 'unset' },
                onClick: () => {
                  Modal.confirm({
                    title: intl.t(
                      'deleteConfirm',
                      'Are you sure you want to delete the following users?            ',
                    ),
                    icon: <ExclamationCircleOutlined />,
                    async onOk() {
                      await handleRemove([record]);
                    },
                    content: (
                      <List<API.UserInfo>
                        dataSource={[record]}
                        rowKey={'id'}
                        renderItem={(item) => (
                          <List.Item>
                            {item.username}
                            {item.fullName
                              ? `(${item.fullName})`
                              : item.email
                              ? `(${item.email})`
                              : ''}
                          </List.Item>
                        )}
                      />
                    ),
                  });
                },
              },
              {
                key: 'disable',
                label: intl.t('button.disable', 'Disable'),
                hidden: record.status === UserStatus.disabled,
                style: { flex: 'unset' },
                onClick: () => {
                  Modal.confirm({
                    title: intl.t(
                      'disableConfirm',
                      'Are you sure you want to disable the following users?',
                    ),
                    icon: <ExclamationCircleOutlined />,
                    async onOk() {
                      await handleDisable([record]);
                      actionRef.current?.reloadAndRest?.();
                    },
                    content: (
                      <List<API.UserInfo>
                        dataSource={[record].filter((user) => user.status !== UserStatus.disabled)}
                        rowKey={'id'}
                        renderItem={(item) => (
                          <List.Item>
                            {item.username}
                            {item.fullName
                              ? `(${item.fullName})`
                              : item.email
                              ? `(${item.email})`
                              : ''}
                          </List.Item>
                        )}
                      />
                    ),
                  });
                },
              },
            ]}
          />,
        ];
      },
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.UserInfo, API.getUsersParams>
        actionRef={actionRef}
        rowKey="id"
        search={false}
        tableAlertRender={false}
        toolbar={{
          search: {
            onSearch: (kws) => {
              setKeywords(kws);
              if (actionRef.current) {
                actionRef.current.setPageInfo?.({
                  ...actionRef.current.pageInfo,
                  current: 1,
                });
                actionRef.current.reload();
              }
            },
          },

          filter: (
            <LightFilter<{ status: { label: string; value: UserStatus; key: string } }>
              onFinish={async ({ status }) => {
                setCurrentStatus(status.value ?? 'all');
                if (actionRef.current) {
                  actionRef.current.setPageInfo?.({
                    ...actionRef.current.pageInfo,
                    current: 1,
                  });
                  actionRef.current?.reload();
                }
                return true;
              }}
            >
              <ProFormSelect<{ value: UserStatus | 'all' }>
                name="status"
                label={intl.t('title.status', 'Status')}
                initialValue={{ value: 'all', label: intl.t('status.all', 'All') }}
                fieldProps={{
                  labelInValue: true,
                }}
                valueEnum={{
                  all: intl.t('status.all', 'All'),
                  ...userStatusEnum,
                }}
              />
            </LightFilter>
          ),
          actions: [
            <Button
              key="create"
              type="primary"
              onClick={() => {
                handleModalVisible(true);
              }}
            >
              <PlusOutlined />
              {intl.t('button.create', 'Create')}
            </Button>,
          ],
        }}
        request={getUsers}
        columns={columns}
        // tableAlertRender={false}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              {intl.t('chosen', 'Chosen')}{' '}
              <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a>{' '}
              {intl.t('item', 'Item(s)')}
            </div>
          }
        >
          <Button
            danger
            onClick={() => {
              Modal.confirm({
                title: intl.t(
                  'deleteConfirm',
                  'Are you sure you want to delete the following users?            ',
                ),
                icon: <ExclamationCircleOutlined />,
                async onOk() {
                  await handleRemove(selectedRowsState);
                  setSelectedRows([]);
                  actionRef.current?.reloadAndRest?.();
                },
                content: (
                  <List<API.UserInfo>
                    dataSource={selectedRowsState}
                    rowKey={'id'}
                    renderItem={(item) => (
                      <List.Item>
                        {item.username}
                        {item.fullName ? `(${item.fullName})` : item.email ? `(${item.email})` : ''}
                      </List.Item>
                    )}
                  />
                ),
              });
            }}
          >
            {intl.t('batchDeletion', 'Batch deletion')}
          </Button>
          {selectedRowsState.filter((user) => user.status !== UserStatus.disabled).length > 0 && (
            <Button
              onClick={() => {
                Modal.confirm({
                  title: intl.t(
                    'disableConfirm',
                    'Are you sure you want to disable the following users?',
                  ),
                  icon: <ExclamationCircleOutlined />,
                  async onOk() {
                    await handleDisable(
                      selectedRowsState.filter((user) => user.status !== UserStatus.disabled),
                    );
                    setSelectedRows([]);
                    actionRef.current?.reloadAndRest?.();
                  },
                  content: (
                    <List<API.UserInfo>
                      dataSource={selectedRowsState.filter(
                        (user) => user.status !== UserStatus.disabled,
                      )}
                      rowKey={'id'}
                      renderItem={(item) => (
                        <List.Item>
                          {item.username}
                          {item.fullName
                            ? `(${item.fullName})`
                            : item.email
                            ? `(${item.email})`
                            : ''}
                        </List.Item>
                      )}
                    />
                  ),
                });
              }}
            >
              {intl.t('batchDisable', 'Batch disable')}
            </Button>
          )}
          {selectedRowsState.filter((user) => user.status !== UserStatus.normal).length > 0 && (
            <Button
              onClick={() => {
                Modal.confirm({
                  title: intl.t(
                    'enableConfirm',
                    'Are you sure you want to enable the following users?',
                  ),
                  icon: <ExclamationCircleOutlined />,
                  async onOk() {
                    await handleEnable(
                      selectedRowsState.filter((user) => user.status !== UserStatus.normal),
                    );
                    setSelectedRows([]);
                    actionRef.current?.reloadAndRest?.();
                  },
                  content: (
                    <List<API.UserInfo>
                      dataSource={selectedRowsState.filter(
                        (user) => user.status !== UserStatus.normal,
                      )}
                      rowKey={'id'}
                      renderItem={(item) => (
                        <List.Item>
                          {item.username}
                          {item.fullName
                            ? `(${item.fullName})`
                            : item.email
                            ? `(${item.email})`
                            : ''}
                        </List.Item>
                      )}
                    />
                  ),
                });
              }}
            >
              {intl.t('batchEnable', 'Batch enable')}
            </Button>
          )}
        </FooterToolbar>
      )}
      <CreateOrUpdateForm
        title={intl.t(
          currentRow ? 'form.title.userUpdate' : 'form.title.userCreate',
          currentRow ? 'Modify user' : 'Add user',
        )}
        onSubmit={async (value) => {
          const success = await (currentRow ? handleUpdate : handleAdd)(value);
          if (success) {
            handleModalVisible(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
          return success;
        }}
        onCancel={() => {
          handleModalVisible(false);
          if (!showDetail) {
            setCurrentRow(undefined);
          }
        }}
        modalVisible={modalVisible}
        values={currentRow}
        parentIntl={intl}
      />

      <Drawer
        width={800}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {showDetail && currentRow?.username && (
          <>
            <ProDescriptions<API.UserInfo>
              column={2}
              title={intl.t('detail.title', 'User Details')}
              request={async () => ({
                data: currentRow || {},
              })}
              params={{
                id: currentRow?.id,
              }}
              className={styles.UserDetails}
              extra={
                <>
                  <a
                    key="grant"
                    onClick={() => {
                      setGranting(true);
                      setDetailActiveTabKey('apps');
                    }}
                    style={{ flex: 'unset' }}
                    hidden={granting}
                  >
                    {intl.t('button.grant', 'Grant')}
                  </a>
                  <a
                    key="save"
                    onClick={async () => {
                      const { createTime, loginTime, extendedData, updateTime, ...vals } =
                        currentRow;
                      const success = await handleUpdate({
                        ...vals,
                        apps: currentUserApps.map((app) => ({ id: app.id, roleId: app.roleId })),
                      });
                      if (success) {
                        setGranting(false);
                      }
                    }}
                    style={{ flex: 'unset' }}
                    hidden={!granting}
                  >
                    {intl.t('button.save', 'Save')}
                  </a>
                  <a
                    key="cancel"
                    onClick={async () => {
                      setGranting(false);
                    }}
                    style={{ flex: 'unset' }}
                    hidden={!granting}
                  >
                    {intl.t('button.cancel', 'Cancel')}
                  </a>
                </>
              }
              columns={columns as ProDescriptionsItemProps<API.UserInfo>[]}
            />
            <Tabs
              activeKey={detailActiveTabKey}
              onChange={(item) => {
                setDetailActiveTabKey(item);
              }}
              items={[
                {
                  label: intl.t('detail.sessions.title', 'Sessions'),
                  key: 'sessions',
                  children: (
                    <ProTable
                      actionRef={sessionRef}
                      toolBarRender={false}
                      request={(params) => {
                        return getUserSessions({ ...params, userId: currentRow.id });
                      }}
                      columns={sesionColumns}
                      rowKey="id"
                      search={false}
                    />
                  ),
                },
                {
                  label: intl.t('detail.apps.title', 'App'),
                  key: 'apps',
                  children: (
                    <GrantView
                      apps={currentUserApps}
                      onChange={setCurrentUserApps}
                      granting={granting}
                      parentIntl={intl}
                    />
                  ),
                },
              ]}
            />
          </>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default UserList;
