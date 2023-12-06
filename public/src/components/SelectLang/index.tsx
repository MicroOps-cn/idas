import { SelectLang as UmiSelectLang, getAllLocales as UmiGetAllLocales } from '@umijs/max';

const disableLangs = ['bn-BD', 'id-ID', 'fa-IR'];

export const allLocales = () => {
  return UmiGetAllLocales().filter((locale) => !disableLangs.includes(locale));
};

const SelectLang: typeof UmiSelectLang = (props) => {
  return (
    <UmiSelectLang
      postLocalesData={(locales) => {
        return locales.filter((locale) => !disableLangs.includes(locale.lang));
      }}
      {...props}
    />
  );
};
export default SelectLang;
