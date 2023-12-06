import React from 'react';

import List from '@/components/List';
import { getPages, patchPages } from '@/services/idas/pages';
import { IntlContext } from '@/utils/intl';
import { MenuRender } from '@/utils/menu';
import { getLocation } from '@/utils/request';
import { PageContainer } from '@ant-design/pro-components';
import { history, useIntl, useModel } from '@umijs/max';

const Page: React.FC<any> = ({}) => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const intl = new IntlContext('pages.pages', useIntl());
  const { pathname } = getLocation();
  return (
    <PageContainer>
      <List
        intl={intl}
        request={{
          list: async (params, options) => {
            const resp = await getPages(params, options);
            return {
              ...resp,
              data: resp.data?.map(({ icon, ...item }) => ({ avatar: icon, ...item })),
            };
          },
          patch: async (body, options) => {
            const ret = await patchPages(body, options);
            if (ret.success) {
              MenuRender.getInstance().render((dynamicMenu) => {
                if (dynamicMenu.length > 0) {
                  setInitialState({
                    ...initialState,
                    dynamicMenu,
                  });
                }
              }, true);
            }
            return ret;
          },
        }}
        onEdit={(entry) => {
          history.push(`${pathname}/${entry.id}`);
        }}
        onCreate={() => {
          history.push(`${pathname}/create`);
        }}
      />
    </PageContainer>
  );
};

export default Page;
