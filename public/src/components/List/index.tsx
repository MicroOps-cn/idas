import { Button, message, Modal, Tag, Typography } from 'antd';
import { isFunction, isString } from 'lodash';
import React, { useState, useRef, useEffect, useCallback } from 'react';

import Avatar from '@/components/Avatar';
import type { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import { ExclamationCircleOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import type { ProListMetas, ProListProps } from '@ant-design/pro-list';
import ProList from '@ant-design/pro-list';
import type { RequestData } from '@ant-design/pro-table';
import type { ActionType } from '@ant-design/pro-table/lib/typing';
import { getLocale } from '@umijs/max';

import styles from './index.less';

interface ListItem {
  id: string;
  avatar?: string;
  name?: string;
  displayName?: string;
  description?: React.ReactNode;
  isDisable?: boolean;
  i18n?: {
    displayName?: Record<string, string>;
    description?: Record<string, string>;
  };
}

type ListRequest<T extends ListItem> = (
  params: {
    pageSize?: number;
    current?: number;
    keywords?: string;
  },
  options?: Record<string, any>,
) => Promise<Partial<RequestData<T>>>;

type PatchRequst = (
  body: { id: string; isDelete?: boolean; isDisable?: boolean }[],
  options?: Record<string, any>,
) => Promise<any>;
type StatusSwitchRequest = (body: { id: string }[], options?: Record<string, any>) => Promise<any>;

interface ListProps<T extends ListItem> extends Omit<ProListProps, 'request' | 'metas'> {
  intl: IntlContext;
  request:
    | ListRequest<T>
    | {
        list: ListRequest<T>;
        patch?: PatchRequst;
        disable?: StatusSwitchRequest;
        enable?: StatusSwitchRequest;
        delete?: StatusSwitchRequest;
      };
  metas?: ProListMetas<T>;
  onEdit?: (item: T) => void;
  onClick?: (item: T) => void;
  onCreate?: () => void;
}

const List = <T extends ListItem>({
  intl,
  request,
  metas,
  onEdit,
  onCreate,
  onClick,
  ...proListProps
}: ListProps<T>) => {
  const defaultPageSize = 20;
  const actionRef = useRef<ActionType>();
  const [keywords, setKeywords] = useState<string>();
  const [isSignalPage, setSignalPage] = useState<boolean>();
  const fetchList = useCallback(
    async (params: {
      pageSize?: number;
      current?: number;
      keyword?: string;
    }): Promise<Partial<RequestData<T>>> => {
      return (isFunction(request) ? request : request.list)({
        pageSize: defaultPageSize,
        ...params,
        keywords: keywords ?? undefined,
      }).then((resp) => {
        if (resp.total === undefined || resp.pageSize === undefined || resp.current === undefined) {
          setSignalPage(true);
        } else if (resp.current <= 1 && resp.total < resp.current * resp.pageSize) {
          setSignalPage(false);
        } else {
          setSignalPage(true);
        }
        return resp;
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
    const handleFunc: StatusSwitchRequest | null =
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
        await handleFunc(isString(id) ? [{ id }] : id.map((itemId) => ({ id: itemId })), { intl });
        hide();
        message.success(success ?? intl.t(`message.operationSuccessd`, 'Operation succeeded.'));
        return true;
      } catch (err) {
        hide();
        if (!(err as RequestError).handled) {
          message.error(
            error ?? intl.t(`message.operationFailed`, 'Operation failed, please try again.'),
          );
        }
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

  return (
    <ProList<T>
      pagination={
        isSignalPage
          ? {
              defaultPageSize: defaultPageSize,
              showSizeChanger: true,
              pageSizeOptions: Array.from(new Set([defaultPageSize, 10, 20, 50, 100])),
            }
          : false
      }
      actionRef={actionRef}
      showActions="always"
      grid={{ gutter: 16, column: 3, xxl: 4, xl: 4, lg: 4, md: 3, sm: 2, xs: 1 }}
      toolbar={{
        search: true,
        onSearch: (kws) => {
          setKeywords(kws);
          actionRef.current?.reload();
        },
        actions: [
          <Button hidden={!onCreate} key="create" type="primary" onClick={onCreate}>
            <PlusOutlined />
            {intl.t('button.create', 'Create')}
          </Button>,
        ],
        settings: [
          {
            icon: <ReloadOutlined onClick={() => actionRef.current?.reload()} />,
            key: 'reload',
            onClick: () => actionRef.current?.reload(),
          },
        ],
      }}
      itemCardProps={{
        className: styles.ListItem,
      }}
      onItem={(item) => {
        return {
          onClick: () => {
            onClick?.(item);
          },
        };
      }}
      metas={{
        title: {
          render: (_, entry) => {
            const i8nDisplayName = (entry.i18n?.displayName ?? {})[getLocale()];
            if (i8nDisplayName) {
              return i8nDisplayName;
            }
            return entry.displayName ?? entry.name;
          },
        },
        subTitle: {
          render: (_, entry) => {
            if (entry.isDisable) {
              return <Tag color="red">{intl.t('button.disabled', 'Disabled')}</Tag>;
            }
            return <></>;
          },
        },
        content: {
          render: (_, entry) => {
            let description: React.ReactNode = (entry.i18n?.description ?? {})[getLocale()];
            if (!description) {
              description = entry.description;
            }
            return (
              <Typography.Paragraph ellipsis={{ rows: 2, tooltip: description }}>
                {description ?? ''}
              </Typography.Paragraph>
            );
          },
        },
        avatar: {
          dataIndex: 'avatar',
          search: false,
          render: (_, entity) => {
            if (isString(entity.avatar)) {
              return <Avatar size="default" src={entity.avatar} />;
            }
            return entity.avatar;
          },
        },
        actions:
          onEdit || handleRemove || handleEnable || handleDisable
            ? {
                cardActionProps: 'actions',
                render: (_, entry) => [
                  onEdit ? (
                    <a
                      onClick={() => {
                        onEdit?.(entry);
                      }}
                      key="edit"
                    >
                      {intl.t('button.edit', 'Edit')}
                    </a>
                  ) : null,
                  handleRemove ? (
                    <a
                      onClick={() => {
                        Modal.confirm({
                          title: intl.t(
                            'confirm.remove',
                            'Are you sure you want to remove this item?',
                          ),
                          icon: <ExclamationCircleOutlined />,
                          onOk() {
                            handleRemove(entry.id);
                          },
                          maskClosable: true,
                        });
                      }}
                      key="remove"
                    >
                      {intl.t('button.remove', 'Remove')}
                    </a>
                  ) : null,
                  entry.isDisable !== undefined &&
                  (entry.isDisable ? handleEnable : handleDisable) ? (
                    <a
                      onClick={() => {
                        Modal.confirm({
                          title: entry.isDisable
                            ? intl.t('confirm.enable', 'Are you sure you want to enable this item?')
                            : intl.t(
                                'confirm.disable',
                                'Are you sure you want to disable this item?',
                              ),
                          icon: <ExclamationCircleOutlined />,
                          onOk() {
                            (entry.isDisable ? handleEnable : handleDisable)?.(entry.id);
                          },
                          maskClosable: true,
                        });
                      }}
                      key="disable_or_enable"
                    >
                      {entry.isDisable
                        ? intl.t('button.enable', 'Enable')
                        : intl.t('button.disable', 'Disable')}
                    </a>
                  ) : undefined,
                ],
              }
            : undefined,
        ...metas,
      }}
      request={fetchList}
      {...proListProps}
    />
  );
};
export default List;
