import { Modal } from 'antd';
import React, { useEffect, useState } from 'react';

import { AvatarUploadField } from '@/components/Avatar';
import { UserStatus } from '@/services/idas/enums';
import { getUserInfo } from '@/services/idas/users';
import { IntlContext } from '@/utils/intl';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { ProFormText, StepsForm } from '@ant-design/pro-form';

import GrantView from './GrantView';

export type FormValueType = Omit<API.UpdateUserRequest, 'id'> & {
  id?: string;
};
export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: FormValueType) => void;
  onSubmit: (values: FormValueType) => Promise<boolean>;
  modalVisible: boolean;
  values?: API.UserInfo;
  title?: React.ReactNode;
  parentIntl: IntlContext;
};

const CreateOrUpdateForm: React.FC<UpdateFormProps> = ({ parentIntl, ...props }) => {
  const intl = new IntlContext('form', parentIntl);
  const [currentStep, setCurrentStep] = useState<number>(0);
  const { values: dfValues, title, modalVisible, onSubmit, onCancel } = props;
  const [userApps, setUserApps] = useState<API.UserApp[]>([]);
  useEffect(() => {
    setCurrentStep(0);
    if (dfValues?.id && modalVisible) {
      getUserInfo({ id: dfValues.id }).then((ret) => {
        if (ret.success && ret.data?.apps) {
          setUserApps(ret.data.apps);
        }
      });
    }
  }, [dfValues, modalVisible]);

  return (
    <StepsForm<FormValueType>
      stepsProps={{
        size: 'small',
      }}
      formProps={{
        preserve: false,
      }}
      current={currentStep}
      onCurrentChange={setCurrentStep}
      stepsFormRender={(dom, submitter) => {
        return (
          <Modal
            width={640}
            styles={{
              body: { padding: '32px 40px 48px' },
            }}
            destroyOnClose
            title={title}
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
        return onSubmit({
          ...values,
          apps: userApps.map((app) => ({ id: app.id, roleId: app.roleId })),
          status: dfValues?.status ? dfValues.status : UserStatus.normal,
          isDelete: dfValues?.isDelete ? dfValues?.isDelete : false,
        });
      }}
    >
      <StepsForm.StepForm<API.UserInfo>
        initialValues={{ ...dfValues }}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 14 }}
        layout={'horizontal'}
        title={intl.t('title.basicConfig', 'Basic')}
        className="basic-view"
      >
        <AvatarUploadField label={intl.t('avatar.label', 'Avatar')} name={'avatar'} />
        <ProFormText hidden={true} name="id" />
        <ProFormText hidden={true} name="storage" />
        <ProFormText
          name="username"
          label={intl.t('userName.label', 'Username')}
          width="md"
          rules={[
            {
              required: true,
              message: intl.t('userName.required', 'Please input username!'),
            },
            {
              pattern: /^[-_A-Za-z0-9]+$/,
              message: intl.t('name.invalid', 'Username format error!'),
            },
          ]}
        />

        <ProFormText name="fullName" label={intl.t('fullName.label', 'FullName')} width="md" />
        <ProFormText name="email" label={intl.t('email.label', 'Email')} width="md" />
        <ProFormText
          name="phoneNumber"
          label={intl.t('phoneNumber.label', 'Telephone number')}
          width="md"
        />
      </StepsForm.StepForm>
      <StepsForm.StepForm
        className="grant-view"
        initialValues={{}}
        title={intl.t('app.title', 'App')}
      >
        <GrantView granting parentIntl={intl} apps={userApps} onChange={setUserApps} />
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default CreateOrUpdateForm;
