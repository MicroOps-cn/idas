import { Alert, Button, Divider, message, Space, Tabs } from 'antd';
import 'antd/es/form/style/index.less';
import { useForm } from 'antd/lib/form/Form';
import classNames from 'classnames';
import { isArray, isNumber } from 'lodash';
import { parse } from 'query-string';
import React, { useEffect, useState } from 'react';
import { useIntl, history, useModel, Link } from 'umi';

import { forgotPasswordPath } from '@/../config/env';
import Footer from '@/components/Footer';
import SelectLang from '@/components/SelectLang';
import type { LoginTypeName, LoginTypeValue } from '@/services/idas/enums';
import { LoginType } from '@/services/idas/enums';
import { sendLoginCaptcha, userLogin as login } from '@/services/idas/user';
import type { LabelValue } from '@/utils/enum';
import { enumToOptions } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import { getApiPath, getPublicPath } from '@/utils/request';
import { LockOutlined, MailOutlined, MobileOutlined, UserOutlined } from '@ant-design/icons';
import { ProFormCaptcha, ProFormCheckbox, ProFormText, LoginForm } from '@ant-design/pro-form';

import VirtualMFADeviceBinding from './components/VirtualMFADeviceBinding';
import styles from './index.less';

const LoginMessage: React.FC<{
  content: string;
  hidden: boolean;
}> = ({ content, hidden }) =>
  !hidden ? (
    <Alert
      style={{
        marginBottom: 24,
      }}
      message={content}
      type="error"
      showIcon
    />
  ) : null;

interface LoginFormComponentProps
  extends React.DetailedHTMLProps<React.HTMLAttributes<HTMLSpanElement>, HTMLSpanElement> {
  loginType: LoginTypeName;
  allows: LoginTypeName[];
  hidden?: boolean;
  children: React.ReactNode;
}

const LoginFormComponent: React.FC<LoginFormComponentProps> = ({
  children,
  allows,
  loginType,
  hidden,
  ...props
}) => {
  const allowTypes = new Set([...allows, allows.map((r) => LoginType[r])]);

  if (allowTypes.has(loginType)) {
    return (
      <span hidden={hidden} {...props}>
        {children}
      </span>
    );
  }
  return <></>;
};

