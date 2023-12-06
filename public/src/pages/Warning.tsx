import { Button, Result } from 'antd';
import React from 'react';

import { IntlContext } from '@/utils/intl';
import { history, useIntl, useSearchParams } from '@umijs/max';

const WarningPage: React.FC = ({}) => {
  const intl = new IntlContext('pages.403', useIntl());
  const [searchParams] = useSearchParams();
  const title = searchParams.get('title');
  const message = searchParams.get('message');

  return (
    <Result
      status="warning"
      title={title}
      subTitle={
        message ??
        intl.t(
          'defaultSubTitle',
          'An unknown error has occurred. Please contact the administrator.',
        )
      }
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          {intl.t('backHome', 'Back Home')}
        </Button>
      }
    />
  );
};

export default WarningPage;
