import { Alert, Button, Empty } from 'antd';
import { QRCodeCanvas } from 'qrcode.react';
import { useEffect, useState } from 'react';
import { useIntl } from 'umi';

import { getTotpSecret as getTOTPSecret } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { LoadingOutlined } from '@ant-design/icons';
import { ProFormText } from '@ant-design/pro-components';
import { GridContent } from '@ant-design/pro-components';

interface VirtualMFABindingProps {
  parentIntl?: IntlContext;
  token?: string;
}

export default ({
  parentIntl = new IntlContext('user.settings.security', useIntl()),
  token: fetchToken,
}: VirtualMFABindingProps) => {
  const intl = new IntlContext('mfa', parentIntl);
  const [secret, setSecret] = useState<API.CreateTOTPSecretResponseData>();
  const [tokenLoading, setTokenLoading] = useState<boolean>(false);
  const fetchTOTPToken = async (token: string) => {
    setSecret(undefined);
    setTokenLoading(true);
    getTOTPSecret({ token })
      .then((resp) => {
        setSecret(resp.data);
      })
      .finally(() => {
        setTokenLoading(false);
      });
  };
  useEffect(() => {
    if (fetchToken) fetchTOTPToken(fetchToken);
  }, [fetchToken]);
  return (
    <GridContent>
      <span>
        <Alert
          style={{ marginBottom: 24 }}
          message={intl.t(
            'code-description',
            'Please obtain two consecutive one-time passwords after scanning and adding MFA and enter them into the input box below.',
          )}
          type="info"
          showIcon
        />
      </span>
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
            <Button
              onClick={() => {
                if (fetchToken) fetchTOTPToken(fetchToken);
              }}
              type="link"
            >
              {intl.t('refresh', 'Refresh')}
            </Button>
          </div>
        </div>
        <div style={{ marginLeft: '50px' }}>
          <ProFormText
            hidden
            width="sm"
            name="token"
            fieldProps={{
              value: secret?.token,
            }}
          />
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
  );
};
