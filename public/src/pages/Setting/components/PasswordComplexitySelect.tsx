import { Button, Select, Space, Tooltip } from 'antd';
import { isNumber } from 'lodash';
import { useEffect, useState } from 'react';

import type { PasswordComplexityName, PasswordComplexityValue } from '@/services/idas/enums';
import { PasswordComplexity } from '@/services/idas/enums';
import { enumToOptions } from '@/utils/enum';
import { IntlContext } from '@/utils/intl';
import { CheckCircleTwoTone, CloseCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';

// enum PasswordComplexityEnum {
//   unsafe = 1,
//   general = 2,
//   safe = 3,
//   verySafe = 4,
// }

interface ModifyPasswordProps {
  parentIntl: IntlContext;
  value?: API.PasswordComplexity;
  onSave: (value: PasswordComplexityValue) => Promise<boolean>;
}

export const PasswordComplexityToolip = ({ parentIntl }: { parentIntl: IntlContext }) => {
  const intl = new IntlContext('password-complexity', parentIntl);
  const options = [
    {
      label: intl.t('option.unsafe', 'Unsafe'),
      desc: intl.t('option.unsafe-description', 'Any character.'),
    },
    {
      label: intl.t('option.general', 'General'),
      desc: intl.t(
        'option.general-description',
        'Composed of at least any two combinations of uppercase letters, lowercase letters, and numbers.',
      ),
    },
    {
      label: intl.t('option.safe', 'Safe'),
      desc: intl.t(
        'option.safe-description',
        'Must include uppercase and lowercase letters and numbers.',
      ),
    },
    {
      label: intl.t('option.very_safe', 'Very Safe'),
      desc: intl.t(
        'option.very_safe-description',
        'Must contain uppercase and lowercase letters, numbers, and special characters.',
      ),
    },
  ];
  return (
    <Tooltip
      placement="bottomLeft"
      overlayStyle={{ maxWidth: 'max-content' }}
      overlayInnerStyle={{ backgroundColor: 'rgba(61, 62, 64, 0.85)' }}
      title={
        <div>
          {options.map((option) => {
            return (
              <li key={option.label}>
                {option.label}: {option.desc}
              </li>
            );
          })}
        </div>
      }
    >
      <QuestionCircleOutlined style={{ color: 'rgba(61, 62, 64, 0.45)', marginInlineStart: 4 }} />
    </Tooltip>
  );
};

export const getPasswordComplexityValue = (c?: API.PasswordComplexity): PasswordComplexityValue => {
  if (isNumber(c)) {
    return c;
  } else if (c !== undefined) {
    return PasswordComplexity[c];
  }
  return 0;
};

export const getPasswordComplexityName = (c?: API.PasswordComplexity): PasswordComplexityName => {
  if (isNumber(c)) {
    return PasswordComplexity[c] as PasswordComplexityName;
  } else if (c !== undefined) {
    return c;
  }
  return PasswordComplexity[0] as PasswordComplexityName;
};

export default ({ parentIntl, value, onSave }: ModifyPasswordProps) => {
  const intl = new IntlContext('password-complexity', parentIntl);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [tmpPasswordComplexity, setTmpPasswordComplexity] = useState<PasswordComplexityValue>(0);

  useEffect(() => {
    if (value !== undefined) setTmpPasswordComplexity(getPasswordComplexityValue(value));
  }, [value]);
  return isEditing ? (
    <>
      <Space.Compact size="small">
        <Select<PasswordComplexity>
          style={{ width: 120 }}
          defaultValue={getPasswordComplexityValue(value)}
          onChange={(val) => {
            setTmpPasswordComplexity(val);
          }}
          options={enumToOptions(PasswordComplexity, intl, 'option')}
        />
        <Button
          size="small"
          onClick={async () => {
            if (await onSave(tmpPasswordComplexity)) {
              setIsEditing(false);
            }
          }}
        >
          <CheckCircleTwoTone />
        </Button>
        <Button
          size="small"
          onClick={() => {
            setIsEditing(false);
          }}
        >
          <CloseCircleOutlined />
        </Button>
      </Space.Compact>
    </>
  ) : (
    <a
      onClick={() => {
        setIsEditing(true);
      }}
    >
      {intl.t('modify', 'Modify')}
    </a>
  );
};
