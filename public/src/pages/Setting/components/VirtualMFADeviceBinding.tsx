import { Button, Empty, Form } from 'antd';
import { QRCodeCanvas } from 'qrcode.react';
import { useState } from 'react';
import { useIntl } from 'umi';

import { bindingTOTP, getTOTPSecret } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { LoadingOutlined } from '@ant-design/icons';
import { ModalForm, ProFormText } from '@ant-design/pro-components';
import { GridContent } from '@ant-design/pro-components';

interface VirtualMFABindingProps {
  parentIntl?: IntlContext;
  currentUser?: API.UserInfo;
}

export default ({
  parentIntl = new IntlContext('user.settings.security', useIntl()),
}: VirtualMFABindingProps) => {
  const intl = new IntlContext('mfa-device', parentIntl);
  const [form] = Form.useForm<API.CreateTOTPRequest>();
  const [secret, setSecret] = useState<API.CreateTOTPSecretResponseData>();
  const [tokenLoading, setTokenLoading] = useState<boolean>(false);
  const fetchToken = async () => {
    setSecret(undefined);
    setTokenLoading(true);
    getTOTPSecret({})
      .then((resp) => {
        setSecret(resp.data);
      })
      .finally(() => {
        setTokenLoading(false);
      });
  };

  return (
    <ModalForm<API.CreateTOTPRequest>
      title={intl.t('bind', 'Bind MFA')}
      trigger={<a>{parentIntl.t('bind', 'Bind')}</a>}
      form={form}
      autoFocusFirstInput
      modalProps={{
        destroyOnClose: true,
        maskClosable: false,
      }}
      onOpenChange={(visible) => {
        if (visible) fetchToken();
      }}
      submitter={{
        submitButtonProps: {
          disabled: !secret,
        },
      }}
      onFinish={async (values) => {
        if (secret?.token) {
          await bindingTOTP({ ...values, token: secret.token });
          return true;
        } else {
          throw new Error('System error: token is null');
        }
      }}
    >
      <GridContent>
        <p>
          {intl.t(
            'code-description',
            'Please obtain two consecutive one-time passwords after scanning and adding MFA and enter them into the input box below.',
          )}
        </p>
        <div style={{ display: 'flex' }}>
          <div>
            <div style={{ width: 220, height: 140 }}>
              {secret?.secret ? (
                <QRCodeCanvas style={{ marginLeft: '40px' }} value={secret?.secret} />
              ) : (
                <Empty
                  description={''}
                  style={{ width: 200, height: 128 }}
                  image={tokenLoading ? <LoadingOutlined /> : undefined}
                />
              )}
            </div>
            <div style={{ display: 'grid', width: '100%' }}>
              <Button onClick={fetchToken} type="link">
                {intl.t('refresh', 'Refresh')}
              </Button>
            </div>
          </div>
          <div style={{ marginLeft: '30px' }}>
            <ProFormText
              rules={[
                {
                  pattern: /[0-9]{6}/,
                  required: true,
                },
              ]}
              width="sm"
              name="firstCode"
              label={intl.t('first-code', 'First code')}
            />
            <ProFormText
              rules={[
                {
                  pattern: /[0-9]{6}/,
                  required: true,
                },
              ]}
              width="sm"
              name="secondCode"
              label={intl.t('second-code', 'Second code')}
            />
          </div>
        </div>
      </GridContent>
    </ModalForm>
  );
};
