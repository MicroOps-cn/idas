import { Button, message, Drawer, Modal, List, Divider, Popconfirm } from 'antd';
import moment from 'moment';
import React, { useState, useRef } from 'react';
import { useIntl } from 'umi';

import { getRoles, createRole, updateRole, deleteRoles } from '@/services/idas/roles';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { FooterToolbar, PageContainer } from '@ant-design/pro-components';
import type { ProDescriptionsItemProps } from '@ant-design/pro-descriptions';
import ProDescriptions from '@ant-design/pro-descriptions';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';

import type { FormValueType } from './components/CreateOrUpdateForm';
import CreateOrUpdateForm from './components/CreateOrUpdateForm';

/**
 * @en-US Add node
 * @zh-CN 添加角色
 * @param fields
 */
const handleAdd = async (fields: FormValueType) => {
  const hide = message.loading('Adding ...');
  try {
    delete fields.id;
    await createRole(fields);
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
 * @zh-CN 更新角色
 *
 * @param fields
 */
const handleUpdate = async (fields: FormValueType) => {
  const hide = message.loading('Configuring');
  try {
    if (fields.id) {
      //@ts-ignore
      await updateRole({ id: fields.id }, fields);
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
 *  Delete node
 * @zh-CN 删除角色
 *
 * @param selectedRows
 */
const handleRemove = async (selectedRows: API.RoleInfo[]) => {
  const hide = message.loading('Deleting ...');
  if (!selectedRows) {
    return true;
  }
  try {
    await deleteRoles(
      selectedRows.map(
        (row): API.DeleteRoleRequest => ({
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

const RoleList: React.FC = () => {
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建/修改窗口的弹窗
   *  */
  const [modalVisible, handleModalVisible] = useState<boolean>(false);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.RoleInfo>();
  const [selectedRowsState, setSelectedRows] = useState<API.RoleInfo[]>([]);
  const [keywords, setKeywords] = useState<string>();

  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */
  const intl = new IntlContext('pages.roles', useIntl());

  const fetchRoles = async (
    params: Omit<API.getRolesParams, 'storage'>,
  ): Promise<API.GetRolesResponse> => {
    let queryParams: API.getRolesParams = { ...params };
    if (keywords) {
      queryParams = { keywords, ...queryParams };
    }
    return getRoles(queryParams);
  };

  const columns: ProColumns<API.RoleInfo>[] = [
    {
      title: intl.t('title.name', 'Role Name'),
      dataIndex: 'name',
    },
    {
      title: intl.t('title.describe', 'Describe'),
      dataIndex: 'describe',
    },
    {
      title: intl.t('title.createTime', 'Create Time'),
      dataIndex: 'create_time',
      render: (_, item) => {
        return moment(item.createTime).locale(intl.locale).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      title: intl.t('title.operate', 'Operate'),
      dataIndex: 'option',
      valueType: 'option',
      render: (_: any, record: API.RoleInfo) => (
        <>
          <a
            key="change"
            onClick={() => {
              handleModalVisible(true);
              setCurrentRow(record);
            }}
          >
            {intl.t('button.change', 'Change')}
          </a>
          <Divider type="vertical" />
          <Popconfirm
            key="delete"
            title={intl.t(
              'delete.popconfirm',
              `Are you sure you want to delete the role named: {name}?`,
              undefined,
              { name: record.name },
            )}
            onConfirm={async () => {
              if (await handleRemove([record])) {
                actionRef.current?.reload();
              }
            }}
          >
            <a>{intl.t('button.delete', 'Delete')}</a>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.RoleInfo, API.getRolesParams>
        actionRef={actionRef}
        rowKey="id"
        search={false}
        toolbar={{
          search: true,
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
        request={fetchRoles}
        columns={columns}
        tableAlertRender={false}
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
                  'delete.confirm',
                  'Are you sure you want to delete the following roles?            ',
                ),
                icon: <ExclamationCircleOutlined />,
                async onOk() {
                  await handleRemove(selectedRowsState);
                  setSelectedRows([]);
                  actionRef.current?.reloadAndRest?.();
                },
                content: (
                  <List<API.RoleInfo>
                    dataSource={selectedRowsState}
                    rowKey={'id'}
                    renderItem={(item) => <List.Item>{item.name}</List.Item>}
                  />
                ),
              });
            }}
          >
            {intl.t('button.batchDeletion', 'Batch deletion')}
          </Button>
        </FooterToolbar>
      )}

      <CreateOrUpdateForm
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
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {showDetail && currentRow?.id && (
          <>
            <ProDescriptions<API.RoleInfo>
              column={2}
              title={intl.t('detail.title', 'Role Details')}
              request={async () => ({
                data: currentRow || {},
              })}
              params={{
                id: currentRow?.id,
              }}
              columns={columns as ProDescriptionsItemProps<API.RoleInfo>[]}
            />
          </>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default RoleList;
