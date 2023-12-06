import { message } from 'antd';
import 'antd/es/form/style/index.less';
import { useEffect, useState } from 'react';
import { useIntl, useParams, history } from 'umi';

import { createApp, getAppInfo, updateApp } from '@/services/idas/apps';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import ProCard from '@ant-design/pro-card';

import CreateOrUpdateForm from './components/CreateOrUpdateForm';

/**
 * @en-US Add node
 * @zh-CN 添加应用
 * @param fields
 */
const handleAdd = async (fields: API.CreateAppRequest) => {
  const hide = message.loading('Adding ...');
  try {
    await createApp({ ...fields });
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
 * @zh-CN 更新应用
 *
 * @param fields
 */
const handleUpdate = async (fields: API.UpdateAppRequest) => {
  const hide = message.loading('Configuring');
  try {
    if (fields.id) {
      await updateApp(
        { id: fields.id },
        {
          id: fields.id,
          name: fields.name,
          status: fields.status,
          description: fields.description,
          grantType: fields.grantType,
          grantMode: fields.grantMode,
          avatar: fields.avatar,
          roles: fields.roles,
          users: fields.users,
          proxy: fields.proxy,
          url: fields.url,
          displayName: fields.displayName,
          i18n: fields.i18n,
        },
      );
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
      message.error('Configuration failed, please try again!');
    }
    return false;
  }
};
const AppEditor: React.FC = ({}) => {
  const { aid } = useParams();
  const intl = new IntlContext('pages.apps', useIntl());
  const [appInfo, setAppInfo] = useState<API.AppInfo | undefined | Record<string, never>>();
  const [loading, setLoading] = useState<boolean>(false);
  const [disabled, setDisable] = useState<boolean>(false);
  const fetchAppInfo = async (appId: string) => {
    try {
      setLoading(true);
      const resp = await getAppInfo({ id: appId });
      if (resp.data) {
        const { data: app } = resp;
        setAppInfo(app);
      }
      setLoading(false);
      return resp;
    } catch (error) {
      setDisable(true);
      if (!(error as RequestError).handled) {
        console.error(`failed to get app info: ${error}`);
      }
      return { success: false };
    }
  };
  useEffect(() => {
    if (aid) {
      fetchAppInfo(aid);
    } else {
      setAppInfo({});
    }
  }, [aid]);
  return (
    // <PageContainer title={aid ? appInfo?.name : false}>
    <ProCard>
      {appInfo ? (
        <CreateOrUpdateForm
          onSubmit={async (value) => {
            const success = await (appInfo?.id ? handleUpdate : handleAdd)(value);
            if (success) {
              history.push(`/apps/${appInfo?.id ?? ''}`);
            }
            return success;
          }}
          values={appInfo?.id ? appInfo : undefined}
          parentIntl={intl}
          loading={loading}
          disabled={disabled}
        />
      ) : null}
    </ProCard>
    // </PageContainer>
  );
};
export default AppEditor;
