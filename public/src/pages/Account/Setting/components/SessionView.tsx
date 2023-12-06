import { message } from 'antd';
import moment from 'moment';
import React, { useRef } from 'react';

import { deleteCurrentUserSession, getCurrentUserSessions } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { LogoutOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';

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
    await deleteCurrentUserSession({
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

const SessionView: React.FC<{ currentUser?: API.UserInfo; parentIntl: IntlContext }> = ({
  currentUser,
  parentIntl,
}) => {
  const intl = new IntlContext('sessions', parentIntl);
  const sessionRef = useRef<ActionType>();
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
  return currentUser ? (
    <ProTable
      actionRef={sessionRef}
      toolBarRender={false}
      request={(params) => {
        return getCurrentUserSessions({ ...params, userId: currentUser.id });
      }}
      columns={sesionColumns}
      rowKey="id"
      search={false}
    />
  ) : null;
};

export default SessionView;
