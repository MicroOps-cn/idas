export default {
  'settings.menuMap.basic': '基本设置',
  'settings.menuMap.security': '安全设置',
  'settings.security.force-enable-mfa': '强制启用多因子认证（MFA）',
  'settings.security.force-enable-mfa-description': '当前状态',
  'settings.security.force-enable-mfa-disabled': '未启用',
  'settings.security.force-enable-mfa-enabled': '已启用',
  'settings.security.password-complexity': '密码复杂度',
  'settings.security.password-complexity-description': '当前密码复杂度要求',
  'settings.security.password-complexity.option.unsafe': '不安全',
  'settings.security.password-complexity.option.general': '一般',
  'settings.security.password-complexity.option.safe': '安全',
  'settings.security.password-complexity.option.very_safe': '非常安全',
  'settings.security.password-complexity.option.unsafe-description': '可以为任意字符。',
  'settings.security.password-complexity.option.general-description':
    '至少由大写字母、小写字母和数字的任意两种组合组成。',
  'settings.security.password-complexity.option.safe-description': '必须包括大小写字母和数字。',
  'settings.security.password-complexity.option.very_safe-description':
    '必须包含大小写字母、数字和特殊字符。',
  'settings.security.password-min-length': '密码最小长度',
  'settings.security.password-min-length-description': '密码长度至少为 {minLen} 位',
  'settings.security.password-min-length-unrestricted': '不限制',
  'settings.security.password-expire-time': '密码过期时间',
  'settings.security.password-expire-time-description':
    '至少 {days} 天修改一次密码，否则帐户自动锁定',
  'settings.security.password-expire-time-unrestricted': '不限制',
  'settings.security.password-failed-lock': '登陆失败自动自动锁定',
  'settings.security.password-failed-lock-description':
    '{min} 分钟内登陆 {fails} 次失败帐户自动锁定。',
  'settings.security.password-failed-lock-unrestricted': '不限制',
  'settings.security.password-history': '历史密码检查策略',
  'settings.security.password-history-description': '禁止使用前 {count} 次密码',
  'settings.security.password-history-unrestricted': '不检查',
  'settings.security.account-inactive-lock': '不活跃帐户自动锁定',
  'settings.security.account-inactive-lock-description': '帐户 {days} 天未登陆自动锁定',
  'settings.security.account-inactive-lock-unrestricted': '不限制',
  'settings.security.login-session-expiration-time': '会话过期时间',
  'settings.security.login-session-expiration-time-description':
    '{loginSessionInactivityHours}小时不活跃自动退出登陆，会话最长保持{loginSessionMaxHours}小时。',
  'settings.security.login-session-expiration-time-unrestricted': '不限制',
  'settings.security.login-session-expiration-time.input.0.tooltip':
    '不活跃会话自动过期时间（该时间周期内如果没有在该平台执行操作则视为不活跃）。',
  'settings.security.login-session-expiration-time.input.0.suffix': '小时',
  'settings.security.login-session-expiration-time.input.1.tooltip': '会话最长保持时间',
  'settings.security.login-session-expiration-time.input.1.suffix': '小时',
  'settings.security.login-session-expiration-tooltip':
    '为了安全起见，如果会话过期时间为0，则会话过期时间将被设置为30天',
  'settings.security.password-complexity.modify': '修改',
  'settings.security.password-min-length.modify': '修改',
  'settings.security.password-expire-time.modify': '修改',
  'settings.security.password-failed-lock.modify': '修改',
  'settings.security.password-history.modify': '修改',
  'settings.security.account-inactive-lock.modify': '修改',
  'settings.security.login-session-expiration-time.modify': '修改',
  'settings.security.password-expire-time.input.suffix': '天',
  'settings.security.password-failed-lock.input.0.tooltip': '连续的密码输入错误数。',
  'settings.security.password-failed-lock.input.0.suffix': '次失败',
  'settings.security.password-failed-lock.input.1.tooltip': '统计周期',
  'settings.security.password-failed-lock.input.1.suffix': '分钟',
  'settings.security.password-history.input.suffix': '天',
  'settings.security.account-inactive-lock.input.suffix': '天',
  'settings.security.login-session-expiration-time.input.suffix': '小时',
  'user.settings.menuMap.basic': '基本设置',
  'user.settings.menuMap.security': '安全设置',
  'user.settings.menuMap.binding': '账号绑定',
  'user.settings.menuMap.notification': '新消息通知',
  'user.settings.basic.avatar': '头像',
  'user.settings.basic.change-avatar': '更换头像',
  'user.settings.basic.email': '邮箱',
  'user.settings.basic.email-message': '请输入您的邮箱!',
  'user.settings.basic.fullname': '姓名',
  'user.settings.basic.profile': '个人简介',
  'user.settings.basic.profile-message': '请输入个人简介!',
  'user.settings.basic.profile-placeholder': '个人简介',
  'user.settings.basic.country': '国家/地区',
  'user.settings.basic.country-message': '请输入您的国家或地区!',
  'user.settings.basic.geographic': '所在省市',
  'user.settings.basic.geographic-message': '请输入您的所在省市!',
  'user.settings.basic.address': '街道地址',
  'user.settings.basic.address-message': '请输入您的街道地址!',
  'user.settings.basic.phone': '联系电话',
  'user.settings.basic.phone-message': '请输入您的联系电话!',
  'user.settings.basic.update': '更新基本信息',
  'user.settings.security.strong': '强',
  'user.settings.security.medium': '中',
  'user.settings.security.weak': '弱',
  'user.settings.security.password': '账户密码',
  'user.settings.security.resetPassword': '重置',
  'user.settings.security.password-description': '当前密码强度',
  'user.settings.security.mfa-phone': '密保手机',
  'user.settings.security.mfa-phone-description': '已绑定手机',
  'user.settings.security.question': '密保问题',
  'user.settings.security.question-description': '未设置密保问题，密保问题可有效保护账户安全',
  'user.settings.security.mfa-email': '安全邮箱',
  'user.settings.security.mfa-email-description': '已绑定邮箱',
  'user.settings.security.mfa-device': 'MFA 设备',
  'user.settings.security.mfa-device-description-unbound':
    '未绑定 MFA 设备，绑定后，可以进行二次确认',
  'user.settings.security.modify': '修改',
  'user.settings.security.set': '设置',
  'user.settings.security.bind': '绑定',
  'user.settings.security.unbind': '解绑',
  'user.settings.binding.taobao': '绑定淘宝',
  'user.settings.binding.taobao-description': '当前未绑定淘宝账号',
  'user.settings.binding.alipay': '绑定支付宝',
  'user.settings.binding.alipay-description': '当前未绑定支付宝账号',
  'user.settings.binding.dingding': '绑定钉钉',
  'user.settings.binding.dingding-description': '当前未绑定钉钉账号',
  'user.settings.binding.bind': '绑定',
  'user.settings.binding.unbind': '解绑',
  'user.settings.notification.password': '账户密码',
  'user.settings.notification.password-description': '其他用户的消息将以站内信的形式通知',
  'user.settings.notification.messages': '系统消息',
  'user.settings.notification.messages-description': '系统消息将以站内信的形式通知',
  'user.settings.notification.todo': '待办任务',
  'user.settings.notification.todo-description': '待办任务将以站内信的形式通知',
  'user.settings.open': '开',
  'user.settings.close': '关',
  'user.settings.basic.username': '用户名',
  'user.settings.security.mfa-device-description-bound': '已绑定',
  'user.settings.security.modify-password.title': '修改密码',
  'user.settings.security.modify-password.old-password': '原密码',
  'user.settings.security.modify-password.new-password': '新密码',
  'user.settings.security.mfa-device.bind': '绑定MFA',
  'user.settings.security.mfa-device.code-description':
    '扫描并添加MFA后，请获得两个连续的一次性密码，并将其输入到下面的输入框中。',
  'user.settings.security.mfa-device.refresh': '刷新',
  'user.settings.security.mfa-device.first-code': '第一个密码',
  'user.settings.security.mfa-device.second-code': '第二个密码',
  'settings.base.version': '当前版本',
  'settings.base.server-version': '后端',
  'settings.base.front-end-version': '前端',
};
