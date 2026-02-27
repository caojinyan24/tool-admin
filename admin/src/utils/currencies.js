// 全部币种配置
export const CURRENCIES = [
  { code: 'USD', name: '美元', symbol: '$' },
  { code: 'CNY', name: '人民币', symbol: '¥' },
  { code: 'EUR', name: '欧元', symbol: '€' },
  { code: 'JPY', name: '日元', symbol: '¥' },
  { code: 'GBP', name: '英镑', symbol: '£' },
  { code: 'AUD', name: '澳元', symbol: 'A$' },
  { code: 'CAD', name: '加元', symbol: 'C$' },
  { code: 'CHF', name: '瑞士法郎', symbol: 'CHF' },
  { code: 'HKD', name: '港币', symbol: 'HK$' },
  { code: 'SGD', name: '新加坡元', symbol: 'S$' },
  { code: 'NZD', name: '新西兰元', symbol: 'NZ$' },
  { code: 'SEK', name: '瑞典克朗', symbol: 'kr' },
  { code: 'NOK', name: '挪威克朗', symbol: 'kr' },
  { code: 'DKK', name: '丹麦克朗', symbol: 'kr' },
  { code: 'PLN', name: '波兰兹罗提', symbol: 'zł' },
  { code: 'CZK', name: '捷克克朗', symbol: 'Kč' },
  { code: 'HUF', name: '匈牙利福林', symbol: 'Ft' },
  { code: 'RUB', name: '俄罗斯卢布', symbol: '₽' },
  { code: 'BRL', name: '巴西雷亚尔', symbol: 'R$' },
  { code: 'MXN', name: '墨西哥比索', symbol: '$' },
  { code: 'ZAR', name: '南非兰特', symbol: 'R' },
  { code: 'INR', name: '印度卢比', symbol: '₹' },
  { code: 'KRW', name: '韩元', symbol: '₩' },
  { code: 'TWD', name: '新台币', symbol: 'NT$' },
  { code: 'THB', name: '泰铢', symbol: '฿' },
  { code: 'MYR', name: '马来西亚林吉特', symbol: 'RM' },
  { code: 'IDR', name: '印尼盾', symbol: 'Rp' },
  { code: 'PHP', name: '菲律宾比索', symbol: '₱' },
  { code: 'VND', name: '越南盾', symbol: '₫' },
  { code: 'AED', name: '阿联酋迪拉姆', symbol: 'د.إ' },
  { code: 'SAR', name: '沙特里亚尔', symbol: 'ر.س' },
  { code: 'ILS', name: '以色列新谢克尔', symbol: '₪' },
  { code: 'TRY', name: '土耳其里拉', symbol: '₺' },
  { code: 'CLP', name: '智利比索', symbol: '$' },
  { code: 'COP', name: '哥伦比亚比索', symbol: '$' },
  { code: 'ARS', name: '阿根廷比索', symbol: '$' },
  { code: 'PEN', name: '秘鲁索尔', symbol: 'S/' },
  { code: 'EGP', name: '埃及镑', symbol: 'E£' },
  { code: 'NGN', name: '尼日利亚奈拉', symbol: '₦' },
  { code: 'MAD', name: '摩洛哥迪拉姆', symbol: 'د.م.' },
  { code: 'QAR', name: '卡塔尔里亚尔', symbol: 'ر.ق' },
  { code: 'KWD', name: '科威特第纳尔', symbol: 'د.ك' },
  { code: 'BHD', name: '巴林第纳尔', symbol: '.د.ب' },
  { code: 'OMR', name: '阿曼里亚尔', symbol: 'ر.ع.' },
  { code: 'JOD', name: '约旦第纳尔', symbol: 'د.أ' },
  { code: 'LBP', name: '黎巴嫩磅', symbol: 'ل.ل' },
  { code: 'PKR', name: '巴基斯坦卢比', symbol: 'Rs' },
  { code: 'LKR', name: '斯里兰卡卢比', symbol: 'Rs' },
  { code: 'NPR', name: '尼泊尔卢比', symbol: 'Rs' },
  { code: 'BDT', name: '孟加拉塔卡', symbol: '৳' },
  { code: 'MMK', name: '缅甸元', symbol: 'Ks' },
  { code: 'KHR', name: '柬埔寨瑞尔', symbol: '៛' },
  { code: 'LAK', name: '老挝基普', symbol: '₭' },
  { code: 'MVR', name: '马尔代夫卢菲亚', symbol: '.ރ' },
  { code: 'BTN', name: '不丹努尔特鲁姆', symbol: 'Nu.' },
  { code: 'AFN', name: '阿富汗尼', symbol: '؋' },
  { code: 'IRR', name: '伊朗里亚尔', symbol: '﷼' },
  { code: 'IQD', name: '伊拉克第纳尔', symbol: 'ع.د' },
  { code: 'YER', name: '也门里亚尔', symbol: 'ر.ي' },
  { code: 'SYR', name: '叙利亚磅', symbol: 'ل.س' },
  { code: 'LYD', name: '利比亚第纳尔', symbol: 'د.ل' },
  { code: 'DZD', name: '阿尔及利亚第纳尔', symbol: 'د.ج' },
  { code: 'TND', name: '突尼斯第纳尔', symbol: 'د.ت' },
  { code: 'GHS', name: '加纳塞地', symbol: '₵' },
  { code: 'KES', name: '肯尼亚先令', symbol: 'Sh' },
  { code: 'UGX', name: '乌干达先令', symbol: 'Sh' },
  { code: 'TZS', name: '坦桑尼亚先令', symbol: 'Sh' },
  { code: 'RWF', name: '卢旺达法郎', symbol: 'Fr' },
  { code: 'BIF', name: '布隆迪法郎', symbol: 'Fr' },
  { code: 'GMD', name: '冈比亚达拉西', symbol: 'D' },
  { code: 'SLL', name: '塞拉利昂利昂', symbol: 'Le' },
  { code: 'LRD', name: '利比里亚元', symbol: '$' },
  { code: 'GNF', name: '几内亚法郎', symbol: 'Fr' },
  { code: 'MLI', name: '马里法郎', symbol: 'Fr' },
  { code: 'CFA', name: '西非法郎', symbol: 'Fr' },
  { code: 'XAF', name: '中非法郎', symbol: 'Fr' },
  { code: 'XOF', name: '西非法郎', symbol: 'Fr' },
  { code: 'MGA', name: '马达加斯加阿里亚里', symbol: 'Ar' },
  { code: 'MUR', name: '毛里求斯卢比', symbol: 'Rs' },
  { code: 'SCR', name: '塞舌尔卢比', symbol: 'Rs' },
  { code: 'MWK', name: '马拉维克瓦查', symbol: 'MK' },
  { code: 'ZMW', name: '赞比亚克瓦查', symbol: 'ZK' },
  { code: 'BWP', name: '博茨瓦纳普拉', symbol: 'P' },
  { code: 'NAD', name: '纳米比亚元', symbol: '$' },
  { code: 'SZL', name: '斯威士兰里兰吉尼', symbol: 'L' },
  { code: 'LSL', name: '莱索托洛蒂', symbol: 'L' },
  { code: 'AZN', name: '阿塞拜疆马纳特', symbol: '₼' },
  { code: 'AMD', name: '亚美尼亚德拉姆', symbol: '֏' },
  { code: 'GEL', name: '格鲁吉亚拉里', symbol: '₾' },
  { code: 'KZT', name: '哈萨克斯坦坚戈', symbol: '₸' },
  { code: 'UZS', name: '乌兹别克斯坦苏姆', symbol: 'so\'m' },
  { code: 'TJS', name: '塔吉克斯坦索莫尼', symbol: 'смн' },
  { code: 'KGS', name: '吉尔吉斯斯坦索姆', symbol: 'сом' },
  { code: 'MNT', name: '蒙古图格里克', symbol: '₮' },
  { code: 'KWD', name: '科威特第纳尔', symbol: 'د.ك' },
  { code: 'BHD', name: '巴林第纳尔', symbol: '.د.ب' },
  { code: 'QAR', name: '卡塔尔里亚尔', symbol: 'ر.ق' },
  { code: 'OMR', name: '阿曼里亚尔', symbol: 'ر.ع.' }
];

// 状态选项
export const STATUS_OPTIONS = [
  { value: 'active', label: '激活' },
  { value: 'inactive', label: '停用' }
];

// 数据源选项
export const SOURCE_OPTIONS = [
  { value: 'central_bank', label: '央行' },
  { value: 'commercial', label: '商业银行' },
  { value: 'foreign_exchange', label: '外汇市场' },
  { value: 'international', label: '国际组织' }
];

// 获取币种名称
export const getCurrencyName = (code) => {
  const currency = CURRENCIES.find(c => c.code === code);
  return currency ? `${currency.name} (${currency.code})` : code;
};

// 获取币种符号
export const getCurrencySymbol = (code) => {
  const currency = CURRENCIES.find(c => c.code === code);
  return currency ? currency.symbol : code;
};