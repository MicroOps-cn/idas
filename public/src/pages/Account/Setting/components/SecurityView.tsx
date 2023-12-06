import { List, message, Modal } from 'antd';
import { Component } from 'react';
import { Link } from 'umi';

import { forgotPasswordPath } from '@/../config/env';
import Switch from '@/components/Switch';
import { patchCurrentUser } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { ExclamationCircleOutlined } from '@ant-design/icons';

import ModifyPassword from './ModifyPassword';
import VirtualMFADeviceBinding from './VirtualMFADeviceBinding';

type Unpacked<T> = T extends (infer U)[] ? U : T;

interface SecurityViewProps {
  parentIntl: IntlContext;
  currentUser?: API.UserInfo;
  reload: () => void;
}

interface SecurityViewState {
  intl: IntlContext;
}
class SecurityView extends Component<SecurityViewProps, SecurityViewState> {
  state: SecurityViewState = {
    intl: new IntlContext('security', this.props.parentIntl),
  };
  constructor(props: Readonly<SecurityViewProps>) {
    super(props);
    this.state = {
      intl: new IntlContext('security', this.props.parentIntl),
    };
  }
  unbindMFADevice(props: SecurityViewProps): () => Promise<boolean> {
    return async () => {
      const hide = message.loading('Unbinding ...');
      try {
        await patchCurrentUser({ totp_as_mfa: false });
        hide();
        message.success('Unbind successfully and will refresh soon');
        props.reload();
        return true;
      } catch (error) {
        console.log(error);
        hide();
        message.error('Unbind failed, please try again');
        return false;
      }
    };
  }
  t = (id: string, defaultMessage?: string): string => {
    const { intl } = this.state;
    return intl?.t(id, defaultMessage) ?? '';
  };
  passwordStrength = {
    strong: <span className="strong">{this.t('strong', 'Strong')}</span>,
    medium: <span className="medium">{this.t('medium', 'Medium')}</span>,
    weak: <span className="weak">{this.t('weak', 'Weak')}</span>,
  };
  getData = () => [
    {
      title: this.t('password'),
      actions: [
        <Link key="reset-password" to={forgotPasswordPath}>
          {this.t('resetPassword', 'Reset')}
        </Link>,
        <ModifyPassword
          key={'modify-password'}
          parentIntl={this.state.intl}
          currentUser={this.props.currentUser}
        />,
      ],
    },
    {
      title: this.t('mfa-phone'),
      description: `${this.t('mfa-phone-description')}: ${
        this.props.currentUser?.phoneNumber ?? ''
      }`,
      actions: [
        <Switch
          request={async () => {
            await patchCurrentUser({
              sms_as_mfa: !this.props.currentUser?.extendedData?.smsAsMFA,
            });
            this.props.reload();
          }}
          checked={this.props.currentUser?.extendedData?.smsAsMFA}
          key={'phone'}
          onChange={() => {}}
          size="small"
        />,
      ],
    },
    {
      title: this.t('mfa-email'),
      description: `${this.t('mfa-email-description')}: ${this.props.currentUser?.email ?? ''}`,
      actions: [
        <Switch
          checked={this.props.currentUser?.extendedData?.emailAsMFA}
          request={async () => {
            await patchCurrentUser({
              email_as_mfa: !this.props.currentUser?.extendedData?.emailAsMFA,
            });
            this.props.reload();
          }}
          key={'email'}
          size="small"
        />,
      ],
    },
    {
      title: this.t('mfa-device'),
      description: this.props.currentUser?.extendedData?.totpAsMFA
        ? this.t('mfa-device-description-bound')
        : this.t('mfa-device-description-unbound'),
      actions: [
        <VirtualMFADeviceBinding
          key="mfa-binding"
          parentIntl={this.state.intl}
          currentUser={this.props.currentUser}
          hidden={this.props.currentUser?.extendedData?.totpAsMFA}
          reload={this.props.reload}
        />,
        <a
          key="mfa-unbind"
          hidden={!this.props.currentUser?.extendedData?.totpAsMFA}
          onClick={() => {
            Modal.confirm({
              title: this.t('unbind-confirm', 'Are you sure you want to unbind the mfa device?'),
              icon: <ExclamationCircleOutlined />,
              onOk: this.unbindMFADevice(this.props),
            });
          }}
        >
          {this.t('unbind', 'Unbind')}
        </a>,
      ],
    },
  ];

  render() {
    const data = this.getData();
    return (
      <>
        <List<Unpacked<typeof data>>
          itemLayout="horizontal"
          dataSource={data}
          renderItem={(item) => (
            <List.Item actions={item.actions}>
              <List.Item.Meta title={item.title} description={item.description} />
            </List.Item>
          )}
        />
      </>
    );
  }
}

export default SecurityView;
