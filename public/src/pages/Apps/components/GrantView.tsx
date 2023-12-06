import { Input, Skeleton, Divider, List, Select } from 'antd';
import React, { useState } from 'react';
import { useEffect } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';

import Avatar from '@/components/Avatar';
import { getUsers as getUsersInfo } from '@/services/idas/users';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';

import styles from '../style.less';

interface AppUserViewProps {
  roles: API.AppRoleInfo[];
  loading?: boolean;
  users: API.UserInfo[];
  onChange?: (users: API.UserInfo[]) => Promise<void> | void;
  granting?: boolean;
  parentIntl: IntlContext;
}

export const AppUserView: React.FC<AppUserViewProps> = ({
  users,
  onChange,
  roles,
  loading,
  granting,
  parentIntl,
}) => {
  const intl = new IntlContext('users', parentIntl);
  return (
    <List<API.UserInfo>
      dataSource={users}
      size={'small'}
      loading={loading}
      renderItem={(item) => (
        <List.Item
          key={item.id}
          actions={
            onChange && granting
              ? [
                  <a
                    key="grant"
                    onClick={() => {
                      onChange(users.filter((user) => user.id != item.id));
                    }}
                  >
                    {intl.t('userList.delete', 'Delete')}
                  </a>,
                ]
              : []
          }
        >
          <List.Item.Meta
            avatar={<Avatar src={`${item.avatar}`} />}
            title={item.fullName ? item.fullName : item.username}
            description={item.email}
          />
          {onChange && granting && roles.length > 0 ? (
            <Select<string>
              onSelect={(val: string) => {
                onChange(
                  users.map((user) => {
                    if (user.id == item.id) {
                      return { ...user, roleId: val };
                    }
                    return user;
                  }),
                );
              }}
              defaultValue={
                item && item.roleId ? item.roleId : roles.find((role) => role.isDefault)?.id
              }
              options={roles.map((role) => ({ key: role.id, value: role.id, label: role.name }))}
            />
          ) : (
            <div>{item && item.role ? item.role : roles.find((role) => role.isDefault)?.name}</div>
          )}
        </List.Item>
      )}
    />
  );
};

interface GrantViewProps {
  roles: API.AppRoleInfo[];
  users: API.UserInfo[];
  onChange: (users: API.UserInfo[]) => Promise<void> | void;
  granting?: boolean;
  loading?: boolean;
  parentIntl: IntlContext;
  type?: 'vertical' | 'horizontal';
}
export const GrantView: React.FC<GrantViewProps> = ({ type = 'vertical', loading, ...props }) => {
  const { roles, onChange, users, granting, parentIntl } = props;
  const intl = new IntlContext('grant', parentIntl);

  const [userListHasMore, setUserListHasMore] = useState<boolean>(true);
  const [userListPageNumber, setUserListPageNumber] = useState<number>(0);
  const [userListKeyworlds, setUserListKeyworlds] = useState<string>();
  const [userList, setUserList] = useState<API.UserInfo[]>([]);

  // const [grantedUserList, setGrantedUserList] = useState<API.User[]>([]);
  const [loadingUserList, setLoadingUserList] = useState<boolean>(false);

  const loadMoreUserList = async (params?: API.getUsersParams) => {
    try {
      setLoadingUserList(true);
      const resp = await getUsersInfo({
        current: userListPageNumber + 1,
        pageSize: 20,
        keywords: userListKeyworlds,
        ...params,
      });
      if (resp && resp.data) {
        const { data: newData, current, pageSize, total } = resp;
        setUserList((oldData) => {
          return [...oldData, ...newData];
        });
        setUserListPageNumber(current);
        if (total < current * pageSize) {
          setUserListHasMore(false);
        }
      } else {
        setUserListHasMore(false);
      }
    } catch (error) {
      if (!(error as RequestError).handled) {
        console.error(`failed to get user list: ${error}`);
      }
    } finally {
      setLoadingUserList(false);
    }
  };

  const loadMoreUserListByKeyworlds = (keywords?: string) => {
    setUserList([]);
    setUserListKeyworlds(keywords);
    setUserListPageNumber(1);
    loadMoreUserList({ keywords, current: 1 });
  };

  useEffect(() => {
    loadMoreUserListByKeyworlds();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div
      style={{
        height: '100%',
        width: '100%',
      }}
    >
      {granting && (
        <>
          <div
            style={{
              height: type === 'vertical' ? 'calc(100% / 2)' : '100%',
              display: granting ? 'block' : 'none',
              width: type === 'vertical' ? '100%' : 'calc((100% - 32px) / 2)',
              float: type === 'horizontal' ? 'left' : 'unset',
            }}
          >
            <Input.Search
              className={styles.SearchInput}
              onSearch={(value) => {
                if (!loadingUserList) loadMoreUserListByKeyworlds(value);
              }}
              loading={loadingUserList}
            />
            <div
              id="scrollableDiv"
              style={{
                height: 'calc( 100% - 32px )',
                overflow: 'auto',
                padding: '0 16px',
                border: '1px solid rgba(140, 140, 140, 0.35)',
              }}
            >
              <InfiniteScroll
                dataLength={userList.length}
                next={loadMoreUserList}
                hasMore={userListHasMore}
                loader={<Skeleton avatar paragraph={{ rows: 1 }} active />}
                endMessage={<Divider plain>End</Divider>}
                scrollableTarget="scrollableDiv"
              >
                <List<API.UserInfo>
                  dataSource={userList.filter((u) => !users.map((u1) => u1.id).includes(u.id))}
                  size={'small'}
                  loading={loadingUserList}
                  renderItem={(item) => (
                    <List.Item
                      key={item.id}
                      actions={[
                        <a
                          key="grant"
                          onClick={() => {
                            onChange([...users, { ...item, role: '' }]);
                          }}
                        >
                          {intl.t('userList.grant', 'Grant')}
                        </a>,
                      ]}
                    >
                      <List.Item.Meta
                        avatar={<Avatar src={`${item.avatar}`} />}
                        title={item.fullName ? item.fullName : item.username}
                        description={item.email}
                      />
                    </List.Item>
                  )}
                />
              </InfiniteScroll>
            </div>
          </div>
          <Divider
            type={type === 'vertical' ? 'horizontal' : 'vertical'}
            style={{
              height: type === 'vertical' ? 'unset' : '100%',
              margin: type === 'vertical' ? '10px 0' : '0 10px',
              float: type === 'horizontal' ? 'left' : 'unset',
              display: granting ? '' : 'none',
            }}
          />
        </>
      )}

      <div
        id="scrollableDiv"
        style={{
          height: type !== 'vertical' || !granting ? '100%' : 'calc((100% - 42px) / 2)',
          width: type === 'horizontal' && granting ? 'calc((100% - 32px) / 2)' : '100%',
          overflow: 'auto',
          padding: '0 16px',
          float: type === 'horizontal' ? 'left' : 'unset',
          border: '1px solid rgba(140, 140, 140, 0.35)',
        }}
      >
        <AppUserView
          parentIntl={parentIntl}
          users={users}
          granting={granting}
          roles={roles}
          loading={loading}
          onChange={onChange}
        />
      </div>
    </div>
  );
};

export default GrantView;
