import { Button, Drawer, List, message, Modal } from 'antd';
import { isFunction, isString } from 'lodash';
import React, { useState, useRef, useEffect, useCallback } from 'react';
import { history } from 'umi';

import type { IntlContext } from '@/utils/intl';
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { FooterToolbar } from '@ant-design/pro-components';
import ProDescriptions from '@ant-design/pro-descriptions';
import type { ActionType, ProColumns, ProTableProps, RequestData } from '@ant-design/pro-table';
import { ProTable } from '@ant-design/pro-table';
import { useLocation } from '@umijs/max';

export interface TableItem extends Record<string, any> {
  id: string;
  createTime?: string;
  updateTime?: string;
  isDisable?: boolean;
}

type ListRequest = (params: {
  pageSize?: number;
  current?: number;
  keywords?: string;
}) => Promise<Partial<RequestData<TableItem>>>;

type PatchRequst = (
  body: { id: string; isDelete?: boolean; isDisable?: boolean }[],
) => Promise<any>;
type StatusSwitchRequest = (body: { id: string }[]) => Promise<any>;

type TableProps = Omit<ProTableProps<TableItem, any>, 'request' | 'metas' | 'columns'> & {
  intl: IntlContext;
  request:
    | ListRequest
    | {
        list: ListRequest;
        patch?: PatchRequst;
        disable?: StatusSwitchRequest;
        enable?: StatusSwitchRequest;
        delete?: StatusSwitchRequest;
      };
  columns: ProColumns<TableItem>[];
};

interface BatchHandleButton {
  onOk?: null | ((items: string | string[]) => Promise<boolean>);
  items?: TableItem[];
  icon?: React.ReactNode;
  title?: React.ReactNode;
  children: React.ReactNode;
  onFinish?: null | ((items: TableItem[]) => Promise<void>);
}

const BatchPopconfirm: React.FC<BatchHandleButton> = ({
  onOk,
  onFinish,
  items,
  children,
  title,
  icon,
}) => {
  if (!onOk || !items || items.length === 0) {
    return null;
  }
  const renderSelectedItem = (item: TableItem) => {
    return <List.Item>{item.id}</List.Item>;
  };
  return (
    <Button
      onClick={() => {
        Modal.confirm({
          title: title,
          icon: icon,
          onOk: async () => {
            if (await onOk(items.map((item) => item.id))) {
              onFinish?.(items);
            }
          },
          content: (
            <List<TableItem> dataSource={items} rowKey={'id'} renderItem={renderSelectedItem} />
          ),
        });
      }}
    >
      {children}
    </Button>
  );
};

