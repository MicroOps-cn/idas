import { Modal, Tree } from 'antd';
import type { DataNode } from 'antd/lib/tree';
import React, { useEffect, useState } from 'react';

import { getPermissions } from '@/services/idas/permissions';
import { IntlContext } from '@/utils/intl';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { ProFormText, StepsForm } from '@ant-design/pro-form';

export type FormValueType = Omit<API.UpdateRoleRequest, 'id'> & {
  id?: string;
};
export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<boolean>;
  modalVisible: boolean;
  values?: API.RoleInfo;
  parentIntl: IntlContext;
};

const CreateOrUpdateForm: React.FC<UpdateFormProps> = ({ parentIntl, ...props }) => {
  const intl = new IntlContext('form', parentIntl);
  const [loading, setLoading] = useState<boolean>(false);
  const { values: dfValues, modalVisible, onSubmit, onCancel } = props;
  const [permissionList, setPermissionList] = useState<API.PermissionInfo[]>([]);
  const [checkedKeys, setCheckedKeys] = useState<string[]>([]);
  const [treeHeight, setTreeHeight] = useState<number>(500);
  const [currentStep, setCurrentStep] = useState<number>(0);

  useEffect(() => {
    const resetTreeHeight = () => {
      if (window.innerHeight > 640) {
        setTreeHeight(window.innerHeight - 440);
      } else {
        setTreeHeight(200);
      }
    };
    resetTreeHeight();
    window.addEventListener('resize', resetTreeHeight);
  }, []);

  useEffect(() => {
    if (modalVisible) {
      setCurrentStep(0);
      setCheckedKeys(
        dfValues?.permission
          ? dfValues.permission.map((item) => {
              return item.id;
            })
          : [],
      );
      setLoading(true);
      getPermissions({ pageSize: 3000 })
        .then((resp) => {
          if (resp.success && resp.data) {
            setPermissionList(resp.data);
          }
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [dfValues, modalVisible]);

  const genPermissionTree = (
    permissions: API.PermissionInfo[] | undefined,
    parent_id: string | undefined,
  ): DataNode[] => {
    if (permissions === undefined) {
      return [];
    }
    const temp: DataNode[] = [];
    const treeArr: API.PermissionInfo[] = permissions;

    treeArr.forEach((item, index) => {
      if (item.parentId === parent_id) {
        const children = genPermissionTree(treeArr, item.id);
        if (children.length > 0) {
          // 递归调用此函数
          temp.push({
            key: treeArr[index].id,
            title: treeArr[index].description ? treeArr[index].description : treeArr[index].name,
            children: children,
          });
        } else {
          temp.push({
            key: treeArr[index].id,
            title: treeArr[index].description ? treeArr[index].description : treeArr[index].name,
          });
        }
      }
    });
    return temp;
  };

  return (
    <StepsForm<FormValueType>
      stepsProps={{
        size: 'small',
      }}
      current={currentStep}
      onCurrentChange={setCurrentStep}
      formProps={{
        preserve: false,
      }}
      stepsFormRender={(dom, submitter) => {
        return (
          <Modal
            width={640}
            confirmLoading={loading}
            bodyStyle={{ padding: '32px 40px 48px' }}
            destroyOnClose
            title={intl.t(
              dfValues ? 'title.roleUpdate' : 'title.roleCreate',
              dfValues ? 'Modify role' : 'Add role',
            )}
            style={{
              maxHeight: 'calc(100vh - 200px)',
            }}
            open={modalVisible}
            footer={submitter}
            onCancel={() => {
              Modal.confirm({
                title: intl.t('cancel?', 'Cancel editing?'),
                icon: <ExclamationCircleOutlined />,
                onOk() {
                  onCancel();
                },
                maskClosable: true,
              });
            }}
          >
            {dom}
          </Modal>
        );
      }}
      onFinish={async (values) => {
        return onSubmit({ ...values, permission: checkedKeys });
      }}
    >
      <StepsForm.StepForm
        initialValues={dfValues}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 14 }}
        layout={'horizontal'}
        title={intl.t('basicConfig', 'Basic information')}
      >
        <ProFormText hidden={true} name="id" />
        <ProFormText
          name="name"
          label={intl.t('name.label', 'Name')}
          width="md"
          rules={[
            {
              required: true,
              message: intl.t('name.required', 'Please input role name!'),
            },
          ]}
        />

        <ProFormText
          name="description"
          label={intl.t('description.label', 'Description')}
          width="md"
        />
      </StepsForm.StepForm>
      <StepsForm.StepForm initialValues={{}} title={intl.t('permission.title', 'Permission')}>
        <Tree
          checkable
          height={treeHeight}
          defaultExpandAll={true}
          checkedKeys={checkedKeys}
          onCheck={(checked) => {
            // @ts-ignore
            setCheckedKeys(checked);
          }}
          treeData={genPermissionTree(permissionList, '')}
        />
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateOrUpdateForm;
