import type { FormInstance } from 'antd';
import { Input, Tabs, message } from 'antd';
import 'antd/dist/antd.css';
import type { DefaultOptionType } from 'antd/lib/select';
import React, { useEffect, useRef, useState } from 'react';

import { AvatarUploadField } from '@/components/Avatar';
import { allLocales } from '@/components/SelectLang';
import { getAppIcons } from '@/services/idas/apps';
import { AppStatus, GrantMode, GrantType } from '@/services/idas/enums';
import type { LabelValue } from '@/utils/enum';
import { enumToOptions } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import { newId } from '@/utils/uuid';
import {
  ProFormGroup,
  ProFormItem,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  StepsForm,
} from '@ant-design/pro-form';

import GrantView from './GrantView';
import ProxySetting from './ProxySetting';
import { RoleView } from './RoleView';

const TextArea = Input.TextArea;

export type FormValueType = API.UpdateAppRequest;
export type UpdateFormProps = {
  onSubmit: (values: FormValueType) => Promise<boolean>;
  values?: Partial<API.AppInfo>;
  title?: React.ReactNode;
  parentIntl: IntlContext;
  loading: boolean;
  disabled: boolean;
};

const CreateOrUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const { values: initialValues, onSubmit, parentIntl, loading, disabled } = props;
  const intl = new IntlContext('form', parentIntl);
  // const [avatar, setAvatar] = useState<UploadFile>();
  const [currentGrantType, setCurrentGrantType] = useState<API.AppMetaGrantType[]>([]);
  const actionRef = useRef<FormInstance<any>>();
  const [grantedUserList, setGrantedUserList] = useState<API.UserInfo[]>([]);
  const [appRoles, setAppRoles] = useState<API.AppRoleInfo[]>([]);
  const [i18n, setI18n] = useState<Required<API.AppI18NOptions>>({
    description: {},
    displayName: {},
  });
  const [proxyConfig, setProxyConfig] = useState<API.AppProxyInfo>({
    domain: '',
    upstream: '',
    urls: [{ id: newId(), method: '*', name: 'default', url: '/' }],
    insecureSkipVerify: false,
    transparentServerName: true,
    hstsOffload: false,
    jwtProvider: false,
    jwtCookieName: '',
    jwtSecret: '',
  });
  const [currentStep, setCurrentStep] = useState<number>(0);
  const groupGrantType: (x: LabelValue[]) => DefaultOptionType[] = (options: LabelValue[]) => {
    const newOptions: (DefaultOptionType & { key: string })[] = [
      {
        key: 'oauth',
        label: intl.t(`grantType.value.oauth`, 'OAuth/OIDC'),
        children: [],
      },
      {
        key: 'other',
        label: intl.t(`grantType.value.other`, 'Other'),
        children: [],
      },
    ];
    const groupKeys: Record<string, string[]> = {
      oauth: ['authorization_code', 'password', 'client_credentials', 'implicit', 'oidc'],
    };
    options.forEach((option: DefaultOptionType) => {
      if (option.key == 'none') {
        option.disabled = true;
      }
      for (const groupKey in groupKeys) {
        if (Object.prototype.hasOwnProperty.call(groupKeys, groupKey)) {
          const keys = groupKeys[groupKey];
          if (keys.indexOf(option.key) >= 0) {
            for (const newOption of newOptions) {
              if (newOption.key === groupKey) {
                if (!newOption.children) {
                  newOption.children = [];
                }
                newOption.children.push(option);
                return;
              }
            }
          }
        }
      }
      for (const newOption of newOptions) {
        if (newOption.key === 'other') {
          if (!newOption.children) {
            newOption.children = [];
          }
          newOption.children.push(option);
          return;
        }
      }
    });
    return newOptions;
  };
  useEffect(() => {
    setCurrentStep(0);
    setAppRoles(
      initialValues?.roles?.map((r) => ({
        id: r.id,
        urls: r.urls ?? [],
        name: r.name,
        isDefault: r.isDefault,
      })) ?? [],
    );
    setGrantedUserList(initialValues?.users ?? []);
    setProxyConfig({
      domain: '',
      upstream: '',
      transparentServerName: true,
      insecureSkipVerify: false,
      hstsOffload: false,
      jwtProvider: false,
      jwtCookieName: '',
      jwtSecret: '',
      ...(initialValues?.proxy ?? {}),
      urls: initialValues?.proxy?.urls ?? [{ id: newId(), method: '*', name: 'default', url: '/' }],
    });
    setCurrentGrantType(initialValues?.grantType ?? []);
    setI18n({
      displayName: {},
      description: {},
      ...(initialValues?.i18n ?? {}),
    });
  }, [initialValues]);

  const hasProxy = (type?: API.AppMetaGrantType[]) => {
    return (type ?? ([] as any)).includes(GrantType.proxy);
  };
  const locales = allLocales();
  return (
    <StepsForm<FormValueType & { _roles: any }>
      stepsProps={{
        size: 'small',
      }}
      formProps={{
        preserve: false,
        disabled: disabled,
        loading: loading,
      }}
      onCurrentChange={async (current) => {
        if (current == 1 && !hasProxy(actionRef.current?.getFieldValue('grantType'))) {
          setCurrentStep(current + (currentStep === 0 ? 1 : -1));
        } else {
          setCurrentStep(current);
        }
      }}
      current={currentStep}
      onFinish={async ({ _roles: _, ...values }) => {
        const { manual } = GrantMode;
        const { grantMode, grantType, status } = values;
        const {
          domain,
          upstream,
          urls: rawUrls,
          transparentServerName,
          insecureSkipVerify,
          hstsOffload,
          jwtProvider,
          jwtSecret,
          jwtCookieName,
        } = proxyConfig;
        const urls = rawUrls.map(({ id, method, name, url, upstream: urlUpstream }) => ({
          id,
          method,
          name,
          url,
          upstream: urlUpstream,
        }));
        return onSubmit({
          ...values,
          users: grantedUserList.map((user) => ({
            id: user.id,
            roleId: user.roleId,
          })),
          status: status !== AppStatus.unknown ? status : AppStatus.normal,
          roles: appRoles.map((role) =>
            hasProxy(grantType) ? role : { ...role, urls: undefined },
          ),
          grantMode: grantMode ?? manual,
          grantType: grantType ?? [GrantType.none as any],
          proxy: hasProxy(grantType)
            ? {
                domain,
                upstream,
                urls,
                insecureSkipVerify,
                transparentServerName,
                hstsOffload,
                jwtProvider,
                jwtSecret,
                jwtCookieName,
              }
            : undefined,
          i18n,
        });
      }}
    >
      <StepsForm.StepForm
        initialValues={initialValues}
        formRef={actionRef}
        layout={'vertical'}
        style={{ maxWidth: 650 }}
        grid={true}
        title={intl.t('basicConfig.title', 'Basic')}
        onFinish={async () => {
          if (initialValues && !initialValues.id) {
            message.error(intl.t('app-id.empty', 'System error, application ID is empty'));
            return false;
          }
          return true;
        }}
      >
        <ProFormText hidden={true} name="id" />
        <AvatarUploadField
          colProps={{ span: 8, sm: 8, xs: 14 }}
          label={intl.t('avatar.label', 'Avatar')}
          name={'avatar'}
          optionsRequest={async (params) => {
            return getAppIcons(params);
          }}
        />
        <ProFormGroup grid colProps={{ span: 16, xs: 14 }}>
          <ProFormText
            name="name"
            label={intl.t('name.label', 'Name')}
            width="md"
            rules={[
              {
                required: true,
                message: intl.t('name.required', 'Please input app name!'),
              },
              {
                pattern: /^[-_A-Za-z0-9]+$/,
                message: intl.t('name.invalid', 'App name format error!'),
              },
            ]}
            disabled={initialValues?.name ? true : false}
          />

          <Tabs
            style={{ width: '100%' }}
            items={[
              {
                label: intl.t('displayName.label', 'Display Name'),
                children: <ProFormText name="displayName" width="md" />,
                key: 'displayName',
              },
              ...locales.map((lang) => ({
                label: intl.t(`displayName.label.${lang}`, lang),
                children: (
                  <Input
                    value={i18n?.displayName?.[lang] ?? ''}
                    onChange={(e) => {
                      setI18n({
                        displayName: { ...i18n.displayName, [lang]: e.target.value },
                        description: i18n.description,
                      });
                    }}
                  />
                ),
                key: `description-${lang}`,
              })),
            ]}
          />
        </ProFormGroup>

        <Tabs
          style={{ width: '100%' }}
          items={[
            {
              label: intl.t('description.label', 'Description'),
              children: (
                <ProFormTextArea colProps={{ span: 24, sm: 24, xs: 14 }} name="description" />
              ),
              key: 'description',
            },
            ...locales.map((lang) => ({
              label: intl.t(`description.label.${lang}`, lang),
              children: (
                <TextArea
                  value={i18n?.description?.[lang] ?? ''}
                  onChange={(e) => {
                    setI18n({
                      description: { ...i18n.description, [lang]: e.target.value },
                      displayName: i18n.displayName,
                    });
                  }}
                />
              ),
              key: `description-${lang}`,
            })),
          ]}
        />
        <ProFormText
          name="url"
          colProps={{ span: 24, sm: 24, xs: 14 }}
          label={intl.t('url.label', 'URL')}
        />
        <ProFormSelect<API.AppMetaGrantType[]>
          colProps={{ span: 12, sm: 12, xs: 14 }}
          name="grantType"
          label={intl.t('grantType.label', 'Grant Type')}
          width="md"
          mode={'tags'}
          options={groupGrantType(enumToOptions(GrantType, parentIntl, 'grantType.value'))}
          fieldProps={{
            onChange: setCurrentGrantType,
          }}
        />
        <ProFormSelect
          colProps={{ span: 12, sm: 12, xs: 14 }}
          name="grantMode"
          label={intl.t('grantMode.label', 'Grant Mode')}
          width="md"
          tooltip={intl.t(
            'grantMode.tooltip',
            'Automatic authorization is only supported when using OAuth2.0 and RADIUS protocols, otherwise manual authorization is only possible.',
          )}
          options={enumToOptions(GrantMode, parentIntl, 'grantMode.value')}
          rules={[
            {
              required: true,
              message: intl.t('grantMode.required', 'Please select Grant Mode!'),
            },
          ]}
        />
      </StepsForm.StepForm>

      <StepsForm.StepForm
        initialValues={{}}
        // labelCol={{ span: 5 }}
        // wrapperCol={{ span: 19 }}
        layout={'vertical'}
        grid={true}
        title={intl.t('proxy.title', 'Proxy')}
        onFinish={async (formVals) => {
          const {
            domain,
            upstream,
            transparentServerName,
            insecureSkipVerify,
            jwtProvider,
            jwtSecret,
            jwtCookieName,
            hstsOffload,
          } = formVals.proxy ?? {};
          setProxyConfig({
            ...proxyConfig,
            domain,
            upstream,
            transparentServerName,
            insecureSkipVerify,
            jwtProvider,
            jwtSecret,
            jwtCookieName,
            hstsOffload,
          });
          const { urls } = proxyConfig;
          // if (!domain) {
          //   message.error(intl.t('proxy.domain.required', 'domain cannot be empty!'));
          //   return false;
          // }
          for (const url of urls) {
            if (!url.name || !url.name.trim()) {
              message.error(intl.t('name.required', 'name cannot be empty!'));
              return false;
            }
            if (!url.method || !url.method.trim()) {
              message.error(intl.t('proxy.method.required', 'method cannot be empty!'));
              return false;
            }
            if (!url.url || !url.url.trim()) {
              message.error(intl.t('proxy.url.required', 'URL cannot be empty!'));
              return false;
            }
          }
          return true;
        }}
      >
        {currentStep === 1 && (
          <ProxySetting dataSource={proxyConfig} setDataSource={setProxyConfig} parentIntl={intl} />
        )}
      </StepsForm.StepForm>
      <StepsForm.StepForm
        initialValues={{}}
        layout={'vertical'}
        title={intl.t('role.title', 'Role')}
      >
        {currentStep === 2 && (
          <ProFormItem
            name={'_roles'}
            rules={[
              {
                validator: () => {
                  if (appRoles.length > 0) {
                    if (appRoles.filter((role) => role.isDefault).length != 1) {
                      return Promise.reject(new Error('No default role specified!'));
                    } else {
                      const appNames = appRoles.map((val) => val.name);
                      if (appNames.some((val, idx) => appNames.includes(val, idx + 1))) {
                        return Promise.reject(new Error('Contains duplicate role names!'));
                      } else if (appNames.some((val) => !Boolean(val))) {
                        return Promise.reject(new Error('Contains empty role names!'));
                      }
                    }
                  }
                  return Promise.resolve();
                },
              },
            ]}
          >
            <RoleView
              value={appRoles}
              urls={hasProxy(currentGrantType) ? proxyConfig.urls : []}
              onChange={setAppRoles}
            />
          </ProFormItem>
        )}
      </StepsForm.StepForm>
      <StepsForm.StepForm
        initialValues={{}}
        layout={'vertical'}
        title={intl.t('user.title', 'User')}
      >
        {currentStep === 3 && (
          <div
            style={{
              height: 'calc( 100vh - 400px )',
            }}
          >
            <GrantView
              users={grantedUserList}
              roles={appRoles}
              onChange={setGrantedUserList}
              granting={true}
              parentIntl={parentIntl}
            />
          </div>
        )}
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateOrUpdateForm;
