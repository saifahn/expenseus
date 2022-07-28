export const mainCategories = {
  unspecified: {
    en_US: 'Unspecified',
  },
  food: {
    en_US: 'Food & Drink',
  },
  daily: {
    en_US: 'Daily',
  },
  transport: {
    en_US: 'Transport',
  },
  utilities: {
    en_US: 'Daily',
  },
  entertainment: {
    en_US: 'Entertainment',
  },
  clothing: {
    en_US: 'Clothing',
  },
  beauty: {
    en_US: 'Beauty',
  },
  travel: {
    en_US: 'Travel',
  },
  home: {
    en_US: 'Home',
  },
  medical: {
    en_US: 'Medical',
  },
  social: {
    en_US: 'Social',
  },
  education: {
    en_US: 'Education',
  },
  housing: {
    en_US: 'Housing',
  },
  insurance: {
    en_US: 'Insurance',
  },
  car: {
    en_US: 'Car',
  },
  other: {
    en_US: 'Other',
  },
};

export type MainCategoryKey = keyof typeof mainCategories;

export const mainCategoryKeys = Object.keys(
  mainCategories,
) as MainCategoryKey[];

type SubCategories = {
  [k: string]: {
    mainCategory: MainCategoryKey;
    en_US: string;
  };
};

export const subcategories: SubCategories = {
  'unspecified.unspecified': {
    mainCategory: 'unspecified',
    en_US: 'Unspecified',
  },
  'food.food': {
    mainCategory: 'food',
    en_US: 'Food & Drink',
  },
  'food.groceries': {
    mainCategory: 'food',
    en_US: 'Groceries',
  },
  'food.eating-out': {
    mainCategory: 'food',
    en_US: 'Eating out',
  },
  'food.delivery': {
    mainCategory: 'food',
    en_US: 'Delivery',
  },
  'daily.daily': {
    mainCategory: 'daily',
    en_US: 'Daily goods',
  },
  'daily.consumables': {
    mainCategory: 'daily',
    en_US: 'Consumables',
  },
  'daily.children': {
    mainCategory: 'daily',
    en_US: 'Child-related',
  },
  'transport.transport': {
    mainCategory: 'transport',
    en_US: 'Transportation',
  },
  'transport.train': {
    mainCategory: 'transport',
    en_US: 'Train',
  },
  'transport.bus': {
    mainCategory: 'transport',
    en_US: 'Bus',
  },
  'transport.taxi': {
    mainCategory: 'transport',
    en_US: 'Taxi',
  },
  'utilities.utilities': {
    mainCategory: 'utilities',
    en_US: 'Utilities',
  },
  'utilities.electricity': {
    mainCategory: 'utilities',
    en_US: 'Electricity',
  },
  'utilities.water': {
    mainCategory: 'utilities',
    en_US: 'Water',
  },
  'utilities.gas': {
    mainCategory: 'utilities',
    en_US: 'Gas',
  },
  'utilities.internet': {
    mainCategory: 'utilities',
    en_US: 'Internet',
  },
  'utilities.phone': {
    mainCategory: 'utilities',
    en_US: 'Phone',
  },
  'entertainment.entertainment': {
    mainCategory: 'entertainment',
    en_US: 'Entertainment',
  },
  'entertainment.film-video': {
    mainCategory: 'entertainment',
    en_US: 'Film and Video',
  },
  'entertainment.books': {
    mainCategory: 'entertainment',
    en_US: 'Books',
  },
  'entertainment.music': {
    mainCategory: 'entertainment',
    en_US: 'Music',
  },
  'entertainment.leisure': {
    mainCategory: 'entertainment',
    en_US: 'Leisure',
  },
  'entertainment.hobbies': {
    mainCategory: 'entertainment',
    en_US: 'Hobbies',
  },
  'entertainment.subscriptions': {
    mainCategory: 'entertainment',
    en_US: 'Subscriptions',
  },
  'clothing.clothing': {
    mainCategory: 'clothing',
    en_US: 'Clothing',
  },
  'clothing.footwear': {
    mainCategory: 'clothing',
    en_US: 'Footwear',
  },
  'beauty.beauty': {
    mainCategory: 'beauty',
    en_US: 'Beauty',
  },
  'beauty.hair': {
    mainCategory: 'beauty',
    en_US: 'Hair',
  },
  'beauty.cosmetics': {
    mainCategory: 'beauty',
    en_US: 'Cosmetics',
  },
  'travel.travel': {
    mainCategory: 'travel',
    en_US: 'Travel',
  },
  'travel.transport': {
    mainCategory: 'travel',
    en_US: 'Transport',
  },
  'travel.accommodation': {
    mainCategory: 'travel',
    en_US: 'Accommodation',
  },
  'home.home': {
    mainCategory: 'home',
    en_US: 'Home',
  },
  'home.furniture': {
    mainCategory: 'home',
    en_US: 'Furniture',
  },
  'home.electronics': {
    mainCategory: 'home',
    en_US: 'Electronics',
  },
  'home.kitchen': {
    mainCategory: 'home',
    en_US: 'Kitchen',
  },
  'medical.medical': {
    mainCategory: 'medical',
    en_US: 'Medical',
  },
  'social.social': {
    mainCategory: 'social',
    en_US: 'Social',
  },
  'social.gifts': {
    mainCategory: 'social',
    en_US: 'Gifts',
  },
  'education.education': {
    mainCategory: 'education',
    en_US: 'Education',
  },
  'education.lessons': {
    mainCategory: 'education',
    en_US: 'Lessons',
  },
  'housing.housing': {
    mainCategory: 'housing',
    en_US: 'Housing',
  },
  'housing.rent': {
    mainCategory: 'housing',
    en_US: 'Rent',
  },
  'housing.mortgage': {
    mainCategory: 'housing',
    en_US: 'Mortgage',
  },
  'insurance.insurance': {
    mainCategory: 'insurance',
    en_US: 'Insurance',
  },
  'car.car': {
    mainCategory: 'car',
    en_US: 'Car',
  },
  'other.other': {
    mainCategory: 'other',
    en_US: 'Other',
  },
};

export type SubcategoryKey = Extract<keyof typeof subcategories, string>;

const categoryKeys = Object.keys(subcategories) as SubcategoryKey[];

type AllCategories = {
  [mainCat in MainCategoryKey]: Array<SubcategoryKey>;
};

// FIXME: this is kind of awful and hard to understand
/**
 * makes an AllCategories from the subcategories and main categories for making
 * tiered options with `optgroup`
 */
export const categories = categoryKeys.reduce((prev, key) => {
  const mainCategory = subcategories[key].mainCategory;
  if (!prev[mainCategory]) {
    prev[mainCategory] = [];
  }
  prev[mainCategory].push(key);
  return prev;
}, {} as AllCategories);

export const categoryNameFromKeyEN = (key: SubcategoryKey) =>
  subcategories[key].en_US;
