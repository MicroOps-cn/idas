import { Input, Space, Tabs, message } from 'antd';
import { RcFile } from 'antd/es/upload';
import type { DefaultOptionType } from 'antd/lib/select';
import { isArray, isFunction, isString, toInteger } from 'lodash';
import React, { useEffect, useState } from 'react';

import { AvatarUploadField } from '@/components/Avatar';
import { allLocales } from '@/components/SelectLang';
import { getAppIcons } from '@/services/idas/apps';
import type { GrantTypeName, GrantTypeValue, JWTSignatureMethodValue } from '@/services/idas/enums';
import { AppStatus, GrantMode, GrantType, JWTSignatureMethod } from '@/services/idas/enums';
import type { LabelValue } from '@/utils/enum';
import { enumToOptions } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import { newId } from '@/utils/uuid';
import { PageLoading, isUrl } from '@ant-design/pro-components';
import type { StepFormProps } from '@ant-design/pro-form';
import {
  ProFormGroup,
  ProFormItem,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  StepsForm,
} from '@ant-design/pro-form';
import { createHashHistory } from '@umijs/max';

import GrantView from './GrantView';
import ProxySetting from './ProxySetting';
import { RoleView } from './RoleView';

const history = createHashHistory();

const TextArea = Input.TextArea;

export type PartialByKeys<T, K = keyof T> = {
  [Q in keyof T as Q extends K ? Q : never]?: T[Q];
} & {
  [Q in keyof T as Q extends K ? never : Q]: T[Q];
};

type RequiredByKeys<T, K extends keyof T> = {
  [P in K]-?: T[P];
} & Pick<T, Exclude<keyof T, K>>;
export type FormValueType = PartialByKeys<
  RequiredByKeys<API.UpdateAppRequest, 'grantType' | 'roles'>,
  'id'
>;

export type UpdateFormProps = {
  onSubmit: (values: FormValueType) => Promise<boolean>;
  values?: Partial<API.AppInfo>;
  title?: React.ReactNode;
  parentIntl: IntlContext;
  loading: boolean;
  disabled: boolean;
};

type EditorAppInfo = PartialByKeys<
  RequiredByKeys<
    Omit<API.AppInfo, 'createTime' | 'updateTime' | 'isDelete'>,
    'proxy' | 'roles' | 'users' | 'status' | 'oAuth2'
  >,
  'id'
>;

