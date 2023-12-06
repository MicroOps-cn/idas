import { Input, Skeleton, Divider, List, Select, Typography } from 'antd';
import React, { useState } from 'react';
import { useEffect } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';

import Avatar from '@/components/Avatar';
import { getApps } from '@/services/idas/apps';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';

import styles from '../index.less';

interface AppUserViewProps {
  loading?: boolean;
  apps: API.UserApp[];
  onChange?: (apps: API.UserApp[]) => Promise<void> | void;
  granting?: boolean;
  parentIntl: IntlContext;
}

export const AppUserView: React.FC<AppUserViewProps> = ({
  onChange,
  loading,
  granting,
  parentIntl,
  apps,
}) => {
  const intl = new IntlContext('apps', parentIntl);
  return (
    <List<API.UserApp>
      dataSource={apps}
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
                      onChange(apps.filter((app) => app.id != item.id));
                    }}
                  >
                    {intl.t('appList.delete', 'Delete')}
                  </a>,
                ]
              : []
          }
        >
          <List.Item.Meta
            avatar={<Avatar src={`${item.avatar}`} />}
            title={item.displayName ?? item.name}
            description={
              <Typography.Paragraph type="secondary" ellipsis={{ tooltip: item.description }}>
                {item.description ?? ''}
              </Typography.Paragraph>
            }
          />
          {onChange && granting && item.roles && item.roles.length > 0 ? (
            <Select<string>
              onSelect={(val: string) => {
                onChange(
                  apps.map((app) => {
                    if (app.id == item.id) {
                      return { ...app, roleId: val };
                    }
                    return app;
                  }),
                );
              }}
              defaultValue={
                item && item.roleId ? item.roleId : item.roles.find((role) => role.isDefault)?.id
              }
              options={item.roles.map((role) => ({
                key: role.id,
                value: role.id,
                label: role.name,
              }))}
            />
          ) : (
            <div>
              {item && item.role ? item.role : item.roles?.find((role) => role.isDefault)?.name}
            </div>
          )}
        </List.Item>
      )}
    />
  );
};

interface GrantViewProps {
  apps: API.UserApp[];
  onChange: (apps: API.UserApp[]) => Promise<void> | void;
  granting?: boolean;
  loading?: boolean;
  parentIntl: IntlContext;
  type?: 'vertical' | 'horizontal';
}
export const GrantView: React.FC<GrantViewProps> = ({ type = 'vertical', loading, ...props }) => {
  const { onChange, apps, granting, parentIntl } = props;
  const intl = new IntlContext('grant', parentIntl);

  const [appListHasMore, setAppListHasMore] = useState<boolean>(true);
  const [appListPageNumber, setAppListPageNumber] = useState<number>(0);
  const [appListKeyworlds, setAppListKeyworlds] = useState<string>();
  const [appList, setAppList] = useState<API.AppInfo[]>([]);

  // const [grantedAppList, setGrantedAppList] = useState<API.User[]>([]);
  const [loadingAppList, setLoadingAppList] = useState<boolean>(false);

  const loadMoreAppList = async (params?: API.getUsersParams) => {
    try {
      setLoadingAppList(true);
      const resp = await getApps({
        current: appListPageNumber + 1,
        pageSize: 20,
        keywords: appListKeyworlds,
        ...params,
      });
      if (resp && resp.data) {
        const { data: newData, current, pageSize, total } = resp;
        setAppList((oldData) => {
          return [...oldData, ...newData];
        });
        setAppListPageNumber(current);
        if (total < current * pageSize) {
          setAppListHasMore(false);
        }
      } else {
        setAppListHasMore(false);
      }
    } catch (error) {
      if (!(error as RequestError).handled) {
        console.error(`failed to get app list: ${error}`);
      }
    } finally {
      setLoadingAppList(false);
    }
  };

  const loadMoreAppListByKeyworlds = (keywords?: string) => {
    setAppList([]);
    setAppListKeyworlds(keywords);
    setAppListPageNumber(1);
    loadMoreAppList({ keywords, current: 1 });
  };

  useEffect(() => {
    loadMoreAppListByKeyworlds();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div
      style={{
        height: type === 'vertical' ? 'calc((100vh - 300px))' : '100%',
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
                if (!loadingAppList) loadMoreAppListByKeyworlds(value);
              }}
              loading={loadingAppList}
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
                dataLength={appList.length}
                next={loadMoreAppList}
                hasMore={appListHasMore}
                loader={<Skeleton avatar paragraph={{ rows: 1 }} active />}
                endMessage={<Divider plain>End</Divider>}
                scrollableTarget="scrollableDiv"
              >
                <List<API.UserApp>
                  dataSource={appList.filter((u) => !apps.map((a1) => a1.id).includes(u.id))}
                  size={'small'}
                  loading={loadingAppList}
                  renderItem={(item) => (
                    <List.Item
                      key={item.id}
                      actions={[
                        <a
                          key="grant"
                          onClick={() => {
                            onChange([...apps, { ...item, roleId: '' }]);
                          }}
                        >
                          {intl.t('appList.grant', 'Grant')}
                        </a>,
                      ]}
                    >
                      <List.Item.Meta
                        avatar={<Avatar src={`${item.avatar}`} />}
                        title={item.displayName ?? item.name}
                        description={
                          <Typography.Paragraph
                            type="secondary"
                            ellipsis={{ tooltip: item.description }}
                          >
                            {item.description ?? ''}
                          </Typography.Paragraph>
                        }
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
          apps={apps.map((app) => ({
            ...(appList.find((a) => a.id == app.id) ?? app),
            roleId: app.roleId,
            role: app.role,
          }))}
          granting={granting}
          loading={loading}
          onChange={onChange}
        />
      </div>
    </div>
  );
};

export default GrantView;
