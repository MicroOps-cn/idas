import { List, Switch, Tooltip } from 'antd';
import type { PrimitiveType } from 'intl-messageformat';
import { Component } from 'react';

import { Input, InputArray } from '@/components/Form/Input';
import { getSecurityConfig, patchSecurityConfig } from '@/services/idas/config';
import { IntlContext } from '@/utils/intl';
import { QuestionCircleOutlined } from '@ant-design/icons';

import PasswordComplexitySelect, {
  getPasswordComplexityName,
  PasswordComplexityToolip,
} from './PasswordComplexitySelect';

type Unpacked<T> = T extends (infer U)[] ? U : T;

interface SecurityViewProps {
  parentIntl: IntlContext;
}

interface SecurityViewState {
  intl: IntlContext;
  setting?: API.RuntimeSecurityConfig;
}
class SecurityView extends Component<SecurityViewProps, SecurityViewState> {
  state: SecurityViewState = {
    intl: new IntlContext('security', this.props.parentIntl),
  };
  constructor(props: Readonly<SecurityViewProps>) {
    getSecurityConfig().then((resp) => {
      this.setState({ setting: resp.data });
    });
    super(props);
    this.state = {
      intl: new IntlContext('security', this.props.parentIntl),
    };
  }
  async handleUpdateSetting(setting: API.PatchSecurityConfigRequest): Promise<boolean> {
    const resp = await patchSecurityConfig(setting);
    getSecurityConfig().then((r) => {
      this.setState({ setting: r.data });
    });
    return resp.success;
  }
  t = (
    id: string,
    defaultMessage?: string,
    description?: string | object,
    values?: Record<string, PrimitiveType>,
  ): string => {
    const { intl } = this.state;
    return intl?.t(id, defaultMessage, description, values) ?? '';
  };
  getData = () => [
    {
      title: this.t('force-enable-mfa'),
      description: (
        <>
          {this.state.setting?.forceEnableMfa
            ? this.t('force-enable-mfa-enabled')
            : this.t('force-enable-mfa-disabled')}
        </>
      ),
      actions: [
        <Switch
          size="small"
          key={'force-enable-mfa'}
          checked={this.state.setting?.forceEnableMfa}
          onChange={(forceEnableMfa) => {
            this.handleUpdateSetting({ forceEnableMfa });
          }}
        />,
      ],
    },
    {
      title: (
        <>
          {this.t('password-complexity')}
          <PasswordComplexityToolip parentIntl={this.state.intl} />
        </>
      ),
      description: (
        <>
          {this.t(
            `password-complexity.option.${getPasswordComplexityName(
              this.state.setting?.passwordComplexity,
            )}-description`,
          )}
        </>
      ),
      actions: [
        <PasswordComplexitySelect
          key={'password-complexity'}
          parentIntl={this.state.intl}
          onSave={async (passwordComplexity) => {
            return this.handleUpdateSetting({ passwordComplexity });
          }}
          value={this.state.setting?.passwordComplexity}
        />,
      ],
    },
    {
      title: this.t('password-min-length'),
      description: (
        <>
          {this.state.setting?.passwordMinLength
            ? this.t('password-min-length-description', '', '', {
                minLen: this.state.setting.passwordMinLength,
              })
            : this.t('password-min-length-unrestricted')}
        </>
      ),
      actions: [
        <Input<number>
          key={'password-min-length'}
          type="number"
          intl={new IntlContext('password-min-length', this.state.intl)}
          onSave={async (passwordMinLength) => {
            return this.handleUpdateSetting({ passwordMinLength });
          }}
          value={this.state.setting?.passwordMinLength}
        />,
      ],
    },
    {
      title: this.t('password-expire-time'),
      description: (
        <>
          {this.state.setting?.passwordExpireTime
            ? this.t('password-expire-time-description', '', '', {
                days: this.state.setting?.passwordExpireTime,
              })
            : this.t('password-expire-time-unrestricted')}
        </>
      ),
      actions: [
        <Input<number>
          key={'password-expire-time'}
          type="number"
          intl={new IntlContext('password-expire-time', this.state.intl)}
          onSave={async (passwordExpireTime) => {
            return this.handleUpdateSetting({ passwordExpireTime });
          }}
          style={{ width: 100 }}
          suffix={'day'}
          value={this.state.setting?.passwordExpireTime}
        />,
      ],
    },
    {
      title: this.t('password-failed-lock'),
      description: (
        <>
          {this.state.setting?.passwordFailedLockDuration &&
          this.state.setting?.passwordFailedLockThreshold
            ? this.t('password-failed-lock-description', '', '', {
                min: this.state.setting?.passwordFailedLockDuration,
                fails: this.state.setting?.passwordFailedLockThreshold,
              })
            : this.t('password-failed-lock-unrestricted')}
        </>
      ),
      actions: [
        <InputArray<number>
          key={'password-failed-lock'}
          type="number"
          intl={new IntlContext('password-failed-lock', this.state.intl)}
          onSave={async (passwordFailedLock) => {
            const [passwordFailedLockThreshold, passwordFailedLockDuration] = passwordFailedLock;
            return this.handleUpdateSetting({
              passwordFailedLockThreshold,
              passwordFailedLockDuration,
            });
          }}
          suffix={['failed', 'minute']}
          style={[{ width: 120 }, { width: 100 }]}
          tooltip={[
            'Number of consecutive password input errors.',
            'The duration of account lockout.',
          ]}
          count={2}
          value={[
            this.state.setting?.passwordFailedLockThreshold ?? 0,
            this.state.setting?.passwordFailedLockDuration ?? 0,
          ]}
        />,
      ],
    },
    {
      title: this.t('password-history'),
      description: (
        <>
          {this.state.setting?.passwordHistory
            ? this.t('password-history-description', '', '', {
                count: this.state.setting?.passwordHistory,
              })
            : this.t('password-history-unrestricted')}
        </>
      ),
      actions: [
        <Input<number>
          key={'password-history'}
          type="number"
          style={{ width: 100 }}
          suffix={'day'}
          intl={new IntlContext('password-history', this.state.intl)}
          onSave={async (passwordHistory) => {
            return this.handleUpdateSetting({ passwordHistory });
          }}
          value={this.state.setting?.passwordHistory}
        />,
      ],
    },
    {
      title: this.t('account-inactive-lock'),
      description: (
        <>
          {this.state.setting?.accountInactiveLock
            ? this.t('account-inactive-lock-description', '', '', {
                days: this.state.setting?.accountInactiveLock,
              })
            : this.t('account-inactive-lock-unrestricted')}
        </>
      ),
      actions: [
        <Input<number>
          key={'account-inactive-lock'}
          type="number"
          style={{ width: 100 }}
          suffix={'day'}
          intl={new IntlContext('account-inactive-lock', this.state.intl)}
          onSave={async (accountInactiveLock) => {
            return this.handleUpdateSetting({ accountInactiveLock });
          }}
          value={this.state.setting?.accountInactiveLock}
        />,
      ],
    },
    {
      title: (
        <>
          {this.t('login-session-expiration-time')}
          <Tooltip
            placement="bottomLeft"
            overlayStyle={{ maxWidth: 'max-content' }}
            overlayInnerStyle={{ backgroundColor: 'rgba(61, 62, 64, 0.85)' }}
            title={this.t('login-session-expiration-tooltip')}
          >
            <QuestionCircleOutlined
              style={{ color: 'rgba(61, 62, 64, 0.45)', marginInlineStart: 4 }}
            />
          </Tooltip>
        </>
      ),
      description: (
        <>
          {this.state.setting?.loginSessionInactivityTime
            ? this.t(
                'login-session-expiration-time-description',
                'Automatically log out after {loginSessionInactivityHours} hours of inactivity, with a maximum session duration of {loginSessionMaxHours} hours.',
                '',
                {
                  loginSessionInactivityHours: this.state.setting?.loginSessionInactivityTime,
                  loginSessionMaxHours: this.state.setting?.loginSessionMaxTime,
                },
              )
            : this.t('login-session-expiration-time-unrestricted')}
        </>
      ),
      actions: [
        <InputArray<number>
          key={'login-session-expiration-time'}
          type="number"
          intl={new IntlContext('login-session-expiration-time', this.state.intl)}
          onSave={async (passwordFailedLock) => {
            const [loginSessionInactivityTime, loginSessionMaxTime] = passwordFailedLock;
            return this.handleUpdateSetting({
              loginSessionInactivityTime,
              loginSessionMaxTime,
            });
          }}
          suffix={['hours', 'hours']}
          style={[{ width: 100 }, { width: 100 }]}
          tooltip={[
            'Session inactive automatic logout time.',
            'The maximum duration of the session.',
          ]}
          count={2}
          value={[
            this.state.setting?.loginSessionInactivityTime ?? 0,
            this.state.setting?.loginSessionMaxTime ?? 0,
          ]}
        />,
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
