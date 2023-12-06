import { Button, Result } from 'antd';
import React from 'react';
import { history, useIntl } from 'umi';

import { IntlContext } from '@/utils/intl';

const NoFoundPage: React.FC = () => {
  const intl = new IntlContext('pages.403', useIntl());
  return (
    <Result
      status="404"
      title="404"
      subTitle={intl.t('subTitle', 'Sorry, the page you visited does not exist.')}
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          {intl.t('backHome', 'Back Home')}
        </Button>
      }
    />
  );
};

export default NoFoundPage;
