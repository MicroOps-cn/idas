import { Button, Result } from 'antd';
import React from 'react';
import { history, useIntl } from 'umi';

import { IntlContext } from '@/utils/intl';

const ForbiddenPage: React.FC = () => {
  const intl = new IntlContext('pages.403', useIntl());
  return (
    <Result
      status="403"
      title="403"
      subTitle={intl.t('subTitle', 'You do not have permission to access this page.')}
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          {intl.t('backHome', 'Back Home')}
        </Button>
      }
    />
  );
};

export default ForbiddenPage;