const Table: React.FC<TableProps> = ({ intl, request, columns, ...proTableProps }) => {
  const actionRef = useRef<ActionType>();
  const [keywords, setKeywords] = useState<string>();
  const [selectedRowsState, setSelectedRows] = useState<TableItem[]>([]);
  const [currentRow, setCurrentRow] = useState<TableItem>();
  const location = useLocation();
  const fetchList = useCallback(
    async (params: {
      pageSize?: number;
      current?: number;
      keyword?: string;
    }): Promise<Partial<RequestData<TableItem>>> => {
      return (isFunction(request) ? request : request.list)({
        ...params,
        keywords: keywords ?? undefined,
      });
    },
    [request, keywords],
  );

  interface createHandlerMessage {
    processing?: React.ReactNode;
    success?: React.ReactNode;
    error?: React.ReactNode;
  }
  const createHandler: (
    handler?: null | StatusSwitchRequest,
    patchParameterRender?: ({ id }: { id: string }) => {
      id: string;
      isDelete?: boolean;
      isDisable?: boolean;
    },
    msg?: createHandlerMessage,
  ) => ((id: string | string[]) => Promise<boolean>) | undefined = (
    handler,
    patchParameterRender,
    { processing, success, error } = {},
  ) => {
    const handleFunc =
      handler ??
      (!isFunction(request) && patchParameterRender
        ? request.patch
          ? async (body: { id: string }[]): Promise<any> => {
              return request.patch?.(body.map(patchParameterRender));
            }
          : null
        : null);

    if (!handleFunc) return undefined;
    return async (id: string | string[]) => {
      const hide = message.loading(processing ?? intl.t(`message.processing`, 'Processing ...'));
      try {
        await handleFunc(isString(id) ? [{ id }] : id.map((itemId) => ({ id: itemId })));
        hide();
        message.success(success ?? intl.t(`message.operationSuccessd`, 'Operation succeeded.'));
        return true;
      } catch (err) {
        hide();
        message.error(
          error ?? intl.t(`message.operationFailed`, 'Operation failed, please try again.'),
        );
        return false;
      }
    };
  };

  const handleRemove = createHandler(
    !isFunction(request) ? request.delete : undefined,
    ({ id }) => ({ id, isDelete: true }),
    {
      processing: intl.t(`message.removing`, 'Removing ...'),
      success: intl.t(`message.removeSuccessd`, 'Remove successfully and will refresh soon'),
      error: intl.t(`message.removeFailed`, 'Remove failed, please try again'),
    },
  );
  const handleDisable = createHandler(
    !isFunction(request) ? request.disable : undefined,
    ({ id }) => ({ id, isDisable: true }),
    {
      processing: intl.t(`message.disabling`, 'Disabling ...'),
      success: intl.t(`message.disableSuccessd`, 'Disabled successfully and will refresh soon'),
      error: intl.t(`message.disableFailed`, 'Disable failed, please try again'),
    },
  );

  const handleEnable = createHandler(
    !isFunction(request) ? request.enable : undefined,
    ({ id }) => ({ id, isDisable: false }),
    {
      processing: intl.t(`message.enabling`, 'Enabling ...'),
      success: intl.t(`message.enableSuccessd`, 'Enabled successfully and will refresh soon'),
      error: intl.t(`message.enableFailed`, 'Enable failed, please try again'),
    },
  );

  useEffect(() => {
    actionRef.current?.reload();
  }, [fetchList]);

  const renderTableColumns = (cs?: ProColumns<TableItem>[]): ProColumns<TableItem>[] => {
    const newColumns = (cs ?? [])?.map((item, idx) => {
      if (idx === 0) {
        return {
          ...item,
          render: (dom: React.ReactNode, entity: TableItem) => {
            return (
              <a
                onClick={() => {
                  setCurrentRow(entity);
                }}
              >
                {dom}
              </a>
            );
          },
        };
      }
      return item;
    });
    return [
      ...newColumns,
      {
        title: intl.t('button.option', 'Option'),
        width: 180,
        key: 'option',
        valueType: 'option',
        render: (_: React.ReactNode, entity: TableItem) => [
          <a
            key="edit"
            onClick={() => {
              history.push(`${location.pathname}/${entity.id}`);
            }}
          >
            {intl.t('button.edit', 'Edit')}
          </a>,
        ],
      },
    ];
  };

  return (
    <>
      <ProTable<TableItem>
        pagination={{
          defaultPageSize: 20,
          showSizeChanger: true,
        }}
        rowKey="id"
        actionRef={actionRef}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
        toolbar={{
          search: true,
          onSearch: (kws) => {
            setKeywords(kws);
            actionRef.current?.reload();
          },
          actions: [
            <Button
              key="create"
              type="primary"
              onClick={() => {
                history.push(`${location.pathname}/create`);
              }}
            >
              <PlusOutlined />
              {intl.t('button.create', 'Create')}
            </Button>,
          ],
        }}
        columns={renderTableColumns(columns)}
        // metas={{
        //   title: { dataIndex: 'name' },
        //   subTitle: {
        //     render: (_, entry) => {
        //       if (entry.isDisable) {
        //         return <Tag color="red">{intl.t('button.disabled', 'Disabled')}</Tag>;
        //       }
        //       // return <Tag color="#5BD8A6">{intl.t('button.Activated', 'Activated')}</Tag>;
        //       return <></>;
        //     },
        //   },
        //   content: {
        //     render: (_, entry) => {
        //       return <Typography.Text ellipsis>{entry.description ?? ''}</Typography.Text>;
        //     },
        //   },
        //   avatar: {
        //     dataIndex: 'avatar',
        //     search: false,
        //     render: (_, entity) => {
        //       if (isString(entity.avatar)) {
        //         return <Avatar size="small" src={entity.avatar} />;
        //       }
        //       return entity.avatar;
        //     },
        //   },
        //   actions: {
        //     cardActionProps: 'actions',
        //     render: (_, entry) => [
        //       <a
        //         onClick={() => {
        //           history.push(`${location.pathname}/${entry.id}`);
        //         }}
        //         key="edit"
        //       >
        //         {intl.t('button.edit', 'Edit')}
        //       </a>,
        //       handleRemove ? (
        //         <a
        //           onClick={() => {
        //             Modal.confirm({
        //               title: intl.t('confirm.remove', 'Are you sure you want to remove this page?'),
        //               icon: <ExclamationCircleOutlined />,
        //               onOk() {
        //                 handleRemove(entry.id);
        //               },
        //               maskClosable: true,
        //             });
        //           }}
        //           key="remove"
        //         >
        //           {intl.t('button.remove', 'Remove')}
        //         </a>
        //       ) : null,
        //       entry.isDisable !== undefined && (entry.isDisable ? handleEnable : handleDisable) ? (
        //         <a
        //           onClick={() => {
        //             Modal.confirm({
        //               title: entry.isDisable
        //                 ? intl.t('button.enable', 'Are you sure you want to enable this page?')
        //                 : intl.t('button.disable', 'Are you sure you want to disable this page?'),
        //               icon: <ExclamationCircleOutlined />,
        //               onOk() {
        //                 handleStatusSwitch(entry.id, entry.isDisable ?? false);
        //               },
        //               maskClosable: true,
        //             });
        //           }}
        //           key="disable_or_enable"
        //         >
        //           {entry.isDisable
        //             ? intl.t('button.enable', 'Enable')
        //             : intl.t('button.disable', 'Disable')}
        //         </a>
        //       ) : undefined,
        //     ],
        //   },
        //   ...metas,
        // }}
        request={async (params) => {
          return fetchList(params);
        }}
        {...proTableProps}
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
          <BatchPopconfirm
            onOk={handleRemove}
            onFinish={async () => {
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
            icon={<ExclamationCircleOutlined />}
            items={selectedRowsState}
            title={intl.t(
              'deleteConfirm',
              'Are you sure you want to delete the following users?            ',
            )}
          >
            {intl.t('batchDeletion', 'Batch deletion')}
          </BatchPopconfirm>

          <BatchPopconfirm
            onOk={handleDisable}
            onFinish={async () => {
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
            icon={<ExclamationCircleOutlined />}
            items={selectedRowsState.filter((item) => item.isDisable === false)}
            title={intl.t(
              'disableConfirm',
              'Are you sure you want to disable the following users?',
            )}
          >
            {intl.t('batchDisable', 'Batch disable')}
          </BatchPopconfirm>

          <BatchPopconfirm
            onOk={handleEnable}
            onFinish={async () => {
              setSelectedRows([]);
              actionRef.current?.reloadAndRest?.();
            }}
            icon={<ExclamationCircleOutlined />}
            items={selectedRowsState.filter((item) => item.isDisable === true)}
            title={intl.t('enableConfirm', 'Are you sure you want to enable the following users?')}
          >
            {intl.t('batchDisable', 'Batch enable')}
          </BatchPopconfirm>
        </FooterToolbar>
      )}
      <Drawer
        open={currentRow !== undefined}
        width={800}
        onClose={() => {
          setCurrentRow(undefined);
        }}
      >
        <ProDescriptions columns={columns as any} dataSource={currentRow} />
      </Drawer>
    </>
  );
};
export default Table;
