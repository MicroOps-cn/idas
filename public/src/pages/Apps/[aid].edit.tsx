import { message } from 'antd';
import { useEffect, useState } from 'react';
import { useIntl, useParams, history } from 'umi';

import { createApp, getAppInfo, updateApp } from '@/services/idas/apps';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import ProCard from '@ant-design/pro-card';

import type { FormValueType, PartialByKeys } from './components/CreateOrUpdateForm';
import CreateOrUpdateForm from './components/CreateOrUpdateForm';

/**
 * @en-US Add node
 * @zh-CN 添加应用
 * @param fields
 */
const handleAdd = async ({
  status: _,
  ...fields
}: API.CreateAppRequest & { id?: string; status?: API.AppMetaStatus }) => {
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
const handleUpdate = async (fields: PartialByKeys<API.UpdateAppRequest, 'id'>) => {
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
          oAuth2: fields.oAuth2,
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
  const [appInfo, setAppInfo] = useState<Partial<API.AppInfo> | undefined>(undefined);
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
      setAppInfo(undefined);
    }
  }, [aid]);
  return (
    // <PageContainer title={aid ? appInfo?.name : false}>
    <ProCard style={{ height: '100%' }}>
      <CreateOrUpdateForm
        onSubmit={async (value: FormValueType) => {
          const success = await (appInfo?.id ? handleUpdate : handleAdd)(value);
          if (success) {
            history.push(`/apps/${appInfo?.id ?? ''}`);
          }
          return success;
        }}
        values={aid ? appInfo : {}}
        parentIntl={intl}
        loading={loading}
        disabled={disabled}
      />
    </ProCard>
    // </PageContainer>
  );
};
export default AppEditor;