const Login: React.FC = ({}) => {
  const { redirect_uri, redirect } = parse(history.location.search);
  const [userLoginState, setUserLoginState] = useState<Omit<API.UserLoginResponse, 'traceId'>>({
    success: true,
  });
  const [loginType, setLoginType] = useState<LoginTypeName>('normal');
  const [token, setToken] = useState<string>();
  const [email, setEmail] = useState<string>();
  const [bindingToken, setBindingToken] = useState<string>();
  const { initialState, setInitialState } = useModel('@@initialState');
  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */
  const intl = new IntlContext('pages.login', useIntl());
  const [form] = useForm();
  const initLoginTypes: LoginTypeName[] = ['normal', 'email', 'sms'];
  const mfaLoginTypes: LoginTypeName[] = ['mfa_totp', 'mfa_email', 'mfa_sms'];

  const [allowLoginTypes, setAllowLoginTypes] = useState<API.LoginType[]>(
    initLoginTypes.map((val) => LoginType[val]),
  );
  const [otherLoginTypes, setOtherLoginTypes] = useState<API.LoginType[]>(
    initLoginTypes.map((val) => LoginType[val]),
  );
  const [oauthLoginTypes, setOAuthLoginTypes] = useState<API.GlobalLoginType[]>([]);
  const handleSetLoginType = (type: LoginType) => {
    if (isNumber(type)) {
      setLoginType(LoginType[type] as LoginTypeName);
    } else {
      setLoginType(type);
    }
  };
  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };

  const [hiddenNormal, setHiddenNormal] = useState<boolean>(false);
  const handleSubmit = async ({
    ...values
  }: API.UserLoginRequest & { phone?: string; email?: string }) => {
    try {
      const msg = await login(
        { ...values, type: loginType, token, bindingToken },
        { skipErrorHandler: true, ignoreError: true },
      );
      if (msg.data?.nextMethod && msg.data.nextMethod.length > 0) {
        setAllowLoginTypes(msg.data.nextMethod);
        handleSetLoginType(msg.data.nextMethod[0] as LoginType);
        setHiddenNormal(true);
        setToken(msg.data.token);
        setEmail(msg.data.email);
      } else if (msg.success) {
        const defaultLoginSuccessMessage = intl.t('success', 'Login succeeded!');
        message.success(defaultLoginSuccessMessage);
        await fetchUserInfo();
        /** 此方法会跳转到 redirect 参数所在的位置 */
        if (redirect_uri) {
          window.location.href = isArray(redirect_uri) ? redirect_uri[0] : redirect_uri;
          return;
        }
        history.push(redirect || '/');
        return;
      } else if (msg.errorCode === 'E0004') {
        message.warning(intl.t(`${msg.errorCode ?? 'normal'}.errorMessage`, msg.errorMessage));
        history.push(`/account/resetPassword?username=${values.username}`);
      } else if (msg.errorCode && ['E0002', 'E0005', 'E0006'].includes(msg.errorCode)) {
        message.warning(intl.t(`${msg.errorCode ?? 'normal'}.errorMessage`, msg.errorMessage));
      }
      // 如果失败去设置用户错误信息
      setUserLoginState(msg);
    } catch (error) {
      setUserLoginState({ success: false });
    }
  };
  const { success, errorMessage, errorCode } = userLoginState;

  const [tabsItem, setTabsItem] = useState<LabelValue[]>(
    enumToOptions(LoginType, intl, 'loginType', (item) =>
      initLoginTypes.includes(item as LoginTypeName),
    ),
  );
  const globalConfig = initialState?.globalConfig ?? null;
  useEffect(() => {
    if (globalConfig) {
      for (let index = 0; index < globalConfig.loginType.length; index++) {
        const t = globalConfig.loginType[index];
        if (t.autoRedirect) {
          let gotoURI = getApiPath(`/api/v1/user/oauth/${t.id}`);
          if (redirect_uri) {
            gotoURI = `${gotoURI}?redirect_uri=${encodeURIComponent(
              isArray(redirect_uri) ? redirect_uri[0] : redirect_uri,
            )}`;
          }
          window.location.href = gotoURI;
          return;
        }
      }
      if (globalConfig.defaultLoginType === LoginType.oauth2) {
        setAllowLoginTypes([]);
        setOtherLoginTypes(
          globalConfig.loginType
            .map((item) => item.type)
            .filter((item) => item !== LoginType.oauth2 && item !== undefined) as LoginTypeValue[],
        );
      } else {
        setAllowLoginTypes(
          globalConfig.loginType
            .map((item) => item.type)
            .filter((item) => item !== LoginType.oauth2 && item !== undefined) as LoginTypeValue[],
        );
      }

      setOAuthLoginTypes(globalConfig.loginType.filter((item) => item.type === LoginType.oauth2));
    }
  }, [globalConfig, redirect_uri]);
  useEffect(() => {
    setTabsItem(
      enumToOptions(
        LoginType,
        intl,
        'loginType',
        (item, val) =>
          allowLoginTypes.includes(item as LoginTypeName) ||
          allowLoginTypes.includes(val as LoginTypeValue),
      ),
    );
    if (allowLoginTypes.length > 0) {
      if (!allowLoginTypes.includes(loginType)) {
        setLoginType(LoginType[allowLoginTypes[0]] as LoginTypeName);
      }
      setHiddenNormal(false);
    } else {
      setHiddenNormal(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [allowLoginTypes]);
  return (
    <div className={styles.container}>
      <div className={styles.lang} data-lang>
        {SelectLang && <SelectLang />}
      </div>
      <div className={styles.content}>
        <div className="login-container">
          <LoginForm<API.UserLoginRequest>
            logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
            title={globalConfig?.title ?? 'IDAS'}
            subTitle={globalConfig?.subTitle ?? 'Identity authentication service'}
            initialValues={{
              autoLogin: true,
            }}
            contentStyle={{ width: 'unset' }}
            form={form}
            actions={
              oauthLoginTypes.length > 0 ? (
                <div className="ant-pro-form-login-page-other">
                  <div key={'loginWith'}>
                    <Divider plain>{intl.t('loginWith', 'Login with')}</Divider>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      justifyContent: 'center',
                      alignItems: 'center',
                      flexDirection: 'column',
                    }}
                  >
                    <Space align="center" size={24} key="loginMethod">
                      {oauthLoginTypes.map((item) => (
                        <div
                          key={item.name}
                          className={classNames(styles.oauthIconButton)}
                          title={intl.t('signInWith', 'Sign in with {name}', '', {
                            name: item.name,
                          })}
                          onClick={() => {
                            let gotoURI = getApiPath(`/api/v1/user/oauth/${item.id}`);
                            if (redirect_uri) {
                              gotoURI = `${gotoURI}?redirect_uri=${encodeURIComponent(
                                isArray(redirect_uri) ? redirect_uri[0] : redirect_uri,
                              )}`;
                            }
                            window.location.href = gotoURI;
                          }}
                        >
                          <img src={item.icon} className={styles.oauthIcon} />
                        </div>
                      ))}
                    </Space>
                  </div>
                  <div className={styles.loginByOtherBtn}>
                    <Button
                      style={{ display: otherLoginTypes.length === 0 ? 'none' : 'unset' }}
                      type="link"
                      onClick={() => {
                        if (globalConfig) {
                          setAllowLoginTypes(
                            globalConfig.loginType
                              .map((item) => item.type)
                              .filter(
                                (item) => item !== LoginType.oauth2 && item !== undefined,
                              ) as LoginTypeValue[],
                          );
                        }
                      }}
                    >
                      {intl.t('loginByOther', 'More login methods')}
                    </Button>
                  </div>
                </div>
              ) : undefined
            }
            submitter={allowLoginTypes.length > 0 ? undefined : false}
            onFinish={async (values) => {
              return handleSubmit(values);
            }}
          >
            <Tabs
              activeKey={loginType}
              onChange={(key) => {
                setLoginType(key as LoginTypeName);
                setUserLoginState({ success: true });
              }}
              className={styles.loginTypeTabs}
              items={tabsItem}
              centered
            />
            <LoginFormComponent
              hidden={hiddenNormal}
              loginType={loginType}
              allows={['normal', ...mfaLoginTypes]}
            >
              <LoginMessage
                hidden={success}
                content={intl.t(`${errorCode ?? 'normal'}.errorMessage`, errorMessage)}
              />
              <ProFormText
                name="username"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined className={styles.prefixIcon} />,
                  autoFocus: true,
                }}
                placeholder={intl.t('username.placeholder', 'Please enter a username')}
                rules={[
                  {
                    required: true,
                    message: intl.t('username.required', 'Please enter a username!'),
                  },
                ]}
              />
              <ProFormText.Password
                name="password"
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined className={styles.prefixIcon} />,
                }}
                placeholder={intl.t('password.placeholder', 'Please input a password')}
                rules={[
                  {
                    required: true,
                    message: intl.t('password.required', 'Please input a password!'),
                  },
                ]}
              />
            </LoginFormComponent>

            <LoginFormComponent loginType={loginType} allows={['mfa_totp']}>
              <LoginMessage
                hidden={success}
                content={intl.t('totp.errorMessage', 'verification code error')}
              />
            </LoginFormComponent>
            <LoginFormComponent loginType={loginType} allows={['email', 'mfa_email']}>
              <LoginMessage
                hidden={success}
                content={intl.t('email.errorMessage', 'Email verification code error')}
              />
              <ProFormText
                fieldProps={{
                  size: 'large',
                  prefix: <MailOutlined className={styles.prefixIcon} />,
                }}
                name="email"
                placeholder={`${intl.t('email.placeholder', 'Please enter your email')} ${
                  email ? `: ${email}` : ''
                }`}
                rules={[
                  {
                    required: true,
                    message: intl.t('email.required', 'Please enter your email!'),
                  },
                ]}
              />
            </LoginFormComponent>
            <LoginFormComponent loginType={loginType} allows={['sms', 'mfa_sms']}>
              <LoginMessage
                hidden={success}
                content={intl.t('sms.errorMessage', 'SMS verification code error')}
              />
              <ProFormText
                fieldProps={{
                  size: 'large',
                  prefix: <MobileOutlined className={styles.prefixIcon} />,
                }}
                name="phone"
                placeholder={intl.t('phoneNumber.placeholder', 'Please enter your phone number')}
                rules={[
                  {
                    required: true,
                    message: intl.t('phoneNumber.required', 'Please enter your phone number!'),
                  },
                  {
                    pattern: /^1\d{10}$/,
                    message: intl.t('phoneNumber.invalid', 'Mobile phone number format error!'),
                  },
                ]}
              />
            </LoginFormComponent>

            <LoginFormComponent
              loginType={loginType}
              allows={['enable_mfa_email', 'enable_mfa_sms', 'enable_mfa_totp']}
              style={{ width: 550, display: 'block' }}
            >
              <Alert
                style={{ marginBottom: 24 }}
                message={intl.t(
                  'mfa.errorMessage',
                  'Because your user is set to enable multiple factor authentication (MFA), you need to enable at least one MFA authentication method.',
                )}
                type="info"
                showIcon
              />
            </LoginFormComponent>

            <LoginFormComponent
              loginType={loginType}
              allows={['enable_mfa_sms']}
              style={{ width: 550, display: 'block' }}
            >
              <Alert
                style={{ marginBottom: 24 }}
                message={intl.t(
                  'enableMfa.smsMessage',
                  'Click Login to automatically enable email as the second authentication factor.',
                )}
                type="info"
                showIcon
              />
            </LoginFormComponent>

            <LoginFormComponent
              loginType={loginType}
              allows={['enable_mfa_email']}
              style={{ width: 550, display: 'block' }}
            >
              <Alert
                style={{ marginBottom: 24 }}
                message={intl.t(
                  'enableMfa.emailMessage',
                  'Click Login to automatically enable email as the second authentication factor.',
                )}
                type="info"
                showIcon
              />
              <LoginMessage
                hidden={success}
                content={intl.t('email.errorMessage', 'Email verification code error')}
              />
              <ProFormText
                fieldProps={{
                  size: 'large',
                  prefix: <MailOutlined className={styles.prefixIcon} />,
                }}
                name="email"
                initialValue={email}
                placeholder={intl.t('email.placeholder', 'Please enter your email')}
                rules={[
                  {
                    required: true,
                    message: intl.t('email.required', 'Please enter your email!'),
                  },
                ]}
              />
            </LoginFormComponent>

            <LoginFormComponent
              loginType={loginType}
              allows={[
                'sms',
                'mfa_sms',
                'mfa_email',
                'email',
                'mfa_totp',
                'enable_mfa_email',
                'enable_mfa_sms',
              ]}
            >
              <ProFormCaptcha
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined className={styles.prefixIcon} />,
                }}
                captchaProps={{
                  size: 'large',
                  hidden: loginType === 'mfa_totp',
                }}
                placeholder={intl.t('captcha.placeholder', 'Please enter the verification code')}
                captchaTextRender={(timing, count) => {
                  if (timing) {
                    return `${count} ${intl.t('getCaptchaSecondText', 'Get verification code')}`;
                  }
                  return intl.t('phone.getVerificationCode', 'Get verification code');
                }}
                name="code"
                rules={[
                  {
                    required: true,
                    message: intl.t('captcha.required', 'Please enter the verification code!'),
                  },
                ]}
                onGetCaptcha={async () => {
                  const req: API.SendLoginCaptchaRequest = { type: loginType };
                  switch (loginType) {
                    case 'enable_mfa_email':
                    case 'mfa_email':
                      form.validateFields(['username']);
                      req.username = form.getFieldValue('username');
                    case 'email':
                      form.validateFields(['email']);
                      req.email = form.getFieldValue('email');
                      break;
                    case 'enable_mfa_sms':
                    case 'mfa_sms':
                      form.validateFields(['username']);
                      req.username = form.getFieldValue('username');
                    case 'sms':
                      form.validateFields(['phone']);
                      req.phone = form.getFieldValue('phone');
                    default:
                      break;
                  }
                  const result = await sendLoginCaptcha(req, { intl });
                  if (!result.success) {
                    return;
                  }
                  console.log(loginType);
                  switch (loginType) {
                    case 'enable_mfa_email':
                    case 'enable_mfa_sms':
                      setBindingToken(result.data?.token);
                      break;
                    default:
                      setToken(result.data?.token);
                      break;
                  }
                  message.success(intl.t('captcha.sent', 'Verification code sent successfully.'));
                }}
              />
            </LoginFormComponent>

            <LoginFormComponent
              style={{ width: 550, display: 'block' }}
              loginType={loginType}
              allows={['enable_mfa_totp']}
            >
              <VirtualMFADeviceBinding token={token} setBindingToken={setBindingToken} />
            </LoginFormComponent>
            <div
              style={{
                marginBottom: 10,
                display: allowLoginTypes.length > 0 ? 'block' : 'none',
              }}
            >
              <ProFormCheckbox noStyle name="autoLogin">
                {intl.t('rememberMe', 'Automatic login')}
              </ProFormCheckbox>
              <LoginFormComponent loginType={loginType} allows={['normal']}>
                <Link
                  style={{
                    float: 'right',
                  }}
                  to={forgotPasswordPath}
                >
                  {intl.t('forgotPassword', 'Forgot password')}
                </Link>
              </LoginFormComponent>
            </div>
          </LoginForm>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Login;
