import { Button, Input, Modal, Form, message } from 'antd';
import { RcFile } from 'antd/es/upload';
import { Component } from 'react';

import { AvatarUpload } from '@/components/Avatar';
import { uploadFile as postFile } from '@/services/idas/files';
import { updateCurrentUser } from '@/services/idas/user';
import type { IntlContext } from '@/utils/intl';

import styles from './BaseView.less';

type BasicUserInfo = Pick<
  API.UserInfo,
  'id' | 'username' | 'avatar' | 'fullName' | 'phoneNumber' | 'email'
>;

/**
 * @en-US Update user
 * @zh-CN 更新用户信息
 * @param fields
 */
const handleUpdate = async ({
  id,
  username,
  avatar,
  fullName,
  phoneNumber,
  email,
}: BasicUserInfo) => {
  const hide = message.loading('Updating ...');
  try {
    await updateCurrentUser({ id, username, avatar, fullName, phoneNumber, email });
    hide();
    message.success('Updated successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Update failed, please try again!');
    return false;
  }
};
type BaseViewProps = {
  currentUser?: API.UserInfo;
  parentIntl: IntlContext;
  reload: () => void;
};

class BaseView extends Component<BaseViewProps> {
  view: HTMLDivElement | undefined = undefined;

  getViewDom = (ref: HTMLDivElement) => {
    this.view = ref;
  };

  handleFinish = async (values: BasicUserInfo) => {
    const { currentUser, reload } = this.props;
    await handleUpdate({ ...currentUser, ...values });
    reload();
  };

  render() {
    const { currentUser } = this.props;

    return (
      <div className={styles.baseView} ref={this.getViewDom}>
        <Form<BasicUserInfo>
          layout="vertical"
          onFinish={this.handleFinish}
          initialValues={currentUser}
          requiredMark
          style={{ display: 'flex' }}
        >
          <div className={styles.left}>
            <Form.Item
              name="username"
              label={this.props.parentIntl.t('basic.username', 'Username')}
            >
              <Input disabled />
            </Form.Item>
            <Form.Item
              name="email"
              label={this.props.parentIntl.t('basic.email')}
              rules={[
                {
                  required: true,
                  message: this.props.parentIntl.t('basic.email-message'),
                },
              ]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="fullName" label={this.props.parentIntl.t('basic.fullname')}>
              <Input />
            </Form.Item>
            <Form.Item name="phoneNumber" label={this.props.parentIntl.t('basic.phone')}>
              <Input />
            </Form.Item>
            <Form.Item>
              <Button htmlType="submit" type="primary">
                {this.props.parentIntl.t('basic.update', 'Update Information')}
              </Button>
            </Form.Item>
          </div>
          <div className={styles.right}>
            <Form.Item name="avatar" label={this.props.parentIntl.t('basic.avatar', 'Avatar')}>
              <AvatarUpload />
            </Form.Item>
            <Modal open={false} />
          </div>
        </Form>
      </div>
    );
  }
}

export default BaseView;
