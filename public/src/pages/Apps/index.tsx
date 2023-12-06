import React from 'react';

import List from '@/components/List';
import { deleteApp, getApps } from '@/services/idas/apps';
import { IntlContext } from '@/utils/intl';
import { PageContainer } from '@ant-design/pro-components';
import { useIntl, history } from '@umijs/max';

const Apps: React.FC = ({}) => {
  const intl = new IntlContext('pages.apps', useIntl());
  return (
    <PageContainer>
      <List
        intl={intl}
        request={{
          list: getApps,
          delete: async (items: { id: string }[], options?: Record<string, any>) => {
            for (const { id } of items) {
              await deleteApp({ id }, options);
            }
          },
        }}
        onEdit={(item) => {
          history.push(`/apps/${item.id}/edit`);
        }}
        onCreate={() => {
          history.push(`/apps/create`);
        }}
        onClick={(item) => {
          history.push(`/apps/${item.id}`);
        }}
      />
    </PageContainer>
  );
};

export default Apps;
