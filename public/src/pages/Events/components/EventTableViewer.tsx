import { DatePicker, Timeline, Tooltip, Typography } from 'antd';
import lodash from 'lodash';
import moment from 'moment';
import type { RangeValue } from 'rc-picker/lib/interface';
import React, { useEffect, useState } from 'react';

import type { IntlContext } from '@/utils/intl';
import ProTable from '@ant-design/pro-table';
import { getLocale } from '@umijs/max';

import styles from './EventTableViewer.less';

moment.locale(getLocale());

const { Paragraph } = Typography;

type TimeRangeValue<DateType = moment.Moment> =
  | Exclude<RangeValue<DateType>, null>
  | (() => Exclude<RangeValue<DateType>, null>);

interface EventLogViewerProps {
  eventId: string;
  expanded: boolean;
  request: (params: API.getEventLogsParams) => Promise<API.GetEventLogsResponse>;
}

const EventLogViewer: React.FC<EventLogViewerProps> = ({ eventId, expanded, request }) => {
  const [dataSource, setDataSource] = useState<API.EventLog[]>();
  useEffect(() => {
    if (expanded && request) {
      request({ eventId }).then((ret) => {
        if (ret.data) {
          setDataSource(ret.data);
        }
      });
    }
  }, [expanded, eventId, request]);
  const parseLog = (log: any): React.ReactChild[] => {
    if (lodash.isString(log)) {
      try {
        return parseLog(JSON.parse(log));
      } catch {
        return [log];
      }
    } else if (lodash.isArray(log)) {
      return log.map((entry) => parseLog(entry)).flat();
    } else if (lodash.isObject(log)) {
      const ret: React.ReactChild[] = [];
      for (const key in log) {
        if (Object.prototype.hasOwnProperty.call(log, key)) {
          const element = log[key as keyof typeof log];
          ret.push(
            <div key={key} className={styles.logline}>
              <div className="log-title">[{key}]</div>
              <Paragraph
                className="log-content"
                ellipsis={{
                  expandable: true,
                }}
                copyable
              >
                {element}
              </Paragraph>
            </div>,
          );
        }
      }
      return ret;
    }
    return [log];
  };
  return (
    <Timeline>
      {dataSource?.map((item) => {
        return (
          <Timeline.Item key={item.id} color="green">
            <div>{moment(item.createTime).format('YYYY-MM-DD hh:mm:ss')}</div>
            {parseLog(item.log)}
            <div className={styles.logtime}>
              {moment(item.createTime).format('YYYY-MM-DD hh:mm:ss')}
            </div>
          </Timeline.Item>
        );
      })}
    </Timeline>
  );
};

interface EventTableViewerProps {
  parentIntl: IntlContext;
  request: (
    params: Pick<
      API.getEventsParams,
      'current' | 'pageSize' | 'endTime' | 'startTime' | 'keywords'
    >,
  ) => Promise<API.GetEventsResponse>;
  logViewerProps: Omit<EventLogViewerProps, 'eventId' | 'expanded'>;
}

