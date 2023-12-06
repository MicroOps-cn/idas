import type { InputProps as AntInputProps } from 'antd';
import { Tooltip } from 'antd';
import { Input as AntInput, Button, Space, message } from 'antd';
import { useEffect, useState } from 'react';

import type { IntlContext } from '@/utils/intl';
import { CheckCircleTwoTone, CloseCircleOutlined } from '@ant-design/icons';

interface InputProps<T extends string | number = string>
  extends Pick<AntInputProps, 'type' | 'placeholder' | 'style'> {
  intl: IntlContext;
  value?: T;
  onSave: (value: T) => Promise<boolean>;
  autoSave?: boolean;
  tooltip?: string;
  suffix?: string;
  prefix?: string;
}
type InputFuncType = <T extends string | number = string>(
  props: InputProps<T>,
) => React.ReactElement;

export const Input: InputFuncType = ({
  intl,
  onSave,
  autoSave,
  value,
  tooltip,
  prefix,
  suffix,
  ...props
}) => {
  const [tempValue, setTempValue] = useState(value);
  const [isEditing, setIsEditing] = useState(false);
  const handleSave = async () => {
    if (tempValue !== undefined && (await onSave(tempValue))) {
      message.info(intl.t('finish', `Update successful.`));
      setIsEditing(false);
    }
  };
  useEffect(() => {
    setTempValue(value);
  }, [value]);
  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      {isEditing ? (
        <Space.Compact size="small">
          <Tooltip title={tooltip ? intl.t('input.tooltip', tooltip) : undefined}>
            <AntInput
              size="small"
              min={1}
              defaultValue={value}
              onChange={(e) => {
                if (props.type === 'number') {
                  setTempValue(parseInt(e.target.value, 10) as any);
                } else {
                  setTempValue(e.target.value as any);
                }
              }}
              onBlur={() => {
                if (autoSave) handleSave();
              }}
              suffix={suffix ? intl.t('input.suffix', suffix) : undefined}
              prefix={prefix ? intl.t('input.prefix', prefix) : undefined}
              {...props}
              style={{ width: 70, ...props.style }}
            />
          </Tooltip>
          <Button size="small" hidden={autoSave} onClick={handleSave}>
            <CheckCircleTwoTone />
          </Button>
          <Button
            size="small"
            onClick={() => {
              setIsEditing(false);
            }}
            hidden={autoSave}
          >
            <CloseCircleOutlined />
          </Button>
        </Space.Compact>
      ) : (
        <a
          style={{ marginLeft: 10 }}
          onClick={() => {
            setIsEditing(true);
          }}
        >
          {intl.t('modify', `Modify`)}
        </a>
      )}
    </div>
  );
};

interface InputArrayProps<T extends string | number = string> extends Pick<AntInputProps, 'type'> {
  intl: IntlContext;
  value?: T[];
  suffix?: string[];
  prefix?: string[];
  tooltip?: string[];
  placeholder?: string[];
  onSave: (value: T[]) => Promise<boolean>;
  autoSave?: boolean;
  count: number;
  style?: React.CSSProperties[];
}
type InputArrayType = <T extends string | number = string>(
  props: InputArrayProps<T>,
) => React.ReactElement;

export const InputArray: InputArrayType = ({
  intl,
  onSave,
  suffix,
  prefix,
  count,
  autoSave,
  style,
  tooltip,
  value,
  type,
}) => {
  const [tempValues, setTempValues] = useState(value ?? []);
  const [isEditing, setIsEditing] = useState(false);
  const handleSave = async () => {
    if (tempValues && (await onSave(tempValues))) {
      message.info(intl.t('finish', `Update successful.`));
      setIsEditing(false);
    }
  };
  const setTempValue = (idx: number, v: any) => {
    setTempValues([...tempValues.slice(0, idx), v, ...tempValues.slice(idx + 1)]);
  };
  const inputGroup = () => {
    const inputs: React.ReactElement[] = [];
    for (let index = 0; index < count; index++) {
      const top = tooltip?.[index] ?? undefined;
      const sfx = suffix?.[index] ?? undefined;
      const pfx = prefix?.[index] ?? undefined;
      const sty = style?.[index] ?? undefined;
      inputs.push(
        <Tooltip
          key={`input-${index}`}
          title={top ? intl.t(`input.${index}.tooltip`, top) : undefined}
        >
          <AntInput
            style={{ width: 70, ...sty }}
            min={1}
            suffix={sfx ? intl.t(`input.${index}.suffix`, sfx) : undefined}
            prefix={pfx ? intl.t(`input.${index}.prefix`, pfx) : undefined}
            defaultValue={value?.[index]}
            type={type}
            onChange={(e) => {
              if (type === 'number') {
                setTempValue(index, parseInt(e.target.value, 10) as any);
              } else {
                setTempValue(index, e.target.value as any);
              }
            }}
            onBlur={() => {
              if (autoSave) handleSave();
            }}
          />
        </Tooltip>,
      );
    }
    return inputs;
  };
  useEffect(() => {
    setTempValues(value ?? []);
  }, [value]);
  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      {isEditing ? (
        <Space.Compact size="small">
          {inputGroup()}
          <Button hidden={autoSave} onClick={handleSave}>
            <CheckCircleTwoTone />
          </Button>
          <Button
            onClick={() => {
              setIsEditing(false);
            }}
            hidden={autoSave}
          >
            <CloseCircleOutlined />
          </Button>
        </Space.Compact>
      ) : (
        <a
          style={{ marginLeft: 10 }}
          onClick={() => {
            setIsEditing(true);
          }}
        >
          {intl.t('modify', `Modify`)}
        </a>
      )}
    </div>
  );
};
