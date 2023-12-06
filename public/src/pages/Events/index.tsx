import React from 'react';

import { getEventLogs, getEvents } from '@/services/idas/events';
import { IntlContext } from '@/utils/intl';
import { PageContainer } from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';

import EventTableViewer from './components/EventTableViewer';

const Events: React.FC = ({}) => {
  const intl = new IntlContext('pages.events', useIntl());
  return (
    <PageContainer>
      <EventTableViewer
        parentIntl={intl}
        request={getEvents}
        logViewerProps={{
          request: getEventLogs,
        }}
      />
    </PageContainer>
  );
};

export default Events;
