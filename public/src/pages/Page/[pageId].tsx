import { isString } from 'lodash';
import React, { useEffect, useMemo, useState } from 'react';

import type { TableItem } from '@/components/Table';
import Table from '@/components/Table';
import { PageFieldType } from '@/services/idas/enums';
import { getPageDatas, patchPageDatas, getPage } from '@/services/idas/pages';
import { IntlContext } from '@/utils/intl';
import { PageContainer } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-table';
import { useIntl, useParams } from '@umijs/max';

const Page: React.FC = ({}) => {
  const { pageId } = useParams();
  const intl = new IntlContext('pages.page', useIntl());

  const [pageConfig, setPageConfig] = useState<API.PageConfig>();
  useEffect(() => {
    getPage({ id: pageId }).then(({ data }) => {
      if (data) {
        setPageConfig(data);
      }
    });
  }, [pageId]);

  const request = useMemo(() => {
    return {
      list: async (params: { pageSize?: number; current?: number; keywords?: string }) => {
        return getPageDatas({ pageId: pageId, ...params }).then(({ data, ...options }) => {
          return {
            ...options,
            data: data?.map(({ data: d, ...meta }) => ({ ...meta, ...d })),
          };
        });
      },
      delete: async (body: { id: string }[]) => {
        return patchPageDatas(
          { pageId: pageId },
          body.map(({ id }) => ({ pageId: pageId, id, isDelete: true })),
        );
      },
    };
  }, [pageId]);

  const renderFields = (fields?: API.FieldConfig[]): ProColumns<TableItem>[] => {
    if (!fields || fields.length === 0) {
      return [];
    }
    return fields.map((field) => {
      const column: ProColumns<TableItem> = {
        title: field.displayName ?? field.name,
        dataIndex: field.name,
        initialValue: field.defaultValue,
        valueEnum: field.valueEnum,
        search: false,
        width: 'md',
      };
      if (field.valueType === PageFieldType.multiSelect || field.valueType === 'multiSelect') {
        column.valueType = 'select';
      } else if (isString(field.valueType)) {
        column.valueType = field.valueType;
      } else {
        column.valueType = PageFieldType[field.valueType] as typeof column.valueType;
      }
      return column;
    });
  };

  return (
    <PageContainer title={pageConfig?.name}>
      <Table
        search={false}
        columns={renderFields(pageConfig?.fields)}
        intl={intl}
        request={request}
      />
    </PageContainer>
  );
};

export default Page;