const CreateOrUpdateForm: React.FC<UpdateFormProps> = ({
  values: initialValues,
  onSubmit,
  parentIntl,
  loading,
  disabled,
}) => {
  const [appInfo, _setAppInfo] = useState<EditorAppInfo>({
    i18n: {
      displayName: {},
      description: {},
    },
    id: undefined,
    status: 'normal',
    proxy: {
      domain: '',
      upstream: '',
      transparentServerName: true,
      insecureSkipVerify: false,
      hstsOffload: false,
      jwtProvider: false,
      jwtCookieName: '',
      jwtSecret: '',
      urls: [{ id: newId(), method: '*', name: 'default', url: '/' }],
    },
    users: [],
    roles: [],
    grantType: [],
    url: '',
    name: '',
    oAuth2: {
      authorizedRedirectUrl: [],
      jwtSignatureKey: '',
      jwtSignatureMethod: JWTSignatureMethod.default,
    },
    grantMode: GrantMode.manual,
  });
  const mergeAppInfo = (info: Partial<EditorAppInfo>, oriAppInfo: EditorAppInfo): EditorAppInfo => {
    if (Object.keys(info).length === 0) {
      return oriAppInfo;
    }
    return {
      ...(oriAppInfo ?? {}),
      ...info,
      i18n: {
        ...(info.i18n ?? oriAppInfo.i18n),
      },
      proxy: {
        ...(info.proxy ?? oriAppInfo.proxy),
        urls: (((info.proxy ?? oriAppInfo.proxy)?.urls ?? []) as API.AppProxyUrl[]).map(
          ({ createTime: _, updateTime: __, ...url }) => url,
        ),
      },
      users: info.users ?? oriAppInfo.users,
      roles: info.roles?.map((role) => ({ ...role, urls: role.urls ?? [] })) ?? oriAppInfo.roles,
    };
  };

  const setAppInfo = (
    info: Partial<EditorAppInfo> | ((prevState: EditorAppInfo) => Partial<EditorAppInfo>),
  ) => {
    _setAppInfo((oriAppInfo) => {
      if (isFunction(info)) {
        return mergeAppInfo(info(oriAppInfo), oriAppInfo);
      }
      return mergeAppInfo(info, oriAppInfo);
    });
  };
  const intl = new IntlContext('form', parentIntl);

  const fixOIDC = (value: GrantTypeValue[], prevState: GrantTypeValue[]): GrantTypeValue[] => {
    const ret = [...value];
    if (
      prevState.includes(GrantType.authorization_code) &&
      !ret.includes(GrantType.authorization_code)
    ) {
      return ret.filter((val) => val !== GrantType.oidc);
    }
    if (ret.includes(GrantType.oidc)) {
      if (!ret.includes(GrantType.authorization_code)) {
        ret.push(GrantType.authorization_code);
      }
    }
    return ret;
  };
  const getGrantTypeValue = (val: API.AppMetaGrantType): GrantTypeValue => {
    if (isString(val)) {
      return GrantType[val as GrantTypeName];
    }
    return val as GrantTypeValue;
  };
  const setCurrentGrantType = (value: React.SetStateAction<API.AppMetaGrantType[]>) => {
    setAppInfo((prevState) => {
      const prevGranType = prevState.grantType?.map(getGrantTypeValue) ?? [];
      if (isArray(value)) {
        return {
          grantType: fixOIDC(value.map(getGrantTypeValue), prevGranType),
        };
      } else {
        return {
          grantType: fixOIDC(value(prevGranType).map(getGrantTypeValue), prevGranType),
        };
      }
    });
  };

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

  const hasProxy = (type?: API.AppMetaGrantType[]) => {
    return (type ?? ([] as any)).includes(GrantType.proxy);
  };
  const locales = allLocales();
  const steps: (StepFormProps & { name: string; allows?: GrantTypeValue[] })[] = [
    {
      layout: 'vertical',
      style: { maxWidth: 650 },
      grid: true,
      title: intl.t('basicConfig.title', 'Basic'),
      name: 'basic',
      onFinish: async ({ grantType: _, ...values }) => {
        if (initialValues?.name && !values.id) {
          message.error(intl.t('app-id.empty', 'System error, application ID is empty'));
          return false;
        }
        setAppInfo(values);
        return true;
      },

      children: (
        <>
          <ProFormText hidden={true} name="id" />
          <AvatarUploadField
            colProps={{ span: 8, sm: 8, xs: 14 }}
            label={intl.t('avatar.label', 'Avatar')}
            name={'avatar'}
            optionsRequest={async (params) => {
              const resp = await getAppIcons(params);
              return { ...resp, data: resp.data ?? undefined };
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
                      defaultValue={appInfo.i18n?.displayName?.[lang] ?? ''}
                      onChange={(e) => {
                        setAppInfo(({ i18n: { displayName, ...i18n } = {}, ...oriAppInfo }) => {
                          return {
                            ...oriAppInfo,
                            i18n: {
                              ...i18n,
                              displayName: { ...displayName, [lang]: e.target.value },
                            },
                          };
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
                    defaultValue={appInfo.i18n?.description?.[lang] ?? ''}
                    onChange={(e) => {
                      setAppInfo(({ i18n: { description, ...i18n } = {}, ...oriAppInfo }) => {
                        return {
                          ...oriAppInfo,
                          i18n: {
                            ...i18n,
                            description: { ...description, [lang]: e.target.value },
                          },
                        };
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
            rules={[
              {
                validator: (__, value) => {
                  if (value && isString(value)) {
                    if (!isUrl(value)) {
                      return Promise.reject(
                        new Error(`${intl.t('invalidURL', 'The value is invalid url')}: ${value}`),
                      );
                    }
                  }
                  return Promise.resolve();
                },
              },
            ]}
          />
          <ProFormSelect<GrantTypeValue[]>
            colProps={{ span: 24, sm: 24, xs: 14 }}
            name="grantType"
            label={intl.t('grantType.label', 'Grant Type')}
            mode={'tags'}
            options={groupGrantType(enumToOptions(GrantType, parentIntl, 'grantType.value'))}
            fieldProps={{
              onChange: setCurrentGrantType,
              value: appInfo.grantType?.map(getGrantTypeValue) ?? [],
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
        </>
      ),
    },
    {
      name: 'proxy',
      allows: [GrantType.proxy],
      layout: 'vertical',
      grid: true,
      title: intl.t('proxy.title', 'Proxy'),
      onFinish: async (formVals) => {
        setAppInfo((oriInfo) => {
          return {
            ...oriInfo,
            proxy: {
              ...(formVals.proxy ?? {}),
              urls: formVals.proxy.urls ?? oriInfo.proxy.urls,
            },
          };
        });
        for (const url of formVals.proxy.urls ?? appInfo.proxy.urls) {
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
      },
      children: (
        <ProxySetting
          dataSource={appInfo.proxy ?? {}}
          setDataSource={(values) => {
            setAppInfo((oriInfo) => {
              if (isFunction(values)) {
                return { ...oriInfo, proxy: { ...oriInfo.proxy, ...values(oriInfo.proxy ?? []) } };
              }
              return { ...oriInfo, proxy: { ...oriInfo.proxy, ...values } };
            });
          }}
          parentIntl={intl}
        />
      ),
    },
    {
      name: 'oAuth2',
      allows: [GrantType.oidc, GrantType.authorization_code],
      title: intl.t('oAuth2.title', 'OAuth2'),
      onFinish: async (values) => {
        setAppInfo(values);
        return true;
      },
      children: (
        <>
          <ProFormTextArea
            name={['oAuth2', 'authorizedRedirectUrl']}
            colProps={{ span: 24, sm: 24, xs: 18 }}
            label={intl.t('authorizedRedirectUrl.label', 'Authorized redirect URLs')}
            hidden={!appInfo.grantType.includes(GrantType.authorization_code)}
            rules={[
              {
                validator: (_, value) => {
                  if (isString(value)) {
                    const values = value
                      .split('\n')
                      .map((val) => val.trim())
                      .filter((val) => val);
                    for (let index = 0; index < values.length; index++) {
                      const val = values[index];
                      if (!isUrl(value)) {
                        return Promise.reject(
                          new Error(`${intl.t('invalidURL', 'The value is invalid url')}: ${val}`),
                        );
                      }
                    }
                  }
                  return Promise.resolve();
                },
              },
            ]}
            convertValue={(value) => {
              if (isArray(value)) {
                return value.join('\n');
              }
              return value;
            }}
            transform={(value, name) => {
              return {
                oAuth2: {
                  [name]: isString(value)
                    ? value
                        .split('\n')
                        .map((val) => val.trim())
                        .filter((val) => val)
                    : value,
                },
              };
            }}
          />
          <ProFormSelect<JWTSignatureMethodValue>
            label={intl.t('jwtSignatureMethod.label', 'Custom JWT signature method')}
            colProps={{ span: 4 }}
            name={['oAuth2', 'jwtSignatureMethod']}
            tooltip={intl.t(
              'jwtSignatureMethod.describe',
              'Customize the method and key (pair) used for issuing JWT',
            )}
            options={enumToOptions(JWTSignatureMethod, parentIntl, 'jwtSignatureMethod.value').sort(
              (item) => toInteger(item.value),
            )}
            onChange={(e) => {
              setAppInfo((oriAppInfo) => {
                return { oAuth2: { ...(oriAppInfo.oAuth2 ?? {}), jwtSignatureMethod: e } };
              });
            }}
          />
          <ProFormTextArea
            label={intl.t('jwtSignatureKey.label', 'JWT signature key')}
            colProps={{ span: 24, sm: 24, xs: 18 }}
            name={['oAuth2', 'jwtSignatureKey']}
            hidden={
              JWTSignatureMethod.default === appInfo.oAuth2?.jwtSignatureMethod ??
              JWTSignatureMethod.default
            }
            placeholder={
              initialValues?.id
                ? intl.t('jwtSignatureKey.placeholder-noChange', '------- No change -------')
                : ''
            }
            tooltip={intl.t('jwtSignatureKey.describe', 'For RSA, this value is the private key')}
          />
        </>
      ),
    },
    {
      name: 'role',
      title: intl.t('role.title', 'Role'),
      children: (
        <ProFormItem
          name={'_roles'}
          rules={[
            {
              validator: () => {
                const roles = appInfo.roles;
                if ((roles.length ?? 0) > 0) {
                  if (roles.filter((role) => role.isDefault).length != 1) {
                    return Promise.reject(new Error('No default role specified!'));
                  } else {
                    const appNames = roles.map((val) => val.name);
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
            value={appInfo.roles}
            urls={hasProxy(appInfo.grantType) ? appInfo.proxy.urls : []}
            onChange={(roles) => setAppInfo({ roles })}
          />
        </ProFormItem>
      ),
    },
    {
      name: 'user',
      title: intl.t('user.title', 'User'),
      style: { width: 600 },
      children: (
        <div
          style={{
            height: 'calc( 100vh - 400px )',
          }}
        >
          <GrantView
            users={appInfo.users}
            roles={appInfo.roles}
            onChange={(users) => {
              return setAppInfo({ users });
            }}
            granting={true}
            parentIntl={parentIntl}
          />
        </div>
      ),
    },
  ];

  const [currentStep, _setCurrentStep] = useState<number>(0);

  const setCurrentStep = (idx: number) => {
    const fixAllowStep = (newIdx: number, oriStep: number): number => {
      if (newIdx < 0) {
        return fixAllowStep(0, oriStep);
      }
      const step = steps[newIdx];
      if (
        step.allows &&
        step.allows.length > 0 &&
        step.allows.findIndex((allow) => appInfo.grantType.includes(allow)) < 0
      ) {
        return fixAllowStep(newIdx + (oriStep < newIdx ? 1 : -1), oriStep);
      }
      return newIdx;
    };
    _setCurrentStep((oriStep: number) => {
      const stepIdx = fixAllowStep(idx, oriStep);
      const step = steps[stepIdx];
      if (history.location.hash !== `#${step.name}`) {
        history.push({ hash: step.name });
      }
      return stepIdx;
    });
  };
  useEffect(() => {
    if (appInfo.grantType && appInfo.grantType.length > 0) {
      const { hash } = history.location;
      setCurrentStep(steps.map((step) => `#${step.name}`).findIndex((step) => step === hash));
    }
  }, [appInfo]);
  useEffect(() => {
    setAppInfo(
      {
        ...initialValues,
        oAuth2: {
          ...(initialValues?.oAuth2 ?? {
            authorizedRedirectUrl: [],
            jwtSignatureMethod: JWTSignatureMethod.default,
          }),
          jwtSignatureKey: '',
        },
      } ?? {},
    );
  }, [initialValues]);
  if (loading || !initialValues || initialValues?.id !== appInfo.id) {
    return <PageLoading />;
  }
  return (
    <StepsForm<FormValueType & { _roles: any }>
      stepsProps={{
        size: 'small',
      }}
      containerStyle={{ height: '100%' }}
      formProps={{
        preserve: false,
        disabled: disabled,
        loading: loading,
      }}
      stepsFormRender={(dom, submitter) => {
        return (
          <div>
            {dom}
            <Space
              style={{ position: 'absolute', bottom: 10, justifyContent: 'center', width: '100%' }}
            >
              {submitter}
            </Space>
          </div>
        );
      }}
      onCurrentChange={setCurrentStep}
      current={currentStep}
      onFinish={async () => {
        return onSubmit({
          ...appInfo,
          users: appInfo.users.map(({ id, roleId }) => ({ id, roleId })),
          status: appInfo.status !== AppStatus.unknown ? appInfo.status : AppStatus.normal,
          roles: appInfo.roles.map(({ id, name, isDefault, urls }) =>
            hasProxy(appInfo.grantType)
              ? { id, name, isDefault, urls }
              : { id, name, isDefault, urls: undefined },
          ),
          proxy: hasProxy(appInfo.grantType) ? appInfo.proxy : undefined,
        });
      }}
    >
      {steps.map(({ name, children, ...props }, idx) => {
        return (
          <StepsForm.StepForm
            key={`step-${name}`}
            initialValues={appInfo}
            syncToInitialValues={false}
            {...props}
          >
            {idx === currentStep ? children : null}
          </StepsForm.StepForm>
        );
      })}
    </StepsForm>
  );
};

export default CreateOrUpdateForm;
