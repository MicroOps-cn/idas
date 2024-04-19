import { List, Space } from 'antd';
import type { PrimitiveType } from 'intl-messageformat';
import { Component } from 'react';

import { IntlContext } from '@/utils/intl';
import { useModel } from '@umijs/max';

type Unpacked<T> = T extends (infer U)[] ? U : T;

interface BaseViewProps {
  parentIntl: IntlContext;
}

declare const buildVersion: string;

export const BaseView = ({ parentIntl }: BaseViewProps) => {
  const { initialState } = useModel('@@initialState');
  const serverVersion = initialState?.globalConfig?.version;
  const intl = new IntlContext('base', parentIntl);

  const data = [
    {
      title: intl.t('version'),
      description: (
        <>
          {serverVersion === buildVersion ? (
            buildVersion
          ) : (
            <Space>
              <span>
                {intl.t('server-version')}: {serverVersion}
              </span>
              |
              <span>
                {intl.t('front-end-version')}: {buildVersion}
              </span>
            </Space>
          )}
        </>
      ),
      actions: undefined,
    },
  ];
  return (
    <>
      <List<Unpacked<typeof data>>
        itemLayout="horizontal"
        dataSource={data}
        renderItem={(item) => (
          <List.Item actions={item.actions}>
            <List.Item.Meta title={item.title} description={item.description} />
          </List.Item>
        )}
      />
    </>
  );
};

export default BaseView;