const EventTableViewer: React.FC<EventTableViewerProps> = ({
  parentIntl: intl,
  logViewerProps,
  request,
}) => {
  const expandedRowRender = (
    record: API.Event,
    _: any,
    __: any,
    expanded: boolean,
  ): React.ReactNode => {
    return <EventLogViewer {...logViewerProps} eventId={record.id} expanded={expanded} />;
  };
  return (
    <ProTable<API.Event, API.currentUserEventsParams>
      columns={[
        {
          title: intl.t('title.eventTime', 'Event Time'),
          hideInSearch: true,
          dataIndex: 'createTime',
          width: 280,
          render: (__, entry) => {
            if (!entry?.createTime) {
              return '';
            }
            const createTime = moment(entry.createTime);
            return `${createTime.format('YYYY-MM-DD HH:mm:ss')} (${createTime.fromNow()})`;
          },
        },
        {
          title: intl.t('title.username', 'Username'),
          dataIndex: 'username',
          width: 180,
          hideInSearch: true,
          render: (_, entity) => {
            if (entity.username) {
              if (entity.userId) {
                return <Tooltip title={entity.userId}>{entity.username}</Tooltip>;
              }
              return entity.username;
            }
            return entity.userId;
          },
        },
        {
          title: intl.t('title.clientIP', 'Client IP'),
          dataIndex: 'clientIp',
          width: 180,
          hideInSearch: true,
        },
        {
          title: intl.t('title.location', 'Location'),
          dataIndex: 'location',
          width: 180,
          hideInSearch: true,
        },
        {
          title: intl.t('title.action', 'Action'),
          dataIndex: 'action',
          hideInSearch: true,
        },
        {
          title: intl.t('title.status', 'Status'),
          hideInSearch: true,
          dataIndex: 'status',
          valueEnum: {
            0: { text: intl.t('status.failed', 'Failed'), status: 'Error' },
            1: { text: intl.t('status.successd', 'Successd'), status: 'Success' },
            false: { text: intl.t('status.failed', 'Failed'), status: 'Error' },
            true: { text: intl.t('status.successd', 'Successd'), status: 'Success' },
          },
          render(dom, entity) {
            if (!entity.status && entity.message) {
              return (
                <Tooltip title={entity.message}>
                  <span>{dom}</span>
                </Tooltip>
              );
            }
            return dom;
          },
        },
        {
          title: intl.t('title.took', 'Took'),
          hideInSearch: true,
          dataIndex: 'took',
          render: (__, entry) => {
            if (!entry.took) {
              return '';
            }
            const took: number = entry.took / 1000 / 1000;
            if (took < 10) {
              return intl.t('took.millisecond', `{took} milliseconds`, '', {
                took: took.toFixed(2),
              });
            } else if (took < 1e3) {
              return intl.t('took.millisecond', `{took} milliseconds`, '', {
                took: took.toFixed(0),
              });
            } else if (took < 10 * 1e3) {
              return intl.t('took.second', `{took} seconds`, '', {
                took: (took / 1000).toFixed(2),
              });
            } else if (took < 120 * 1e3) {
              return intl.t('took.second', `{took} seconds`, '', {
                took: (took / 1000).toFixed(0),
              });
            } else {
              return moment.duration(took).humanize();
            }
          },
        },
        {
          title: intl.t('title.keywords', `Keywords`),
          hideInTable: true,
          dataIndex: 'keywords',
          hideInSearch: false,
        },
        {
          title: intl.t('title.timeRange', 'Time Range'),
          hideInTable: true,
          dataIndex: 'timeRange',
          valueType: 'dateTimeRange',
          hideInSearch: false,
          search: {
            transform: (value: string[]) => {
              return { startTime: moment(value[0]).format(), endTime: moment(value[1]).format() };
            },
          },
          initialValue: [moment().startOf('day'), moment().endOf('day')],
          formItemProps: {
            style: { width: 500 },
          },
          renderFormItem() {
            const ranges: { text: string; value: TimeRangeValue<moment.Moment> }[] = [
              {
                text: intl.t('timeRange.today', 'Today'),
                value: [moment().startOf('day'), moment().endOf('day')],
              },
              {
                text: intl.t('timeRange.lastDay', 'Last day'),
                value: [moment().add(-1, 'd'), moment()],
              },
              {
                text: intl.t('timeRange.last3Days', 'Last 3 Days'),
                value: [moment().add(-3, 'd'), moment()],
              },
              {
                text: intl.t('timeRange.lastWeek', 'Last Week'),
                value: [moment().add(-1, 'w'), moment()],
              },
              {
                text: intl.t('timeRange.lastMonth', 'Last Month'),
                value: [moment().add(-1, 'M'), moment()],
              },
              {
                text: intl.t('timeRange.thisWeek', 'This week'),
                value: [moment().startOf('week'), moment().endOf('week')],
              },
              {
                text: intl.t('timeRange.thisMonth', 'This month'),
                value: [moment().startOf('month'), moment().endOf('month')],
              },
            ];
            return (
              <DatePicker.RangePicker
                showTime
                {...{
                  ranges: ranges.reduce(
                    (acc, item) => ((acc[item.text] = item.value), acc),
                    {} as Record<string, TimeRangeValue<moment.Moment>>,
                  ),
                }}
              />
            );
          },
        },
      ]}
      request={async (params) => {
        return request(params);
      }}
      rowKey={'id'}
      expandable={{
        expandedRowRender,
      }}
      pagination={{ pageSize: 10, showSizeChanger: true, pageSizeOptions: [5, 10, 20, 50, 100] }}
    />
  );
};

export default EventTableViewer;
