import { Card, Row, message } from 'antd';
import React, { useRef } from 'react';

import { AvatarUploadField } from '@/components/Avatar';
import { createPage, getPage, updatePage } from '@/services/idas/pages';
import { IntlContext } from '@/utils/intl';
import type { RequestError } from '@/utils/request';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  ProForm,
  ProFormGroup,
  ProFormText,
  ProFormTextArea,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { useIntl, useParams } from '@umijs/max';

import type { FieldConfig } from './components/FieldsConfig';
import FieldsConfig from './components/FieldsConfig';

/**
 * @en-US Add page
 * @zh-CN 添加页面
 * @param fields
 */
const handleAdd = async (fields: API.PageConfig) => {
  const hide = message.loading('Adding ...');
  try {
    await createPage(fields);
    hide();
    message.success('Added successfully');
    return true;
  } catch (error) {
    hide();
    if (!(error as RequestError).handled) {
      message.error('Adding failed, please try again!');
    }
    return false;
  }
};

/**
 * @en-US Update Page
 * @zh-CN 更新页面
 * @param fields
 */
const handleUpdate = async (fields: API.PageConfig) => {
  const hide = message.loading('Updating ...');
  try {
    await updatePage({ id: fields.id }, fields);
    hide();
    message.success('Updated successfully');
    return true;
  } catch (error) {
    hide();
    if (!(error as RequestError).handled) {
      message.error('Update failed, please try again!');
    }
    return false;
  }
};
const PageForm: React.FC = ({}) => {
  const { pageId } = useParams();
  const intl = new IntlContext('pages.pages', useIntl());
  const formRef = useRef<ProFormInstance<API.PageConfig>>();
  const handleSubmit = pageId === 'create' ? handleAdd : handleUpdate;
  return (
    <Card
      style={{
        height: '100%',
      }}
      bodyStyle={{ height: '100%' }}
    >
      <ProForm<API.PageConfig>
        style={{ maxWidth: 800, margin: '8px auto auto' }}
        formRef={formRef}
        request={
          pageId !== 'create'
            ? async () => {
                const { data, errorCode, errorMessage } = await getPage({ id: pageId });
                if (data) {
                  return data;
                }
                throw new Error(`Page configuration query failed: [${errorCode}]${errorMessage}`);
              }
            : undefined
        }
        onFinish={async (value) => {
          handleSubmit({
            ...value,
            fields: (value.fields as FieldConfig[])?.map(
              ({ id, index, ...item }: FieldConfig) => item,
            ),
          });
        }}
      >
        <ProFormText hidden width="md" name="id" label="id" />
        <ProFormSwitch hidden width="md" name="isDisable" label="isDisable" />
        <Row>
          <ProFormGroup spaceProps={{ style: { display: 'inline-block' } }}>
            <ProFormText
              width="md"
              name="name"
              label={intl.t('name.label', 'Page Name')}
              placeholder={intl.t('name.placeholder', 'Please input page name')}
            />
            <ProFormTextArea
              width="lg"
              name="description"
              label={intl.t('description.label', 'Description')}
              placeholder={intl.t('description.placeholder', 'Please input page description')}
            />
          </ProFormGroup>
          <AvatarUploadField
            label={intl.t('icon.label', 'Icon')}
            name={'icon'}
            formItemProps={{
              style: { marginLeft: 100 },
            }}
          />
        </Row>
        <FieldsConfig label={intl.t('fields.label', 'Fields')} name={'fields'} parentIntl={intl} />
      </ProForm>
    </Card>
  );
};
export default PageForm;
