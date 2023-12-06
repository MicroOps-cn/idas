import type { SwitchProps as AntSwitchProps } from 'antd';
import { Switch as AntSwitch } from 'antd';
import React, { useEffect, useState } from 'react';

interface SwitchProps extends AntSwitchProps {
  request: () => any;
}

const Switch: React.FC<SwitchProps> = ({ request, checked, onChange, ...props }) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [value, setValue] = useState<boolean | undefined>(checked);
  useEffect(() => {
    setValue(checked);
  }, [checked]);

  return (
    <AntSwitch
      onChange={async (status, event) => {
        try {
          onChange?.(status, event);
          setLoading(true);
          await request();
          setValue(status);
        } catch {
          setValue(!status);
        } finally {
          setLoading(false);
        }
      }}
      loading={loading}
      checked={value}
      {...props}
    />
  );
};

export default Switch;
