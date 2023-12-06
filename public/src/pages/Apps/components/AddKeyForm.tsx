import { Alert, Input, message, Modal, Tag } from 'antd';
import React, { useEffect, useState } from 'react';

import { createAppKey } from '@/services/idas/apps';
import { IntlContext } from '@/utils/intl';

interface AddKeyFormProps {
  visible: boolean;
  onClose: () => void;
  app?: API.AppInfo;
  parentIntl: IntlContext;
}

const AddKeyForm: React.FC<AddKeyFormProps> = ({ visible, onClose, app, parentIntl }) => {
  const intl = new IntlContext('keypair.form', parentIntl);
  const [createSuccessed, setCreateSuccessed] = useState<boolean>(false);
  const [appKey, setAppKey] = useState<API.AppKeyInfo | undefined>();
  const [name, setName] = useState<string>('');
  useEffect(() => {
    if (visible && app) {
      setCreateSuccessed(false);
      setAppKey(undefined);
      setName('');
    }
  }, [app, visible]);
  return (
    <>
      <Modal
        title={appKey ? intl.t('title.success', 'Success') : intl.t('title', 'Add App Key')}
        open={visible}
        onCancel={onClose}
        cancelButtonProps={{
          hidden: Boolean(appKey),
        }}
        onOk={async () => {
          if (app?.id && !createSuccessed) {
            if (!name) {
              message.error(intl.t('name.empty', 'please input name!'));
              return;
            }
            const resp = await createAppKey({ appId: app.id }, { name: name, appId: app.id });
            if (resp.success) {
              setAppKey(resp.data);
              setCreateSuccessed(true);
              return;
            }
            return;
          }
          onClose();
        }}
      >
        <Input
          hidden={createSuccessed}
          value={name}
          onChange={(e) => {
            setName(e.target.value);
          }}
          required={!createSuccessed}
          placeholder={intl.t('name.placeholder', 'Please Input name or description of app key')}
        />
        {appKey && (
          <div>
            <Alert
              description={intl.t(
                'used.tips',
                'If used for applications such as OAuth, please use Access Secret. If used for radius, please use Secret Hash.',
              )}
              type="info"
              showIcon
            />
            {intl.t('tips', 'App key pair, please keep it properly.')}:
            <ul style={{ padding: '0px' }}>
              <li>
                Access Key: <Tag>{appKey.key}</Tag>
              </li>
              <li>
                Access Secret: <Tag>{appKey.privateKey}</Tag>
              </li>
              <li>
                Secret Hash: <Tag>{appKey.secret}</Tag>
              </li>
            </ul>
          </div>
        )}
      </Modal>
    </>
  );
};
export default AddKeyForm;
