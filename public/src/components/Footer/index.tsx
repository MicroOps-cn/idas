import { useIntl, useModel } from 'umi';

import type { FooterProps } from '@ant-design/pro-components';
import { DefaultFooter } from '@ant-design/pro-components';

const Footer: React.FC<FooterProps> = (props) => {
  const intl = useIntl();
  const currentYear = new Date().getFullYear();

  const { initialState } = useModel('@@initialState');
  const copyright =
    initialState?.globalConfig?.copyright ??
    intl.formatMessage({
      id: 'app.copyright.produced',
      defaultMessage: 'Wiseasy',
    });
  return <DefaultFooter {...props} copyright={`${currentYear} ${copyright}`} />;
};

export default Footer;
