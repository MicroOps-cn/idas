import { Card, Divider, Radio, message, Switch, Form, Row, Space } from 'antd';
import type { FormLayout } from 'antd/lib/form/Form';
import { isString } from 'lodash';
import React, { useEffect, useState } from 'react';

import { PageFieldType } from '@/services/idas/enums';
import { createPageData, getPage, getPageData, updatePageData } from '@/services/idas/pages';
import { IntlContext } from '@/utils/intl';
import type { ProFormColumnsType } from '@ant-design/pro-form';
import { BetaSchemaForm } from '@ant-design/pro-form';
import { useIntl, history, useParams } from '@umijs/max';

const PageDataForm: React.FC = ({}) => {
  const { pageId, id } = useParams();
  const intl = new IntlContext('pages.page', useIntl());
  const [columns, setColumns] = useState<ProFormColumnsType[]>([]);
  const [formLayoutType, setFormLayoutType] = useState<FormLayout>('horizontal');
  const [grid, setGrid] = useState(true);

  const renderField = (field: API.FieldConfig): ProFormColumnsType => {
    const column: ProFormColumnsType = {
      title: field.displayName ?? field.name,
      dataIndex: field.name,
      initialValue: field.defaultValue,
      valueEnum: field.valueEnum,
      width: 'md',
    };
    if (field.valueType === PageFieldType.multiSelect || field.valueType === 'multiSelect') {
      column.valueType = 'select';
    } else if (isString(field.valueType)) {
      column.valueType = field.valueType;
    } else {
      column.valueType = PageFieldType[field.valueType] as any;
    }
    return column;
  };

  useEffect(() => {
    getPage({ id: pageId }).then(({ data }) => {
      if (data && data.fields) {
        setColumns(data.fields.map(renderField));
      }
    });
  }, [id, pageId]);
  const handleSubmit = async (data: Omit<API.PageData, 'createTime' | 'updateTime'>) => {
    const hide = message.loading(
      pageId === 'create'
        ? intl.t('message.adding', 'Adding ...')
        : intl.t('message.updating', 'Updating ...'),
    );
    try {
      if (id === 'create') {
        await createPageData({ pageId: pageId }, {
          ...data,
          id: undefined,
        } as API.createPageDataParams);
      } else {
        await updatePageData({ pageId: pageId, id: data.id }, data);
      }
      hide();
      message.success(
        pageId === 'create'
          ? intl.t('message.added.success', 'Added successfully')
          : intl.t('message.updated.success', 'Updated successfully'),
      );
      return true;
    } catch (error) {
      hide();
      message.error(
        pageId === 'create'
          ? intl.t('message.add.failed', 'Add successfully')
          : intl.t('message.update.failed', 'Update successfully'),
      );
      return false;
    }
  };
  return (
    <Card
      style={{
        height: '100%',
      }}
      bodyStyle={{ height: '100%' }}
    >
      <Row>
        <Space size="large">
          <Form.Item label={intl.t('layout', 'Layout of labels')}>
            <Radio.Group
              defaultValue={formLayoutType}
              onChange={(e) => {
                setFormLayoutType(e.target.value);
              }}
              optionType="button"
              options={['horizontal', 'vertical', 'inline']}
            />
          </Form.Item>
          <Form.Item label={intl.t('gird', 'Gird')}>
            <Switch defaultChecked={grid} onChange={setGrid} />
          </Form.Item>
        </Space>
      </Row>
      <Divider style={{ margin: '0 0 24px 0' }} />
      <BetaSchemaForm<Record<string, any>>
        layoutType={'Form'}
        grid={grid}
        colProps={{
          span: 12,
        }}
        layout={formLayoutType}
        rowProps={{
          gutter: [16, 16],
        }}
        columns={columns}
        request={
          id !== 'create'
            ? async () => {
                const { data, errorCode, errorMessage } = await getPageData({ pageId, id });
                if (data?.data) {
                  return data.data;
                }
                throw new Error(`Page configuration query failed: [${errorCode}]${errorMessage}`);
              }
            : undefined
        }
        onFinish={async (value) => {
          if (
            await handleSubmit({
              data: value,
              id: id,
              pageId: pageId,
            })
          ) {
            history.push('./');
          }
        }}
      />
    </Card>
  );
};
export default PageDataForm;
