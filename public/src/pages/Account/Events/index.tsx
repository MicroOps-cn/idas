import React from 'react';

import defaultSettings from '@/../config/defaultSettings';
import { getActions } from '@/components/RightContent';
import EventTableViewer from '@/pages/Events/components/EventTableViewer';
import { currentUserEventLogs, currentUserEvents } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { getPublicPath } from '@/utils/request';
import { ProLayout } from '@ant-design/pro-components';
import { useModel, useIntl, history } from '@umijs/max';

const Events: React.FC = ({}) => {
  const intl = new IntlContext('pages.events', useIntl());
  const {
    initialState = {
      settings: { navTheme: 'light' },
      globalConfig: null,
      currentUser: { role: undefined },
    },
  } = useModel('@@initialState');
  // @ts-expect-error
  const { navTheme = 'dark' } = initialState.settings;

  const globalConfig = initialState?.globalConfig ?? null;
  return (
    <ProLayout
      logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
      title={globalConfig?.title ?? defaultSettings.title}
      onMenuHeaderClick={() => {
        history.push('/');
      }}
      layout="top"
      navTheme={navTheme}
      actionsRender={() => getActions(initialState?.currentUser?.role)}
    >
      <EventTableViewer
        parentIntl={intl}
        request={currentUserEvents}
        logViewerProps={{
          request: currentUserEventLogs,
        }}
      />
    </ProLayout>
  );
};

export default Events;
