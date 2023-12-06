import { Button, message } from 'antd';
import 'antd/es/form/style/index.less';
import React, { useState } from 'react';
import { useIntl, useModel } from 'umi';

import defaultSettings from '@/../config/defaultSettings';
import Footer from '@/components/Footer';
import SelectLang from '@/components/SelectLang';
import { forgotPassword } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { getPublicPath } from '@/utils/request';
import { MailOutlined, UserOutlined } from '@ant-design/icons';
import { ProFormText, LoginForm } from '@ant-design/pro-form';

import styles from './index.less';

const ForgotPassword: React.FC = () => {
  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */
  const intl = new IntlContext('pages.forgotPassword', useIntl());
  const handleForgotPassword = async (values: API.ForgotUserPasswordRequest): Promise<boolean> => {
    try {
      const msg = await forgotPassword({ ...values });
      if (msg.success) {
        const defaultSuccessMessage = intl.t(
          'success',
          'The email was sent successfully, please check it.',
        );
        message.success(defaultSuccessMessage);
        return true;
      }
    } catch (error) {
      console.error(error);
    }
    return false;
  };
  const [loading, setLoading] = useState<boolean>(false);
  const { initialState } = useModel('@@initialState');
  const globalConfig = initialState?.globalConfig ?? null;
  return (
    <div className={styles.container}>
      <div className={styles.lang} data-lang>
        {SelectLang && <SelectLang />}
      </div>
      <div className={styles.content}>
        <LoginForm<API.ForgotUserPasswordRequest>
          logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
          title={globalConfig?.title ?? defaultSettings.title}
          subTitle={<> </>}
          initialValues={{
            autoLogin: true,
          }}
          submitter={{
            render: (props) => {
              return (
                <Button loading={loading} onClick={props.submit} block type="primary">
                  {intl.t('button.submit', 'Submit!')}
                </Button>
              );
            },
          }}
          onFinish={async (values) => {
            setLoading(true);
            await handleForgotPassword(values);
            setLoading(false);
          }}
        >
          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              disabled: loading,
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
          <ProFormText
            name="email"
            fieldProps={{
              size: 'large',
              disabled: loading,
              prefix: <MailOutlined className={styles.prefixIcon} />,
            }}
            placeholder={intl.t('email.placeholder', 'Please enter email address')}
            rules={[
              {
                required: true,
                message: intl.t('email.required', 'Please enter email address!'),
              },
            ]}
          />
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default ForgotPassword;
