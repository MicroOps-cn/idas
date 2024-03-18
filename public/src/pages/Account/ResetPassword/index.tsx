import { Button, message } from 'antd';
import 'antd/es/form/style/index.less';
import type { StoreValue } from 'antd/lib/form/interface';
import { parse } from 'query-string';
import React, { useState } from 'react';
import { history, useIntl, useModel } from 'umi';

import defaultSettings from '@/../config/defaultSettings';
import { loginPath, forgotPasswordPath } from '@/../config/env';
import Footer from '@/components/Footer';
import SelectLang from '@/components/SelectLang';
import { resetPassword } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { getPublicPath } from '@/utils/request';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { ProFormText, LoginForm } from '@ant-design/pro-form';
import { Link } from '@umijs/max';

import styles from './index.less';

const ResetPassword: React.FC = () => {
  const rootIntl = useIntl();
  const intl = new IntlContext('pages.resetPassword', rootIntl);
  const handleResetPassword = async (values: API.ResetUserPasswordRequest): Promise<boolean> => {
    try {
      // 登录
      const msg = await resetPassword({ ...values });
      if (msg.success) {
        const defaultSuccessMessage = intl.t(
          'success',
          'The password is reset successfully. Please login again with the new password.',
        );
        message.success(defaultSuccessMessage);
        return true;
      }
    } catch (error) {
      console.error(error);
    }
    return false;
  };
  const query = parse(history.location.search);
  const [loading, setLoading] = useState<boolean>(false);
  const { initialState } = useModel('@@initialState');
  const globalConfig = initialState?.globalConfig ?? null;
  return (
    <div className={styles.container}>
      <div className={styles.lang} data-lang>
        {SelectLang && <SelectLang />}
      </div>
      <div className={styles.content}>
        <LoginForm<API.ResetUserPasswordRequest>
          logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
          title={globalConfig?.title ?? defaultSettings.title}
          subTitle={<> </>}
          initialValues={{
            username: query.username,
            token: query.token,
          }}
          submitter={{
            render: (submitProps) => {
              return (
                <Button loading={loading} onClick={submitProps.submit} block type="primary">
                  {intl.t('resetPassword', 'Reset Password')}
                </Button>
              );
            },
          }}
          onFinish={async (values) => {
            setLoading(true);
            if (
              await handleResetPassword({
                newPassword: values.newPassword,
                oldPassword: values.oldPassword,
                userId: query.userId,
                token: query.token,
                username: query.username,
              })
            ) {
              history.push(loginPath);
            }
            setLoading(false);
          }}
        >
          <ProFormText
            fieldProps={{
              value: query.username,
              size: 'large',
              disabled: Boolean(query.username),
              prefix: <UserOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.t('username.placeholder', 'Please enter your username')}
            rules={[
              {
                required: true,
                message: intl.t('username.required', 'Please enter your username!'),
              },
            ]}
          />
          <ProFormText.Password
            name="oldPassword"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={styles.prefixIcon} />,
            }}
            hidden={Boolean(query.token)}
            placeholder={intl.t('oldPassword.placeholder', 'Please enter current password')}
            rules={[
              {
                required: !query.token,
                message: intl.t('oldPassword.required', 'Please enter current password!'),
              },
            ]}
          />
          <ProFormText.Password
            name="newPassword"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.t('password.placeholder', 'Please enter a new password')}
            rules={[
              {
                required: true,
                message: intl.t('password.required', 'Please enter a new password!'),
              },
            ]}
          />
          <ProFormText.Password
            name="newPasswordConfirm"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.t('confirmPassword.placeholder', 'Confirm new password.')}
            rules={[
              {
                required: true,
                message: intl.t('confirmPassword.required', 'Confirm new password!'),
              },
              ({ getFieldValue }) => ({
                validator: (_, value: StoreValue) => {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(
                    new Error('The two passwords that you entered do not match!'),
                  );
                },
              }),
            ]}
          />
          <Link
            style={{
              float: 'right',
            }}
            to={forgotPasswordPath}
          >
            {intl.t('forgotPassword', 'Forgot password')}
          </Link>
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default ResetPassword;
