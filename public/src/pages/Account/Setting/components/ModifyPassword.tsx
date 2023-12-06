import { Form, message } from 'antd';

import { resetPassword } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { ModalForm, ProFormText } from '@ant-design/pro-components';
import { GridContent } from '@ant-design/pro-components';

interface ModifyPasswordProps {
  parentIntl: IntlContext;
  currentUser?: API.UserInfo;
}

export default ({ parentIntl, currentUser }: ModifyPasswordProps) => {
  const intl = new IntlContext('modify-password', parentIntl);
  const [form] = Form.useForm<API.ResetUserPasswordRequest>();

  return (
    <ModalForm<API.ResetUserPasswordRequest & { newPassword2?: string }>
      title={intl.t('title', 'Modify Password')}
      trigger={<a>{parentIntl.t('modify', 'Modify')}</a>}
      form={form}
      autoFocusFirstInput
      modalProps={{
        destroyOnClose: true,
        maskClosable: false,
      }}
      width={600}
      onFinish={async ({ newPassword2, ...values }) => {
        if (newPassword2 != values.newPassword) {
          message.error('The passwords entered two times are inconsistent.');
        }
        const resp = await resetPassword({ ...values, username: currentUser?.username });
        if (resp.success) {
          message.info('Password updated successfully!');
          return true;
        }
        return false;
      }}
    >
      <GridContent>
        <ProFormText.Password
          rules={[
            {
              required: true,
            },
          ]}
          width="lg"
          name="oldPassword"
          label={intl.t('old-password', 'Old password')}
        />
        <ProFormText.Password
          rules={[
            {
              required: true,
            },
          ]}
          width="lg"
          name="newPassword"
          label={intl.t('new-password', 'New password')}
        />
        <ProFormText.Password
          rules={[
            {
              required: true,
            },
          ]}
          width="lg"
          name="newPassword2"
          label={intl.t('new-password', 'New password')}
        />
      </GridContent>
    </ModalForm>
  );
};
